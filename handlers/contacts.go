package handlers

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	"github.com/cloudsrc/api.awaymail.v1.go/src/domain/redis"
	googleWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	sess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/session"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	svc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/contacts"
)

type Contacts struct {
	App           *libs.App
	ContactsSvc   svc.Service
	GoogleWrapper googleWrapper.Wrapper
	RedisRepo     redis.Repository
}

//SetRouter : routers for casbu
func (c *Contacts) SetRouter() {
	sess := sess.Session{App: c.App, GoogleWrapper: c.GoogleWrapper, RedisRepo: c.RedisRepo}
	v3 := c.App.Engine.Group("/v3")
	v3.POST("/contacts", sess.CheckToken, c.addContactsV3)
	v3.PUT("/contacts/:id", sess.CheckToken, c.updateContactsV3)
	v3.DELETE("/contacts/:id", sess.CheckToken, c.deleteContactsV3)
	v3.GET("/contacts", sess.CheckToken, c.getContactsListV3)
	v3.GET("/contacts/recent", sess.CheckToken, c.getRecentContactsListV3)
}

func (c *Contacts) addContactsV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status:  http.StatusOK,
		Message: "success",
	}
	var req svc.ContactAddReq
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := c.ContactsSvc.AddContacts(ctxSess, req, ctxSess.UserSession)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}

	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (c *Contacts) updateContactsV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status:  http.StatusOK,
		Message: "success",
	}
	var req svc.ContactUpdateReq
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	req.ID = ctx.Params.ByName("id")
	if len(req.ID) == 0 {
		res.Message = "missing contact id"
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := c.ContactsSvc.UpdateContacts(ctxSess, req, ctxSess.UserSession)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}

	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (c *Contacts) deleteContactsV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status:  http.StatusOK,
		Message: "success",
	}
	contactsID := ctx.Params.ByName("id")
	if len(contactsID) == 0 {
		res.Message = "missing contact id"
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = contactsID
	ctxSess.Lv1("Incoming message")

	err := c.ContactsSvc.DeleteContacts(ctxSess, contactsID, ctxSess.UserSession)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.SendResponse(ctxSess, ctx)
}

func (c *Contacts) getContactsListV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status:  http.StatusOK,
		Message: "success",
	}
	ctxSess.Request = ctxSess.UserSession.UserID
	ctxSess.Lv1("Incoming message")

	resp, err := c.ContactsSvc.GetContactsList(ctxSess, ctxSess.UserSession)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}
	if reflect.ValueOf(resp).IsNil() {
		resp = []struct{}{}
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (c *Contacts) getRecentContactsListV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status:  http.StatusOK,
		Message: "success",
	}
	ctxSess.Request = ctxSess.UserSession.UserID
	ctxSess.Lv1("Incoming message")

	resp, err := c.ContactsSvc.GetRecentContactsList(ctxSess, ctxSess.UserSession)
	if err != nil {
		res.Status, res.Message = errHandler(err)
		res.SendResponse(ctxSess, ctx)
		return
	}
	if reflect.ValueOf(resp).IsNil() {
		resp = []struct{}{}
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}
