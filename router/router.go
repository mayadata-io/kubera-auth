package router

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/oauth/providers"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	v1 "github.com/mayadata-io/kubera-auth/versionedController/v1"
	"github.com/mayadata-io/kubera-auth/versionedController/v1/configuration"
	"github.com/mayadata-io/kubera-auth/versionedController/v1/email"
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
		email.New(),
	}
	unauthenticatedLinks = map[string][]string{
		"/v1" + v1.TokenRoute:         {http.MethodPost, http.MethodGet},
		"/v1" + v1.EmailRoute:         {http.MethodGet},
		"/v1" + "/oauth":              {http.MethodGet},
		"/v1" + v1.ConfigurationRoute: {http.MethodGet},
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
	routerV1 := router.Group("/v1")
	routerV1.Use(Middleware)
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

	v1.Server.SocialLoginRequest(c, user, u)
}

//Middleware ...
func Middleware(c *gin.Context) {

	auth := c.Request.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	for path, methods := range unauthenticatedLinks {
		if path == c.Request.URL.Path {
			for _, method := range methods {
				if method == c.Request.Method {
					return
				}
			}
		}
	}

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}

	if token == "" {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Token",
		})
		return
	}

	userInfo, err := v1.Server.GetUserFromToken(token)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Set(types.UserInfoKey, userInfo)
}
