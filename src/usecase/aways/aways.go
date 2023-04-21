package aways

import (
	"fmt"
	"net/http"
	"time"

	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/aways"
	"github.com/cloudsrc/api.awaymail.v1.go/src/domain/redis"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	awayRepo  domain.Repository
	redisRepo redis.Repository
}

func New(awayRepo domain.Repository, redisRepo redis.Repository) Service {
	if awayRepo == nil {
		panic("contacts repository is nil")
	}

	return &service{
		awayRepo:  awayRepo,
		redisRepo: redisRepo,
	}
}

func (s *service) CreateAway(ctxSess *ctxSess.Context, user models.UserSession, req *CreateAwayReq) (resp *CreateAwayResp, err error) {
	if req.ActivateAllow.IsZero() || req.DeactivateAllow.IsZero() {
		err = constants.ErrorInvalidRequest
		ctxSess.ErrorMessage = err.Error()
		return
	}

	if len(req.Repeat) > 0 {
		req.Repeat = utils.RemoveDuplicateStr(req.Repeat)
		for _, each := range req.Repeat {
			if notDay := utils.ValidateWeekday(each); !notDay {
				err = constants.NewError(http.StatusBadRequest, fmt.Sprintf("invlidate date: %s", each))
				ctxSess.ErrorMessage = err.Error()
				return
			}
		}
	}

	// check empty string from an array
	allowedContacts := RemoveEmptyStrings(req.AllowedContacts)
	allowedKeywords := RemoveEmptyStrings(req.AllowedKeywords)

	away, err := s.awayRepo.CreateAway(&domain.Away{
		ID:              primitive.NewObjectID(),
		UserID:          user.UserID,
		Title:           req.Title,
		IsEnabled:       req.IsEnabled,
		ActivateAllow:   req.ActivateAllow,
		DeactivateAllow: req.DeactivateAllow,
		Repeat:          req.Repeat,
		AllDay:          req.AllDay,
		Message:         req.Message,
		KeyDuration:     req.KeyDuration,
		AllowedContacts: allowedContacts,
		AllowedKeywords: allowedKeywords,
	})
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	resp = &CreateAwayResp{
		ID:              away.ID,
		Title:           away.Title,
		IsEnabled:       away.IsEnabled,
		ActivateAllow:   away.ActivateAllow,
		DeactivateAllow: away.DeactivateAllow,
		Repeat:          away.Repeat,
		AllDay:          away.AllDay,
		Message:         away.Message,
		KeyDuration:     away.KeyDuration,
		AllowedContacts: away.AllowedContacts,
		AllowedKeywords: away.AllowedKeywords,
	}

	return
}

func (s *service) GetAwayList(ctxSess *ctxSess.Context, user models.UserSession, req EnableAwayReq) (resp []*CreateAwayResp, err error) {
	entity, err := s.awayRepo.GetAwayList(user.UserID, req.IsEnabled)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDataNotFound
		return
	}

	for _, each := range entity {
		resp = append(resp, &CreateAwayResp{
			ID:              each.ID,
			Title:           each.Title,
			IsEnabled:       each.IsEnabled,
			ActivateAllow:   each.ActivateAllow,
			DeactivateAllow: each.DeactivateAllow,
			Repeat:          each.Repeat,
			AllDay:          each.AllDay,
			Message:         each.Message,
			KeyDuration:     each.KeyDuration,
			AllowedContacts: each.AllowedContacts,
			AllowedKeywords: each.AllowedKeywords,
		})
	}

	return
}

func (s *service) GetAway(ctxSess *ctxSess.Context, awayID string) (resp *CreateAwayResp, err error) {
	docID, err := primitive.ObjectIDFromHex(awayID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDataNotFound
		return
	}
	entity, err := s.awayRepo.GetAway(docID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDataNotFound
		return
	}

	resp = &CreateAwayResp{
		ID:              entity.ID,
		Title:           entity.Title,
		IsEnabled:       entity.IsEnabled,
		ActivateAllow:   entity.ActivateAllow,
		DeactivateAllow: entity.DeactivateAllow,
		Repeat:          entity.Repeat,
		AllDay:          entity.AllDay,
		Message:         entity.Message,
		KeyDuration:     entity.KeyDuration,
		AllowedContacts: entity.AllowedContacts,
		AllowedKeywords: entity.AllowedKeywords,
	}

	return
}

