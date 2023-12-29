import json
import logging

import grpc_
from kafka_ import NotificationConsumer
from notification_bot.messagesender import NotificationSender


def main():
    # Load configuration from a JSON file
    with open('config.json', 'r') as file:
        config = json.load(file)
    # Enable logging
    logging.basicConfig(format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.DEBUG)
    logging.getLogger("httpx").setLevel(logging.DEBUG)

    logging.info("Loading configuration from 'config.json'...")
    # Initialize the Telegram notification sender
    logging.info("Configuring notification sender...")
    def get_chat_ids(chat_id): return grpc_.get_chat_ids(chat_id, config['grpc']['get_chat_ids'])
    notification_sender = NotificationSender(config['telegram']['token'], get_chat_ids)
    # Initialize the Kafka message processor
    logging.info("Initializing Kafka message processor...")
    processor = NotificationConsumer(bootstrap_servers=config['kafka']['bootstrap_servers'],
                                     group_id=config['kafka']['group_id'],
                                     auto_offset_reset=config['kafka']['auto_offset_reset'],
                                     topic=config['kafka']['topic'], message_handler=notification_sender.send_message)
    # Configure Kafka consumer
    logging.info("Configuring Kafka consumer...")
    processor.configure_consumer()
    # Start listening to Kafka messages in a separate thread
    logging.info("Start listening to Kafka messages ...")
    processor.consume_message()
    logging.info("Terminated")


if __name__ == '__main__':
    main()
