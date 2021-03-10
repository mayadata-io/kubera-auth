package signup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/pkg/models"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
)

//SignupController is the type the request in which the request will be parsed
type SignupController struct {
	controller.GenericController
	routePath string
}

// Model defines the json struct in which the request will be parsed
type Model struct {
	UnverifiedEmail string `json:"unverified_email"`
	Password        string `json:"password"`
	Name            string `json:"name"`
}

// New creates a new LoginUser
func New() *SignupController {
	return &SignupController{
		routePath: controller.SignupRoute,
	}
}

//Post registers a user in database taking minimal details needed
func (signupController *SignupController) Post(c *gin.Context) {
	signuplModel := &Model{}
	err := c.BindJSON(signuplModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Unable to parse JSON",
		})
		return
	}

	newUser := &models.UserCredentials{
		UserName: &signuplModel.UnverifiedEmail,
		Name:     &signuplModel.Name,
		Password: &signuplModel.Password,
	}

	// First create the user then immediately send a verification link to his email
	controller.Server.SelfSignupUser(c, newUser)
}

// Register will rsgister this controller to the specified router
func (signupController *SignupController) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, signupController, signupController.routePath)
}