func (s *service) UpdateAway(ctxSess *ctxSess.Context, user models.UserSession, req *UpdateAwayReq) (resp *CreateAwayResp, err error) {
	entity, err := s.awayRepo.GetAway(req.ID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDataNotFound
		return
	}

	if entity.UserID != user.UserID {
		err = constants.ErrorNotAuthorized
		ctxSess.ErrorMessage = err.Error()
		return
	}

	// check empty string from an array
	allowedContacts := RemoveEmptyStrings(req.AllowedContacts)
	allowedKeywords := RemoveEmptyStrings(req.AllowedKeywords)

	entity.Title = req.Title
	entity.IsEnabled = req.IsEnabled
	entity.ActivateAllow = req.ActivateAllow
	entity.DeactivateAllow = req.DeactivateAllow
	entity.Repeat = req.Repeat
	entity.AllDay = req.AllDay
	entity.Message = req.Message
	entity.KeyDuration = req.KeyDuration
	entity.AllowedContacts = allowedContacts
	entity.AllowedKeywords = allowedKeywords

	err = s.awayRepo.UpdateAway(entity)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	resp = &CreateAwayResp{
		ID:              entity.ID,
		Title:           entity.Title,
		IsEnabled:       entity.IsEnabled,
		ActivateAllow:   entity.ActivateAllow,
		DeactivateAllow: entity.DeactivateAllow,
		Repeat:          entity.Repeat,
		AllDay:          entity.AllDay,
		Message:         entity.Message,
		KeyDuration:     entity.KeyDuration,
		AllowedContacts: entity.AllowedContacts,
		AllowedKeywords: entity.AllowedKeywords,
	}

	return
}

func (s *service) EnableAway(ctxSess *ctxSess.Context, user models.UserSession, isEnabled bool, awayID primitive.ObjectID) (resp *CreateAwayResp, err error) {
	entity, err := s.awayRepo.GetAway(awayID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDataNotFound
		return
	}

	if entity.UserID != user.UserID {
		err = constants.ErrorNotAuthorized
		ctxSess.ErrorMessage = err.Error()
		return
	}

	// check empty string from an array
	allowedContacts := RemoveEmptyStrings(entity.AllowedContacts)
	allowedKeywords := RemoveEmptyStrings(entity.AllowedKeywords)

	err = s.awayRepo.UpdateAway(&domain.Away{
		ID:              entity.ID,
		UserID:          entity.UserID,
		Title:           entity.Title,
		IsEnabled:       isEnabled,
		ActivateAllow:   entity.ActivateAllow,
		DeactivateAllow: entity.DeactivateAllow,
		Repeat:          entity.Repeat,
		AllDay:          entity.AllDay,
		Message:         entity.Message,
		KeyDuration:     entity.KeyDuration,
		AllowedContacts: allowedContacts,
		AllowedKeywords: allowedKeywords,
	})
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	resp = &CreateAwayResp{
		ID:              entity.ID,
		Title:           entity.Title,
		IsEnabled:       entity.IsEnabled,
		ActivateAllow:   entity.ActivateAllow,
		DeactivateAllow: entity.DeactivateAllow,
		Repeat:          entity.Repeat,
		AllDay:          entity.AllDay,
		Message:         entity.Message,
		KeyDuration:     entity.KeyDuration,
		AllowedContacts: allowedContacts,
		AllowedKeywords: allowedKeywords,
	}

	key := fmt.Sprintf(constants.AwayModeStart, ctxSess.UserSession.Email)
	if isEnabled {
		s.redisRepo.SetKey(key, time.Now().Unix())
	} else {
		s.redisRepo.DelKey(key)
	}

	s.redisRepo.DelKey(fmt.Sprintf(constants.ListBreakthrough, ctxSess.UserSession.Email))

	return
}

func (s *service) EnableAllAway(ctxSess *ctxSess.Context, user models.UserSession, isEnabled bool) (total int64, err error) {
	total, err = s.awayRepo.UpdateAwayMode(isEnabled, user.UserID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	key := fmt.Sprintf(constants.AwayModeStart, ctxSess.UserSession.Email)
	if isEnabled {
		s.redisRepo.SetKey(key, time.Now().Unix())
	} else {
		s.redisRepo.DelKey(key)
	}

	s.redisRepo.DelKey(fmt.Sprintf(constants.ListBreakthrough, ctxSess.UserSession.Email))

	return
}

func (s *service) DeleteAway(ctxSess *ctxSess.Context, user models.UserSession, awayID string) (err error) {
	docID, err := primitive.ObjectIDFromHex(awayID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDataNotFound
		return
	}
	entity, err := s.awayRepo.GetAway(docID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	if entity.UserID != user.UserID {
		err = constants.ErrorNotAuthorized
		ctxSess.ErrorMessage = err.Error()
		return
	}

	err = s.awayRepo.DeleteAway(entity.ID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	return
}
