package login

import (
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/oauth/providers"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
)

//LoginUser is the type the request in which the request will be parsed
type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// New creates a new LoginUser
func New() *LoginUser {
	return &LoginUser{}
}

func init() {
	if os.Getenv("PORTAL_URL") == "" {
		log.Fatal("Environment variables PORTAL_URL are not set")
	}
}

// Login lets a user login into the kubera-core
func (loginUser *LoginUser) Login(c *gin.Context) {
	err := c.BindJSON(loginUser)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	controller.Server.LocalLoginRequest(c, loginUser.Username, loginUser.Password)
	return
}

/* SocialLogin will be triggered on GET request on the same path as Login alogn with a "auth_type" parameter
** so as to identify the type of login user is up to. This has to be triggered through a href request
** so that the user is able to be redirected to provider page for login.
** Javascript Get Request can block the redirection of user */
func (loginUser *LoginUser) SocialLogin(c *gin.Context) {

	authType := c.Query("auth_type")
	switch models.AuthType(authType) {
	case models.GithubAuth:
		{
			githubURL := controller.Server.GithubConfig.AuthCodeURL(types.GithubState)
			c.Redirect(http.StatusFound, githubURL)
		}
	default:
		{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unknown Authentication Type",
			})
		}
	}
}

//CallbackRequest will be triggered by the provider automatically after the login
func (loginUser *LoginUser) CallbackRequest(c *gin.Context) {

	var user *models.UserCredentials
	var err error
	u := types.PortalURL + "/login?"
	values := url.Values{}
	state := c.Query("state")

	switch state {
	case types.GithubState:
		{
			user, err = providers.GetGithubUser(c)
			if err != nil {
				log.Errorln("Error getting user from github", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				c.Redirect(http.StatusFound, u)
				return
			}
		}
	default:
		{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "state Invalid",
			})
			c.Redirect(http.StatusFound, u)
			return
		}
	}

	ti, err := controller.Server.SocialLoginRequest(c, user)
	if err != nil {
		log.Errorln("Error logging in", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Redirect(http.StatusFound, u)
		return
	}
	values.Set("access_token", ti.GetAccess())
	c.Redirect(http.StatusFound, u+values.Encode())
}
