package users

import (
	"fmt"
	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"github.com/cloudsrc/api.awaymail.v1.go/src/domain/redis"
	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/users"
	"github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/session"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	config        config.Config
	userRepo      domain.Repository
	googleWrapper google.Wrapper
	redisRepo     redis.Repository
}

func New(cfg config.Config, userRepo domain.Repository, googleWrapper google.Wrapper, redisRepo redis.Repository) Service {
	if userRepo == nil {
		panic("contacts repository is nil")
	}

	if googleWrapper == nil {
		panic("google wrapper is nil")
	}

	return &service{
		config:        cfg,
		userRepo:      userRepo,
		googleWrapper: googleWrapper,
		redisRepo:     redisRepo,
	}
}

func (s *service) GetAuthToken(ctxSess *ctxSess.Context, req *GetAuthTokenReq, agent string) (resp GetAuthTokenResp, newUserID primitive.ObjectID, err error) {
	userID, email, authToken, name, photo, isNew, err := s.getToken(ctxSess, agent, req)
	if err != nil {
		return
	}

	if isNew {
		newUserID = userID
	}

	encRefreshToken, err := utils.Encrypt(s.config.Salt, req.RefreshToken)
	resp.Token, resp.Expiry = session.NewBearerToken(userID, email, req.Provider, encRefreshToken, authToken, name, photo, req.SwiftToken, req.DevTokenKey)

	key := fmt.Sprintf(constants.UserSession, userID.Hex(), req.SwiftToken)
	s.redisRepo.SetKey(key, userID.Hex())
	//if os.Getenv("GOLANG_ENV") != "prod" {
	//	key = fmt.Sprintf(constants.UserTokenDev, userID.Hex())
	//	s.redisRepo.SetKey(key, req.DevTokenKey)
	//}

	return
}

func (s *service) Logout(ctxSess *ctxSess.Context) {
	if ctxSess.UserSession.SwiftToken != "" {
		entity, _ := s.userRepo.GetUserByID(ctxSess.UserSession.UserID)
		if entity == nil {
			return
		}
		for _, each := range entity.SwiftToken {
			if each == ctxSess.UserSession.SwiftToken {
				removeSwiftToken(entity.SwiftToken, each)
				s.userRepo.UpdateSwiftToken(ctxSess.UserSession.Email, entity.SwiftToken)
			}
		}
	}

	key := fmt.Sprintf(constants.UserSession, ctxSess.UserSession.UserID.Hex(), ctxSess.UserSession.SwiftToken)
	s.redisRepo.DelKey(key)
}

func (s *service) GetUserProfile(ctxSess *ctxSess.Context) (resp *GetUserProfileResp, err error) {
	user := ctxSess.UserSession
	if len(user.UserID.Hex()) == 0 {
		err = constants.ErrorDataNotFound
		return
	}
	resp = &GetUserProfileResp{
		Email: user.Email,
		Name:  user.Name,
		Photo: user.Photo,
	}

	return
}
