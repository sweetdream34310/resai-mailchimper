package mongo

import (
	"errors"
	"time"

	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/users"
	dbMongo "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	mo "go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type usersRepository struct {
	db      *mgo.Database
	mongoDB dbMongo.Client
}

const usersCollectionName string = "users"

func NewUsers(db *mgo.Database, mongoDB dbMongo.Client) *usersRepository {
	r := &usersRepository{db: db, mongoDB: mongoDB}
	if r.db == nil {
		panic("mongo client is nil")
	}

	if err := mongoDB.Collection(usersCollectionName); err != nil {
		panic("collection casbu not found")
	}
	return r
}

func (r *usersRepository) AddUser(entity *domain.User) (*domain.User, error) {
	entity.ID, _ = primitive.ObjectIDFromHex(bson.NewObjectId().Hex())
	entity.LastLogin = time.Now().UTC()
	entity.CreatedAt = time.Now().UTC()
	entity.Active = 1
	opts := &dbMongo.Options{CollectionName: usersCollectionName, Upsert: true}
	if _, err := r.mongoDB.InsertOne(*entity, opts); err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *usersRepository) GetUser(email string) (entity *domain.User, err error) {
	opts := &dbMongo.Options{}
	opts.CollectionName = usersCollectionName
	err = r.mongoDB.FindOne(dbMongo.M{"email": email}, &entity, opts)
	if err == mo.ErrNoDocuments {
		err = nil
	}
	return
}

func (r *usersRepository) GetUserByID(id primitive.ObjectID) (entity *domain.User, err error) {
	opts := &dbMongo.Options{}
	opts.CollectionName = usersCollectionName
	err = r.mongoDB.FindOne(dbMongo.M{"_id": id}, &entity, opts)
	if err == mo.ErrNoDocuments {
		err = nil
	}
	return
}

func (r *usersRepository) GetAllUser() (entity []*domain.User, err error) {
	opts := &dbMongo.Options{}
	opts.CollectionName = usersCollectionName
	err = r.mongoDB.Find(dbMongo.M{"active": 1}, &entity, opts)
	if err == mo.ErrNoDocuments {
		err = nil
	}
	return
}

func (r *usersRepository) UpdateUser(email, authToken, refreshToken, name, photo string, swiftToken []string) (err error) {
	opts := &dbMongo.Options{CollectionName: usersCollectionName, Upsert: true}
	change := dbMongo.M{"$set": dbMongo.M{
		"last_login":    time.Now().UTC(),
		"auth_token":    authToken,
		"refresh_token": refreshToken,
		"name":          name,
		"photo":         photo,
		"swift_token":   swiftToken,
	},
	}
	entity := &domain.User{}
	err = r.mongoDB.FindOneAndUpdate(dbMongo.M{"email": email}, change, &entity, opts)
	if err != nil {
		if errors.Is(err, mo.ErrNoDocuments) {
			err = nil
			return
		}
		return
	}
	return
}

func (r *usersRepository) UpdateSwiftToken(email string, swiftToken []string) (err error) {
	opts := &dbMongo.Options{CollectionName: usersCollectionName, Upsert: true}
	change := dbMongo.M{"$set": dbMongo.M{
		"swift_token": swiftToken,
	},
	}
	entity := &domain.User{}
	err = r.mongoDB.FindOneAndUpdate(dbMongo.M{"email": email}, change, &entity, opts)
	if err != nil {
		if errors.Is(err, mo.ErrNoDocuments) {
			err = nil
			return
		}
		return
	}
	return
}

func (r *usersRepository) UpdateFlagArchive(email string, archive bool) (err error) {
	opts := &dbMongo.Options{CollectionName: usersCollectionName, Upsert: true}
	change := dbMongo.M{"$set": dbMongo.M{
		"archive": archive,
	},
	}
	entity := &domain.User{}
	err = r.mongoDB.FindOneAndUpdate(dbMongo.M{"email": email}, change, &entity, opts)
	if err != nil {
		if errors.Is(err, mo.ErrNoDocuments) {
			err = nil
			return
		}
		return
	}
	return
}
