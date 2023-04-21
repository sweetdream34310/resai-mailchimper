package messages

import "gopkg.in/mgo.v2/bson"

type Repository interface {
	FindInbox(selector interface{}) (*Message, error)
	InsertInbox(*Message) error
	DeleteArchiveByEmail(email string) (err error)
	InsertSentBox(*Message) error
	GetInbox(skip int, userID bson.ObjectId) ([]Message, error)
}
