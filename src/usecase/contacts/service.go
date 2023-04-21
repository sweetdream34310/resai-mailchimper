package contacts

import (
	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/contacts"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
)

type Service interface {
	AddContacts(ctxSess *ctxSess.Context, req ContactAddReq, user models.UserSession) (interface{}, error)
	UpdateContacts(ctxSess *ctxSess.Context, req ContactUpdateReq, user models.UserSession) (*domain.Contacts, error)
	DeleteContacts(ctxSess *ctxSess.Context, id string, user models.UserSession) (err error)
	GetContactsList(ctxSess *ctxSess.Context, user models.UserSession) (interface{}, error)
	GetRecentContactsList(ctxSess *ctxSess.Context, user models.UserSession) (interface{}, error)
}
