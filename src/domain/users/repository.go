package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	AddUser(entity *User) (*User, error)
	GetUser(email string) (entity *User, err error)
	GetUserByID(id primitive.ObjectID) (entity *User, err error)
	GetAllUser() (entity []*User, err error)
	UpdateUser(email, authToken, refreshToken, name, photo string, swiftToken []string) (err error)
	UpdateSwiftToken(email string, swiftToken []string) (err error)
	UpdateFlagArchive(email string, archive bool) (err error)
}
