package mongo

import (
	dbMongo "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	domainMsg "github.com/cloudsrc/api.awaymail.v1.go/src/domain/messages"
)

type messagesRepository struct {
	db      *mgo.Database
	mongoDB dbMongo.Client
}

const (
	inboxCollectionName string = "messages.inbox"
	sentCollectionName  string = "messages.sent"
)

func NewMessages(db *mgo.Database, mongoDB dbMongo.Client) *messagesRepository {
	r := &messagesRepository{db: db, mongoDB: mongoDB}
	if r.db == nil {
		panic("mongo client is nil")
	}

	if err := mongoDB.Collection(inboxCollectionName); err != nil {
		panic("collection casbu not found")
	}
	return r
}

func (r *messagesRepository) FindInbox(selector interface{}) (entity *domainMsg.Message, err error) {
	if err = r.db.C(inboxCollectionName).Find(selector).One(&entity); err != nil {
		return
	}
	return
}

func (r *messagesRepository) InsertInbox(entity *domainMsg.Message) (err error) {
	if _, err = r.db.C(inboxCollectionName).Upsert(bson.M{"message_id": entity.MessageID}, &entity); err != nil {
		return err
	}
	return nil
}

func (r *messagesRepository) DeleteArchiveByEmail(email string) (err error) {
	_, err = r.db.C(inboxCollectionName).RemoveAll(bson.M{"to": email})
	return
}

func (r *messagesRepository) InsertSentBox(entity *domainMsg.Message) (err error) {
	if _, err = r.db.C(sentCollectionName).Upsert(bson.M{"message_id": entity.MessageID}, &entity); err != nil {
		return err
	}
	return nil
}

func (r *messagesRepository) GetInbox(skip int, userID bson.ObjectId) (entity []domainMsg.Message, err error) {
	var data []domainMsg.Message
	pipeline := []bson.M{}
	pipeline = []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$sort": bson.D{{"received_at", -1}}},
		{"$skip": skip * 10},
		{"$limit": 10}}
	err = r.db.C(inboxCollectionName).Pipe(pipeline).All(&data)
	return data, err
}
