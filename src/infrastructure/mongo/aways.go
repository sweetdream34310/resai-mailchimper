package mongo

import (
	"errors"

	dbMongo "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mo "go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/aways"
)

type awaysRepository struct {
	db      *mgo.Database
	mongoDB dbMongo.Client
}

const awaysCollectionName string = "aways"

func NewAways(db *mgo.Database, mongoDB dbMongo.Client) *awaysRepository {
	r := &awaysRepository{db: db, mongoDB: mongoDB}
	if r.db == nil {
		panic("mongo client is nil")
	}

	if err := mongoDB.Collection(awaysCollectionName); err != nil {
		panic("collection casbu not found")
	}
	return r
}

func (a *awaysRepository) CreateAway(req *domain.Away) (entity *domain.Away, err error) {
	opts := &dbMongo.Options{}
	opts.CollectionName = awaysCollectionName
	if err = a.mongoDB.FindOne(dbMongo.M{
		"$and": []dbMongo.M{
			{"title": req.Title},
			{"user_id": req.UserID},
		}},
		entity, opts); err == mo.ErrNoDocuments {
		options := &dbMongo.Options{CollectionName: awaysCollectionName, Upsert: true}
		if _, err = a.mongoDB.InsertOne(*req, options); err != nil {
			return nil, err
		}
	}
	return req, nil
}

func (a *awaysRepository) GetAwayList(userID primitive.ObjectID, enable bool) (entity []*domain.Away, err error) {
	filter := bson.M{"user_id": userID}
	if enable {
		filter = bson.M{"$and": []bson.M{
			{"user_id": userID},
			{"is_enabled": enable},
		}}
	}
	opts := &dbMongo.Options{}
	opts.CollectionName = awaysCollectionName
	err = a.mongoDB.Find(filter, &entity, opts)
	return
}

func (a *awaysRepository) GetAway(awayID primitive.ObjectID) (entity *domain.Away, err error) {
	opts := &dbMongo.Options{}
	opts.CollectionName = awaysCollectionName
	err = a.mongoDB.FindOne(dbMongo.M{"_id": awayID}, &entity, opts)
	return
}

func (a *awaysRepository) UpdateAway(req *domain.Away) (err error) {

	opts := &dbMongo.Options{CollectionName: awaysCollectionName, Upsert: true}
	change := dbMongo.M{"$set": dbMongo.M{
		"title":            req.Title,
		"repeat":           req.Repeat,
		"all_day":          req.AllDay,
		"activate_allow":   req.ActivateAllow,
		"deactivate_allow": req.DeactivateAllow,
		"is_enabled":       req.IsEnabled,
		"message":          req.Message,
		"key_duration":     req.KeyDuration,
		"allowed_contacts": req.AllowedContacts,
		"allowed_keywords": req.AllowedKeywords,
	}}
	_, _, err = a.mongoDB.UpdateOne(dbMongo.M{"_id": req.ID}, change, opts)
	if err != nil {
		if errors.Is(err, mo.ErrNoDocuments) {
			err = nil
			return
		}
		return
	}
	return
}

func (a *awaysRepository) UpdateAwayMode(isEnabled bool, userID primitive.ObjectID) (total int64, err error) {
	opts := &dbMongo.Options{CollectionName: awaysCollectionName, Upsert: true}
	change := dbMongo.M{"$set": dbMongo.M{
		"is_enabled": isEnabled,
	}}
	total, err = a.mongoDB.UpdateMany(bson.M{"user_id": userID}, change, opts)
	return
}

func (a *awaysRepository) DeleteAway(awayID primitive.ObjectID) (err error) {
	opts := &dbMongo.Options{}
	opts.CollectionName = awaysCollectionName
	_, err = a.mongoDB.DeleteOne(dbMongo.M{"_id": awayID}, opts)
	return
}
