package emailmanager

import (
	"time"

	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/manager/jwtmanager"
	"github.com/mayadata-io/kubera-auth/pkg/generates"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/types"
)

func SendVerificationEmail(accessGenerate *generates.JWTAccessGenerate, userInfo *models.PublicUserInfo) error {
	tgr := &jwtmanager.TokenGenerateRequest{
		UserInfo:       userInfo,
		AccessTokenExp: time.Minute * types.VerificationLinkExpirationTimeUnit,
	}

	tokenInfo, err := jwtmanager.GenerateAuthToken(accessGenerate, tgr, models.TokenVerify)
	if err != nil {
		return err
	}

	link := types.PortalURL + "/api/auth/v1/email?access=" + tokenInfo.Access

	buf, err := generates.GetEmailBody(userInfo.Name, link)
	if err != nil {
		log.Error("Error occurred while getting email body for user: " + userInfo.UID + "error: " + err.Error())
		return err
	}

	err = generates.SendEmail(userInfo.UnverifiedEmail, "Email Verification", buf.String())
	if err != nil {
		log.Error("Error occurred while sending email for user: " + userInfo.UID + "error: " + err.Error())
		return err
	}

	return nil
}
