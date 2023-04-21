package libs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/database"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
	bucketGCS "github.com/cloudsrc/api.awaymail.v1.go/libs/gcp/storage"
	dbMongo "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mongo"
)

type (
	// App stucture
	App struct {
		Logger  *logrus.Logger
		Engine  *gin.Engine
		Config  config.Config
		DB      *mongoClient
		DBMongo dbMongo.Client
		Rabbit  *RabbitClient
		Redis   *RedisClient
		Router  *gin.RouterGroup
		Bucket  *bucketGCS.Connections
	}
)

func NewApp() *App {

	config := config.Load()

	log := logrus.New()

	if os.Getenv("GOLANG_ENV") == "prod" {
		log.Level = logrus.InfoLevel
	} else {
		log.Level = logrus.DebugLevel
	}

	log.Formatter = &logrus.JSONFormatter{}

	engine := gin.New()

	if os.Getenv("GOLANG_ENV") == "test" {

		gin.SetMode(gin.TestMode)

	} else {
		// enabling cors
		engine.Use(cors.Middleware(cors.Config{
			Origins:         "*",
			Methods:         "GET, PUT, POST, DELETE, OPTIONS",
			RequestHeaders:  "Origin, Authorization, Content-Type, X-Requested-With, Access-Control-Allow-Origin, X-Token-Key",
			ExposedHeaders:  "",
			MaxAge:          50 * time.Second,
			ValidateHeaders: false,
		}))

		engine.Use(gzip.Gzip(gzip.DefaultCompression))

		// if os.Getenv("GOLANG_ENV") != "local" {
		//engine.Use(ginrus.Ginrus(log, time.RFC3339, true))
		//engine.Use(RequestLoggerMiddleware())
		// }
		engine.Use(RequestIdMiddleware())

		engine.Use(gin.Recovery())
	}

	mgoClient := newMongoClient(config)

	db := database.New(config)

	if err := mgoClient.MongoDB.C("users").EnsureIndex(mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		Background: true,
		Sparse:     true,
		Name:       "users",
	}); err != nil {
		panic(err)
	}

	if err := mgoClient.MongoDB.C("messages.inbox").EnsureIndex(mgo.Index{
		Key:        []string{"message_id"},
		Unique:     true,
		Background: true,
		Sparse:     true,
		Name:       "messages.inbox",
	}); err != nil {
		panic(err)
	}

	if err := mgoClient.MongoDB.C("messages.sent").EnsureIndex(mgo.Index{
		Key:        []string{"message_id"},
		Unique:     true,
		Background: true,
		Sparse:     true,
		Name:       "messages.sent",
	}); err != nil {
		panic(err)
	}

	if err := mgoClient.MongoDB.C("aways").EnsureIndex(mgo.Index{
		Key:        []string{"title"},
		Unique:     false,
		Background: true,
		Sparse:     true,
		Name:       "aways",
	}); err != nil {
		panic(err)
	}

	if err := mgoClient.MongoDB.C("contacts").EnsureIndex(mgo.Index{
		Key:        []string{"email"},
		Unique:     false,
		Background: true,
		Sparse:     true,
		Name:       "contacts",
	}); err != nil {
		panic(err)
	}

	return &App{
		Logger:  log,
		Engine:  engine,
		Config:  config,
		DB:      mgoClient,
		DBMongo: db,
		Rabbit:  newRabbitclient(config),
		Redis:   newRedisClient(config),
		Bucket:  bucketGCS.NewConnection(config),
	}
}

// Shutdown : Closes all connection.
func (app *App) Shutdown(srv *http.Server) error {

	app.Logger.Println("Shuting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	app.Redis.Close()

	if err := srv.Shutdown(ctx); err != nil {
		cancel()
		return err
	}

	defer cancel()

	return nil
}

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		fmt.Println(string(body))
		fmt.Println(c.Request.Header)
		c.Next()
	}
}
