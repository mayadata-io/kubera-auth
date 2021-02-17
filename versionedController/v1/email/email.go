package email

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
)

type EmailController struct {
	controller.GenericController
	routePath string
	model     *Model
}

//Model ...
type Model struct {
	Email string `json:"email,omitempty"`
}

// New creates a new User
func New() *EmailController {
	return &EmailController{
		routePath: controller.EmailRoute,
		model:     &Model{},
	}
}

//Post send the verification link to an email address
func (emailController *EmailController) Post(c *gin.Context) {
	err := c.BindJSON(emailController.model)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Unable to parse JSON",
		})
		return
	}

	controller.Server.SendVerificationLink(c, emailController.model.Email)
}

//Get verifies the email by a link
func (emailController *EmailController) Get(c *gin.Context) {
	token := c.Query("access")

	jwtUserInfo, err := controller.Server.GetUserFromToken(token)
	if err != nil {
		log.Error("Error occurred while parsing jwt token error: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Set("userInfo", jwtUserInfo)

	controller.Server.VerifyEmail(c)
}

// Register will rsgister this controller to the specified router
func (emailController *EmailController) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, emailController, emailController.routePath)
}
