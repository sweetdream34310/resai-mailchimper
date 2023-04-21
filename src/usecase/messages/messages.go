package messages

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/aways"
	domainCts "github.com/cloudsrc/api.awaymail.v1.go/src/domain/contacts"
	domainMsg "github.com/cloudsrc/api.awaymail.v1.go/src/domain/messages"
	"github.com/cloudsrc/api.awaymail.v1.go/src/domain/redis"
	userDomain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/users"
	chatgptWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/chatgpt"
	"github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/gaurun"
	googleWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	mailWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/mail"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"github.com/mailgun/mailgun-go/v4"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/gmail/v1"
	"gopkg.in/mgo.v2"
)

type service struct {
	config         config.Config
	awayRepo       domain.Repository
	contactsRepo   domainCts.Repository
	messageRepo    domainMsg.Repository
	userRepo       userDomain.Repository
	googleWrapper  googleWrapper.Wrapper
	chatgptWrapper chatgptWrapper.Wrapper
	gaurunWrapper  gaurun.Wrapper
	redisRepo      redis.Repository
	rabbit         *libs.RabbitClient
	mailGun        *mailgun.MailgunImpl
}

func New(config config.Config, awayRepo domain.Repository, messageRepo domainMsg.Repository, userRepo userDomain.Repository,
	contactsRepo domainCts.Repository, googleWrapper googleWrapper.Wrapper, chatgptWrapper chatgptWrapper.Wrapper, gaurunWrapper gaurun.Wrapper, redisRepo redis.Repository,
	rabbit *libs.RabbitClient, mailGun *mailgun.MailgunImpl) Service {
	if awayRepo == nil {
		panic("contacts repository is nil")
	}
	if googleWrapper == nil {
		panic("google wrapper is nil")
	}
	if chatgptWrapper == nil {
		panic("chatgpt wrapper is nil")
	}
	if mailGun == nil {
		panic("mailGun is nil")
	}

	return &service{
		config:         config,
		awayRepo:       awayRepo,
		messageRepo:    messageRepo,
		userRepo:       userRepo,
		contactsRepo:   contactsRepo,
		googleWrapper:  googleWrapper,
		chatgptWrapper: chatgptWrapper,
		gaurunWrapper:  gaurunWrapper,
		redisRepo:      redisRepo,
		rabbit:         rabbit,
		mailGun:        mailGun,
	}
}

func (s *service) SentMessage(ctxSess *ctxSess.Context, user models.UserSession, req SendMessageReq) (err error) {
	message := googleWrapper.Messsage{
		To:             req.To,
		Message:        req.Message,
		Subject:        req.Subject,
		AttachmentsURL: req.AttachmentsURL,
		Email:          user.Email,
		ThreadID:       req.ThreadID,
		MessageID:      req.MessageID,
	}
	if len(req.AttachmentsURL) == 0 {
		err = s.googleWrapper.SendMessage(ctxSess, message)
	} else {
		err = s.googleWrapper.SendMessageWithAttachment(ctxSess, message)
	}
	if err != nil {
		err = constants.ErrorGeneral
		return
	}

	key := fmt.Sprintf("recent:contact:%s", user.UserID.Hex())
	go s.redisRepo.PushCache(key, req.To)

	return
}

func (s *service) UpdateMessage(ctxSess *ctxSess.Context, messageId string, req UpdateRequest) (resp *gmail.Message, err error) {
	resp, err = s.googleWrapper.UpdateMessage(ctxSess, messageId, req.IsRead)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	//getMessageReq := &MessageQueue{
	//	MessageID: messageId,
	//	UserIdHex: ctxSess.UserSession.UserID.Hex(),
	//	UserId:    ctxSess.UserSession.UserID,
	//	IsRead:    *req.IsRead,
	//}
	//go s.GetMessage(ctxSess, getMessageReq)
	return
}

func (s *service) DeleteMessage(ctxSess *ctxSess.Context, messageId string) (err error) {
	err = s.googleWrapper.DeleteMessage(ctxSess, messageId)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	//cQuery := "user:inbox:" + ctxSess.UserSession.UserID.Hex() + ":" + messageId
	//err = s.redisRepo.DelCache(cQuery, "message", "body_text", "from", "subject", "received_at")
	//if err != nil {
	//	ctxSess.ErrorMessage = err.Error()
	//}
	return
}

