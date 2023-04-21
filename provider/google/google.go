package google

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/people/v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	"github.com/cloudsrc/api.awaymail.v1.go/models"
)

const endpointProfile string = "https://www.googleapis.com/oauth2/v2/userinfo"

var scopes = []string{
	people.ContactsReadonlyScope,
	gmail.MailGoogleComScope,
	calendar.CalendarReadonlyScope,
}

var (
	LABEL_INBOX            = "INBOX"
	LABEL_SENT             = "SENT"
	LABEL_UNREAD           = "UNREAD"
	USER_PER_DAY_API_LIMIT = 86400
)

// Provider : Provider
type Provider struct {
	Config        *config.Config
	Srv           *gmail.Service
	Redis         libs.RedisClient
	MongoClient   *mgo.Database
	Rabbit        *libs.RabbitClient
	User          models.User
	providerName  string
	agent         string
	nextPageToken string
}

type Messsage struct {
	To             string   `json:"to"`
	Subject        string   `json:"subject"`
	Message        string   `json:"message"`
	AttachmentsURL []string `json:"attachments_url,omitempty"`
}

// New : new provider config
func New(config *config.Config, rclient libs.RedisClient, mclient *mgo.Database, rabbitClient *libs.RabbitClient) *Provider {
	return &Provider{
		Redis:       rclient,
		Config:      config,
		MongoClient: mclient,
		Rabbit:      rabbitClient,
	}
}

// Name is the name used to retrieve this provider later.
func (p *Provider) Name() string {
	return p.providerName
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *Provider) SetName(name string) {
	p.providerName = name
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *Provider) SetAgent(agent string) {
	p.agent = agent
}

// GetUserProfile is to retrive user profile information
func (p *Provider) GetUserProfile() error {
	clientID := p.Config.Google.Ios.ClientID
	clientSecret := p.Config.Google.Ios.ClientSecret
	if p.agent == "website" {
		clientID = p.Config.Google.Website.ClientID
		clientSecret = p.Config.Google.Website.ClientSecret
	}
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
	token := &oauth2.Token{RefreshToken: p.User.RefreshToken}
	tok, err := config.TokenSource(context.TODO(), token).Token()
	if err != nil {
		return err
	}
	p.User.Self = "me"
	if tok.Valid() {
		srv, err := gmail.New(config.Client(context.Background(), tok))
		if err != nil {
			return err
		}
		p.Srv = srv
		profile, err := p.Srv.Users.GetProfile(p.User.Self).Do()
		if err != nil {
			return err
		}
		p.User = models.User{
			DB:           p.MongoClient,
			Email:        profile.EmailAddress,
			AuthToken:    tok.AccessToken,
			RefreshToken: tok.RefreshToken,
		}
		if err := p.User.GetUser(); err != nil {
			return err
		}
		go p.Rabbit.Consume("inbox_processor_"+p.User.ID.Hex(), p.ProcessInbox)
		go p.Rabbit.Consume("sent_processor_"+p.User.ID.Hex(), p.ProcessSent)
	} else {
		return errors.New("token not valid")
	}
	return nil
}

// ValidateToken : ValidateToken
func (p *Provider) ValidateToken(token string) error {
	p.User.RefreshToken = token
	if err := p.GetUserProfile(); err != nil {
		return err
	}
	return nil
}

// GetUserInbox : get user inbox using Gmail API
func (p *Provider) GetUserInbox() error {
	var messageList []string
	skip := "initial"
	p.Redis.CreateIndex(fmt.Sprintf("user:inbox:%s", p.User.ID.Hex()),
		fmt.Sprintf(`ON hash PREFIX 1 user:inbox:%s: SCHEMA `+
			`body_text TEXT NOSTEM `+
			`from TEXT NOSTEM `+
			`subject TEXT NOSTEM `+
			`received_at NUMERIC SORTABLE`,
			p.User.ID.Hex()))
	for skip != "" {
		mList, _ := p.getMessagesList(LABEL_INBOX)
		for _, ii := range mList {
			cQuery := "user:inbox:" + p.User.ID.Hex() + ":" + ii.Id
			messageList = append(messageList, cQuery)
			if res := p.Redis.GetCache(cQuery, "message"); res == nil {
				if err := p.Rabbit.Publish("inbox_processor_"+p.User.ID.Hex(), map[string]interface{}{
					"message_id": ii.Id,
				}); err != nil {
					return err
				}
			}
		}
		skip = p.nextPageToken
	}
	p.removeDiffMessages(difference(p.Redis.GetKeysPrefix("user:inbox:"+p.User.ID.Hex()+":*"), messageList))
	return nil
}

