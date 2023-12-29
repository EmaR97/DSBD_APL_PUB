package entity

import (
	"encoding/base64"
	"golang.org/x/exp/rand"
)

type Token struct {
	Id     string `json:"_id" bson:"_id"`
	UserID string `json:"user_id" bson:"user_id"`
}

func (u Token) ID() string {
	return u.Id
}

func NewToken(userID string, tokenLength int) (*Token, error) {
	token, err := generateToken(tokenLength)
	if err != nil {
		return nil, err
	}
	return &Token{
		UserID: userID,
		Id:     token, // Add other fields if necessary (e.g., expiration time)
	}, nil
}
func generateToken(tokenLength int) (string, error) {
	tokenBytes := make([]byte, tokenLength)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token, nil
}
