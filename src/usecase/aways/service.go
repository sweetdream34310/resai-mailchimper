package aways

import (
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	CreateAway(ctxSess *ctxSess.Context, user models.UserSession, req *CreateAwayReq) (resp *CreateAwayResp, err error)
	GetAwayList(ctxSess *ctxSess.Context, user models.UserSession, req EnableAwayReq) (resp []*CreateAwayResp, err error)
	GetAway(ctxSess *ctxSess.Context, awayID string) (resp *CreateAwayResp, err error)
	UpdateAway(ctxSess *ctxSess.Context, user models.UserSession, req *UpdateAwayReq) (resp *CreateAwayResp, err error)
	EnableAway(ctxSess *ctxSess.Context, user models.UserSession, enable bool, awayID primitive.ObjectID) (resp *CreateAwayResp, err error)
	EnableAllAway(ctxSess *ctxSess.Context, user models.UserSession, isEnabled bool) (total int64, err error)
	DeleteAway(ctxSess *ctxSess.Context, user models.UserSession, awayID string) (err error)
}