// GetUserSent : get user inbox using Gmail API
func (p *Provider) GetUserSent() error {
	var messageList []string
	skip := "initial"
	p.Redis.CreateIndex(fmt.Sprintf("user:sent:%s", p.User.ID.Hex()),
		fmt.Sprintf(`ON hash PREFIX 1 user:sent:%s: SCHEMA `+
			`body_text TEXT NOSTEM `+
			`from TEXT NOSTEM `+
			`subject TEXT NOSTEM `+
			`sent NUMERIC SORTABLE`,
			p.User.ID.Hex()))
	for skip != "" {
		mList, _ := p.getMessagesList(LABEL_SENT)
		for _, ii := range mList {
			cQuery := "user:sent:" + p.User.ID.Hex() + ":" + ii.Id
			messageList = append(messageList, cQuery)
			if res := p.Redis.GetCache(cQuery, "message"); res == nil {
				if err := p.Rabbit.Publish("sent_processor_"+p.User.ID.Hex(), map[string]interface{}{
					"message_id": ii.Id,
				}); err != nil {
					return err
				}
			}
		}
		skip = p.nextPageToken
	}
	p.removeDiffSentMessages(difference(p.Redis.GetKeysPrefix("user:sent:"+p.User.ID.Hex()+":*"), messageList))
	return nil
}

// UpdateMessage : update user message
func (p *Provider) UpdateMessage(id string, isRead *bool) (interface{}, error) {
	var req gmail.ModifyMessageRequest
	if !*isRead {
		req.AddLabelIds = append(req.AddLabelIds, LABEL_UNREAD)
	} else {
		req.RemoveLabelIds = append(req.AddLabelIds, LABEL_UNREAD)
	}
	res, err := p.Srv.Users.Messages.Modify(p.User.Self, id, &req).Do()
	if err != nil {
		return nil, err
	}
	return p.GetMessage(res.Id)
}

// GetMessage : get message for the user based on the id provided.
func (p *Provider) GetMessage(id string) (interface{}, error) {
	msg := models.Message{
		DB:        p.MongoClient,
		MessageID: id,
		UserID:    p.User.ID,
	}
	cQuery := "user:inbox:" + p.User.ID.Hex() + ":" + id
	mess, err := p.Srv.Users.Messages.Get(p.User.Self, id).Do()
	if err != nil {
		return nil, err
	}
	isRead := true
	for _, label := range mess.LabelIds {
		if label == LABEL_UNREAD {
			isRead = false
		}
	}
	received := time.Unix(0, mess.InternalDate*1000000).UTC()
	msg.Received = &received
	msg.IsRead = isRead
	p.parseMessagePart(mess.Payload, &msg)
	msg.Insert("messages.inbox", bson.M{"message_id": id})
	encoded, _ := json.Marshal(msg)
	if msg.From == nil {
		p.Redis.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", "noreply_example.com", "received_at", mess.InternalDate)
		return msg, nil
	}
	email, _ := mail.ParseAddress(*msg.From)
	if email == nil {
		p.Redis.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", "noreply_example.com", "received_at", mess.InternalDate)
		return msg, nil
	}
	p.Redis.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", strings.ReplaceAll(email.Address, "@", "_"), "received_at", mess.InternalDate)
	return msg, nil
}
func (p *Provider) GetSentMessage(id string) (interface{}, error) {
	msg := models.Message{
		DB:        p.MongoClient,
		MessageID: id,
		UserID:    p.User.ID,
	}
	cQuery := "user:sent:" + p.User.ID.Hex() + ":" + id
	mess, err := p.Srv.Users.Messages.Get(p.User.Self, id).Do()
	if err != nil {
		return nil, err
	}
	sent := time.Unix(0, mess.InternalDate*1000000).UTC()
	msg.Sent = &sent
	p.parseMessagePart(mess.Payload, &msg)
	msg.Insert("messages.sent", bson.M{"message_id": id})
	encoded, _ := json.Marshal(msg)
	if msg.From == nil {
		p.Redis.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", "noreply_example.com", "sent", mess.InternalDate)
		return msg, nil
	}
	email, _ := mail.ParseAddress(*msg.From)
	if email == nil {
		p.Redis.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", "noreply_example.com", "sent", mess.InternalDate)
		return msg, nil
	}
	p.Redis.SetCache(cQuery, "message", encoded, "body_text", msg.BodyText, "subject", msg.Subject, "from", strings.ReplaceAll(email.Address, "@", "_"), "sent", mess.InternalDate)
	return msg, nil
}

// SendMessage : send the messages using Gmail API
func (p *Provider) SendMessage(m interface{}) error {
	mess := m.(Messsage)
	var message gmail.Message
	msg := fmt.Sprintf("From: %s\r\n"+
		"reply-to: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n%s", p.User.Self, p.User.Email, mess.To, mess.Subject, mess.Message)
	message.Raw = base64.StdEncoding.EncodeToString([]byte(msg))
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	if _, err := p.Srv.Users.Messages.Send(p.User.Self, &message).Do(); err != nil {
		return err
	}
	return nil
}

