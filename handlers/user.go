package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/people/v1"

	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	"github.com/cloudsrc/api.awaymail.v1.go/queue"
	"github.com/cloudsrc/api.awaymail.v1.go/src/domain/redis"
	googleWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	sess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/session"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	svc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/users"
)

type User struct {
	App           *libs.App
	UserSvc       svc.Service
	Queue         queue.Messages
	GoogleWrapper googleWrapper.Wrapper
	RedisRepo     redis.Repository
}

var scopes = []string{
	people.ContactsReadonlyScope,
	gmail.GmailModifyScope,
	calendar.CalendarReadonlyScope,
}

func (u *User) SetRouter() {
	sess := sess.Session{App: u.App, GoogleWrapper: u.GoogleWrapper, RedisRepo: u.RedisRepo}
	v3 := u.App.Engine.Group("/v3")
	v3.POST("/session", u.getAuthTokenV3)
	v3.GET("/user/profile", sess.CheckToken, u.getProfile)
	v3.GET("/logout", sess.CheckToken, u.logout)
}

func (u *User) getAuthTokenV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	agent := ctx.Request.Header.Get("X-agent")
	req := &svc.GetAuthTokenReq{}
	if err := ctx.ShouldBindBodyWith(req, binding.JSON); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, newUserID, err := u.UserSvc.GetAuthToken(ctxSess, req, agent)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}
	if newUserID.Hex() != "000000000000000000000000" {
		go u.Queue.RunOneID(newUserID)
	}

	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (u *User) getProfile(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	ctxSess.Lv1("Incoming message")

	resp, err := u.UserSvc.GetUserProfile(ctxSess)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}

	res.Data = resp
	res.Message = "success"
	res.SendResponse(ctxSess, ctx)
}

func (u *User) logout(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	ctxSess.Lv1("Incoming message")

	u.UserSvc.Logout(ctxSess)
	res.Message = "success"
	res.SendResponse(ctxSess, ctx)
}

func errHandler(err error) (code int, message string) {
	if errResp, ok := err.(*constants.ApplicationError); ok {
		code = errResp.Code
		message = errResp.Message
	} else {
		errHandler(constants.ErrorGeneral)
	}
	return
}
