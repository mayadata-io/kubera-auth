package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	v1 "github.com/mayadata-io/kubera-auth/versionedController/v1"
	"github.com/mayadata-io/kubera-auth/versionedController/v1/login"
	"github.com/mayadata-io/kubera-auth/versionedController/v1/password"
	"github.com/mayadata-io/kubera-auth/versionedController/v1/user"
)

const (
	loginRoute          = "/login"
	githubLoginRoute    = "/oauth"
	updatePasswordRoute = "/update/password"
	resetPasswordRoute  = "/reset/password"
	createRoute         = "/create"
	updateDetailsRoute  = "/update/details"
	getUsersRoute       = "/users"
	logoutRoute         = "/logout"
)

var (
	userController     v1.UserController     = user.New()
	loginController    v1.LoginController    = login.New()
	passwordController v1.PasswordController = password.New()
)

// New will create a new routes
func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.EnableJsonDecoderDisallowUnknownFields()
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AddAllowHeaders("Access-Control-Allow-Origin", "Authorization")
	config.AllowAllOrigins = true

	router.Use(cors.New(config))

	v1 := router.Group("v1")
	{
		v1.GET(getUsersRoute, userController.GetAllUsers)
		v1.POST(logoutRoute, userController.Logout)
		v1.POST(loginRoute, loginController.Login)
		v1.GET(loginRoute, loginController.SocialLogin)
		v1.GET(githubLoginRoute, loginController.CallbackRequest)
		v1.POST(updatePasswordRoute, passwordController.Update)
		v1.POST(createRoute, userController.Create)
		v1.POST(updateDetailsRoute, userController.UpdateUserDetails)
		v1.POST(resetPasswordRoute, passwordController.Reset)
	}

	return router
}