func (s *service) ArchiveInbox(ctxSess *ctxSess.Context, req ArchiveReq) {
	_ = s.userRepo.UpdateFlagArchive(ctxSess.UserSession.Email, req.Archive)
	if req.Archive {
		s.processArchive(ctxSess)
	}
	return
}

func (s *service) DeleteArchiveInbox(ctxSess *ctxSess.Context) (err error) {
	err = s.messageRepo.DeleteArchiveByEmail(ctxSess.UserSession.Email)
	return
}

func (s *service) GetUserInbox(ctxSess *ctxSess.Context, user models.UserSession, req InboxReq) (resp InboxResp, err error) {
	go func() {
		u, _ := s.userRepo.GetUser(ctxSess.UserSession.Email)
		if u != nil {
			if u.Archive {
				s.processArchive(ctxSess)
			}
		}
	}()

	query := "*"
	var subjectQuery []string

	if len(strings.TrimSpace(req.Search)) > 0 {
		subjectQuery = append(subjectQuery, strings.ReplaceAll(req.Search, "@", "_"))
		query = "@from|subject|body_text:(" + strings.Join(subjectQuery, "|") + ")"
	}

	var offset int64
	if req.Page > 0 {
		req.Page = req.Page - 1
	}
	offset = req.Page * 10
	resp.Messages = s.redisRepo.SearchIndex(fmt.Sprintf("user:inbox:%s",
		user.UserID.Hex()),
		query, "SORTBY", "received_at", "DESC", "RETURN", 1,
		"message", "LIMIT", offset, 10)

	total := s.redisRepo.AggregateIndex(fmt.Sprintf("user:inbox:%s",
		user.UserID.Hex()),
		query)

	totalEmails := total
	totalPage := int64(math.Ceil(float64(totalEmails) / 10))

	resp.Page = Page{
		CurrentPage: req.Page + 1,
		TotalPage:   totalPage,
		TotalEmails: totalEmails,
	}

	return
}

func (s *service) GetUserSent(ctxSess *ctxSess.Context, user models.UserSession, req InboxSentReq) (resp SentEmailResp, err error) {
	go s.getUserSentBox(ctxSess, user)

	query := "*"
	var userQuery []string
	var subjectQuery []string
	if len(strings.TrimSpace(req.Search)) > 0 {
		userQuery = append(userQuery, strings.ReplaceAll(req.Search, "@", "_"))
		subjectQuery = append(subjectQuery, req.Search)
		query = "(@from:(" + strings.Join(userQuery, "|") + ")" + " | @subject|body_text:(" + strings.Join(subjectQuery, "|") + "))"
	}

	var offset int64
	if req.Page > 0 {
		req.Page = req.Page - 1
	}
	offset = req.Page * 10
	resp.Messages = s.redisRepo.SearchIndex(fmt.Sprintf("user:sent:%s",
		user.UserID.Hex()),
		query, "SORTBY", "sent", "DESC", "RETURN", 1,
		"message", "LIMIT", offset, 10)

	total := s.redisRepo.AggregateIndex(fmt.Sprintf("user:sent:%s",
		user.UserID.Hex()),
		query)

	totalEmails := total
	totalPage := int64(math.Ceil(float64(totalEmails) / 10))

	resp.Page = Page{
		CurrentPage: req.Page + 1,
		TotalPage:   totalPage,
		TotalEmails: totalEmails,
	}

	return resp, nil
}

