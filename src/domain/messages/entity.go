package messages

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	UserID            primitive.ObjectID `bson:"user_id" json:"-"`
	MessageID         string             `bson:"message_id,omitempty" json:"message_id"`
	OriginalMessageID string             `bson:"original_message_id,omitempty" json:"original_message_id"`
	ThreadID          string             `bson:"thread_id,omitempty" json:"thread_id"`
	Subject           string             `bson:"subject,omitempty" json:"subject,omitempty"`
	From              *string            `bson:"from,omitempty" json:"from,omitempty"`
	To                *string            `bson:"to,omitempty" json:"to,omitempty"`
	Received          *time.Time         `bson:"received_at,omitempty" json:"received_at,omitempty"`
	Sent              *time.Time         `bson:"sent,omitempty" json:"sent,omitempty"`
	BodyText          string             `bson:"body_text,omitempty" json:"body_text,omitempty"`
	BodyHTML          string             `bson:"body_html,omitempty" json:"body_html,omitempty"`
	IsRead            bool               `bson:"is_read,omitempty" json:"is_read"`
	IsReply           bool               `bson:"is_reply,omitempty" json:"is_reply"`
	MimeFlow          []string           `bson:"mime_flow,omitempty" json:"mime_flow"`
	Labels            []string           `bson:"labels,omitempty" json:"labels"`
	Snippet           string             `bson:"snippet,omitempty" json:"snippet"`
	Attachments       []Attachment       `bson:"attachments,omitempty" json:"attachments"`
}

type Attachment struct {
	AttachmentID   string `json:"attachment_id"`
	AttachmentSize int64  `json:"attachment_size"`
	FileName       string `json:"filename"`
	MimeType       string `json:"mime_type"`
}
