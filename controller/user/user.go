package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/controller"
	"github.com/mayadata-io/kubera-auth/pkg/models"
)

// User is a type to be accepted as input
type User models.UserCredentials

// New creates a new User
func New() *User {
	return &User{}
}

// Logout lets a user login into the kubera-core
func (user *User) Logout(c *gin.Context) {
	controller.Server.LogoutRequest(c)
	return
}

// UpdateUserDetails updates a user details
func (user *User) UpdateUserDetails(c *gin.Context) {
	err := c.BindJSON(user)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	userModel := models.UserCredentials(*user)
	controller.Server.UpdateUserDetailsRequest(c, &userModel)
	return
}

//Create creates a user, request should be sent by admin
func (user *User) Create(c *gin.Context) {
	err := c.BindJSON(user)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	userModel := models.UserCredentials(*user)
	userModel.Kind = models.LocalAuth
	controller.Server.CreateRequest(c, &userModel)
	return
}

//GetAllUsers responds with a list of users
func (user *User) GetAllUsers(c *gin.Context) {
	controller.Server.GetUsersRequest(c)
	return
}
