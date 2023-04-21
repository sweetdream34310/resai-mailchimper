package google

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/people/v1"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
)

type provider struct {
	config      config.Config
	oauthConfig *oauth2.Config
}

func New(config config.Config) Wrapper {
	defaultConfig := &oauth2.Config{
		ClientID:     config.Google.Ios.ClientID,
		ClientSecret: config.Google.Ios.ClientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  constants.GoogleAuthURL,
			TokenURL: constants.GoogleTokenURL,
		},
	}
	return &provider{
		config:      config,
		oauthConfig: defaultConfig,
	}
}

var scopes = []string{
	people.UserinfoProfileScope,
	people.ContactsReadonlyScope,
	gmail.MailGoogleComScope,
	calendar.CalendarReadonlyScope,
}

const (
	meSelf       = "me"
	peopleMeSelf = "people/me"
)

func (p *provider) ValidateToken(ctxSess *ctxSess.Context, agent, refreshToken string) (token *oauth2.Token, err error) {
	p.oauthConfig.ClientID, p.oauthConfig.ClientSecret, err = p.getClientIDSecret(agent)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		return
	}
	tok := &oauth2.Token{RefreshToken: refreshToken}
	if token, err = p.oauthConfig.TokenSource(context.TODO(), tok).Token(); err != nil {
		//ctxSess.ErrorMessage = err.Error()
		err = constants.ErrorInvalidRequest
		return
	}
	return
}

func (p *provider) GetProfile(ctxSess *ctxSess.Context, agent, refreshToken string) (resp UserProfile, err error) {
	token, err := p.ValidateToken(ctxSess, agent, refreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, refreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return
		}
	}
	p.oauthConfig.ClientID, p.oauthConfig.ClientSecret, err = p.getClientIDSecret(agent)
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		return
	}

	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		return
	}
	profile, err := gmailSvc.Users.GetProfile(meSelf).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		return
	}
	peopleSvc, err := people.New(p.oauthConfig.Client(context.Background(), token))

	resp = UserProfile{
		Email:        profile.EmailAddress,
		AuthToken:    token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	me, err := peopleSvc.People.Get(peopleMeSelf).PersonFields("names,photos").Do()
	if err == nil {
		resp.Name = me.Names[0].DisplayName
		resp.Photo = me.Photos[0].Url
	}

	return
}

func (p *provider) GetMessagesList(ctxSess *ctxSess.Context, req GetMessageRequest) (*gmail.ListMessagesResponse, error) {
	var (
		response *gmail.ListMessagesResponse
		err      error
	)
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return nil, err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))

	list := gmailSvc.Users.Messages.List(meSelf).
		MaxResults(req.Limit).
		Q(req.Q)
	if len(req.Label) != 0 {
		list.LabelIds(req.Label)
	}
	if req.NextPageToken != "" {
		list.PageToken(req.NextPageToken)
	}
	response, err = list.Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		return nil, err
	}

	return response, nil
}

func (p *provider) GetMessages(ctxSess *ctxSess.Context, messageId string) (*gmail.Message, error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return nil, err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))

	msg, err := gmailSvc.Users.Messages.Get(meSelf, messageId).Fields().Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return msg, err
}

func (p *provider) GetThreadMessages(ctxSess *ctxSess.Context, threadId string) (*gmail.Thread, error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return nil, err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))

	msg, err := gmailSvc.Users.Threads.Get(meSelf, threadId).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return msg, err
}

func (p *provider) UpdateMessage(ctxSess *ctxSess.Context, id string, isRead *bool) (message *gmail.Message, err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return nil, err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))

	var req gmail.ModifyMessageRequest
	if !*isRead {
		req.AddLabelIds = append(req.AddLabelIds, constants.LABEL_UNREAD)
	} else {
		req.RemoveLabelIds = append(req.AddLabelIds, constants.LABEL_UNREAD)
	}
	res, err := gmailSvc.Users.Messages.Modify(meSelf, id, &req).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
		return nil, err
	}

	return p.GetMessages(ctxSess, res.Id)
}

func (p *provider) DeleteMessage(ctxSess *ctxSess.Context, messageId string) (err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))

	err = gmailSvc.Users.Messages.Delete(meSelf, messageId).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return
}

func (p *provider) ModifyMessage(ctxSess *ctxSess.Context, req *gmail.BatchModifyMessagesRequest) (err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))

	err = gmailSvc.Users.Messages.BatchModify(meSelf, req).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return
}

func (p *provider) GetUserLabel(ctxSess *ctxSess.Context, messageId string) (res *gmail.Label, err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return nil, err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))

	res, err = gmailSvc.Users.Labels.Get(meSelf, messageId).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return
}

func (p *provider) GetLabelList(ctxSess *ctxSess.Context) (res []*gmail.Label, err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return nil, err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))
	out, err := gmailSvc.Users.Labels.List(meSelf).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	res = out.Labels
	return
}

func (p *provider) CreateLabel(ctxSess *ctxSess.Context, req *gmail.Label) (res *gmail.Label, err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return nil, err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))
	res, err = gmailSvc.Users.Labels.Create(meSelf, req).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return
}

func (p *provider) DeleteLabel(ctxSess *ctxSess.Context, labelID string) (err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))
	err = gmailSvc.Users.Labels.Delete(meSelf, labelID).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return
}

func (p *provider) PatchLabel(ctxSess *ctxSess.Context, labelID string, req *gmail.Label) (res *gmail.Label, err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return nil, err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))
	res, err = gmailSvc.Users.Labels.Patch(meSelf, labelID, req).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return
}

func (p *provider) GetAttachmentDetail(ctxSess *ctxSess.Context, messageID, attachmentID string) (res *gmail.MessagePartBody, err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return nil, err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))
	res, err = gmailSvc.Users.Messages.Attachments.Get(meSelf, messageID, attachmentID).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return
}

