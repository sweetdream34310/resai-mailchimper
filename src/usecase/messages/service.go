package messages

import (
	"google.golang.org/api/gmail/v1"

	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
)

type Service interface {
	ArchiveInbox(ctxSess *ctxSess.Context, req ArchiveReq)
	DeleteArchiveInbox(ctxSess *ctxSess.Context) (err error)
	GetUserInbox(ctxSess *ctxSess.Context, user models.UserSession, req InboxReq) (resp InboxResp, err error)
	GetUserSent(ctxSess *ctxSess.Context, user models.UserSession, req InboxSentReq) (resp SentEmailResp, err error)
	GetMessage(ctxSess *ctxSess.Context, queue *MessageQueue) (resp interface{}, err error)
	GetSentMessage(ctxSess *ctxSess.Context, queue *MessageQueue) (resp interface{}, err error)
	SentMessage(ctxSess *ctxSess.Context, user models.UserSession, req SendMessageReq) (err error)
	UpdateMessage(ctxSess *ctxSess.Context, messageId string, req UpdateRequest) (resp *gmail.Message, err error)
	DeleteMessage(ctxSess *ctxSess.Context, messageId string) (err error)
	TagLabelUserMessage(ctxSess *ctxSess.Context, req TagLabelReq) (err error)
	TagLabelArchiveMessage(ctxSess *ctxSess.Context, req TagLabelArchiveReq) (err error)
	GetAttachment(ctxSess *ctxSess.Context, req AttachmentRequest) (resp AttachmentResponse, err error)
	GetThreadByID(ctxSess *ctxSess.Context, threadId string) (resp ThreadMessageResp, err error)

	GetUserMessage(ctxSess *ctxSess.Context, req UserMessageReq) (resp UserMessageResp, err error)
	GetUserMessageDetail(ctxSess *ctxSess.Context, messageID, threadID string) (resp *UserMessage, err error)
	GetUserMessageByID(ctxSess *ctxSess.Context, messageID string) (resp *UserMessage, err error)

	BreakthroughNotification(ctxSess *ctxSess.Context, req BreakthroughReq) (err error)
	PushNotification(ctxSess *ctxSess.Context, in *PushNotification) (err error)
}
