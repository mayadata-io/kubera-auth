package password

import (
	"net/http"

	log "github.com/golang/glog"

	"github.com/gin-gonic/gin"
	"github.com/mayadata-io/kubera-auth/controller"
)

//Password ...
type Password struct {
	Username    string `json:"username,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

// New creates a new User
func New() *Password {
	return &Password{}
}

//Update updates the password of the concerned user
func (password *Password) Update(c *gin.Context) {
	err := c.BindJSON(password)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	controller.Server.UpdatePasswordRequest(c, password.OldPassword, password.NewPassword)
	return
}

//Reset updates the password of concerned user ggiven that request should be sent by admin
func (password *Password) Reset(c *gin.Context) {
	err := c.BindJSON(password)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	controller.Server.ResetPasswordRequest(c, password.NewPassword, password.Username)
	return
}
