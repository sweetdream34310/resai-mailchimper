package mongo

import (
	"context"
	"errors"
	"time"

	mo "go.mongodb.org/mongo-driver/mongo"
	moOpts "go.mongodb.org/mongo-driver/mongo/options"
)

type (
	//M map filter type
	M map[string]interface{}
	//Client mongo interface for mongo client connection
	Client interface {
		Disconnect() error
		DB(dbName string) *connection
		Collection(collectionName string) (err error)
		GetContext() context.Context
		SetContext(c context.Context)
		WithTimeout(timeSec time.Duration) context.CancelFunc

		Find(filter interface{}, outputVal interface{}, opts ...*Options) (err error)
		FindOne(filter interface{}, outputVal interface{}, opts ...*Options) error
		FindOneAndUpdate(filter, update interface{}, outputVal interface{}, opts ...*Options) error
		InsertOne(document interface{}, opts ...*Options) (insertedID string, err error)

		UpdateMany(filter, update interface{}, opts ...*Options) (modifiedCount int64, err error)
		UpdateOne(filter, update interface{}, opts ...*Options) (upsertedID string, modifiedCount int64, err error)

		CountDocuments(filter interface{}, opts ...*Options) (total int64, err error)

		DeleteOne(filter interface{}, opts ...*Options) (DeletedCount int64, err error)
		DeleteMany(filter interface{}, opts ...*Options) (DeletedCount int64, err error)

		Aggregate(pipeline interface{}, outputVal interface{}, opts ...*Options) (err error)
	}

	connection struct {
		ctx             context.Context
		client          *mo.Client
		collection      *mo.Collection
		multiCollection map[string]*mo.Collection
		dbName          string
		// debug           bool
	}
)

func Connect(ctx context.Context, URI string, opts ...ClientOptions) (mongoClient Client, err error) {
	clientOptions := moOpts.Client().ApplyURI(URI)

	if len(opts) > 0 {
		if opts[0].MaxConnIdleTime > 0 {
			clientOptions.SetMaxConnIdleTime(opts[0].MaxConnIdleTime)
		}
		if opts[0].MaxPoolSize > 0 {
			clientOptions.SetMaxPoolSize(opts[0].MaxPoolSize)
		}
		if opts[0].MinPoolSize > 0 {
			clientOptions.SetMinPoolSize(opts[0].MinPoolSize)
		}
		clientOptions.SetAuth(moOpts.Credential(opts[0].Auth))
	}

	// Connect to MongoDB
	moClient, err := mo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	//Check the connection
	err = moClient.Ping(ctx, nil)
	if err != nil {
		return
	}

	mongoClient = &connection{
		ctx:    ctx,
		client: moClient,
	}
	return
}

func (conn *connection) DB(dbName string) *connection {
	conn.dbName = dbName
	conn.multiCollection = map[string]*mo.Collection{}
	return conn
}

func (conn *connection) Collection(collectionName string) (err error) {
	if conn.dbName == "" {
		err = errors.New("DB required")
		return
	}
	if _, ok := conn.multiCollection[collectionName]; ok {
		return
	}
	conn.collection = conn.client.Database(conn.dbName).Collection(collectionName)
	conn.multiCollection[collectionName] = conn.collection
	return
}

// Disconnect terminate connection with mongo client
func (conn *connection) Disconnect() (err error) {
	err = conn.client.Disconnect(conn.ctx)
	return
}

// GetContext get connection context
func (conn *connection) GetContext() context.Context {
	return conn.ctx
}

// SetContext set connection context
func (conn *connection) SetContext(c context.Context) {
	conn.ctx = c
}

// WithTimeout set timeout based on context
func (conn *connection) WithTimeout(timeSec time.Duration) context.CancelFunc {
	ctx, cancel := context.WithTimeout(conn.ctx, timeSec*time.Second)
	conn.ctx = ctx
	return cancel
}
