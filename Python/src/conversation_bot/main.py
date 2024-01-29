import json
import logging
import os
import threading

from mongoengine import connect
from telegram import Update
from telegram.ext import Application

from conversation_bot.cam_sub_manager import CamSubHandler
from conversation_bot.login_manager import LoginHandler
from grpc_ import run_grpc_server, get_cam_ids, SubscriptionServiceServicer
from http_.auth_request import AuthClient
from mongo import Subscription


def main() -> None:
    # Load configuration from a JSON file
    config_path = os.environ.get('CONFIG_PATH', 'config.json')
    with open(config_path, 'r') as file:
        config = json.load(file)

    # Enable logging
    logging.basicConfig(format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.DEBUG)
    logging.getLogger("httpx").setLevel(logging.DEBUG)

    logging.info("Loading configuration from 'config.json'...")

    # Create a thread for the gRPC server
    logging.info("Starting gRPC server thread...")
    stop_event = threading.Event()
    SubscriptionServiceServicer.callback = lambda cam_id: [subscription.id.chat_id for subscription in
                                                           Subscription.get_by_cam_id(cam_id) if
                                                           subscription.to_notify()]
    grpc_server_thread = threading.Thread(target=run_grpc_server, args=(config['grpc']['get_chat_ids'], stop_event))
    grpc_server_thread.start()
    logging.info("gRPC server thread started. " + config['grpc']['get_chat_ids'])

    # Set up the Telegram bot application
    logging.info("Configuring subscription manager...")
    sub_manager = Application.builder().token(config['telegram']['token']).build()

    # Database connection
    logging.info("Connecting to the database...")
    connect(config["mongo"]["dbname"], host=config["mongo"]["host"], port=int(config["mongo"]["port"]),
            username=config["mongo"]["username"], password=config["mongo"]["password"])
    logging.info("Database connection established.")

    # Handlers
    logging.info("Setting up conversation handlers...")
    cam_sub_handler = CamSubHandler(lambda user_id: get_cam_ids(user_id, config['grpc']['get_cam_ids']))
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
    stop_event.set()
    grpc_server_thread.join()  # Wait for the server thread to finish
    logging.info("Terminated")


if __name__ == '__main__':
    main()
