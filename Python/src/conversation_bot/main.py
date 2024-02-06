import logging
import threading
import time
from random import random

import google
from google.protobuf.json_format import Parse, MessageToJson
from telegram import Update
from telegram.ext import Application

from conversation_bot.cam_sub_manager import CamSubHandler
from conversation_bot.login_manager import LoginHandler
from http_.auth_request import AuthClient
from kafka_ import MyConsumer
from kafka_.producer import MyProducer
from message import cam_pb2, subscription_pb2
from mongo import Subscription
from utility import enable_log_load_conf, connect_to_mongo_db


def main() -> None:
    config = enable_log_load_conf()

    logging.info("Configuring Kafka response producer...")
    response_producer = MyProducer(bootstrap_servers=config['kafka']['bootstrap_servers'])

    logging.info("Configuring Kafka request consumer...")
    request_consumer = MyConsumer(bootstrap_servers=config['kafka']['bootstrap_servers'], topic_handlers={
        "user_to_notify": lambda msg: get_user_to_notify(msg, response_producer)}, group_id=config['kafka']['group_id'],
                                  auto_offset_reset=config['kafka']['auto_offset_reset'])
    # Start listening to Kafka messages in a separate thread
    logging.info("Start listening to Kafka messages ...")
    processor_thread = threading.Thread(target=request_consumer.loop_consume_message)
    processor_thread.start()

    # Database connection
    connect_to_mongo_db(dbname_=config["mongo"]["dbname"], host_=config["mongo"]["host"],
                        port=int(config["mongo"]["port"]), username_=config["mongo"]["username"],
                        password_=config["mongo"]["password"])

    # Initialize the Kafka message consumer
    logging.info("Initializing Kafka message consumer...")
    response_consumer = MyConsumer(bootstrap_servers=config['kafka']['bootstrap_servers'],
                                   topic_handlers={"conversation_bot": None}, group_id=config['kafka']['group_id'],
                                   auto_offset_reset=config['kafka']['auto_offset_reset'])

    # Set up the Telegram bot application
    logging.info("Setting up conversation handlers...")
    logging.info("Configuring subscription manager...")
    sub_manager = Application.builder().token(config['telegram']['token']).build()
    # Handlers
    cam_sub_handler = CamSubHandler(lambda userId: get_user_id(userId, response_producer, response_consumer))
    sub_conv_handler = cam_sub_handler.cam_sub_conv_handler()
    auth_client = AuthClient(config['auth'])
    login_handler = LoginHandler(sub_conv_handler, auth_client.login)
    login_conv_handler = login_handler.login_conv_handler()
    sub_manager.add_handler(login_conv_handler)

    # Start polling for incoming Telegram updates
    logging.info("Start polling for Telegram updates ...")
    sub_manager.run_polling(allowed_updates=Update.ALL_TYPES)

    logging.info("Shutting down ...")
    logging.info("Stopping Grpc listener ...")
    request_consumer.stop()
    response_consumer.close()
    processor_thread.join()  # Wait for the server thread to finish
    request_consumer.close()
    logging.info("Terminated")


def get_user_to_notify(msg, producer: MyProducer):
    request = Parse(msg, subscription_pb2.CamIdRequest())  # return
    logging.info("request.cam_id:" + request.cam_id)
    subs = Subscription.get_by_cam_id(request.cam_id)
    chat_ids = [subscription.id.chat_id for subscription in subs if subscription.to_notify()]
    print(chat_ids)
    response = subscription_pb2.ChatIdsResponse(chat_ids=chat_ids, request_id=request.request_id)
    producer.produce_message(request.response_topic, MessageToJson(response))


def get_user_id(user_id, producer: MyProducer, consumer: MyConsumer, request_topic="cam_ids",
                response_topic="conversation_bot"):
    request_id = str(random())
    request = cam_pb2.UserIdRequest(user_id=user_id, request_id=request_id, response_topic=response_topic)
    producer.produce_message(request_topic, request.SerializeToString())
    start_time = time.time()
    while time.time() - start_time < 5:
        message = consumer.listen_for_response(5)
        if not message:
            return []
        try:
            response = Parse(message, cam_pb2.CamIdsResponse())  # return
            if request_id == response.request_id:
                return response.cam_ids
        except google.protobuf.message.DecodeError as e:
            logging.error(f"Error decoding message: {e}")
            return []


if __name__ == '__main__':
    main()
