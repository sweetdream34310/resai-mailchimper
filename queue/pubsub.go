package queue

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"github.com/cloudsrc/api.awaymail.v1.go/src/domain/redis"
	googleWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/gcp"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils"
	messagesSvc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/messages"
)

type pubSub struct {
	svc           *gcp.PubSubCon
	googleWrapper googleWrapper.Wrapper
	redisRepo     redis.Repository
	messagesSvc   messagesSvc.Service
}

func NewPubSub(config config.Config, redisRepo redis.Repository, googleWrapper googleWrapper.Wrapper, messagesSvc messagesSvc.Service) {
	u := pubSub{svc: gcp.GetConnection(config), redisRepo: redisRepo, googleWrapper: googleWrapper, messagesSvc: messagesSvc}
	subscription := u.svc.Connection.Subscription(config.Gpubsub.Subscribe)
	go subscription.Receive(u.svc.Context, u.notification)
}

func (u *pubSub) notification(ctx context.Context, msg *pubsub.Message) {
	msg.Ack()
	sess := createSession(utils.GenerateThreadId(), "pub/sub notification")
	in := &messagesSvc.PushNotification{}

	err := json.Unmarshal(msg.Data, &in)
	if err != nil {
		sess.ErrorMessage = err.Error()
		sess.Lv4()
		return
	}
	sess.Request = in
	sess.Lv1("Incoming message pub/sub notification")

	sess.UserSession.Email = in.EmailAddress

	go u.messagesSvc.PushNotification(sess, in)
}
