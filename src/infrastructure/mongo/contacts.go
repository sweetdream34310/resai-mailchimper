package mongo

import (
	"errors"

	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/contacts"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	dbMongo "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mo "go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

type contactsRepository struct {
	db      *mgo.Database
	mongoDB dbMongo.Client
}

const contactsCollectionName string = "contacts"

func NewContacts(db *mgo.Database, mongoDB dbMongo.Client) *contactsRepository {
	r := &contactsRepository{db: db, mongoDB: mongoDB}
	if r.db == nil {
		panic("mongo client is nil")
	}

	if err := mongoDB.Collection(contactsCollectionName); err != nil {
		panic("collection casbu not found")
	}
	return r
}

func (r *contactsRepository) Add(entity *domain.Contacts) (*domain.Contacts, error) {
	opts := &dbMongo.Options{}
	opts.CollectionName = contactsCollectionName
	if err := r.mongoDB.FindOne(dbMongo.M{
		"$and": []dbMongo.M{
			{"email": entity.Email},
			{"user_id": entity.UserID},
		}},
		entity, opts); err == mo.ErrNoDocuments {
		options := &dbMongo.Options{CollectionName: contactsCollectionName, Upsert: true}
		if _, err = r.mongoDB.InsertOne(*entity, options); err != nil {
			return nil, err
		}
	}
	return entity, nil
}

func (r *contactsRepository) Update(entity *domain.Contacts) (out *domain.Contacts, err error) {
	opts := &dbMongo.Options{CollectionName: contactsCollectionName, Upsert: true}
	change := dbMongo.M{"$set": dbMongo.M{
		"name":  entity.Name,
		"email": entity.Email,
	}}
	err = r.mongoDB.FindOneAndUpdate(dbMongo.M{"_id": entity.ID}, change, &out, opts)
	if err != nil {
		if errors.Is(err, mo.ErrNoDocuments) {
			err = nil
			return
		}
		return
	}
	return
}

func (r *contactsRepository) Get(id string) (entity *domain.Contacts, err error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	opts := &dbMongo.Options{}
	opts.CollectionName = contactsCollectionName
	err = r.mongoDB.FindOne(dbMongo.M{"_id": docID}, &entity, opts)
	if err == mo.ErrNoDocuments {
		err = nil
	}
	return
}

func (r *contactsRepository) Delete(id string) (err error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	opts := &dbMongo.Options{}
	opts.CollectionName = contactsCollectionName
	_, err = r.mongoDB.DeleteOne(dbMongo.M{"_id": docID}, opts)
	return
}

func (r *contactsRepository) GetList(user models.UserSession) (entity []*domain.Contacts, err error) {
	opts := &dbMongo.Options{}
	opts.CollectionName = contactsCollectionName
	err = r.mongoDB.Find(dbMongo.M{"user_id": user.UserID}, &entity, opts)
	if err == mo.ErrNoDocuments {
		err = nil
	}
	return
}
