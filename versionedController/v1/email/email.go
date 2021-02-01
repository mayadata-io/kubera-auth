package email

import (
	"net/http"

	log "github.com/golang/glog"

	"github.com/gin-gonic/gin"
	"github.com/mayadata-io/kubera-auth/pkg/generates"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/types"
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
			"message": "Unable to parse JSON",
		})
		return
	}

	jwtUser, _ := c.Get(types.UserInfoKey)
	jwtUserInfo := jwtUser.(*models.PublicUserInfo)

	link := controller.Server.GenerateVerificationLink(c, emailController.model.Email)
	if link == "" {
		return
	}

	buf, err := generates.GetEmailBody(jwtUserInfo.Name, link)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to parse email template",
		})
		return
	}

	err = generates.SendEmail(emailController.model.Email, "Email Verification", buf.String())
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to send email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification Email sent",
	})
}

//Get verifies the email by a link
func (emailController *EmailController) Get(c *gin.Context) {
	token := c.Query("access")

	jwtUserInfo, err := controller.Server.Manager.ParseToken(token)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{
			"message": "Invalid Token",
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