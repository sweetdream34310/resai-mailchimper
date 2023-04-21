package contacts

import (
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
)

type Repository interface {
	Add(entity *Contacts) (*Contacts, error)
	Get(id string) (*Contacts, error)
	GetList(user models.UserSession) ([]*Contacts, error)
	Update(entity *Contacts) (*Contacts, error)
	Delete(id string) error
}
