package database

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/mongo"
)

type DB struct {
	mongo.Client
}

func New(cfg config.Config) *DB {
	uri := fmt.Sprintf("mongodb://%v/%v", cfg.MongoDB.Hosts[0], cfg.MongoDB.Database)
	credential := mongo.Credential{
		AuthSource: "admin",
		Username:   cfg.MongoDB.User,
		Password:   cfg.MongoDB.Password,
	}
	client, err := mongo.Connect(context.Background(), uri, mongo.ClientOptions{
		MaxPoolSize:     5,
		MinPoolSize:     2,
		MaxConnIdleTime: 30 * time.Second,
		Auth:            credential,
	})
	if err != nil {
		panic("failed to connect db mongo")
	}
	conn := client.DB(cfg.MongoDB.Database)
	return &DB{Client: conn}
}
