package models

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User : This defines the user structure.
type User struct {
	ID           bson.ObjectId `json:"_id" bson:"_id"`
	PhoneNumber  string        `bson:"phone_number,omitempty" json:"phone_number"`
	Active       int           `bson:"active" json:"active"`
	LastLogin    time.Time     `bson:"last_login" json:"-"`
	CreatedAt    time.Time     `bson:"created_at" json:"-"`
	FullName     string        `bson:"full_name" json:"full_name"`
	DB           *mgo.Database `bson:"-" json:"-"`
	SwiftToken   string        `bson:"swift_token" json:"-"`
	RefreshToken string        `bson:"refresh_token" json:"-"`
	AuthToken    string        `bson:"auth_token" json:"-"`
	Email        string        `bson:"email" json:"email"`
	Self         string        `bson:"-" json:"-"`
}

// GetUser : get user profile.
func (u *User) GetUser() error {
	db := u.DB
	authToken := u.AuthToken
	refreshToken := u.RefreshToken
	if err := u.DB.C("users").Find(
		bson.M{"email": u.Email},
	).One(&u); err == mgo.ErrNotFound {
		u.ID = bson.NewObjectId()
		u.LastLogin = time.Now().UTC()
		u.CreatedAt = time.Now().UTC()
		u.Active = 1
		if err := u.DB.C("users").Insert(u); err != nil {
			return err
		}
	}
	if u.Active == 0 {
		return errors.New("user not allowed to access the system")
	}
	u.DB = db
	if _, err := u.DB.C("users").Find(bson.M{"email": u.Email}).Apply(mgo.Change{
		Update: bson.M{"$set": bson.M{
			"last_login":    time.Now().UTC(),
			"auth_token":    authToken,
			"refresh_token": refreshToken,
		}},
		ReturnNew: true,
	}, nil); err != nil {
		return err
	}
	return nil
}
