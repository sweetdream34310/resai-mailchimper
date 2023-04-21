package gcp

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
)

type PubSubCon struct {
	Connection *pubsub.Client
	Context    context.Context
}

var (
	connection *PubSubCon
	mutex      sync.Mutex
)

func GetConnection(config config.Config) *PubSubCon {
	if connection == nil {
		mutex.Lock()
		defer mutex.Unlock()
		connection = newConnection(config)
	}

	return connection
}

func newConnection(config config.Config) *PubSubCon {
	ctx := context.Background()
	client := new(pubsub.Client)
	client, err := pubsub.NewClient(ctx, config.Gpubsub.ProjectName)
	if err != nil {
		filename := config.Gpubsub.CredentialPath
		client, err = pubsub.NewClient(ctx, config.Gpubsub.ProjectName, option.WithCredentialsFile(filename))
		if err != nil {
			log.Panic("got an error while connecting pubsub, ", zap.Error(err))
		}
	}

	return &PubSubCon{Connection: client, Context: ctx}
}
