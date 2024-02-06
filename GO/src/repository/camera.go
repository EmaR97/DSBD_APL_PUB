package repository

import (
	"CamMonitoring/src/entity"
	"CamMonitoring/src/utility"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type CameraRepository struct {
	MongoDBRepository[entity.Camera]
}

// GetAllByUser retrieves cameras an existing user in MongoDB.
func (r *CameraRepository) GetAllByUser(userId string) ([]entity.Camera, error) {
	var cameras []entity.Camera
	log.Printf(userId)
	cursor, err := r.collection.Find(context.Background(), bson.M{"user_id": userId})
	if err != nil {
		// Handle error
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			utility.ErrorLog().Printf("Error closing cursor: %v", err)
		}
	}(cursor, context.Background())

	if err := cursor.All(context.Background(), &cameras); err != nil {
		// Handle error
	}
	return cameras, err
}

// GetAllCamIdsByUser retrieves cameras for a given user from MongoDB using the callback function.
func (r *CameraRepository) GetAllCamIdsByUser(userId string) ([]string, error) {
	cameras, err := r.GetAllByUser(userId)
	if err != nil {
		return nil, fmt.Errorf("GetAllByUserCallback not set")
	}
	var camIds []string
	for _, camera := range cameras {
		camIds = append(camIds, camera.Id)
	}
	return camIds, err
}
