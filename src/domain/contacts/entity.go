package contacts

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contacts struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"-"`
	Name   string             `bson:"name" json:"name"`
	Email  string             `bson:"email" json:"email"`
}
