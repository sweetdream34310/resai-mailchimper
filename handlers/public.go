package handlers

import (
	"fmt"
	messagesSvc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/messages"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/satori/uuid"

	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
)

type (
	Public struct {
		App         *libs.App
		MessagesSvc messagesSvc.Service
	}
)

//SetRouter : routers for casbu
func (p *Public) SetRouter() {
	p.App.Router.Static("/file", "./public")
	p.App.Router.POST("/file", p.uploadFile)
	p.App.Router.GET("/breakthrough/notification", p.breakthrough)
	p.App.Router.GET("/breakthrough/amp/notification", p.ampbreakthrough)
}

func (p *Public) ampbreakthrough(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	var req messagesSvc.BreakthroughReq
	if err := ctx.Bind(&req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}

	go func() {
		er := p.MessagesSvc.BreakthroughNotification(ctxSess, req)
		if er != nil {
			ctxSess.ErrorMessage = er.Error()
			ctxSess.Lv4()
		}
	}()

	res.SendResponse(ctxSess, ctx)
}

func (p *Public) breakthrough(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	var req messagesSvc.BreakthroughReq
	if err := ctx.Bind(&req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}

	go func() {
		er := p.MessagesSvc.BreakthroughNotification(ctxSess, req)
		if er != nil {
			ctxSess.ErrorMessage = er.Error()
			ctxSess.Lv4()
		}
	}()

	ctx.Redirect(http.StatusFound, "https://www.awaymail.net/breakthrough")
}

func (p *Public) uploadFile(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	ctxSess.Lv1("Incoming message")
	bucketConn := p.App.Bucket
	bucket := bucketConn.Connection.Bucket(bucketConn.BucketName)
	_, err := bucket.Attrs(bucketConn.Context)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusInternalServerError
		res.SendResponse(ctxSess, ctx)
		return
	}
	f, uploadedFile, err := ctx.Request.FormFile("file")
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusInternalServerError
		res.SendResponse(ctxSess, ctx)
		return
	}
	fileName := uuid.NewV4().String() + filepath.Ext(uploadedFile.Filename)
	sw := bucket.Object(fileName).NewWriter(bucketConn.Context)
	if _, err = io.Copy(sw, f); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusInternalServerError
		res.SendResponse(ctxSess, ctx)
		return
	}
	if err = sw.Close(); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusInternalServerError
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = map[string]interface{}{
		"url": fmt.Sprintf("https://%v.storage.googleapis.com/%v", bucketConn.BucketName, sw.Attrs().Name),
	}
	res.SendResponse(ctxSess, ctx)
}
