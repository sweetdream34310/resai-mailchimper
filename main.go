package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudsrc/api.awaymail.v1.go/handlers"
	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	"github.com/cloudsrc/api.awaymail.v1.go/queue"
	"github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/chatgpt"
	"github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/gaurun"
	"github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	"github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/mongo"
	"github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/redis"
	awaysSvc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/aways"
	contactsSvc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/contacts"
	labelSvc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/label"
	messagesSvc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/messages"
	userSvc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/users"
	"github.com/mailgun/mailgun-go/v4"
)

type servers struct {
	Host string
	Port uint
}

func httpAddress(s servers) string {
	return s.Host + ":" + fmt.Sprintf("%d", s.Port)
}

func main() {

	app := libs.NewApp()

	startServer(app)
}

func startServer(app *libs.App) {

	s := servers{
		Host: app.Config.Server.Host,
		Port: app.Config.Server.Port,
	}

	app.Router = app.Engine.Group("/v2")

	contactsRepo := mongo.NewContacts(app.DB.MongoDB, app.DBMongo)
	userRepo := mongo.NewUsers(app.DB.MongoDB, app.DBMongo)
	awayRepo := mongo.NewAways(app.DB.MongoDB, app.DBMongo)
	messagesRepo := mongo.NewMessages(app.DB.MongoDB, app.DBMongo)
	redisRepo := redis.NewRedisInbox(app.Redis)

	googleWrapper := google.New(app.Config)
	chatgptWrapper := chatgpt.Setup(app.Config)
	gaurunWrapper := gaurun.NewPushNotification(app.Config)
	mg := mailgun.NewMailgun(app.Config.Mailgun.Domain, app.Config.Mailgun.APIKey)

	contactsSvc := contactsSvc.New(app.Config, contactsRepo, redisRepo)
	userSvc := userSvc.New(app.Config, userRepo, googleWrapper, redisRepo)
	awaysSvc := awaysSvc.New(awayRepo, redisRepo)
	labelSvc := labelSvc.New(googleWrapper)
	messagesSvc := messagesSvc.New(app.Config, awayRepo, messagesRepo, userRepo, contactsRepo, googleWrapper, chatgptWrapper, gaurunWrapper, redisRepo, app.Rabbit, mg)

	queueHandler := queue.Messages{
		App:         app,
		MessagesSvc: messagesSvc,
		UserRepo:    userRepo,
	}
	queueHandler.Run()

	queue.NewPubSub(app.Config, redisRepo, googleWrapper, messagesSvc)

	messageHandler := handlers.Message{
		App:           app,
		MessagesSvc:   messagesSvc,
		GoogleWrapper: googleWrapper,
		RedisRepo:     redisRepo,
	}
	messageHandler.SetRouter()
	publicHandler := handlers.Public{App: app, MessagesSvc: messagesSvc}
	publicHandler.SetRouter()
	awayHandler := handlers.Away{
		App:           app,
		AwaySvc:       awaysSvc,
		GoogleWrapper: googleWrapper,
		RedisRepo:     redisRepo,
	}
	awayHandler.SetRouter()
	userHandler := handlers.User{
		App:           app,
		UserSvc:       userSvc,
		Queue:         queueHandler,
		GoogleWrapper: googleWrapper,
		RedisRepo:     redisRepo,
	}
	userHandler.SetRouter()
	contactsHandler := handlers.Contacts{
		App:           app,
		ContactsSvc:   contactsSvc,
		GoogleWrapper: googleWrapper,
		RedisRepo:     redisRepo,
	}
	contactsHandler.SetRouter()
	labelHandler := handlers.Label{
		App:           app,
		LabelSvc:      labelSvc,
		GoogleWrapper: googleWrapper,
		RedisRepo:     redisRepo,
	}
	labelHandler.SetRouter()

	srv := &http.Server{
		Addr:    httpAddress(s),
		Handler: app.Engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatalf("listen: %s\n", err)
		}
	}()

	app.Logger.Info("server started listening on : ", httpAddress(s))

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := app.Shutdown(srv); err != nil {
		app.Logger.Fatal("Server forced to shutdown:", err)
	}

	app.Logger.Println("Server exiting")
}
