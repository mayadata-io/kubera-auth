package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/mayadata-io/kubera-auth/controller"
	"github.com/mayadata-io/kubera-auth/controller/login"
	"github.com/mayadata-io/kubera-auth/controller/password"
	"github.com/mayadata-io/kubera-auth/controller/user"
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
	userController     controller.UserController     = user.New()
	loginController    controller.LoginController    = login.New()
	passwordController controller.PasswordController = password.New()
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

	// Handle the request for chaos-schedule
	router.GET(getUsersRoute, userController.GetAllUsers)
	router.POST(logoutRoute, userController.Logout)
	router.POST(loginRoute, loginController.Login)
	router.GET(loginRoute, loginController.SocialLogin)
	router.GET(githubLoginRoute, loginController.CallbackRequest)
	router.POST(updatePasswordRoute, passwordController.Update)
	router.POST(createRoute, userController.Create)
	router.POST(updateDetailsRoute, userController.UpdateUserDetails)
	router.POST(resetPasswordRoute, passwordController.Reset)
	return router
}