func (s *service) GetMessage(ctxSess *ctxSess.Context, req *MessageQueue) (resp interface{}, err error) {
	cQuery := "user:inbox:" + req.UserId.Hex() + ":" + req.MessageID
	mess, err := s.googleWrapper.GetMessages(ctxSess, req.MessageID)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	received := time.Unix(0, mess.InternalDate*1000000).UTC()
	msg := &domainMsg.Message{
		MessageID: req.MessageID,
		ThreadID:  mess.ThreadId,
		UserID:    req.UserId,
		Received:  &received,
		Labels:    mess.LabelIds,
		Snippet:   mess.Snippet,
	}
	msg.IsRead = true
	for _, label := range mess.LabelIds {
		if label == constants.LABEL_UNREAD {
			msg.IsRead = false
			break
		}
	}

	thread, err := s.googleWrapper.GetThreadMessages(ctxSess, mess.ThreadId)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		return
	}
	if len(thread.Messages) > 1 {
		msg.IsReply = true
	}

	entity := &domainMsg.Message{
		MessageID: req.MessageID,
		ThreadID:  mess.ThreadId,
		UserID:    req.UserId,
		Received:  &received,
		Labels:    mess.LabelIds,
		Snippet:   mess.Snippet,
	}
	entity.Snippet, _ = utils.Encrypt(s.config.Salt, entity.Snippet)
	s.parseMessagePart(mess.Payload, msg, entity)
	err = s.messageRepo.InsertInbox(entity)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		ctxSess.Lv4()
	}
	encoded, _ := json.Marshal(msg)
	if msg.From == nil {
		s.redisRepo.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", "noreply_example.com", "received_at", mess.InternalDate)
		return msg, nil
	}
	email, _ := mail.ParseAddress(*msg.From)
	if email == nil {
		s.redisRepo.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", "noreply_example.com", "received_at", mess.InternalDate)
		return msg, nil
	}
	s.redisRepo.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", strings.ReplaceAll(email.Address, "@", "_"), "received_at", mess.InternalDate)
	return msg, nil
}

func (s *service) GetSentMessage(ctxSess *ctxSess.Context, req *MessageQueue) (resp interface{}, err error) {
	cQuery := "user:sent:" + req.UserId.Hex() + ":" + req.MessageID
	mess, err := s.googleWrapper.GetMessages(ctxSess, req.MessageID)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	sent := time.Unix(0, mess.InternalDate*1000000).UTC()
	msg := &domainMsg.Message{
		MessageID: req.MessageID,
		ThreadID:  mess.ThreadId,
		UserID:    req.UserId,
		Sent:      &sent,
		Labels:    mess.LabelIds,
		Snippet:   mess.Snippet,
	}
	entity := msg
	entity.Snippet, _ = utils.Encrypt(s.config.Salt, entity.Snippet)
	s.parseMessagePart(mess.Payload, msg, entity)
	_ = s.messageRepo.InsertSentBox(entity)
	encoded, _ := json.Marshal(msg)
	if msg.From == nil {
		s.redisRepo.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", "noreply_example.com", "sent", mess.InternalDate)
		return msg, nil
	}
	email, _ := mail.ParseAddress(*msg.From)
	if email == nil {
		s.redisRepo.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", "noreply_example.com", "sent", mess.InternalDate)
		return msg, nil
	}
	s.redisRepo.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", strings.ReplaceAll(email.Address, "@", "_"), "sent", mess.InternalDate)
	return msg, nil
}

func (s *service) GetUserLabel(ctxSess *ctxSess.Context, labelId string) (resp *gmail.Label, err error) {
	resp, err = s.googleWrapper.GetUserLabel(ctxSess, labelId)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	return
}

func (s *service) TagLabelUserMessage(ctxSess *ctxSess.Context, req TagLabelReq) (err error) {
	for k, v := range req.LabelIds {
		if strings.Contains(strings.ToUpper(v), "LABEL") {
			req.LabelIds[k] = strings.Title(strings.ToLower(v))
		} else {
			req.LabelIds[k] = strings.ToUpper(v)
		}
	}
	err = s.googleWrapper.ModifyMessage(ctxSess, &gmail.BatchModifyMessagesRequest{
		AddLabelIds: req.LabelIds,
		Ids:         req.MessageIds,
	})
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	return
}

func (s *service) TagLabelArchiveMessage(ctxSess *ctxSess.Context, req TagLabelArchiveReq) (err error) {
	archive, err := s.checkLabelArchive(ctxSess)
	if err != nil {
		return
	}

	err = s.googleWrapper.ModifyMessage(ctxSess, &gmail.BatchModifyMessagesRequest{
		AddLabelIds:    []string{archive.Id},
		RemoveLabelIds: []string{constants.LABEL_INBOX},
		Ids:            req.MessageIds,
	})
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	return
}

