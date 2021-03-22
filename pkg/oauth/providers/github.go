package providers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"
	"github.com/google/go-github/github"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
	"golang.org/x/oauth2"
)

// getUserFromToken Returns the user information from the token
func getGitHubUser(c *gin.Context, token *oauth2.Token) (*models.UserCredentials, error) {
	ctx := c.Request.Context()
	ts := oauth2.StaticTokenSource(token)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	githubUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	githubUserEmails, _, err := client.Users.ListEmails(ctx, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	currTime := time.Now()
	user := models.UserCredentials{
		Name:         githubUser.GetName(),
		Kind:         models.GithubAuth,
		Role:         models.RoleUser,
		State:        models.StateActive,
		LoggedIn:     true,
		SocialAuthID: strconv.FormatInt(*githubUser.ID, 10),
		CreatedAt:    &currTime,
	}

	for _, githubUserEmail := range githubUserEmails {
		if githubUserEmail.Primary != nil && githubUserEmail.Email != nil {
			user.Email = githubUserEmail.GetEmail()
			user.OnBoardingState = models.BoardingStateEmailVerified
			break
		}
	}

	return &user, err
}

// GetGithubUser gives the details of the user fetched as from github
func GetGithubUser(c *gin.Context) (*models.UserCredentials, error) {
	token, err := controller.Server.GithubConfig.GetToken(c)
	if err != nil {
		log.Errorln("Error getting token from github", err)
		return nil, err
	}

	githubUser, err := getGitHubUser(c, token)
	if err != nil {
		log.Errorln("Error getting user from github", err)
		return nil, err
	}

	return githubUser, nil
}
