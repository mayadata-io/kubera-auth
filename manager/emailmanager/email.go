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

func SendVerificationEmail(accessGenerate *generates.JWTAccessGenerate, userInfo *models.PublicUserInfo, emailType EmailType) error {
	tgr := &jwtmanager.TokenGenerateRequest{
		UserInfo:       userInfo,
		AccessTokenExp: time.Minute * types.VerificationLinkExpirationTimeUnit,
	}

	tokenInfo, err := jwtmanager.GenerateAuthToken(accessGenerate, tgr, models.TokenEmail)
	if err != nil {
		return err
	}

	var email string
	// var emailType EmailType
	var buf *bytes.Buffer

	switch emailType {
	case VerificationEmail:
		email = userInfo.UnverifiedEmail
		link := types.PortalURL + "/api/auth/v1/email?access=" + tokenInfo.Access
		buf, err = generates.GetEmailBody(types.VerificationEmailTemplatePath, userInfo.Name, link, "")
		if err != nil {
			log.Error("Error occurred while getting email body for user: " + userInfo.UID + "error: " + err.Error())
			return err
		}
	case ResetPasswordEmail:
		email = userInfo.Email
		link := types.PortalURL + "/change-password?access=" + tokenInfo.Access
		retriggerLink := types.PortalURL + "/api/auth/v1/password?email=" + email
		buf, err = generates.GetEmailBody(types.ResetPasswordEmailTemplatePath, userInfo.Name, link, retriggerLink)
		if err != nil {
			log.Error("Error occurred while getting email body for user: " + userInfo.UID + "error: " + err.Error())
			return err
		}
	}

	err = generates.SendEmail(email, "Email Verification", buf.String())
	if err != nil {
		log.Error("Error occurred while sending email for user: " + userInfo.UID + "error: " + err.Error())
		return err
	}

	return nil
}
