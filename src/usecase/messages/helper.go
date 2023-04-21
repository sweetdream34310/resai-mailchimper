package messages

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/gmail/v1"
	"gopkg.in/mgo.v2"

	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/aways"
	domainMsg "github.com/cloudsrc/api.awaymail.v1.go/src/domain/messages"
	googleWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
)

func (s *service) getUserInbox(ctxSess *ctxSess.Context, user models.UserSession) (err error) {
	_ = s.redisRepo.CreateInboxIndex(user)

	var messageList []string
	skip := "initial"

	key := fmt.Sprintf(constants.UserArchiveSession, ctxSess.UserSession.UserID.Hex())
	s.redisRepo.SetKey(key, ctxSess.UserSession.UserID.Hex())

	for skip != "" {
		var mList *gmail.ListMessagesResponse
		if s.userLimitCheck(user.UserID.Hex()) {
			if skip == "initial" {
				skip = ""
			}
			mList, err = s.googleWrapper.GetMessagesList(ctxSess, googleWrapper.GetMessageRequest{skip, constants.LABEL_INBOX, "", constants.DefaultLimit})
			if err != nil {
				ctxSess.Lv4()
			}
			for _, ii := range mList.Messages {
				cQuery := "user:inbox:" + user.UserID.Hex() + ":" + ii.Id
				messageList = append(messageList, cQuery)
				if res := s.redisRepo.GetCache(cQuery, "message"); res == nil {
					if err = s.rabbit.Publish("inbox_processor_"+user.UserID.Hex(), map[string]interface{}{
						"message_id":    ii.Id,
						"user_id_hex":   user.UserID.Hex(),
						"user_id":       user.UserID,
						"email":         user.Email,
						"auth_token":    ctxSess.UserSession.AuthToken,
						"refresh_token": ctxSess.UserSession.RefreshToken,
					}); err != nil {
						ctxSess.ErrorMessage = err.Error()
						ctxSess.Lv4()
						return
					}
				}
			}
		}
		skip = mList.NextPageToken
	}
	s.removeDiffMessages(difference(s.redisRepo.GetKeysPrefix("user:inbox:"+user.UserID.Hex()+":*"), messageList))
	s.redisRepo.DelKey(key)
	return
}

func (s *service) getUserSentBox(ctxSess *ctxSess.Context, user models.UserSession) (err error) {
	_ = s.redisRepo.CreateSentIndex(user)

	var messageList []string
	skip := "initial"

	for skip != "" {
		var mList *gmail.ListMessagesResponse
		if s.userLimitCheck(user.UserID.Hex()) {
			if skip == "initial" {
				skip = ""
			}
			mList, err = s.googleWrapper.GetMessagesList(ctxSess, googleWrapper.GetMessageRequest{skip, constants.LABEL_SENT, "", constants.DefaultLimit})
			if err != nil {
				ctxSess.Lv4()
			}
			for _, ii := range mList.Messages {
				cQuery := "user:sent:" + user.UserID.Hex() + ":" + ii.Id
				messageList = append(messageList, cQuery)
				if res := s.redisRepo.GetCache(cQuery, "message"); res == nil {
					if err = s.rabbit.Publish("sent_processor_"+user.UserID.Hex(), map[string]interface{}{
						"message_id":    ii.Id,
						"user_id_hex":   user.UserID.Hex(),
						"user_id":       user.UserID,
						"email":         user.Email,
						"auth_token":    ctxSess.UserSession.AuthToken,
						"refresh_token": ctxSess.UserSession.RefreshToken,
					}); err != nil {
						ctxSess.ErrorMessage = err.Error()
						ctxSess.Lv4()
						return
					}
				}
			}
		}
		skip = mList.NextPageToken
	}
	s.removeDiffMessages(difference(s.redisRepo.GetKeysPrefix("user:sent:"+user.UserID.Hex()+":*"), messageList))
	return
}

func (s *service) removeDiffMessages(diff []string) {
	if len(diff) != 0 {
		for _, d := range diff {
			s.redisRepo.DelCache(d, "message", "body_text", "from", "subject", "received_at")
		}
	}
}

func (s *service) awayModeNotification(ctxSess *ctxSess.Context, emailFrom string) bool {
	notifyList, err := s.contactsRepo.GetList(ctxSess.UserSession)
	if err != nil && err != mgo.ErrNotFound {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return false
	}

	if notifyList != nil {
		for _, v := range notifyList {
			email, _ := utils.Decrypt(s.config.Salt, v.Email)
			if strings.Contains(strings.ToLower(emailFrom), strings.ToLower(email)) {
				return true
			}
		}
	}

	return false
}

