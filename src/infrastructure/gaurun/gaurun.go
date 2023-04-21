package gaurun

import (
	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/rest"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
)

type PushNotification struct {
	config config.Config
}

func NewPushNotification(config config.Config) *PushNotification {
	return &PushNotification{config: config}
}

func (p *PushNotification) Send(ctxSess *ctxSess.Context, payload Notification) error {
	gaurunURL := p.config.Gaurun.Url
	//if ctxSess.UserSession.DevTokenKey {
	//	gaurunURL = p.config.Gaurun.UrlDev
	//}

	if _, err := rest.Request("POST", gaurunURL, &rest.RequestOptions{
		Payload: payload,
	}); err != nil {
		ctxSess.ErrorMessage = err.Error()
		return err
	}
	return nil
}
