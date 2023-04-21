package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	"github.com/cloudsrc/api.awaymail.v1.go/src/domain/redis"
	googleWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	sess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/session"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	labelSvc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/label"
)

type Label struct {
	App           *libs.App
	LabelSvc      labelSvc.Service
	GoogleWrapper googleWrapper.Wrapper
	RedisRepo     redis.Repository
}

func (m *Label) SetRouter() {
	sess := sess.Session{App: m.App, GoogleWrapper: m.GoogleWrapper, RedisRepo: m.RedisRepo}
	v3 := m.App.Engine.Group("/v3")
	v3.GET("/label", sess.CheckToken, m.findAllLabels)
	v3.GET("/label/:id", sess.CheckToken, m.getLabel)
	v3.POST("/label", sess.CheckToken, m.createLabel)
	v3.PATCH("/label/:id", sess.CheckToken, m.patchLabel)
	v3.DELETE("/label/:id", sess.CheckToken, m.deleteLabel)
}

func (m *Label) getLabel(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	labelID := strings.ToUpper(ctx.Params.ByName("id"))
	if strings.Contains(labelID, "LABEL") {
		labelID = strings.Title(strings.ToLower(labelID))
	}

	ctxSess.Request = labelID
	ctxSess.Lv1("Incoming message")

	resp, err := m.LabelSvc.GetLabel(ctxSess, labelID)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (m *Label) findAllLabels(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	ctxSess.Lv1("Incoming message")

	resp, err := m.LabelSvc.FindAllLabels(ctxSess)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (m *Label) createLabel(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	req := &labelSvc.CreateLabelReq{}
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.Bind(req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := m.LabelSvc.CreateLabel(ctxSess, req)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (m *Label) patchLabel(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	req := &labelSvc.PatchLabelReq{}
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.Bind(req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	labelID := strings.ToUpper(ctx.Params.ByName("id"))
	if strings.Contains(labelID, "LABEL") {
		labelID = strings.Title(strings.ToLower(labelID))
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := m.LabelSvc.PatchLabel(ctxSess, labelID, req)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (m *Label) deleteLabel(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	labelID := strings.ToUpper(ctx.Params.ByName("id"))
	if strings.Contains(labelID, "LABEL") {
		labelID = strings.Title(strings.ToLower(labelID))
	}
	ctxSess.Request = labelID
	ctxSess.Lv1("Incoming message")

	err := m.LabelSvc.DeleteLabel(ctxSess, labelID)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = map[string]interface{}{
		"message": "success",
	}
	res.SendResponse(ctxSess, ctx)
}
