package messages

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InboxReq struct {
	Page   int64  `json:"page" form:"page"`
	Search string `json:"search" form:"search"`
}

type InboxSentReq struct {
	Page   int64  `json:"page" form:"page"`
	Search string `json:"search" form:"search"`
}

type InboxResp struct {
	Messages []interface{} `json:"messages"`
	Page     Page          `json:"page"`
}

type SentEmailResp struct {
	Messages []interface{} `json:"messages"`
	Page     Page          `json:"page"`
}

type Page struct {
	CurrentPage int64 `json:"current_page"`
	TotalPage   int64 `json:"total_page"`
	TotalEmails int64 `json:"total_emails"`
}

type MessageQueue struct {
	MessageID    string             `json:"message_id"`
	UserIdHex    string             `json:"user_id_hex"`
	UserId       primitive.ObjectID `json:"user_id"`
	Email        string             `json:"email"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	IsRead       bool               `json:"is_read"`
	IsReply      bool               `json:"is_reply"`
}

type SendMessageReq struct {
	To             string   `json:"to" binding:"required"`
	Message        string   `json:"message" binding:"required"`
	Subject        string   `json:"subject" binding:"required"`
	AttachmentsURL []string `json:"attachments_url"`
	ThreadID       string   `json:"thread_id"`
	MessageID      string   `json:"message_id"`
}

type UpdateRequest struct {
	IsRead *bool `json:"is_read" binding:"required"`
}

type UserMessageReq struct {
	Label         string `json:"label" form:"label"`
	NextPageToken string `json:"nextPageToken" form:"nextPageToken"`
	PrevPageToken string `json:"prevPageToken" form:"prevPageToken"`
	Search        string `json:"search" form:"search"`
}

type UserMessageResp struct {
	Messages         []UserMessage `json:"messages"`
	NextPageToken    string        `json:"nextPageToken"`
	CurrentPageToken string        `json:"currentPageToken"`
	PrevPageToken    string        `json:"prevPageToken"`
	TotalEmails      int64         `json:"totalEmails"`
}

type ThreadMessageResp struct {
	Messages []UserMessage `json:"messages"`
}

type UserMessage struct {
	MessageID            string       `bson:"message_id,omitempty" json:"message_id"`
	LabelIds             []string     `bson:"labels_id,omitempty" json:"labels_id,omitempty"`
	OriginalMessageID    string       `bson:"original_message_id,omitempty" json:"original_message_id"`
	ThreadID             string       `bson:"thread_id,omitempty" json:"thread_id"`
	HistoryId            uint64       `json:"history_id,omitempty"`
	Subject              string       `bson:"subject,omitempty" json:"subject,omitempty"`
	From                 *string      `bson:"from,omitempty" json:"from,omitempty"`
	To                   *string      `bson:"to,omitempty" json:"to,omitempty"`
	Received             *time.Time   `bson:"received_at,omitempty" json:"received_at,omitempty"`
	Sent                 *time.Time   `bson:"sent,omitempty" json:"sent,omitempty"`
	BodyText             string       `bson:"body_text,omitempty" json:"body_text,omitempty"`
	BodyRawText          string       `json:"body_raw_text"`
	BodyHTML             string       `bson:"body_html,omitempty" json:"body_html,omitempty"`
	IsRead               bool         `bson:"is_read,omitempty" json:"is_read"`
	IsReply              bool         `bson:"is_reply,omitempty" json:"is_reply"`
	MimeFlow             []string     `bson:"mime_flow,omitempty" json:"mime_flow"`
	Labels               []string     `bson:"labels,omitempty" json:"labels"`
	Snippet              string       `bson:"snippet,omitempty" json:"snippet"`
	Attachments          []Attachment `bson:"attachments,omitempty" json:"attachments"`
	Thumbnail            []byte       `json:"thumbnail"`
	Summary              string       `json:"summary"`
	FlagBreakthrough     bool         `json:"flag_breakthrough"`
	CountReply           int          `json:"count_reply"`
	AwayContainsKeywords []string     `json:"away_contains_keywords"`
}

type Attachment struct {
	AttachmentID   string `json:"attachment_id"`
	AttachmentSize int64  `json:"attachment_size"`
	FileName       string `json:"filename"`
	MimeType       string `json:"mime_type"`
}

type AttachmentRequest struct {
	MessageID    string `json:"message_id" form:"message_id" binding:"required"`
	AttachmentID string `json:"attachment_id" form:"attachment_id" binding:"required"`
	FileName     string `json:"filename" form:"filename" binding:"required"`
	MimeType     string `json:"mime_type" form:"mime_type" binding:"required"`
}

type AttachmentResponse struct {
	FileByte []byte `json:"file_byte"`
}

type TagLabelReq struct {
	MessageIds []string `json:"message_ids" binding:"required"`
	LabelIds   []string `json:"label_ids" binding:"required"`
}

type TagLabelArchiveReq struct {
	MessageIds []string `json:"message_ids" binding:"required"`
}

type PushNotification struct {
	EmailAddress string `json:"emailAddress"`
	HistoryId    uint64 `json:"historyId"`
}

type ArchiveReq struct {
	Archive bool `json:"archive"`
}

type BreakthroughReq struct {
	Email      string `json:"email" form:"email"`
	XRequestID string `json:"id" form:"id"`
}
