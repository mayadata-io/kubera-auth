package email

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/pkg/types"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
)

type EmailController struct {
	controller.GenericController
	routePath string
}

//Model ...
type Model struct {
	UnverifiedEmail string `json:"unverified_email,omitempty"`
	Resend          bool   `json:"resend,omitempty"`
	Restore         bool   `json:"restore,omitempty"`
}

// New creates a new User
func New() *EmailController {
	return &EmailController{
		routePath: controller.EmailRoute,
	}
}

//Post send the verification link to an email address
func (emailController *EmailController) Post(c *gin.Context) {
	emailModel := &Model{}
	err := c.BindJSON(emailModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Unable to parse JSON",
		})
		return
	}

	if emailModel.Restore {
		controller.Server.RestoreEmail(c)
	} else {
		controller.Server.SendVerificationLink(c, emailModel.Resend, emailModel.UnverifiedEmail)
	}
}

//Get verifies the email by a link
func (emailController *EmailController) Get(c *gin.Context) {
	token := c.Query("access")
	redirectURL := types.PortalURL + "/verified-email"

	jwtUserCredentials, err := controller.Server.GetUserFromToken(token)
	if err != nil {
		log.Error("Error occurred while parsing jwt token error: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		// Redirecting user back to UI if access token is not valid
		c.Redirect(http.StatusPermanentRedirect, redirectURL)
		return
	}
	c.Set(types.JWTUserCredentialsKey, jwtUserCredentials)

	controller.Server.VerifyEmail(c, redirectURL)
}

// Register will register this controller to the specified router
func (emailController *EmailController) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, emailController, emailController.routePath)
}
