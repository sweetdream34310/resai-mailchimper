package mail

import (
	"fmt"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email struct {
}

func replaceBody(body, date, url, email, emailName string) string {
	body = strings.ReplaceAll(body, "{{.Email.awaymail.date}}", date)
	body = strings.ReplaceAll(body, "{{.Email.awaymail.breakthrough}}", url)
	body = strings.ReplaceAll(body, "{{.Email.awaymail.email}}", email)
	body = strings.ReplaceAll(body, "{{.Email.awaymail.name}}", emailName)
	return body
}

func ampEmailV2(emailTo, date, ampUrl, url, email, emailName string) (*mail.SGMailV3, error) {
	address := "away@awaymail.com"
	name := "Awaymail"
	from := mail.NewEmail(name, address)
	subject := fmt.Sprintf("Awaymail - Breakthrough Email from %s", email)
	address = emailTo
	name = ""
	text, err := os.ReadFile("./resources/email.txt")
	if err != nil {
		return nil, err
	}
	html, err := os.ReadFile("./resources/email.html")
	if err != nil {
		return nil, err
	}
	amp, err := os.ReadFile("./resources/email.amp.html")
	if err != nil {
		return nil, err
	}
	to := mail.NewEmail(name, address)
	content := mail.NewContent("text/plain", replaceBody(string(text), date, url, email, emailName))
	contentHtml := mail.NewContent("text/html", replaceBody(string(html), date, url, email, emailName))
	contentAmp := mail.NewContent("text/x-amp-html", replaceBody(string(amp), date, ampUrl, email, emailName))

	m := mail.NewV3MailInit(from, subject, to, content, contentAmp, contentHtml)
	//var m *mail.SGMailV3
	//if !strings.Contains(emailTo, "@icloud.com") {
	//	m = mail.NewV3MailInit(from, subject, to, content, contentHtml, contentAmp)
	//} else {
	//	m = mail.NewV3MailInit(from, subject, to, content, contentHtml)
	//}
	var dis bool
	mtrack := &mail.ClickTrackingSetting{
		Enable:     &dis,
		EnableText: &dis,
	}
	mOTrack := &mail.OpenTrackingSetting{
		Enable: &dis,
	}
	tracSetting := &mail.TrackingSettings{
		ClickTracking: mtrack,
		OpenTracking:  mOTrack,
	}
	m.TrackingSettings = tracSetting
	return m, nil
}

func (e *Email) SendEmail(emailTo, date, ampUrl, url, email, emailName string) (int, string, error) {
	client := sendgrid.NewSendClient("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")
	emailV3, err := ampEmailV2(emailTo, date, ampUrl, url, email, emailName)
	if err != nil {
		return 0, "", err
	}
	response, err := client.Send(emailV3)
	if err != nil {
		return 0, "", err
	}

	return response.StatusCode, response.Body, nil
}
