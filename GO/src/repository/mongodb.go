package repository

import (
	"CamMonitoring/src/_interface"
	"CamMonitoring/src/service/storage"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "golang.org/x/exp/constraints"
	"reflect"
)

// MongoDBRepository is an implementation of Repository that interacts with MongoDB.
type MongoDBRepository[Entity _interface.Entity] struct {
	database   *mongo.Database
	collection *mongo.Collection
}

// NewMongoDBRepository creates a new instance of MongoDBRepository.
func NewMongoDBRepository[Entity _interface.Entity](
	databaseService *storage.MongoDBService, databaseName string, collectionName *string,
) *MongoDBRepository[Entity] {
	if collectionName == nil {
		entityType := reflect.TypeOf((*Entity)(nil)).Elem()
		typeName := entityType.Name()
		collectionName = &typeName
	}
	database := databaseService.GetDatabase(databaseName)
	collection := database.Collection(*collectionName)
	return &MongoDBRepository[Entity]{
		database:   database,
		collection: collection,
	}
}

// GetByID retrieves a user by ID from MongoDB.
func (r *MongoDBRepository[Entity]) GetByID(id string) (Entity, error) {
	filter := bson.M{"_id": id}
	var entity Entity
	err := r.collection.FindOne(context.Background(), filter).Decode(&entity)
	if err != nil {
		return entity, err
	}
	return entity, nil
}

// Add adds a new user to MongoDB.
func (r *MongoDBRepository[Entity]) Add(entity Entity) error {
	//println(entity)
	_, err := r.collection.InsertOne(context.Background(), entity)
	return err
}

// Update updates an existing user in MongoDB.
func (r *MongoDBRepository[Entity]) Update(entity Entity) error {
	filter := bson.M{"_id": entity.ID()}
	update := bson.M{"$set": entity}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

// GetAll retrieves all users from MongoDB.
func (r *MongoDBRepository[Entity]) GetAll() ([]Entity, error) {
	var entities []Entity
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &entities)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *MongoDBRepository[Entity]) Delete(id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(context.Background(), filter)
	return err

}
