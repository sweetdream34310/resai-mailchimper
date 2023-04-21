package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	DB        *mgo.Database `bson:"-" json:"-"`
	UserID    bson.ObjectId `bson:"user_id" json:"-"`
	MessageID string        `bson:"message_id,omitempty" json:"message_id"`
	Subject   string        `bson:"subject,omitempty" json:"subject,omitempty"`
	From      *string       `bson:"from,omitempty" json:"from,omitempty"`
	To        *string       `bson:"to,omitempty" json:"to,omitempty"`
	Received  *time.Time    `bson:"received_at,omitempty" json:"received_at,omitempty"`
	Sent      *time.Time    `bson:"sent,omitempty" json:"sent,omitempty"`
	BodyText  string        `bson:"body_text,omitempty" json:"body_text,omitempty"`
	BodyHTML  string        `bson:"body_html,omitempty" json:"body_html,omitempty"`
	IsRead    bool          `bson:"is_read,omitempty" json:"is_read"`
	MimeFlow  []string      `bson:"mime_flow,omitempty" json:"mime_flow"`
}

func (m *Message) Find(collection string, selector interface{}) (interface{}, error) {
	if err := m.DB.C(collection).
		Find(selector).
		One(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Message) Insert(collection string, selector interface{}) error {
	if _, err := m.DB.C(collection).
		Upsert(selector, &m); err != nil {
		return err
	}
	return nil
}

func (m *Message) GetInbox(collection string, skip int) ([]Message, error) {
	var data []Message
	pipeline := []bson.M{}
	pipeline = []bson.M{
		{"$match": bson.M{"user_id": m.UserID}},
		{"$sort": bson.D{{"received_at", -1}}},
		{"$skip": skip * 10},
		{"$limit": 10}}
	err := m.DB.C(collection).Pipe(pipeline).All(&data)
	return data, err
}