func (s *service) awayModeNotificationV2(ctxSess *ctxSess.Context, emailFrom string, awayList []*domain.Away, bodyText string) bool {
	for _, v := range awayList {
		if v.IsEnabled {
			for _, allowContact := range v.AllowedContacts {
				if strings.TrimSpace(allowContact) == "" {
					continue
				}
				if allowContact == emailFrom {
					return true
				}
			}
			for _, allowKeyword := range v.AllowedKeywords {
				if strings.TrimSpace(allowKeyword) == "" {
					continue
				}
				if strings.Contains(bodyText, allowKeyword) {
					return true
				}
			}
		}
	}
	return false
}

func (s *service) createQuerySearch(ctxSess *ctxSess.Context, awayList []*domain.Away, search, listBreakthrough string) (query string, err error) {
	notifyList, errs := s.contactsRepo.GetList(ctxSess.UserSession)
	if errs != nil && errs != mgo.ErrNotFound {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}

	var q strings.Builder
	q.WriteString("{")

	if notifyList != nil {
		for k, v := range notifyList {
			notifyList[k].Email, _ = utils.Decrypt(s.config.Salt, v.Email)
			notifyList[k].Name, _ = utils.Decrypt(s.config.Salt, v.Name)
		}
		for _, each := range notifyList {
			q.WriteString("from:")
			q.WriteString(each.Email)
			q.WriteString(" ")
		}
	}

	for _, v := range awayList {
		if !v.IsEnabled {
			continue
		}
		if len(v.AllowedContacts) > 0 {
			for _, val := range v.AllowedContacts {
				q.WriteString("from:")
				q.WriteString(val)
				q.WriteString(" ")
			}
		}
		if len(v.AllowedKeywords) > 0 {
			for _, val := range v.AllowedKeywords {
				q.WriteString(val)
				q.WriteString(" ")
			}
		}
	}

	if len(listBreakthrough) > 0 {
		ls := strings.Split(listBreakthrough, " ")
		for _, v := range ls {
			if len(v) == 0 {
				continue
			}
			q.WriteString("rfc822msgid:")
			q.WriteString(v)
			q.WriteString(" ")
		}
	}
	q.WriteString("}")
	if q.String() == "{}" {
		err = constants.ErrorDataNotFound
		return
	} else {
		startTime := s.redisRepo.GetKey(fmt.Sprintf(constants.AwayModeStart, ctxSess.UserSession.Email))
		_, er := strconv.ParseInt(startTime, 10, 64)
		if er == nil {
			q.WriteString("after:")
			q.WriteString(startTime)
			q.WriteString(" ")
		}
	}

	if len(search) > 0 {
		q.WriteString(" ")
		q.WriteString(search)
	}

	query = q.String()
	return
}

func (s *service) createQueryInbox(ctxSess *ctxSess.Context, user models.UserSession, search string) (query string, err error) {
	notifyList, errs := s.contactsRepo.GetList(user)
	if errs != nil && errs != mgo.ErrNotFound {
		ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorDatabase
		return
	}
	for k, v := range notifyList {
		notifyList[k].Email, _ = utils.Decrypt(s.config.Salt, v.Email)
		notifyList[k].Name, _ = utils.Decrypt(s.config.Salt, v.Name)
	}
	var userQuery []string
	var subjectQuery []string
	for _, each := range notifyList {
		if len(strings.TrimSpace(search)) > 0 {
			if strings.Contains(each.Email, strings.ToLower(search)) {
				userQuery = append(userQuery, strings.ReplaceAll(each.Email, "@", "_"))
			}
		} else {
			userQuery = append(userQuery, strings.ReplaceAll(each.Email, "@", "_"))
		}
	}
	if len(userQuery) == 0 {
		for _, each := range notifyList {
			userQuery = append(userQuery, strings.ReplaceAll(each.Email, "@", "_"))
		}
	}
	query = "@from:(" + strings.Join(userQuery, "|") + ")"
	if len(strings.TrimSpace(search)) > 0 {
		query = query + " " + "@subject|body_text:(" + strings.Join(append(subjectQuery, search), "|") + ")"
	}
	return
}

