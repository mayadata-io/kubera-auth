package password

import (
	"net/http"

	log "github.com/golang/glog"

	"github.com/gin-gonic/gin"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
)

type PasswordController struct {
	controller.GenericController
	routePath string
}

// New creates a new User
func New() *PasswordController {
	return &PasswordController{
		routePath: controller.PasswordRoute,
	}
}

//Put updates the password of the concerned user
func (password *PasswordController) Put(c *gin.Context) {
	type model struct {
		NewPassword string `json:"new_password,omitempty"`
	}

	passwordModel := &model{}
	err := c.BindJSON(passwordModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}
	controller.Server.UpdatePasswordRequest(c, passwordModel.NewPassword)
}

//Get sends an email for resetting the password of a user to the concerned user's email
func (password *PasswordController) Get(c *gin.Context) {
	email := c.Query("email")
	controller.Server.ForgotPasswordRequest(c, email)
}

// Register will register this controller to the specified router
func (password *PasswordController) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, password, password.routePath)
}
