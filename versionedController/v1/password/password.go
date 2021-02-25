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

//Model ...
type Model struct {
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

// New creates a new User
func New() *PasswordController {
	return &PasswordController{
		routePath: controller.PasswordRoute,
	}
}

//Put updates the password of the concerned user
func (password *PasswordController) Put(c *gin.Context) {
	passwordModel := &Model{}
	err := c.BindJSON(passwordModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	controller.Server.UpdatePasswordRequest(c, passwordModel.OldPassword, passwordModel.NewPassword)
	return
}

// Register will rsgister this controller to the specified router
func (password *PasswordController) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, password, password.routePath)
}
