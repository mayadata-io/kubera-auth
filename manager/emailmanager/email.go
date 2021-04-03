package emailmanager

import (
	"bytes"
	"time"

	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/manager/jwtmanager"
	"github.com/mayadata-io/kubera-auth/pkg/generates"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/types"
)

type EmailType string

const (
	VerificationEmail  EmailType = "Verification"
	ResetPasswordEmail EmailType = "Reset"
)

func SendEmail(accessGenerate *generates.JWTAccessGenerate, userInfo *models.PublicUserInfo, emailType EmailType) error {
	tgr := &jwtmanager.TokenGenerateRequest{
		UserInfo:       userInfo,
		AccessTokenExp: time.Minute * types.VerificationLinkExpirationTimeUnit,
	}

	tokenInfo, err := jwtmanager.GenerateAuthToken(accessGenerate, tgr, models.TokenEmail)
	if err != nil {
		return err
	}

	var email string
	var buf *bytes.Buffer
	var subject string

	switch emailType {
	case VerificationEmail:
		email = userInfo.UnverifiedEmail
		templateVar := generates.TemplateVariables{
			Username: userInfo.Name,
			Link:     types.PortalURL + "/api/auth/v1/email?access=" + tokenInfo.Access,
		}
		subject = "Email Verification"

		buf, err = generates.GetEmailBody(types.VerificationEmailTemplatePath, templateVar)
		if err != nil {
			log.Error("Error occurred while getting email body for user: " + userInfo.UID + "error: " + err.Error())
			return err
		}
	case ResetPasswordEmail:
		email = userInfo.Email
		templateVar := generates.TemplateVariables{
			Username:      userInfo.Name,
			Link:          types.PortalURL + "/change-password?access=" + tokenInfo.Access,
			RetriggerLink: types.PortalURL + "/api/auth/v1/password?access=" + tokenInfo.Access,
		}
		subject = "Password Reset"

		buf, err = generates.GetEmailBody(types.ResetPasswordEmailTemplatePath, templateVar)
		if err != nil {
			log.Error("Error occurred while getting email body for user: " + userInfo.UID + "error: " + err.Error())
			return err
		}
	}

	err = generates.SendEmail(email, subject, buf.String())
	if err != nil {
		log.Error("Error occurred while sending email for user: " + userInfo.UID + "error: " + err.Error())
		return err
	}

	return nil
}
