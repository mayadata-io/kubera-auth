package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/pkg/models"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
)

// UserController is a type to be accepted as input
type UserController struct {
	controller.GenericController
	routePath string
}

// New creates a new User
func New() *UserController {
	return &UserController{
		routePath: controller.UserRoute,
	}
}

// Put updates a user details
func (user *UserController) Put(c *gin.Context) {
	userModel := &models.UserCredentials{}
	err := c.BindJSON(userModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	controller.Server.UpdateUserDetailsRequest(c, userModel)
}

//Patch updates the password of concerned user given that request should be sent by admin
func (user *UserController) Patch(c *gin.Context) {
	userModel := &models.UserCredentials{}
	err := c.BindJSON(userModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}
	controller.Server.ResetPasswordRequest(c, userModel.Password, userModel.UserName)
}

//Post creates a user, request should be sent by admin
func (user *UserController) Post(c *gin.Context) {
	userModel := &models.UserCredentials{}
	err := c.BindJSON(userModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	controller.Server.CreateRequest(c, userModel)
}

// Get will respond with a particular user or all users
func (user *UserController) Get(c *gin.Context) {
	// Get all users
	controller.Server.GetUsersRequest(c)
}

// GetByUID will respond with a particular user or all users
func (user *UserController) GetByUID(c *gin.Context) {
	userID := c.Param("userID")
	controller.Server.GetUserByUID(c, userID)
}

// GetByUsername will respond with a particular user or all users
func (user *UserController) GetByUsername(c *gin.Context) {
	userID := c.Param("username")
	controller.Server.GetUserByUserName(c, userID)
}

// Register will rsgister this controller to the specified router
func (user *UserController) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, user, user.routePath)
	router.GET(user.routePath+"/uid/:userID", user.GetByUID)
	router.GET(user.routePath+"/username/:username", user.GetByUsername)
}