func (p *provider) SendMessage(ctxSess *ctxSess.Context, req Messsage) (err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))

	var message gmail.Message
	msg := fmt.Sprintf("From: %s\r\n"+
		"reply-to: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n%s", meSelf, req.Email, req.To, req.Subject, req.Message)

	if len(req.ThreadID) > 0 {
		message.ThreadId = req.ThreadID
		msg = fmt.Sprintf("From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"In-Reply-To: %s\r\n"+
			"References: %s\r\n"+
			"\r\n%s", meSelf, req.To, req.Subject, req.MessageID, req.MessageID, req.Message)
	}

	message.Raw = base64.StdEncoding.EncodeToString([]byte(msg))
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	if _, err = gmailSvc.Users.Messages.Send(meSelf, &message).Do(); err != nil {
		ctxSess.ErrorMessage = err.Error()
		return err
	}
	return nil
}

func (p *provider) SendMessageWithAttachment(ctxSess *ctxSess.Context, req Messsage) (err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))

	boundary := randStr(32, "alphanum")
	var message gmail.Message
	var attachmentStr string
	for count, attachmentURL := range req.AttachmentsURL {
		url := strings.Split(attachmentURL, "/")
		filename := url[len(url)-1]
		filePath := "./public/" + filename
		fileBytes, errs := ioutil.ReadFile(filePath)
		if errs != nil {
			if er := downloadFile(attachmentURL, filePath); er != nil {
				ctxSess.ErrorMessage = er.Error()
				return er
			}
			fileBytes, _ = ioutil.ReadFile(filePath)
		}
		fileMIMEType := http.DetectContentType(fileBytes)
		fileData := base64.StdEncoding.EncodeToString(fileBytes)
		var astr string
		if count == 0 {
			astr = "Content-Type: " + fileMIMEType + "; name=" + string('"') + filename + string('"') + " \n" +
				"MIME-Version: 1.0\n" +
				"Content-Transfer-Encoding: base64\n" +
				"Content-Disposition: attachment; filename=" + string('"') + filename + string('"') + " \n\n" +
				chunkSplit(fileData, 76, "\n") +
				"--" + boundary + "\n"
		} else {
			astr = "Content-Type: " + fileMIMEType + "; name=" + string('"') + filename + string('"') + " \n" +
				"MIME-Version: 1.0\n" +
				"Content-Transfer-Encoding: base64\n" +
				"Content-Disposition: attachment; filename=" + string('"') + filename + string('"') + " \n\n" +
				chunkSplit(fileData, 76, "\n")
		}
		attachmentStr = attachmentStr + astr
	}

	messageBody := []byte("Content-Type: multipart/mixed; boundary=" + boundary + " \n" +
		"MIME-Version: 1.0\n" +
		"to: " + req.To + "\n" +
		"subject: " + req.Subject + "\n\n" +
		"--" + boundary + "\n" +
		"Content-Type: text/plain; charset=" + string('"') + "UTF-8" + string('"') + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: 7bit\n\n" +
		req.Message + "\n\n" +
		"--" + boundary + "\n" +
		attachmentStr +
		"--" + boundary + "--")
	if len(req.ThreadID) > 0 {
		message.ThreadId = req.ThreadID
		messageBody = []byte("Content-Type: multipart/mixed; boundary=" + boundary + " \n" +
			"MIME-Version: 1.0\n" +
			"to: " + req.To + "\n" +
			"subject: " + req.Subject + "\n\n" +
			"In-Reply-To: " + req.MessageID + "\n\n" +
			"References: " + req.MessageID + "\n\n" +
			"--" + boundary + "\n" +
			"Content-Type: text/plain; charset=" + string('"') + "UTF-8" + string('"') + "\n" +
			"MIME-Version: 1.0\n" +
			"Content-Transfer-Encoding: 7bit\n\n" +
			req.Message + "\n\n" +
			"--" + boundary + "\n" +
			attachmentStr +
			"--" + boundary + "--")
	}
	message.Raw = base64.StdEncoding.EncodeToString([]byte(messageBody))
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)
	if _, err = gmailSvc.Users.Messages.Send(meSelf, &message).Do(); err != nil {
		ctxSess.ErrorMessage = err.Error()
		return err
	}
	return nil
}

func (p *provider) WatchPushNotification(ctxSess *ctxSess.Context, refreshToken string, req *gmail.WatchRequest) (res *gmail.WatchResponse, err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, refreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, refreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return nil, err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))
	res, err = gmailSvc.Users.Watch(meSelf, req).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return
}

func (p *provider) StopWatchPushNotification(ctxSess *ctxSess.Context, refreshToken string) (err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, refreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, refreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))
	err = gmailSvc.Users.Stop(meSelf).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return
}

func (p *provider) GetHistory(ctxSess *ctxSess.Context, historyID uint64) (resp *gmail.ListHistoryResponse, err error) {
	token, err := p.ValidateToken(ctxSess, constants.IosAgent, ctxSess.UserSession.RefreshToken)
	if err != nil {
		token, err = p.ValidateToken(ctxSess, constants.WebAgent, ctxSess.UserSession.RefreshToken)
		if err != nil {
			err = errors.New(fmt.Sprintf("token GetMessages is nil: %s ", err.Error()))
			ctxSess.ErrorMessage = err.Error()
			return nil, err
		}
	}
	gmailSvc, err := gmail.New(p.oauthConfig.Client(context.Background(), token))
	resp, err = gmailSvc.Users.History.List(meSelf).HistoryTypes("messageAdded").StartHistoryId(historyID).Do()
	if err != nil {
		ctxSess.ErrorMessage = err.Error()
	}
	return
}
