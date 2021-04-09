package signup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/pkg/models"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
)

//SignupController is the extension to GenericController which contains the path of this endpoint too.
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
	signupModel := &Model{}
	err := c.BindJSON(signupModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Unable to parse JSON",
		})
		return
	}

	newUser := &models.UserCredentials{
		UserName:        signupModel.UnverifiedEmail,
		Name:            signupModel.Name,
		Password:        signupModel.Password,
		UnverifiedEmail: signupModel.UnverifiedEmail,
	}

	// First create the user then immediately send a verification link to his email
	controller.Server.SelfSignupUser(c, newUser)
}

// Register will register this controller to the specified router
func (signupController *SignupController) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, signupController, signupController.routePath)
}
