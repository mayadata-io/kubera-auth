package generates

import (
	"bytes"
	"html/template"

	"gopkg.in/gomail.v2"

	"github.com/mayadata-io/kubera-auth/pkg/types"
)

type TemplateVariables struct {
	Link          string
	Username      string
	RetriggerLink string
}

var (
	emailHost   = "smtp.gmail.com"
	emailPort   = 465
	contentType = "text/html"
	keySubject  = "Subject"
	keyTo       = "To"
	keyFrom     = "From"
)

// SendEmail will send the email to the defined email.
func SendEmail(sendTo, subject string, body string) error {
	dialer, message := configureEmail()
	message.SetHeader(keyTo, sendTo)
	message.SetHeader(keySubject, subject)
	message.SetBody(contentType, body)
	message.Embed(types.TemplatePath + types.KuberaPortalImagePath)
	message.Embed(types.TemplatePath + types.MayadataLogoImagePath)
	message.Embed(types.TemplatePath + types.BackgroundEmailImagePath)
	return dialer.DialAndSend(message)
}

/*
	configureEmail will initialize the configuration for sending email.
	It will configure the Host, Port, Username, Password and from.
*/
func configureEmail() (*gomail.Dialer, *gomail.Message) {
	dialer := gomail.NewDialer(emailHost, emailPort, types.EmailUsername, types.EmailPassword)
	dialer.SSL = true
	message := gomail.NewMessage()
	message.SetHeader(keyFrom, types.EmailUsername)
	return dialer, message
}

// GetEmailBody forms the html template body of email
func GetEmailBody(emailTemplateFilePath string, templateVar TemplateVariables) (*bytes.Buffer, error) {
	t, err := template.ParseFiles(types.TemplatePath + emailTemplateFilePath)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, templateVar)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
