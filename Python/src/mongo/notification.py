from mongoengine import Document, StringField, IntField, DoesNotExist


# Define the MongoEngine model
class Notification(Document):
    # Autogenerated unique ID
    # _id field is autogenerated by default if not specified explicitly
    # You can also specify your own custom ID field if needed
    cam_id = StringField(required=True)
    timestamp = IntField(required=True)
    link = StringField(required=True)

    def __str__(self):
        return f"Notification(cam_id={self.cam_id}, timestamp={self.timestamp}, link={self.link})"


# Save incoming message and return the auto-generated ID
def save_incoming_message(cam_id, timestamp, link):
    notification = Notification(cam_id=cam_id, timestamp=timestamp, link=link)
    notification.save()
    return str(notification.id)


# Retrieve a document by its auto-generated ID
def get_notification_by_id(notification_id):
    try:
        notification = Notification.objects.get(id=notification_id)
        return notification
    except DoesNotExist:
        return None  # Return None if the document does not exist