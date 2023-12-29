from mongoengine import StringField

from mongo.base_document import BaseDocument


class User(BaseDocument):
    meta = {'collection': 'User'}

    _id = StringField(primary_key=True)
    password = StringField(required=True)
