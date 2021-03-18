package login

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
	"golang.org/x/oauth2"
)

//LoginController is the type the request in which the request will be parsed
type LoginController struct {
	controller.GenericController
	routePath string
}

// Model defines the json struct in which the request will be parsed
type Model struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// New creates a new LoginUser
func New() *LoginController {
	return &LoginController{
		routePath: controller.TokenRoute,
	}
}

func init() {
	if os.Getenv("PORTAL_URL") == "" {
		log.Fatal("Environment variables PORTAL_URL are not set")
	}
}

// Post lets a user login into the kubera-core
func (login *LoginController) Post(c *gin.Context) {
	loginModel := &Model{}
	err := c.BindJSON(loginModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}
	controller.Server.LocalLoginRequest(c, loginModel.Username, loginModel.Password)
}

// Get will be triggered on GET request on the same path as Login along with a "auth_type" parameter
// so as to identify the type of login user is up to. This has to be triggered through a href request
// so that the user is able to be redirected to provider page for login.
// Javascript Get Request can block the redirection of user
func (login *LoginController) Get(c *gin.Context) {
	authType := c.Query("auth_type")
	switch models.AuthType(authType) {
	case models.GithubAuth:
		if !controller.Server.Config.DisableGithubAuth {
			githubURL := controller.Server.GithubConfig.AuthCodeURL(types.GithubState)
			c.Redirect(http.StatusFound, githubURL)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Authentication type not allowed",
			})
		}
	case models.GoogleAuth:
		if !controller.Server.Config.DisableGoogleAuth {
			googleURL := controller.Server.GoogleConfig.AuthCodeURL(types.GoogleState,
				oauth2.SetAuthURLParam("include_granted_scopes", "true"), oauth2.AccessTypeOnline)
			log.Info(googleURL)
			c.Redirect(http.StatusFound, googleURL)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Authentication type not allowed",
			})
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unknown Authentication Type",
		})
	}
}

// Delete lets a user logout of the kubera-core
func (login *LoginController) Delete(c *gin.Context) {
	controller.Server.LogoutRequest(c)
}

// Register will rsgister this controller to the specified router
func (login *LoginController) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, login, login.routePath)
}
