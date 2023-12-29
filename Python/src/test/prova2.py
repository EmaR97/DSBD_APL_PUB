from mongoengine import connect, Document, EmbeddedDocument, StringField, IntField, EmbeddedDocumentField

# Connect to MongoDB
connect('my_database', host='localhost', port=27017)


# Define an embedded document with primary_key=True
class Address(EmbeddedDocument):
    street = StringField()
    city = StringField()
    state = StringField()
    zip_code = StringField()


# Define a MongoDB document model with an embedded document as _id
class Person(Document):
    _id = EmbeddedDocumentField(Address, primary_key=True)
    name = StringField(required=True)
    age = IntField(required=True)


# Create a person with an embedded address as _id
john_address = Address(street="123 Main St", city="Anytown", state="CA", zip_code="12345")
john = Person(_id=john_address, name="John Doe", age=25)
john.save()

# Query the database using the embedded address as _id
retrieved_john = Person.objects(_id=john_address).first()
print("Person Information:")
print(f"Name: {retrieved_john.name}")
print(f"Age: {retrieved_john.age}")