func (s *service) validateAwayMode(awayList []*domain.Away) (flagEnableAway bool, endDate string) {
	if len(awayList) > 0 {
		now := time.Now().UTC()
		for _, away := range awayList {
			if away.IsEnabled {
				endDate = fmt.Sprintf("%s 23:59:59", away.DeactivateAllow.Format("2006-01-02"))
				if away.AllDay {
					start, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s 00:00:00", away.ActivateAllow.Format("2006-01-02")))
					end, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s 23:59:59", away.DeactivateAllow.Format("2006-01-02")))
					if now.After(start) && now.Before(end) {
						flagEnableAway = true
						break
					} else if len(away.Repeat) > 0 {
						for _, eachDay := range away.Repeat {
							if eachDay == now.Weekday().String() {
								flagEnableAway = true
								break
							}
						}
					}
				} else {
					if now.After(away.ActivateAllow) && now.Before(away.DeactivateAllow) {
						flagEnableAway = true
						break
					} else if len(away.Repeat) > 0 {
						start, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", now.Format("2006-01-02"), away.ActivateAllow.Format("15:04:05")))
						end, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", now.Format("2006-01-02"), away.DeactivateAllow.Format("15:04:05")))
						if now.After(start) && now.Before(end) {
							flagEnableAway = true
							break
						}
					}
				}
			}
		}
	}
	return
}

func difference(slice1 []string, slice2 []string) []string {
	var diff []string
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}
	return diff
}

func (s *service) userLimitCheck(key string) bool {
	val := s.redisRepo.GetKey(key)
	if val == "" {
		s.redisRepo.SetKey(key, 1, "EX", 86400)
		return true
	}
	totalUserAPILimit, _ := strconv.Atoi(val)
	if totalUserAPILimit < constants.USER_PER_DAY_API_LIMIT {
		s.redisRepo.IncrKey(key)
		return true
	}
	return false
}

func (s *service) parseMessagePart(origmsg *gmail.MessagePart, mess *domainMsg.Message, entity *domainMsg.Message) {
	mess.MimeFlow = append(mess.MimeFlow, origmsg.MimeType)
	entity.MimeFlow = append(entity.MimeFlow, origmsg.MimeType)
	for _, ii := range origmsg.Headers {
		if ii.Name == "Subject" {
			encSubject, _ := utils.Encrypt(s.config.Salt, ii.Value)
			entity.Subject = encSubject
			mess.Subject = ii.Value
		}
		if ii.Name == "From" {
			encFrom, _ := utils.Encrypt(s.config.Salt, ii.Value)
			entity.From = &encFrom
			mess.From = &ii.Value
		}
		if ii.Name == "To" {
			email, _ := mail.ParseAddress(ii.Value)
			entity.To = &email.Address
			mess.To = &ii.Value
		}
		if strings.ToLower(ii.Name) == "message-id" {
			entity.OriginalMessageID = ii.Value
			mess.OriginalMessageID = ii.Value
		}
	}
	if len(origmsg.Filename) != 0 {
		encFilename, _ := utils.Encrypt(s.config.Salt, origmsg.Filename)
		encMimeType, _ := utils.Encrypt(s.config.Salt, origmsg.MimeType)
		entity.Attachments = append(entity.Attachments, domainMsg.Attachment{
			FileName:       encFilename,
			MimeType:       encMimeType,
			AttachmentID:   origmsg.Body.AttachmentId,
			AttachmentSize: origmsg.Body.Size,
		})
		mess.Attachments = append(mess.Attachments, domainMsg.Attachment{
			FileName:       origmsg.Filename,
			MimeType:       origmsg.MimeType,
			AttachmentID:   origmsg.Body.AttachmentId,
			AttachmentSize: origmsg.Body.Size,
		})
	}

	if strings.HasPrefix(origmsg.MimeType, "multipart") {
		for _, ii := range origmsg.Parts {
			s.parseMessagePart(ii, mess, entity)
		}
	}
	if origmsg.MimeType == "text/plain" {
		sBodyText, _ := base64.URLEncoding.DecodeString(origmsg.Body.Data)
		encBodyText, _ := utils.Encrypt(s.config.Salt, string(sBodyText))
		entity.BodyText = encBodyText
		mess.BodyText = string(sBodyText)
	} else if origmsg.MimeType == "text/html" {
		sBodyHTML, _ := base64.URLEncoding.DecodeString(origmsg.Body.Data)
		encBodyHTML, _ := utils.Encrypt(s.config.Salt, string(sBodyHTML))
		entity.BodyHTML = encBodyHTML
		mess.BodyHTML = string(sBodyHTML)
	}
}

func (s *service) getUserMessageList(ctxSess *ctxSess.Context, req UserMessageReq) (mList *gmail.ListMessagesResponse, err error) {
	mList, err = s.googleWrapper.GetMessagesList(ctxSess, googleWrapper.GetMessageRequest{req.NextPageToken, req.Label, req.Search, constants.MinimumLimit})
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return
}

