package user

import (
	"encoding/json"
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
	type model struct {
		UnverifiedEmail string `json:"unverified_email"`
		Company         string `json:"company"`
		CompanyRole     string `json:"company_role"`
		Name            string `json:"name"`
	}

	requestModel := &model{}
	err := c.BindJSON(requestModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	userCredentialsModel := &models.UserCredentials{}

	userBytes, err := json.Marshal(requestModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}
	err = json.Unmarshal(userBytes, userCredentialsModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	controller.Server.UpdateUserDetailsRequest(c, userCredentialsModel)
}

//Patch updates the password of concerned user given that request should be sent by admin
func (user *UserController) Patch(c *gin.Context) {
	type model struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	requestModel := &model{}
	err := c.BindJSON(requestModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}
	controller.Server.ResetPasswordRequest(c, requestModel.Password, requestModel.UserName)
}

//Post creates a user, request should be sent by admin
func (user *UserController) Post(c *gin.Context) {
	type model struct {
		UnverifiedEmail string `json:"unverified_email"`
		Company         string `json:"company"`
		CompanyRole     string `json:"company_role"`
		Name            string `json:"name"`
		UserName        string `json:"username"`
		Password        string `json:"password"`
	}

	requestModel := &model{}
	err := c.BindJSON(requestModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	userCredentialsModel := &models.UserCredentials{}

	userBytes, err := json.Marshal(requestModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}
	err = json.Unmarshal(userBytes, userCredentialsModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	controller.Server.CreateRequest(c, userCredentialsModel)
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

// Register will register this controller to the specified router
func (user *UserController) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, user, user.routePath)
	router.GET(user.routePath+"/uid/:userID", user.GetByUID)
	router.GET(user.routePath+"/username/:username", user.GetByUsername)
}
