package gaurun

type NotifExtend struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type NotifPayload struct {
	Token            []string      `json:"token"`
	Platform         int           `json:"platform"`
	Message          string        `json:"message,omitempty"`
	Sound            string        `json:"sound,omitempty"`
	Category         string        `json:"category,omitempty"`
	Badge            int           `json:"badge,omitempty"`
	Expiry           int           `json:"expiry,omitempty"`
	ContentAvailable bool          `json:"content_available,omitempty"`
	MutableContent   bool          `json:"mutable_content,omitempty"`
	PushType         string        `json:"push_type,omitempty"`
	Extend           []NotifExtend `json:"extend,omitempty"`
}

type Notification struct {
	Notifications []NotifPayload `json:"notifications"`
}
