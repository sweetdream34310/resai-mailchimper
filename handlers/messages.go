package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	"github.com/cloudsrc/api.awaymail.v1.go/models"
	"github.com/cloudsrc/api.awaymail.v1.go/provider"
	"github.com/cloudsrc/api.awaymail.v1.go/provider/google"
	"github.com/cloudsrc/api.awaymail.v1.go/src/domain/redis"
	googleWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	sess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/session"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	messagesSvc "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/messages"
)

type (
	Message struct {
		App           *libs.App
		MessagesSvc   messagesSvc.Service
		GoogleWrapper googleWrapper.Wrapper
		RedisRepo     redis.Repository
	}
	InboxRequest struct {
		Skip string `form:"skip"`
	}
	SendRequest struct {
		To             string   `json:"to" binding:"required"`
		Message        string   `json:"message" binding:"required"`
		Subject        string   `json:"subject" binding:"required"`
		AttachmentsURL []string `json:"attachments_url"`
	}
	UpdateRequest struct {
		IsRead *bool `json:"is_read" binding:"required"`
	}
)

//SetRouter : routers for casbu
func (m *Message) SetRouter() {
	//session := Session{App: m.App}
	sess := sess.Session{App: m.App, GoogleWrapper: m.GoogleWrapper, RedisRepo: m.RedisRepo}
	v3 := m.App.Engine.Group("/v3")
	v3.POST("/archive", sess.CheckToken, m.archiveInbox)
	v3.GET("/archive/inbox", sess.CheckToken, m.getUserInboxV3)
	v3.DELETE("/archive/inbox", sess.CheckToken, m.deleteArchiveInbox)
	//v3.GET("/sent", sess.CheckToken, m.getUserSentV3)
	//v3.GET("/message/:id", sess.CheckToken, m.getMessageV3)
	v3.GET("/sent/message/:id", sess.CheckToken, m.getSentMessageV3)
	v3.POST("/send", sess.CheckToken, m.sendMessageV3)
	v3.PATCH("/message/:id", sess.CheckToken, m.updateMessageV3)
	v3.DELETE("/message/:id", sess.CheckToken, m.deleteMessageV3)
	v3.GET("/user/detail/message/:id", sess.CheckToken, m.GetUserMessageByID)
	v3.GET("/user/message/:label", sess.CheckToken, m.getUserMessage)
	v3.GET("/user/thread/message/:id", sess.CheckToken, m.getThreadMessage)
	v3.GET("/message/attachment", sess.CheckToken, m.getAttachment)
	v3.POST("/user/message/tag", sess.CheckToken, m.addTagMessage)
	v3.POST("/inbox/archive", sess.CheckToken, m.tagArchiveMessage)

	//m.App.Router.GET("/inbox", session.CheckToken, m.getUserInbox)
	//m.App.Router.GET("/sent", session.CheckToken, m.getUserSent)
	//m.App.Router.GET("/message/:id", session.CheckToken, m.getMessage)
	//m.App.Router.POST("/send", session.CheckToken, m.sendMessage)
	//m.App.Router.PATCH("/message/:id", session.CheckToken, m.updateMessage)
}

