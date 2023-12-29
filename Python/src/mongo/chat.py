from mongoengine import StringField, IntField, LongField

from mongo.base_document import BaseDocument


class Chat(BaseDocument):
    meta = {'collection': 'Chat'}

    _id = LongField(primary_key=True)
    user_id = StringField(required=True)
