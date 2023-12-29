from mongoengine import Document, DoesNotExist, StringField


class BaseDocument(Document):
    meta = {'abstract': True}
    _id = StringField(primary_key=True)

    @classmethod
    def get_by_id(cls, _id):
        try:
            obj = cls.objects.get(_id=_id)
            return obj
        except DoesNotExist:
            return None
