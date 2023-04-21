package queue

import (
	"encoding/json"
	"fmt"

	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/users"
	Logger "github.com/cloudsrc/api.awaymail.v1.go/src/shared/logger"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	messagesSvc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/messages"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Messages struct {
	App         *libs.App
	MessagesSvc messagesSvc.Service
	UserRepo    domain.Repository
}

func (m *Messages) Run() {
	users, _ := m.UserRepo.GetAllUser()
	for _, each := range users {
		go m.App.Rabbit.Consume("inbox_processor_"+each.ID.Hex(), m.processInbox)
		go m.App.Rabbit.Consume("sent_processor_"+each.ID.Hex(), m.ProcessSent)
	}
}

func (m *Messages) RunOneID(id primitive.ObjectID) {
	go m.App.Rabbit.Consume("inbox_processor_"+id.Hex(), m.processInbox)
	go m.App.Rabbit.Consume("sent_processor_"+id.Hex(), m.ProcessSent)
}

func (m *Messages) processInbox(data []byte) error {
	ctxSess := createSession(utils.GenerateThreadId(), "inbox_processor")
	message := &messagesSvc.MessageQueue{}
	if err := json.Unmarshal(data, &message); err != nil {
		ctxSess.ErrorMessage = err.Error()
		ctxSess.Lv4()
		return err
	}
	ctxSess.UserSession = models.UserSession{
		UserID:       message.UserId,
		Email:        message.Email,
		RefreshToken: message.RefreshToken,
		AuthToken:    message.AccessToken,
	}
	ctxSess.Request = message
	ctxSess.Lv1(fmt.Sprintf("Incoming message %s", message.Email))
	if _, err := m.MessagesSvc.GetMessage(ctxSess, message); err != nil {
		ctxSess.Lv4()
		return err
	}
	ctxSess.Lv4()
	return nil
}

func (m *Messages) ProcessSent(data []byte) error {
	ctxSess := createSession(utils.GenerateThreadId(), "sent_processor")
	message := &messagesSvc.MessageQueue{}
	if err := json.Unmarshal(data, &message); err != nil {
		ctxSess.ErrorMessage = err.Error()
		ctxSess.Lv4()
		return err
	}
	ctxSess.UserSession = models.UserSession{
		UserID:       message.UserId,
		Email:        message.Email,
		RefreshToken: message.RefreshToken,
		AuthToken:    message.AccessToken,
	}
	ctxSess.Request = message
	ctxSess.Lv1(fmt.Sprintf("Incoming message %s", message.Email))
	if _, err := m.MessagesSvc.GetSentMessage(ctxSess, message); err != nil {
		ctxSess.Lv4()
		return err
	}
	ctxSess.Lv4()
	return nil
}

func createSession(threadID, topic string) *ctxSess.Context {
	return ctxSess.New(Logger.GetLogger()).
		SetAppName("api.awaymail.v1.go").
		SetAppVersion("0.0").
		SetURL(topic).
		SetXRequestID(threadID).
		SetMethod("QUEUE")
}
