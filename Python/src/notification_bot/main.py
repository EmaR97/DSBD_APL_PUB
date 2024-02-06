import logging

from google.protobuf.json_format import MessageToJson

from kafka_ import MyConsumer, MyProducer
from message import subscription_pb2
from notification_bot.messagesender import NotificationSender
from utility import enable_log_load_conf, connect_to_mongo_db


def main():
    config = enable_log_load_conf()
    # Database connection
    connect_to_mongo_db(dbname_=config["mongo"]["dbname"], host_=config["mongo"]["host"],
                        port=int(config["mongo"]["port"]), username_=config["mongo"]["username"],
                        password_=config["mongo"]["password"])

    logging.info("Initializing query sender and receiver over Kafka...")
    request_producer = MyProducer(bootstrap_servers=config['kafka']['bootstrap_servers'])

    def get_chat_ids_request(cam_id: str, notification_id: str):
        request = subscription_pb2.CamIdRequest(cam_id=cam_id, request_id=notification_id,
                                                response_topic=config['kafka']['topic'])
        request_encoded = MessageToJson(request)
        request_producer.produce_message(config['kafka']['request_topic'], request_encoded, cam_id)

    notification_sender = NotificationSender(config['telegram']['token'], get_chat_ids_request)

    logging.info("Initializing Kafka notification consumer...")
    notification_consumer = MyConsumer(bootstrap_servers=config['kafka']['bootstrap_servers'],
                                       topic_handlers={config['kafka']['topic']: notification_sender.handle_msg},
                                       group_id=config['kafka']['group_id'],
                                       auto_offset_reset=config['kafka']['auto_offset_reset'])
    try:
        logging.info("Start listening for notification on Kafka...")
        notification_consumer.loop_consume_message()
    finally:
        logging.info("Closing connection")
        notification_consumer.close()


if __name__ == '__main__':
    main()
