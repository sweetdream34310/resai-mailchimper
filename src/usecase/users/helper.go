package users

import (
	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/users"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/people/v1"
)

var scopes = []string{
	people.ContactsReadonlyScope,
	gmail.GmailModifyScope,
	calendar.CalendarReadonlyScope,
}

func (s *service) getToken(ctxSess *ctxSess.Context, agent string, req *GetAuthTokenReq) (userID primitive.ObjectID, email, authToken, name, photo string, isNew bool, err error) {
	switch req.Provider {
	case constants.GmailClient:
		userID, email, authToken, name, photo, isNew, err = s.gmailToken(ctxSess, agent, req)
	default:
		err = constants.ErrorClientNotSupported
		ctxSess.ErrorMessage = err.Error()
	}

	return
}

func (s *service) gmailToken(ctxSess *ctxSess.Context, agent string, req *GetAuthTokenReq) (userID primitive.ObjectID, email, authToken, name, photo string, isNew bool, err error) {
	profile, err := s.googleWrapper.GetProfile(ctxSess, agent, req.RefreshToken)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}

	encAuthToken, err := utils.Encrypt(s.config.Salt, profile.AuthToken)
	encRefreshToken, err := utils.Encrypt(s.config.Salt, profile.RefreshToken)
	encName, err := utils.Encrypt(s.config.Salt, profile.Name)
	encPhoto, err := utils.Encrypt(s.config.Salt, profile.Photo)

	entity, err := s.userRepo.GetUser(profile.Email)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	if entity == nil {
		usr := &domain.User{
			Email:        profile.Email,
			AuthToken:    encAuthToken,
			RefreshToken: encRefreshToken,
			Name:         encName,
			Photo:        encPhoto,
		}
		if req.SwiftToken != "" {
			usr.SwiftToken = []string{req.SwiftToken}
		}
		entity, err = s.userRepo.AddUser(usr)
		if err == nil {
			isNew = true
		}
	} else {
		var flag bool
		for _, each := range entity.SwiftToken {
			if each == req.SwiftToken {
				flag = true
				break
			}
		}

		if !flag && req.SwiftToken != "" {
			entity.SwiftToken = append(entity.SwiftToken, req.SwiftToken)
			err = s.userRepo.UpdateUser(profile.Email, encAuthToken, encRefreshToken, encName, encPhoto, entity.SwiftToken)
		}
	}
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}
	userID = entity.ID
	email = profile.Email
	authToken = encAuthToken
	name = profile.Name
	photo = profile.Photo

	return
}

func removeSwiftToken(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
