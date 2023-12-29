package repository

import (
	"CamMonitoring/src/entity"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type TokenRepository struct {
	MongoDBRepository[entity.Token]
}

// GetByUser updates an existing user in MongoDB.
func (r *TokenRepository) GetByUser(userId string) (entity.Token, error) {
	filter := bson.M{"user_id": userId}
	var token entity.Token
	err := r.collection.FindOne(context.Background(), filter).Decode(&token)
	if err != nil {
		return token, err
	}
	return token, nil
}

// RemoveToken removes the token from the database
func (r *TokenRepository) RemoveToken(userID string) error {
	// Implement the logic to remove the token from the database
	return r.Delete(userID)
}

// StoreToken stores the token in the database
func (r *TokenRepository) StoreToken(tokenEntity entity.Token) error {
	oldToken, err := r.GetByUser(tokenEntity.UserID)
	if err == nil {
		if err := r.Delete(oldToken.Id); err != nil {
			return fmt.Errorf("error deleting old token: %s", err)
		}
	}
	return r.Add(tokenEntity)
}

// IsValidToken checks if the token is valid by querying the database
func (r *TokenRepository) IsValidToken(token string) (bool, error) {
	// Query the database to check if the token is present and valid
	// Implement the logic based on your database structure and requirements
	_, err := r.GetByID(token)
	if err != nil {
		return false, fmt.Errorf("invalid token: %s", err)
	}
	return true, nil
}
