package providers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
	"golang.org/x/oauth2"
)

const (
	userInfo = "https://www.googleapis.com/userinfo/v2/me"
)

type account struct {
	// LastName = `family_name` of user
	LastName string `json:"family_name"`
	// Name is the full-name of the user
	Name string `json:"name"`
	// Picture is the link of the Google Profile picture of the user
	Picture string `json:"picture"`
	Locale  string `json:"locale"`
	Email   string `json:"email"`
	// FirstName = `given_name` in json
	FirstName string `json:"given_name"`
	// ID has to be a string because it overflows an unsigned long integer
	ID string `json:"id"`
	// Hd is the hosted G Suite domain of the user. Provided only if the user
	// belongs to a hosted domain.
	Hd string `json:"hd,omitempty"`
	// VerifiedEmail is the value of the email verification status
	VerifiedEmail bool `json:"verified_email"`
}

// getGoogleUser gets the user model based on the structure
func getGoogleUser(c *gin.Context, token *oauth2.Token) (*models.UserCredentials, error) {
	ctx := c.Request.Context()
	// use oauth2.ReuseTokenSource when access_type is set to offline, refreshing every hour
	ts := oauth2.StaticTokenSource(token)
	tc := oauth2.NewClient(ctx, ts)
	resp, _ := tc.Get(userInfo)
	var gUser account
	if resp.StatusCode == http.StatusOK {
		_ = json.NewDecoder(resp.Body).Decode(&gUser)
	}
	user := models.UserCredentials{
		Name:         gUser.Name,
		Kind:         models.GmailAuth,
		Role:         models.RoleUser,
		Email:        gUser.Email,
		State:        models.StateActive,
		SocialAuthID: gUser.ID,
		LoggedIn:     true,
		Photo:        gUser.Picture,
	}
	return &user, nil
}

// GetGoogleUser gives the details of the user fetched as from github
func GetGoogleUser(c *gin.Context) (*models.UserCredentials, error) {
	token, err := controller.Server.GoogleConfig.GetToken(c)
	if err != nil {
		log.Errorln("Error getting token from google", err)
		return nil, err
	}
	user, err := getGoogleUser(c, token)
	if err != nil {
		log.Errorln("Error getting user from google", err)
		return nil, err
	}
	return user, nil
}
