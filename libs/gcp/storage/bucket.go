package bucket

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"google.golang.org/api/option"
)

type Connections struct {
	Connection *storage.Client
	Context    context.Context
	BucketName string
}

var (
	connection *Connections
)

func NewConnection(config config.Config) *Connections {
	ctx := context.Background()
	client := new(storage.Client)
	//client, err := storage.NewClient(ctx)
	//if err != nil {
	filename := config.Gcs.CredentialPath
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(filename))
	if err != nil {
		panic(err)
	}
	//}

	bucketName := config.Gcs.BucketName
	bucket := client.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		panic(err)
	}

	return &Connections{Connection: client, Context: ctx, BucketName: bucketName}
}
