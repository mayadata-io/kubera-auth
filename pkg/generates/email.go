package generates

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"
	"gopkg.in/gomail.v2"

	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/types"
)

type TemplateVariables struct {
	Link     string
	Username string
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
	message.Embed("./versionedController/v1/email/kuberaPortal.png")
	message.Embed("./versionedController/v1/email/mayadata-logo.png")
	message.Embed("./versionedController/v1/email/bg-kubera-email.png")
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
func GetEmailBody(c *gin.Context, link string) *bytes.Buffer {

	jwtUser, _ := c.Get("userInfo")
	jwtUserInfo := jwtUser.(*models.PublicUserInfo)

	t, err := template.ParseFiles("./versionedController/v1/email/emailTemplate.html")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to parse email template",
		})
		return nil
	}

	templateVar := TemplateVariables{
		Username: jwtUserInfo.Name,
		Link:     link,
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, templateVar)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to modify email template",
		})
		return nil
	}

	return buf
}
