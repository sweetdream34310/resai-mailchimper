package users

import (
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	GetAuthToken(ctxSess *ctxSess.Context, req *GetAuthTokenReq, agent string) (resp GetAuthTokenResp, newUserID primitive.ObjectID, err error)
	GetUserProfile(ctxSess *ctxSess.Context) (resp *GetUserProfileResp, err error)
	Logout(ctxSess *ctxSess.Context)
}