func (p *Provider) SendMessageWithAttachment(m interface{}) error {
	mess := m.(Messsage)
	boundary := randStr(32, "alphanum")
	var message gmail.Message
	var attachmentStr string
	for count, attachmentURL := range mess.AttachmentsURL {
		url := strings.Split(attachmentURL, "/")
		filename := url[len(url)-1]
		filePath := "./public/" + filename
		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			if err := downloadFile(attachmentURL, filePath); err != nil {
				return err
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
		"to: " + mess.To + "\n" +
		"subject: " + mess.Subject + "\n\n" +

		"--" + boundary + "\n" +
		"Content-Type: text/plain; charset=" + string('"') + "UTF-8" + string('"') + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: 7bit\n\n" +
		mess.Message + "\n\n" +
		"--" + boundary + "\n" +

		attachmentStr +
		"--" + boundary + "--")
	message.Raw = base64.StdEncoding.EncodeToString([]byte(messageBody))
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)
	if _, err := p.Srv.Users.Messages.Send(p.User.Self, &message).Do(); err != nil {
		return err
	}
	return nil
}

func (p *Provider) getMessagesList(label string) ([]*gmail.Message, error) {
	var (
		response *gmail.ListMessagesResponse
		err      error
	)
	if p.userLimitCheck(p.User.ID.Hex()) {
		if p.nextPageToken == "" {
			response, err = p.Srv.Users.Messages.List(p.User.Self).
				LabelIds(label).Do()
			if err != nil {
				return nil, err
			}
		} else {
			response, err = p.Srv.Users.Messages.List(p.User.Self).
				LabelIds(label).
				PageToken(p.nextPageToken).Do()
			if err != nil {
				return nil, err
			}
		}
	}
	p.nextPageToken = response.NextPageToken
	return response.Messages, nil
}

func (p *Provider) parseMessagePart(origmsg *gmail.MessagePart, mess *models.Message) {
	mess.MimeFlow = append(mess.MimeFlow, origmsg.MimeType)
	for _, ii := range origmsg.Headers {
		if ii.Name == "Subject" {
			mess.Subject = ii.Value
		}
		if ii.Name == "From" {
			mess.From = &ii.Value
		}
	}
	if strings.HasPrefix(origmsg.MimeType, "multipart") {
		for _, ii := range origmsg.Parts {
			p.parseMessagePart(ii, mess)
		}
	}
	if origmsg.MimeType == "text/plain" {
		sBodyText, _ := base64.URLEncoding.DecodeString(origmsg.Body.Data)
		mess.BodyText = string(sBodyText)
	} else if origmsg.MimeType == "text/html" {
		sBodyHTML, _ := base64.URLEncoding.DecodeString(origmsg.Body.Data)
		mess.BodyHTML = string(sBodyHTML)
	}
}

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	var strBytes = make([]byte, strSize)
	_, _ = rand.Read(strBytes)
	for k, v := range strBytes {
		strBytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(strBytes)
}

func chunkSplit(body string, limit int, end string) string {
	var charSlice []rune

	// push characters to slice
	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	var result = ""

	for len(charSlice) >= 1 {
		// convert slice/array back to string
		// but insert end at specified limit
		result = result + string(charSlice[:limit]) + end

		// discard the elements that were copied over to result
		charSlice = charSlice[limit:]

		// change the limit
		// to cater for the last few words in
		if len(charSlice) < limit {
			limit = len(charSlice)
		}
	}
	return result
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("file or url not found")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	//Write the bytes to the fiel
	if _, err := io.Copy(file, response.Body); err != nil {
		return err
	}
	return nil
}

func (p *Provider) removeDiffMessages(diff []string) {
	if len(diff) != 0 {
		for _, d := range diff {
			p.Redis.DelCache(d, "message", "body_text", "from", "subject", "received_at")
		}
	}
}

func (p *Provider) removeDiffSentMessages(diff []string) {
	if len(diff) != 0 {
		for _, d := range diff {
			p.Redis.DelCache(d, "message", "body_text", "from", "subject", "posted_at")
		}
	}
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

func (p *Provider) userLimitCheck(key string) bool {
	val := p.Redis.Getkey(key)
	if val == "" {
		p.Redis.SetKey(key, 1, "EX", 86400)
		return true
	}
	totalUserAPILimit, _ := strconv.Atoi(val)
	if totalUserAPILimit < USER_PER_DAY_API_LIMIT {
		p.Redis.Incrkey(key)
		return true
	}
	return false
}

func (p *Provider) ProcessInbox(data []byte) error {
	var message struct {
		MessageID string `json:"message_id"`
	}
	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	if _, err := p.GetMessage(message.MessageID); err != nil {
		return err
	}
	return nil
}

func (p *Provider) ProcessSent(data []byte) error {
	var message struct {
		MessageID string `json:"message_id"`
	}
	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	if _, err := p.GetSentMessage(message.MessageID); err != nil {
		return err
	}
	return nil
}
