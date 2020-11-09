package configuration

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mayadata-io/kubera-auth/pkg/k8s"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
)

type Controller struct {
	controller.GenericController
	routePath string
	model     *Model
}

//Model ...
type Model struct {
	ClientID     string `json:"GITHUB_CLIENT_ID"`
	ClientSecret string `json:"GITHUB_CLIENT_SECRET"`
}

// New creates a new User
func New() *Controller {
	return &Controller{
		routePath: controller.CredentialsRoute,
		model:     &Model{},
	}
}

//Put updates the password of the concerned user
func (configurationController *Controller) Put(c *gin.Context) {
	err := c.BindJSON(configurationController.model)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	credentialType := c.Query("type")
	fmt.Println("credentials", credentialType)
	var data map[string]string

	switch credentialType {
	case string(models.GithubAuth):
		{
			data = map[string]string{
				types.GITHUB_CLIENT_ID:     configurationController.model.ClientID,
				types.GITHUB_CLIENT_SECRET: configurationController.model.ClientSecret,
			}
		}
	default:
		{
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	secret, err := k8s.ClientSet.CoreV1().Secrets(types.DefaultNamespace).Get(c.Request.Context(), types.Crednetials, metav1.GetOptions{})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error storing credentials : %v", err)
		log.Errorln("Error getting secret ", err)
		return
	}

	secret.StringData = data
	secret, err = k8s.ClientSet.CoreV1().Secrets(types.DefaultNamespace).Update(c.Request.Context(), secret, metav1.UpdateOptions{})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error storing credentials : %v", err)
		log.Errorln("Error updating secret ", err)
		return
	}

	c.String(http.StatusOK, "Credentials Saved sucessfully")
}

// Register will rsgister this controller to the specified router
func (configurationController *Controller) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, configurationController, configurationController.routePath)
}
