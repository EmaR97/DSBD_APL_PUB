import time
from datetime import datetime
from typing import List

from mongoengine import (EmbeddedDocument, StringField, BooleanField, EmbeddedDocumentField, LongField, IntField, )

from mongo.base_document import BaseDocument


class Subscription(BaseDocument):
    meta = {'collection': 'Subscription'}

    class Id(EmbeddedDocument):
        chat_id = LongField()
        cam_id = StringField()

    class Specs(EmbeddedDocument):
        interval = IntField(0, default=0)  # minutes
        period = StringField()  # 22:00-06:00

    _id = EmbeddedDocumentField(Id, primary_key=True)
    is_subscribed = BooleanField(default=True)
    specs = EmbeddedDocumentField(Specs)
    last_notification = LongField(default=0)  # time of the last notification send

    @classmethod
    def get_by_chat_and_cam_id(cls, chat_id, cam_id) -> 'Subscription':
        _id = cls.Id(chat_id=chat_id, cam_id=cam_id)
        return cls.get_by_id(_id)

    @classmethod
    def get_by_cam_id(cls, cam_id) -> List['Subscription']:
        return cls.objects(_id__cam_id=cam_id)

    def has_interval_elapsed(self):
        if self.specs.interval > 0:
            # Check if the current time is within the specified interval
            current_time = time.time()/60
            return current_time - self.last_notification > self.specs.interval
        return True

    def is_within_period(self):
        if self.specs.period:
            # Check if the current time is within the specified period
            current_time = datetime.now().time()
            start_time, end_time = map(lambda x: datetime.strptime(x, "%H:%M").time(), self.specs.period.split('-'))
            return start_time <= current_time <= end_time
        return True

    def to_notify(self):
        to_notify = self.is_subscribed and self.has_interval_elapsed() and self.is_within_period()
        if to_notify:
            self.last_notification = time.time()/60
            self.save()
        return to_notify
