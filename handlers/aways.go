package handlers

import (
	"net/http"
	"time"

	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	"github.com/cloudsrc/api.awaymail.v1.go/models"
	"github.com/cloudsrc/api.awaymail.v1.go/src/domain/redis"
	googleWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	sess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/session"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	svc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/aways"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Away struct {
		App           *libs.App
		AwaySvc       svc.Service
		GoogleWrapper googleWrapper.Wrapper
		RedisRepo     redis.Repository
	}
	Request struct {
		ID              string           `json:"_id,omitempty"`
		Title           string           `json:"title" binding:"required"`
		ActivateAllow   time.Time        `json:"activate_allow"`
		DeactivateAllow time.Time        `json:"deactivate_allow"`
		Repeat          []string         `json:"repeat"`
		IsEnabled       *bool            `json:"is_enabled"`
		AllDay          *bool            `json:"all_day"`
		KeyDuration     string           `json:"key_duration"`
		AllowedUsers    []models.Allowed `json:"allowed_users"`
		AllowedSubjects []models.Allowed `json:"allowed_subjects"`
	}
)

// SetRouter : routers for casbu
func (a *Away) SetRouter() {
	//session := Session{App: a.App}
	sess := sess.Session{App: a.App, GoogleWrapper: a.GoogleWrapper, RedisRepo: a.RedisRepo}
	v3 := a.App.Engine.Group("/v3")
	v3.POST("/away", sess.CheckToken, a.createAwayV3)
	v3.GET("/aways", sess.CheckToken, a.getUserAwaysV3)
	v3.GET("/away/:id", sess.CheckToken, a.getUserAwayV3)
	v3.PUT("/away", sess.CheckToken, a.updateAwayV3)
	v3.PATCH("/away/:id", sess.CheckToken, a.enableAwayV3)
	v3.PATCH("/away", sess.CheckToken, a.enableAwayModeV3)
	v3.DELETE("/away/:id", sess.CheckToken, a.deleteAwayV3)
}

func (a *Away) createAwayV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	req := &svc.CreateAwayReq{}
	if err := ctx.ShouldBindBodyWith(req, binding.JSON); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := a.AwaySvc.CreateAway(ctxSess, ctxSess.UserSession, req)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}

	res.Data = resp

	// send response
	res.SendResponse(ctxSess, ctx)

}

func (a *Away) getUserAwayV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	awayID := ctx.Params.ByName("id")
	ctxSess.Request = awayID
	ctxSess.Lv1("Incoming message")

	resp, err := a.AwaySvc.GetAway(ctxSess, awayID)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (a *Away) getUserAwaysV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	var req svc.EnableAwayReq
	if err := ctx.Bind(&req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = ctxSess.UserSession.UserID
	ctxSess.Lv1("Incoming message")

	resp, err := a.AwaySvc.GetAwayList(ctxSess, ctxSess.UserSession, req)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (a *Away) deleteAwayV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}

	awayID := ctx.Params.ByName("id")
	ctxSess.Request = awayID
	ctxSess.Lv1("Incoming message")

	if err := a.AwaySvc.DeleteAway(ctxSess, ctxSess.UserSession, awayID); err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Message = "success"
	res.SendResponse(ctxSess, ctx)
}

func (a *Away) updateAwayV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}

	req := &svc.UpdateAwayReq{}
	if err := ctx.ShouldBindBodyWith(req, binding.JSON); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := a.AwaySvc.UpdateAway(ctxSess, ctxSess.UserSession, req)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (a *Away) enableAwayModeV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	req := &svc.EnableAwayReq{}
	err := ctx.ShouldBindBodyWith(req, binding.JSON)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	total, err := a.AwaySvc.EnableAllAway(ctxSess, ctxSess.UserSession, req.IsEnabled)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Message = "success"
	res.Data = map[string]interface{}{
		"total": total,
	}
	res.SendResponse(ctxSess, ctx)
}

func (a *Away) enableAwayV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	awayID, _ := primitive.ObjectIDFromHex(ctx.Params.ByName("id"))
	req := &svc.EnableAwayReq{}
	if err := ctx.ShouldBindBodyWith(req, binding.JSON); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := a.AwaySvc.EnableAway(ctxSess, ctxSess.UserSession, req.IsEnabled, awayID)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}