func (s *service) GetUserMessageDetail(ctxSess *ctxSess.Context, messageID, threadID string) (resp *UserMessage, err error) {
	resp, err = s.getUserMessage(ctxSess, messageID, threadID, "")
	return
}

func (s *service) GetUserMessage(ctxSess *ctxSess.Context, req UserMessageReq) (resp UserMessageResp, err error) {
	if strings.ToLower(req.Label) == strings.ToLower(constants.LABEL_ARCHIVE) {
		archive, errs := s.checkLabelArchive(ctxSess)
		if errs != nil {
			return
		}

		req.Label = archive.Id
	}

	resp.Messages = []UserMessage{}
	awayList, err := s.awayRepo.GetAwayList(ctxSess.UserSession.UserID, false)
	if err != nil && err != mgo.ErrNotFound {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	flagAwayMode, _ := s.validateAwayMode(awayList)
	var listBreakthrough string
	if flagAwayMode {
		listBreakthrough = s.redisRepo.GetKey(fmt.Sprintf(constants.ListBreakthrough, ctxSess.UserSession.Email))
		req.Search, err = s.createQuerySearch(ctxSess, awayList, req.Search, listBreakthrough)
		if err != nil {
			err = nil
			return
		}
	}

	if req.Label == constants.LABEL_SNOOZED {
		req.Label = ""
		req.Search = "in:snoozed " + req.Search
	}

	listMessages, err := s.getUserMessageList(ctxSess, req)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	var msgList []UserMessage
	eg := errgroup.Group{}

	label := &gmail.Label{}
	if !flagAwayMode && len(req.Search) == 0 {
		eg.Go(func() error {
			label, err = s.GetUserLabel(ctxSess, req.Label)
			return err
		})
	}
	if req.Label == "" {
		label.MessagesTotal = int64(len(listMessages.Messages))
	}

	for _, each := range listMessages.Messages {
		eg.Go(func() error {
			msg, errs := s.getUserMessage(ctxSess, each.Id, each.ThreadId, strings.ToUpper(req.Label))
			if errs == nil {
				if len(awayList) > 0 {
					now := time.Now().UTC()
					for _, away := range awayList {
						if away.IsEnabled {
							if away.AllDay {
								for _, valKey := range away.AllowedKeywords {
									if strings.ContainsAny(msg.BodyRawText, valKey) {
										msg.AwayContainsKeywords = append(msg.AwayContainsKeywords, valKey)
									}
								}
							} else {
								if now.After(away.ActivateAllow) && now.Before(away.DeactivateAllow) {
									for _, valKey := range away.AllowedKeywords {
										if strings.Contains(msg.BodyRawText, valKey) {
											msg.AwayContainsKeywords = append(msg.AwayContainsKeywords, valKey)
										}
									}
								} else if len(away.Repeat) > 0 {
									for _, valKey := range away.AllowedKeywords {
										if strings.ContainsAny(msg.BodyRawText, valKey) {
											msg.AwayContainsKeywords = append(msg.AwayContainsKeywords, valKey)
										}
									}
								}
							}
						}
					}
				}

				for _, v := range msg.Attachments {
					thumbnail, errThumbnail := s.GetAttachment(ctxSess, AttachmentRequest{
						MessageID:    msg.MessageID,
						AttachmentID: v.AttachmentID,
						FileName:     v.FileName,
						MimeType:     v.MimeType,
					})
					if errThumbnail == nil {
						msg.Thumbnail = thumbnail.FileByte
					}
				}
				msgList = append(msgList, *msg)
			}
			return errs
		})
		time.Sleep(5 * time.Millisecond)
	}
	if errs := eg.Wait(); errs != nil {
		ctxSess.ErrorMessage = errs.Error()
		errs = constants.ErrorGeneral
		return
	}

	ls := strings.Split(listBreakthrough, " ")
	for _, each := range listMessages.Messages {
		for _, eachMsg := range msgList {
			if each.Id == eachMsg.MessageID {
				for _, v := range ls {
					if eachMsg.OriginalMessageID == v {
						eachMsg.FlagBreakthrough = true
						break
					}
				}
				resp.Messages = append(resp.Messages, eachMsg)
				break
			}
		}
	}

	resp.NextPageToken = listMessages.NextPageToken
	resp.CurrentPageToken = req.NextPageToken
	resp.PrevPageToken = req.PrevPageToken
	if label != nil {
		resp.TotalEmails = label.MessagesTotal
	}

	return
}

func (s *service) GetThreadByID(ctxSess *ctxSess.Context, threadId string) (resp ThreadMessageResp, err error) {
	thread, err := s.googleWrapper.GetThreadMessages(ctxSess, threadId)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorGeneral
		return
	}

	for _, each := range thread.Messages {
		sent := time.Unix(0, each.InternalDate*1000000).UTC()
		msg := &UserMessage{
			MessageID: each.Id,
			ThreadID:  each.ThreadId,
			Sent:      &sent,
			Labels:    each.LabelIds,
			Snippet:   each.Snippet,
			HistoryId: each.HistoryId,
		}

		s.parseUserMessagePart(each.Payload, msg)

		resp.Messages = append(resp.Messages, *msg)
	}

	return
}

