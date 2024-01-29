import logging

import google.protobuf
import telegram
from google.protobuf.json_format import Parse

from message import notification_pb2


class NotificationSender:
    def __init__(self, bot_token_, chat_ids_callback):
        self.bot = telegram.Bot(bot_token_)
        self.get_chat_ids = chat_ids_callback

    async def send_message(self, msg, ):
        m = {'key': msg.key(), 'value': msg.value(), 'partition': msg.partition(), 'offset': msg.offset()}
        logging.debug(f'Received message: {m}')
        try:
            n = Parse(msg.value(), notification_pb2.Notification())  # return
        except google.protobuf.message.DecodeError as e:
            logging.error(f"Error decoding message: {e}")
            return
        try:
            chat_ids = self.get_chat_ids(n.cam_id)
        except Exception as e:
            logging.error(f"Error get chat ids: {e}")
            return
        for chat_id in chat_ids if chat_ids else []:
            logging.debug(f'Sending message to chat: {chat_id}')
            try:
                await self.bot.send_message(chat_id=chat_id, text=str(n))
            except Exception as e:
                logging.error(f'Error sending message to chat: {chat_id}: {e}')
