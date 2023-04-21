package google

import (
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"

	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
)

type Wrapper interface {
	ValidateToken(ctxSess *ctxSess.Context, agent, refreshToken string) (token *oauth2.Token, err error)
	GetProfile(ctxSess *ctxSess.Context, agent, refreshToken string) (resp UserProfile, err error)
	GetMessagesList(ctxSess *ctxSess.Context, req GetMessageRequest) (*gmail.ListMessagesResponse, error)
	GetMessages(ctxSess *ctxSess.Context, messageId string) (message *gmail.Message, err error)
	GetThreadMessages(ctxSess *ctxSess.Context, threadId string) (*gmail.Thread, error)
	SendMessage(ctxSess *ctxSess.Context, req Messsage) (err error)
	SendMessageWithAttachment(ctxSess *ctxSess.Context, req Messsage) (err error)
	UpdateMessage(ctxSess *ctxSess.Context, id string, isRead *bool) (message *gmail.Message, err error)
	DeleteMessage(ctxSess *ctxSess.Context, messageId string) (err error)
	ModifyMessage(ctxSess *ctxSess.Context, req *gmail.BatchModifyMessagesRequest) (err error)
	GetAttachmentDetail(ctxSess *ctxSess.Context, messageID, attachmentID string) (res *gmail.MessagePartBody, err error)
	GetHistory(ctxSess *ctxSess.Context, historyID uint64) (resp *gmail.ListHistoryResponse, err error)

	//label
	GetUserLabel(ctxSess *ctxSess.Context, messageId string) (res *gmail.Label, err error)
	GetLabelList(ctxSess *ctxSess.Context) (res []*gmail.Label, err error)
	CreateLabel(ctxSess *ctxSess.Context, req *gmail.Label) (res *gmail.Label, err error)
	DeleteLabel(ctxSess *ctxSess.Context, labelID string) (err error)
	PatchLabel(ctxSess *ctxSess.Context, labelID string, req *gmail.Label) (res *gmail.Label, err error)

	WatchPushNotification(ctxSess *ctxSess.Context, refreshToken string, req *gmail.WatchRequest) (res *gmail.WatchResponse, err error)
	StopWatchPushNotification(ctxSess *ctxSess.Context, refreshToken string) (err error)
}