func (m *Message) getUserMessage(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	var req messagesSvc.UserMessageReq
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.Bind(&req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	req.Label = strings.ToUpper(ctx.Params.ByName("label"))
	if strings.Contains(req.Label, "LABEL") {
		req.Label = strings.Title(strings.ToLower(req.Label))
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := m.MessagesSvc.GetUserMessage(ctxSess, req)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (m *Message) getThreadMessage(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	threadId := strings.ToUpper(ctx.Params.ByName("id"))
	ctxSess.Request = threadId
	ctxSess.Lv1("Incoming message")

	resp, err := m.MessagesSvc.GetThreadByID(ctxSess, threadId)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (m *Message) getAttachment(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	var req messagesSvc.AttachmentRequest
	res := libs.ResponseFile{
		Status: http.StatusOK,
	}
	if err := ctx.Bind(&req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendFleResponse(ctx, "")
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := m.MessagesSvc.GetAttachment(ctxSess, req)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendFleResponse(ctx, "")
		return
	}
	res.Data = resp.FileByte
	//ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", req.FileName))
	res.SendFleResponse(ctx, req.MimeType)
}

func (m *Message) addTagMessage(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	var req messagesSvc.TagLabelReq
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.Bind(&req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	err := m.MessagesSvc.TagLabelUserMessage(ctxSess, req)
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

func (m *Message) tagArchiveMessage(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	var req messagesSvc.TagLabelArchiveReq
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.Bind(&req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	err := m.MessagesSvc.TagLabelArchiveMessage(ctxSess, req)
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

func (m *Message) archiveInbox(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	var req messagesSvc.ArchiveReq
	if err := ctx.Bind(&req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Lv1("Incoming message")

	m.MessagesSvc.ArchiveInbox(ctxSess, req)

	res.SendResponse(ctxSess, ctx)
}

func (m *Message) deleteArchiveInbox(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	ctxSess.Lv1("Incoming message")

	err := m.MessagesSvc.DeleteArchiveInbox(ctxSess)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}

	res.SendResponse(ctxSess, ctx)
}

func (m *Message) getUserInboxV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	var req messagesSvc.InboxReq
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.Bind(&req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := m.MessagesSvc.GetUserInbox(ctxSess, ctxSess.UserSession, req)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	if req.Page == 0 {
		req.Page += 1
	}
	ctxSess.Put("responsePage", resp.Page)
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (m *Message) getUserSentV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	var req messagesSvc.InboxSentReq
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := m.MessagesSvc.GetUserSent(ctxSess, ctxSess.UserSession, req)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	if req.Page == 0 {
		req.Page += 1
	}
	ctxSess.Put("responsePage", resp.Page)
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (m *Message) getMessageV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	messageID := ctx.Params.ByName("id")

	req := &messagesSvc.MessageQueue{
		MessageID: messageID,
		UserIdHex: ctxSess.UserSession.UserID.Hex(),
		UserId:    ctxSess.UserSession.UserID,
		IsRead:    true,
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := m.MessagesSvc.GetMessage(ctxSess, req)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (m *Message) getSentMessageV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	messageID := ctx.Params.ByName("id")

	req := &messagesSvc.MessageQueue{
		MessageID: messageID,
		UserIdHex: ctxSess.UserSession.UserID.Hex(),
		UserId:    ctxSess.UserSession.UserID,
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	resp, err := m.MessagesSvc.GetSentMessage(ctxSess, req)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}

func (m *Message) sendMessageV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	var req messagesSvc.SendMessageReq
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil { //validation error
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	err := m.MessagesSvc.SentMessage(ctxSess, ctxSess.UserSession, req)
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

func (m *Message) updateMessageV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	var req messagesSvc.UpdateRequest
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	messageID := ctx.Params.ByName("id")
	ctxSess.Request = req
	ctxSess.Lv1("Incoming message")

	_, err := m.MessagesSvc.UpdateMessage(ctxSess, messageID, req)
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

func (m *Message) deleteMessageV3(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	messageID := ctx.Params.ByName("id")
	ctxSess.Request = messageID
	ctxSess.Lv1("Incoming message")

	err := m.MessagesSvc.DeleteMessage(ctxSess, messageID)
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

func (m *Message) getUserInbox(ctx *gin.Context) {
	var req InboxRequest
	var response []interface{}
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.ShouldBind(&req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(nil, ctx)
		return
	}
	var user models.User
	prov, _ := ctx.Get("provider")
	provi := prov.(provider.Provider)
	p := ctx.Request.Header.Get("X-Client")
	switch p {
	case "gmail":
		user = prov.(*google.Provider).User
	}
	away := models.Away{
		DB:     m.App.DB.MongoDB,
		UserID: user.ID,
	}
	aways, err := away.GetAways("aways")
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(nil, ctx)
		return
	}
	go provi.GetUserInbox()
	query := "*"
	var userQuery []string
	var subjectQuery []string
	if len(aways) > 0 {
		for _, away := range aways {
			if len(away.AllowedUsers) > 0 && *away.IsEnabled {
				for _, user := range away.AllowedUsers {
					if *user.Activated {
						if user.Name == ".*" {
							goto doneProcessing
						}
						userQuery = append(userQuery, strings.ReplaceAll(user.Name, "@", "_"))
					}
				}
			}
			if len(away.AllowedSubjects) > 0 && *away.IsEnabled {
				for _, subject := range away.AllowedSubjects {
					if *subject.Activated {
						if subject.Name == ".*" {
							goto doneProcessing
						}
						subjectQuery = append(subjectQuery, subject.Name)
					}
				}
			}
		}
		if len(subjectQuery) > 0 {
			query = "@subject|body_text:(" + strings.Join(subjectQuery, "|") + ")"
		}
		if len(userQuery) > 0 && len(subjectQuery) > 0 {
			query = "@from:(" + strings.Join(userQuery, "|") + ")" + " " + query
		} else {
			if len(userQuery) > 0 {
				query = "@from:(" + strings.Join(userQuery, "|") + ")"
			}
		}
	}
doneProcessing:
	skip, _ := strconv.Atoi(req.Skip)
	response = m.App.Redis.SearchIndex(fmt.Sprintf("user:inbox:%s",
		user.ID.Hex()),
		query, "SORTBY", "received_at", "DESC", "RETURN", 1,
		"message", "LIMIT", skip, 10)
	if response == nil {
		inbox := models.Message{
			DB:     m.App.DB.MongoDB,
			UserID: user.ID,
		}
		inboxArr, _ := inbox.GetInbox("messages.inbox", skip)
		res.Data = map[string]interface{}{
			"messages": inboxArr,
			"skip":     strconv.Itoa(skip + 10),
		}
		res.SendResponse(nil, ctx)
		return
	}
	res.Data = map[string]interface{}{
		"messages": response,
		"skip":     strconv.Itoa(skip + 10),
	}
	res.SendResponse(nil, ctx)
}

func (m *Message) getUserSent(ctx *gin.Context) {
	var req InboxRequest
	var response []interface{}
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.ShouldBind(&req); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(nil, ctx)
		return
	}
	var user models.User
	prov, _ := ctx.Get("provider")
	provi := prov.(provider.Provider)
	p := ctx.Request.Header.Get("X-Client")
	switch p {
	case "gmail":
		user = prov.(*google.Provider).User
	}
	away := models.Away{
		DB:     m.App.DB.MongoDB,
		UserID: user.ID,
	}
	aways, err := away.GetAways("aways")
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(nil, ctx)
		return
	}
	go provi.GetUserSent()
	query := "*"
	var userQuery []string
	var subjectQuery []string
	if len(aways) > 0 {
		for _, away := range aways {
			if len(away.AllowedUsers) > 0 && *away.IsEnabled {
				for _, user := range away.AllowedUsers {
					if *user.Activated {
						if user.Name == ".*" {
							goto doneProcessing
						}
						userQuery = append(userQuery, strings.ReplaceAll(user.Name, "@", "_"))
					}
				}
			}
			if len(away.AllowedSubjects) > 0 && *away.IsEnabled {
				for _, subject := range away.AllowedSubjects {
					if *subject.Activated {
						if subject.Name == ".*" {
							goto doneProcessing
						}
						subjectQuery = append(subjectQuery, subject.Name)
					}
				}
			}
		}
		if len(subjectQuery) > 0 {
			query = "@subject|body_text:(" + strings.Join(subjectQuery, "|") + ")"
		}
		if len(userQuery) > 0 && len(subjectQuery) > 0 {
			query = "@from:(" + strings.Join(userQuery, "|") + ")" + " " + query
		} else {
			if len(userQuery) > 0 {
				query = "@from:(" + strings.Join(userQuery, "|") + ")"
			}
		}
	}
doneProcessing:
	skip, _ := strconv.Atoi(req.Skip)
	response = m.App.Redis.SearchIndex(fmt.Sprintf("user:sent:%s",
		user.ID.Hex()),
		query, "SORTBY", "sent", "DESC", "RETURN", 1,
		"message", "LIMIT", skip, 10)
	if response == nil {
		inbox := models.Message{
			DB:     m.App.DB.MongoDB,
			UserID: user.ID,
		}
		inboxArr, _ := inbox.GetInbox("messages.inbox", skip)
		res.Data = map[string]interface{}{
			"messages": inboxArr,
			"skip":     strconv.Itoa(skip + 10),
		}
		res.SendResponse(nil, ctx)
		return
	}
	res.Data = map[string]interface{}{
		"messages": response,
		"skip":     strconv.Itoa(skip + 10),
	}
	res.SendResponse(nil, ctx)
}

func (m *Message) sendMessage(ctx *gin.Context) {
	var req SendRequest
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil { //validation error
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(nil, ctx)
		return
	}
	prov, _ := ctx.Get("provider")
	provi := prov.(provider.Provider)
	message := google.Messsage{
		To:             req.To,
		Message:        req.Message,
		Subject:        req.Subject,
		AttachmentsURL: req.AttachmentsURL,
	}
	if len(req.AttachmentsURL) == 0 {
		if err := provi.SendMessage(message); err != nil {
			res.Message = err.Error()
			res.Status = http.StatusBadRequest
			res.SendResponse(nil, ctx)
			return
		}
	} else {
		if err := provi.SendMessageWithAttachment(message); err != nil {
			res.Message = err.Error()
			res.Status = http.StatusBadRequest
			res.SendResponse(nil, ctx)
			return
		}
	}
	res.Data = map[string]interface{}{
		"message": "success",
	}
	res.SendResponse(nil, ctx)
}

func (m *Message) getMessage(ctx *gin.Context) {
	res := libs.Response{
		Status: http.StatusOK,
	}
	messageID := ctx.Params.ByName("id")
	prov, _ := ctx.Get("provider")
	provi := prov.(provider.Provider)
	message, err := provi.GetMessage(messageID)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(nil, ctx)
		return
	}
	res.Data = message
	res.SendResponse(nil, ctx)
}

func (m *Message) updateMessage(ctx *gin.Context) {
	var req UpdateRequest
	res := libs.Response{
		Status: http.StatusOK,
	}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(nil, ctx)
		return
	}
	messageID := ctx.Params.ByName("id")
	prov, _ := ctx.Get("provider")
	provi := prov.(provider.Provider)
	message, err := provi.UpdateMessage(messageID, req.IsRead)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(nil, ctx)
		return
	}
	res.Data = message
	res.SendResponse(nil, ctx)
}

func (m *Message) GetUserMessageByID(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	res := libs.Response{
		Status: http.StatusOK,
	}
	messageID := ctx.Params.ByName("id")
	ctxSess.Request = messageID
	ctxSess.Lv1("Incoming message")

	resp, err := m.MessagesSvc.GetUserMessageByID(ctxSess, messageID)
	if err != nil {
		res.Message = err.Error()
		res.Status = http.StatusBadRequest
		res.SendResponse(ctxSess, ctx)
		return
	}
	res.Data = resp
	res.SendResponse(ctxSess, ctx)
}
