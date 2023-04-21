package contacts

import (
	"fmt"
	"strings"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/contacts"
	"github.com/cloudsrc/api.awaymail.v1.go/src/domain/redis"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	config       config.Config
	contactsRepo domain.Repository
	redisRepo    redis.Repository
}

func New(config config.Config, contactsRepo domain.Repository, redisRepo redis.Repository) Service {
	if contactsRepo == nil {
		panic("contacts repository is nil")
	}
	if redisRepo == nil {
		panic("redis repository is nil")
	}

	return &service{
		config:       config,
		contactsRepo: contactsRepo,
		redisRepo:    redisRepo,
	}
}

func (s *service) AddContacts(ctxSess *ctxSess.Context, req ContactAddReq, user models.UserSession) (resp interface{}, err error) {
	if len(strings.TrimSpace(req.Name)) == 0 {
		req.Name = req.Email
	}
	encEmail, err := utils.Encrypt(s.config.Salt, req.Email)
	encName, err := utils.Encrypt(s.config.Salt, req.Name)
	entity, err := s.contactsRepo.Add(&domain.Contacts{
		ID:     primitive.NewObjectID(),
		Email:  encEmail,
		Name:   encName,
		UserID: user.UserID,
	})
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	entity.Name = req.Name
	entity.Email = req.Email

	resp = entity

	return
}

func (s *service) GetContactsList(ctxSess *ctxSess.Context, user models.UserSession) (resp interface{}, err error) {
	ls, err := s.contactsRepo.GetList(user)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDataNotFound
		return
	}

	for k, v := range ls {
		ls[k].Email, _ = utils.Decrypt(s.config.Salt, v.Email)
		ls[k].Name, _ = utils.Decrypt(s.config.Salt, v.Name)
	}
	resp = ls

	return
}

func (s *service) GetRecentContactsList(ctxSess *ctxSess.Context, user models.UserSession) (resp interface{}, err error) {
	key := fmt.Sprintf("recent:contact:%s", user.UserID.Hex())
	resp, err = s.redisRepo.GetPushCache(key, 0, 15)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDataNotFound
		return
	}

	resp = utils.RemoveDuplicateStr(resp.([]string))

	return
}

func (s *service) UpdateContacts(ctxSess *ctxSess.Context, req ContactUpdateReq, user models.UserSession) (resp *domain.Contacts, err error) {
	contact, err := s.contactsRepo.Get(req.ID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDataNotFound
		return
	}

	decEmail, _ := utils.Decrypt(s.config.Salt, contact.Email)
	if req.Email != decEmail {
		err = constants.ErrorEmailNotMatch
		ctxSess.ErrorMessage = err.Error()
		return
	}

	encEmail, err := utils.Encrypt(s.config.Salt, req.Email)
	encName, err := utils.Encrypt(s.config.Salt, req.Name)
	id, _ := primitive.ObjectIDFromHex(req.ID)
	entity, err := s.contactsRepo.Update(&domain.Contacts{
		Email:  encEmail,
		Name:   encName,
		UserID: user.UserID,
		ID:     id,
	})
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	entity.Name = req.Name
	entity.Email = req.Email

	resp = entity

	return
}

func (s *service) DeleteContacts(ctxSess *ctxSess.Context, id string, user models.UserSession) (err error) {
	contact, err := s.contactsRepo.Get(id)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	if contact.UserID != user.UserID {
		err = constants.ErrorNotAuthorized
		ctxSess.ErrorMessage = err.Error()
		return
	}

	err = s.contactsRepo.Delete(id)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}

	return
}
