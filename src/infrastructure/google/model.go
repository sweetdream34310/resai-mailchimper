package google

type UserProfile struct {
	Email        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
	AuthToken    string `json:"auth_token"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
}

type Messsage struct {
	Email          string   `json:"email"`
	To             string   `json:"to"`
	Subject        string   `json:"subject"`
	Message        string   `json:"message"`
	ThreadID       string   `json:"thread_id,omitempty"`
	MessageID      string   `json:"message_id,omitempty"`
	AttachmentsURL []string `json:"attachments_url,omitempty"`
}

type GetMessageRequest struct {
	NextPageToken string `json:"nextPageToken"`
	Label         string `json:"label"`
	Q             string `json:"q"`
	Limit         int64  `json:"limit"`
}
