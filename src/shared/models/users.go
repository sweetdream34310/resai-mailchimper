package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSession struct {
	UserID       primitive.ObjectID `json:"user_id"`
	Email        string             `json:"email"`
	RefreshToken string             `json:"refresh_token"`
	AuthToken    string             `json:"auth_token"`
	Name         string             `json:"name"`
	Photo        string             `json:"photo"`
	SwiftToken   string             `json:"swift_token"`
	DevTokenKey  bool               `json:"dev_token_key"`
}
