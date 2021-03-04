package providers

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
	"golang.org/x/oauth2"
)

// getUserFromToken gets the user model based on the structure
func getUserFromToken(c *gin.Context, token *oauth2.Token) (*models.UserCredentials, error) {
	ctx := c.Request.Context()
	ts := oauth2.StaticTokenSource(token)
	tc := oauth2.NewClient(ctx, ts)
	// client := google.DefaultClient()

	googleUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	currTime := time.Now()
	user := models.UserCredentials{
		Name:         googleUser.Name,
		Kind:         models.GmailAuth,
		Role:         models.RoleUser,
		State:        models.StateActive,
		SocialAuthID: googleUser.ID,
		LoggedIn:     &types.TrueValue,
		CreatedAt:    &currTime,
	}

	return &user, err
}

// GetGoogleUser gives the details of the user fetched as from github
func GetGoogleUser(c *gin.Context) (*models.UserCredentials, error) {
	token, err := controller.Server.GoogleConfig.GetToken(c)
	if err != nil {
		log.Errorln("Error getting token from github", err)
		return nil, err
	}

	githubUser, err := getUserFromToken(c, token)
	if err != nil {
		log.Errorln("Error getting user from github", err)
		return nil, err
	}

	return githubUser, nil
}