func (s *service) GetAttachment(ctxSess *ctxSess.Context, req AttachmentRequest) (resp AttachmentResponse, err error) {
	attDetail, err := s.googleWrapper.GetAttachmentDetail(ctxSess, req.MessageID, req.AttachmentID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorGeneral
		return
	}
	byt, err := base64.URLEncoding.DecodeString(attDetail.Data)
	resp = AttachmentResponse{
		FileByte: byt,
	}
	return
}

func (s *service) BreakthroughNotification(ctxSess *ctxSess.Context, req BreakthroughReq) (err error) {
	defer ctxSess.Lv4()

	key := fmt.Sprintf(constants.Breakthrough, req.Email, req.XRequestID)
	msgID := s.redisRepo.GetKey(key)
	if msgID == "" {
		err = constants.ErrorDataNotFound
		return
	}

	user, _ := s.userRepo.GetUser(req.Email)
	if user == nil {
		return
	}
	decRefreshToken, _ := utils.Decrypt(s.config.Salt, user.RefreshToken)
	ctxSess.UserSession.RefreshToken = decRefreshToken
	ctxSess.UserSession.UserID = user.ID

	mess, err := s.googleWrapper.GetMessages(ctxSess, msgID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		return
	}

	//if os.Getenv("GOLANG_ENV") != "prod" && os.Getenv("GOLANG_ENV") != "dev" {
	//	key := fmt.Sprintf(constants.UserTokenDev, user.ID.Hex())
	//	userTokenDev := s.redisRepo.GetKey(key)
	//	ctxSess.UserSession.DevTokenKey, _ = strconv.ParseBool(userTokenDev)
	//}
	badge, errs := s.googleWrapper.GetUserLabel(ctxSess, constants.LABEL_UNREAD)
	if errs != nil {
		ctxSess.ErrorMessage = errs.Error()
		errs = constants.ErrorGeneral
		return
	}

	t2 := ctxSess.Lv2("Push Notification msgID: ", msgID, "start")
	msg, _ := s.GetUserMessageDetail(ctxSess, msgID, mess.ThreadId)
	if msg == nil || msg.To == nil {
		err = constants.ErrorDataNotFound
		return
	}

	list := s.redisRepo.GetKey(fmt.Sprintf(constants.ListBreakthrough, req.Email))
	if strings.TrimSpace(list) == "" {
		list = msg.OriginalMessageID
	} else {
		list = fmt.Sprintf("%s %s", list, msg.OriginalMessageID)
	}
	s.redisRepo.SetKey(fmt.Sprintf(constants.ListBreakthrough, req.Email), list)

	var flagInbox bool
	for _, eachLabel := range msg.Labels {
		if eachLabel == "INBOX" {
			flagInbox = true
			break
		}
	}
	if !flagInbox {
		return
	}

	var from string
	if msg.From != nil {
		from = *msg.From
	}

	notif := gaurun.NotifPayload{
		Token:    user.SwiftToken,
		Platform: 1,
		Message:  fmt.Sprintf("%s\n%s\n%s", from, msg.Subject, msg.Snippet),
		Badge:    int(badge.ThreadsUnread),
		Expiry:   10,
		Sound:    "NEW CASBU.wav",
		Extend: []gaurun.NotifExtend{
			{
				Key: "message_id",
				Val: msg.MessageID,
			},
			{
				Key: "email",
				Val: req.Email,
			},
		},
		Category:       "new",
		MutableContent: true,
	}

	notification := gaurun.Notification{}
	notification.Notifications = append(notification.Notifications, notif)
	if len(notification.Notifications) > 0 {
		if err = s.gaurunWrapper.Send(ctxSess, notification); err != nil {
			ctxSess.ErrorMessage = err.Error()
			errs = constants.ErrorGeneral
			ctxSess.Lv3(t2, "Push Notification msgID: ", msgID, "end with error", "notification", notification)
			return
		}
	}
	ctxSess.Lv3(t2, "Push Notification msgID: ", msgID, "end", "notification", notification)

	s.redisRepo.DelKey(key)

	return
}

