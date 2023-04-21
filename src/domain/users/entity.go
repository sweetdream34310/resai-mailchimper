package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	PhoneNumber  string             `bson:"phone_number,omitempty" json:"phone_number"`
	Active       int                `bson:"active" json:"active"`
	AwayMode     bool               `bson:"away_mode" json:"away_mode"`
	Archive      bool               `bson:"archive" json:"archive"`
	LastLogin    time.Time          `bson:"last_login" json:"-"`
	CreatedAt    time.Time          `bson:"created_at" json:"-"`
	FullName     string             `bson:"full_name" json:"full_name"`
	SwiftToken   []string           `bson:"swift_token" json:"-"`
	RefreshToken string             `bson:"refresh_token" json:"-"`
	AuthToken    string             `bson:"auth_token" json:"-"`
	Email        string             `bson:"email" json:"email"`
	Name         string             `bson:"name" json:"name"`
	Photo        string             `bson:"photo" json:"photo"`
}
