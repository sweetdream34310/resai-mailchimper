package provider

// Provider needs to be implemented for each 3rd party authentication provider
// e.g. Facebook, Twitter, etc...
type Provider interface {
	Name() string
	SetName(string)
	SetAgent(string)
	ValidateToken(string) error
	GetUserProfile() error
	GetUserInbox() error
	GetUserSent() error
	GetMessage(string) (interface{}, error)
	UpdateMessage(string, *bool) (interface{}, error)
	SendMessage(interface{}) error
	SendMessageWithAttachment(interface{}) error
}