func (s *service) PushNotification(ctxSess *ctxSess.Context, in *PushNotification) (err error) {
	defer ctxSess.Lv4()
	user, _ := s.userRepo.GetUser(in.EmailAddress)
	if user == nil {
		return
	}
	decRefreshToken, _ := utils.Decrypt(s.config.Salt, user.RefreshToken)
	ctxSess.UserSession.RefreshToken = decRefreshToken
	ctxSess.UserSession.UserID = user.ID

	//if os.Getenv("GOLANG_ENV") != "prod" && os.Getenv("GOLANG_ENV") != "dev" {
	//	key := fmt.Sprintf(constants.UserTokenDev, user.ID.Hex())
	//	userTokenDev := s.redisRepo.GetKey(key)
	//	ctxSess.UserSession.DevTokenKey, _ = strconv.ParseBool(userTokenDev)
	//}

	awayList, err := s.awayRepo.GetAwayList(ctxSess.UserSession.UserID, false)
	if err != nil && err != mgo.ErrNotFound {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	historyKey := fmt.Sprintf(constants.UserHistory, ctxSess.UserSession.Email)
	hs := s.redisRepo.GetKey(historyKey)
	hsID, err := strconv.ParseUint(hs, 0, 64)
	if err != nil {
		hsID = 0
	}
	if hsID == in.HistoryId {
		err = constants.ErrorDoublePushNotification
		ctxSess.ErrorMessage = err.Error()
		return
	}
	s.redisRepo.SetKey(fmt.Sprintf(constants.UserHistory, in.EmailAddress), in.HistoryId)

	watchKey := fmt.Sprintf(constants.UserWatch, in.EmailAddress)
	watchSess := s.redisRepo.GetKey(watchKey)
	if watchSess == "" {
		err = s.googleWrapper.StopWatchPushNotification(ctxSess, ctxSess.UserSession.RefreshToken)
		if err != nil {
			ctxSess.ErrorMessage = err.Error()
			ctxSess.Lv4()
		}
		watchRes, _ := s.googleWrapper.WatchPushNotification(ctxSess, ctxSess.UserSession.RefreshToken, &gmail.WatchRequest{
			LabelIds:  []string{"INBOX"},
			TopicName: fmt.Sprintf("projects/%s/topics/%s", s.config.Gpubsub.ProjectName, s.config.Gpubsub.Topic),
		})
		if watchRes == nil {
			return
		}
		s.redisRepo.SetKey(watchKey, watchRes.HistoryId, "EX", 86400)
		s.redisRepo.SetKey(fmt.Sprintf(constants.UserHistory, in.EmailAddress), watchRes.HistoryId)
		hsID = watchRes.HistoryId
	}

	lsHistory, err := s.googleWrapper.GetHistory(ctxSess, hsID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		return
	}

	if lsHistory.History != nil {
		badge, errs := s.googleWrapper.GetUserLabel(ctxSess, constants.LABEL_UNREAD)
		if errs != nil {
			ctxSess.ErrorMessage = errs.Error()
			errs = constants.ErrorGeneral
			return
		}
		for _, each := range lsHistory.History {
			notification := gaurun.Notification{}
			t2 := ctxSess.Lv2("Push Notification HistoryID: ", each.Id, "len msgID: ", len(each.MessagesAdded), "start")
			for _, eachMsg := range each.MessagesAdded {
				if eachMsg.Message != nil {
					msgID := eachMsg.Message.Id
					msg, _ := s.GetUserMessageDetail(ctxSess, msgID, eachMsg.Message.ThreadId)
					if msg == nil || msg.To == nil {
						continue
					}
					var flagInbox bool
					for _, eachLabel := range msg.Labels {
						if eachLabel == "INBOX" {
							flagInbox = true
							break
						}
					}
					if !flagInbox {
						continue
					}

					emailFrom := strings.Split(*msg.From, " ")
					eFrom := emailFrom[len(emailFrom)-1]
					if string(eFrom[0]) == "<" && string(eFrom[len(eFrom)-1]) == ">" {
						eFrom = eFrom[1:]
						eFrom = eFrom[:len(eFrom)-1]
						msg.From = &eFrom
					}
					flagAwayMode, endDate := s.validateAwayMode(awayList)
					if flagAwayMode {
						if !s.awayModeNotificationV2(ctxSess, *msg.From, awayList, msg.BodyText) {
							userProfile, _ := s.googleWrapper.GetProfile(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
							m := new(mailWrapper.Email)
							ampBreakthroughURL := fmt.Sprintf("%s/v2/breakthrough/amp/notification?email=%s&id=%s", s.config.Application.URL, in.EmailAddress, ctxSess.XRequestID)
							breakthroughURL := fmt.Sprintf("%s/v2/breakthrough/notification?email=%s&id=%s", s.config.Application.URL, in.EmailAddress, ctxSess.XRequestID)

							s.redisRepo.SetKey(fmt.Sprintf(constants.Breakthrough, in.EmailAddress, ctxSess.XRequestID), msgID)
							resCode, body, er := m.SendEmail(*msg.From, endDate, ampBreakthroughURL, breakthroughURL, *msg.To, userProfile.Name)
							if er != nil {
								ctxSess.ErrorMessage = er.Error()
								ctxSess.Lv4()
								continue
							}
							if resCode != http.StatusAccepted {
								ctxSess.ErrorMessage = body
								ctxSess.Lv4()
								continue
							}

							continue
						}
					}

					var from string
					if msg.From != nil {
						from = *msg.From
					}

					notif := gaurun.NotifPayload{
						Token:    user.SwiftToken,
						Platform: 1,
						Message:  fmt.Sprintf("%s\n%s\n%s", from, msg.Subject, msg.Snippet),
						Badge:    int(badge.ThreadsUnread),
						Expiry:   10,
						Sound:    "NEW CASBU.wav",
						Extend: []gaurun.NotifExtend{
							{
								Key: "message_id",
								Val: msg.MessageID,
							},
							{
								Key: "email",
								Val: in.EmailAddress,
							},
						},
						Category:       "new",
						MutableContent: true,
					}
					notification.Notifications = append(notification.Notifications, notif)
				}
			}
			if len(notification.Notifications) > 0 {
				if err = s.gaurunWrapper.Send(ctxSess, notification); err != nil {
					ctxSess.ErrorMessage = err.Error()
					errs = constants.ErrorGeneral
					ctxSess.Lv3(t2, "Push Notification HistoryID: ", each.Id, "end with error", "notification", notification)
					return
				}
			}
			ctxSess.Lv3(t2, "Push Notification HistoryID: ", each.Id, "end", "notification", notification)
		}
	}

	return
}

func (s *service) GetUserMessageByID(ctxSess *ctxSess.Context, messageID string) (resp *UserMessage, err error) {
	mess, err := s.googleWrapper.GetMessages(ctxSess, messageID)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}

	received := time.Unix(0, mess.InternalDate*1000000).UTC()
	msg := &UserMessage{
		MessageID: messageID,
		ThreadID:  mess.ThreadId,
		Received:  &received,
		Labels:    mess.LabelIds,
		Snippet:   mess.Snippet,
	}

	s.parseUserMessagePart(mess.Payload, msg)
	return msg, nil
}
