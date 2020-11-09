package router

import (
	"net/http"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/oauth/providers"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	v1 "github.com/mayadata-io/kubera-auth/versionedController/v1"
	"github.com/mayadata-io/kubera-auth/versionedController/v1/configuration"
	"github.com/mayadata-io/kubera-auth/versionedController/v1/login"
	"github.com/mayadata-io/kubera-auth/versionedController/v1/password"
	"github.com/mayadata-io/kubera-auth/versionedController/v1/user"
)

const (
	githubLoginRoute = "/oauth"
	healthCheckRoute = "/health"
)

var (
	controllers = []v1.Controller{
		login.New(),
		user.New(),
		password.New(),
		configuration.New(),
	}
)

func registerControllers(router *gin.RouterGroup) {
	for _, controller := range controllers {
		controller.Register(router)
	}
}

// New will create a new routes
func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.EnableJsonDecoderDisallowUnknownFields()
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AddAllowHeaders("Access-Control-Allow-Origin", "Authorization")
	config.AllowAllOrigins = true

	router.Use(cors.New(config))

	v1.InitializeServer()
	routerV1 := router.Group("v1")
	{
		routerV1.GET(githubLoginRoute, CallbackRequest)
		routerV1.GET(healthCheckRoute, HealthCheck)
	}
	registerControllers(routerV1)

	return router
}

// HealthCheck will respond with the server status
func HealthCheck(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

//CallbackRequest will be triggered by the provider automatically after the login
func CallbackRequest(c *gin.Context) {

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

	ti, err := v1.Server.SocialLoginRequest(c, user)
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
