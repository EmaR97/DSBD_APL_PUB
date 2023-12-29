from mongoengine import StringField, DoesNotExist

from mongo.base_document import BaseDocument


class Cam(BaseDocument):
    meta = {'collection': 'Camera'}

    _id = StringField(primary_key=True)
    user_id = StringField(required=True)

    @classmethod
    def get_by_user_id(cls, user_id):
        try:
            cams = cls.objects(user_id=user_id)
            return list(cams)
        except DoesNotExist:
            return None