func (s *service) getUserMessage(ctxSess *ctxSess.Context, messageID string, threadID string, label string) (msg *UserMessage, err error) {
	mess := &gmail.Message{}
	thread, err := s.googleWrapper.GetThreadMessages(ctxSess, threadID)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		return
	}
	// check if threadID is same with messageID
	for _, v := range thread.Messages {
		if v.Id == messageID {
			mess = v
			break
		}
	}
	if mess.Id == "" {
		err = constants.ErrorDataNotFound
		return
	}

	sent := time.Unix(0, mess.InternalDate*1000000).UTC()
	msg = &UserMessage{
		MessageID: mess.Id,
		ThreadID:  mess.ThreadId,
		Sent:      &sent,
		Labels:    mess.LabelIds,
		Snippet:   mess.Snippet,
		HistoryId: mess.HistoryId,
	}
	if label != "" {
		for _, v := range thread.Messages {
			for _, val := range v.LabelIds {
				if val == label {
					msg.CountReply++
				}
			}
		}
	}

	msg.IsRead = true
	for _, label := range mess.LabelIds {
		if label == constants.LABEL_UNREAD {
			msg.IsRead = false
			break
		}
	}

	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		return
	}
	if len(thread.Messages) > 1 {
		msg.IsReply = true
	}

	s.parseUserMessagePart(mess.Payload, msg)

	return
}

func (s *service) parseUserMessagePart(origmsg *gmail.MessagePart, mess *UserMessage) {
	mess.MimeFlow = append(mess.MimeFlow, origmsg.MimeType)
	for _, ii := range origmsg.Headers {
		if ii.Name == "Subject" {
			mess.Subject = ii.Value
		}
		if ii.Name == "From" {
			mess.From = &ii.Value
		}
		if ii.Name == "To" {
			mess.To = &ii.Value
		}
		if strings.ToLower(ii.Name) == "message-id" {
			mess.OriginalMessageID = ii.Value
		}
	}
	if len(origmsg.Filename) != 0 {
		mess.Attachments = append(mess.Attachments, Attachment{
			FileName:       origmsg.Filename,
			MimeType:       origmsg.MimeType,
			AttachmentID:   origmsg.Body.AttachmentId,
			AttachmentSize: origmsg.Body.Size,
		})
	}

	if strings.HasPrefix(origmsg.MimeType, "multipart") {
		for _, ii := range origmsg.Parts {
			s.parseUserMessagePart(ii, mess)
		}
	}
	if origmsg.MimeType == "text/plain" {
		sBodyText, _ := base64.URLEncoding.DecodeString(origmsg.Body.Data)
		mess.BodyText = string(sBodyText)
		mess.BodyRawText = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(strings.TrimSuffix(mess.BodyText, "\r\n"), "")
	} else if origmsg.MimeType == "text/html" {
		sBodyHTML, _ := base64.URLEncoding.DecodeString(origmsg.Body.Data)
		mess.BodyHTML = string(sBodyHTML)
		re := regexp.MustCompile(`<[^>]*>`)
		mess.BodyText = re.ReplaceAllString(strings.TrimSuffix(mess.BodyHTML, "\r\n"), "")
		mess.BodyRawText = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(mess.BodyText, "")
	}
	//if strings.TrimSpace(mess.BodyRawText) != "" {
	//	mess.Summary = s.callChatGpt(mess.BodyRawText, 0)
	//}
}

func (s *service) callChatGpt(context string, retry int) (summary string) {
	chatgptResp, err := s.chatgptWrapper.CreateSummary(context)
	if err == nil {
		if len(chatgptResp.Choices) > 0 {
			summary = chatgptResp.Choices[0].Message.Content
			if strings.Contains(summary, "As an AI language model") && retry < 5 {
				return s.callChatGpt(context, retry+1)
			}
		}
	}
	return
}

func (s *service) processArchive(ctxSess *ctxSess.Context) {
	key := fmt.Sprintf(constants.UserArchiveSession, ctxSess.UserSession.UserID.Hex())
	flagArchive := s.redisRepo.GetKey(key)
	if flagArchive == "" {
		go s.getUserInbox(ctxSess, ctxSess.UserSession)
	}
}

func (s *service) checkLabelArchive(ctxSess *ctxSess.Context) (archive *gmail.Label, err error) {
	ls, err := s.googleWrapper.GetLabelList(ctxSess)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	var flagArchive bool
	for _, eachLabel := range ls {
		if eachLabel.Name == constants.LABEL_AwayARCHIVE {
			flagArchive = true
			archive = eachLabel
			break
		}
	}

	if !flagArchive {
		archive, err = s.googleWrapper.CreateLabel(ctxSess, &gmail.Label{Name: constants.LABEL_AwayARCHIVE})
		if err != nil {
			ctxSess.ErrorMessage = err.Error()
			err = constants.ErrorGeneral
			return
		}
	}

	return
}
