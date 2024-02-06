import asyncio
import logging

import telegram
from google.protobuf.json_format import Parse, ParseError

from message import notification_pb2, subscription_pb2
from mongo.notification import save_incoming_message, get_notification_by_id


class NotificationSender:
    def __init__(self, bot_token_, chat_ids_callback):
        self.bot = telegram.Bot(bot_token_)
        self.get_chat_ids = chat_ids_callback

    def handle_msg(self, msg):
        try:
            notification = Parse(msg, notification_pb2.Notification())  # return
            self.receive_notification(notification)
            return
        except ParseError:
            pass
        try:
            response = Parse(msg, subscription_pb2.ChatIdsResponse())  # return
            self.send_notifiaction(response)
        except ParseError as e:
            logging.error(f"Error decoding message: {e}")

    def receive_notification(self, notification):
        notification_id = save_incoming_message(notification.cam_id, notification.timestamp, notification.link)
        self.get_chat_ids(notification.cam_id, notification_id)

    def send_notifiaction(self, response):
        notification = get_notification_by_id(response.request_id)
        logging.debug(f'Sending chat to notify: {response.chat_ids}')
        for chat_id in response.chat_ids if response.chat_ids else []:
            logging.debug(f'Sending message to chat: {chat_id}')
            try:
                asyncio.get_event_loop().run_until_complete(
                    self.bot.send_message(chat_id=chat_id, text=str(notification)))
            except Exception as e:
                logging.error(f'Error sending message to chat: {chat_id}: {e}')
