package configuration

import (
	"net/http"
	"strconv"

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
	ClientID     *string `json:"GITHUB_CLIENT_ID,omitempty"`
	ClientSecret *string `json:"GITHUB_CLIENT_SECRET,omitempty"`
	EnableGithub *bool   `json:"ENABLE_GITHUB,omitempty"`
}

// New creates a new User
func New() *Controller {
	return &Controller{
		routePath: controller.ConfigurationRoute,
		model:     &Model{},
	}
}

//Put updates the password of the concerned user
func (configurationController *Controller) Put(c *gin.Context) {

	userInfo, err := controller.Server.GetUserFromToken(c.Request)
	if err != nil || userInfo == nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	} else if userInfo.GetRole() != models.RoleAdmin {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	err = c.BindJSON(configurationController.model)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	if configurationController.model.ClientID != nil && configurationController.model.EnableGithub != nil {
		credentialType := c.Query("type")
		var data map[string]string

		switch credentialType {
		case string(models.GithubAuth):
			{
				data = map[string]string{
					types.GITHUB_CLIENT_ID:     *configurationController.model.ClientID,
					types.GITHUB_CLIENT_SECRET: *configurationController.model.ClientSecret,
				}
			}
		default:
			{
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		secret, err := k8s.ClientSet.CoreV1().Secrets(types.DefaultNamespace).Get(c.Request.Context(), types.Credentials, metav1.GetOptions{})
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

		controller.Server.GithubConfig.ClientID = *configurationController.model.ClientID
		controller.Server.GithubConfig.ClientSecret = *configurationController.model.ClientSecret

	} else if configurationController.model.EnableGithub != nil {
		cm, err := k8s.ClientSet.CoreV1().ConfigMaps(types.DefaultNamespace).Get(c.Request.Context(), types.DefaultConfigMap, metav1.GetOptions{})
		if err != nil {
			c.String(http.StatusInternalServerError, "Error enabling github login : %v", err)
			log.Errorln("Error getting configmap ", err)
			return
		}

		cm.Data[types.DISABLE_GITHUBAUTH] = strconv.FormatBool(!*configurationController.model.EnableGithub)
		cm, err = k8s.ClientSet.CoreV1().ConfigMaps(types.DefaultNamespace).Update(c.Request.Context(), cm, metav1.UpdateOptions{})
		if err != nil {
			c.String(http.StatusInternalServerError, "Error enabling github login : %v", err)
			log.Errorln("Error updating configmap ", err)
			return
		}

		controller.Server.Config.DisableGithubAuth = !*configurationController.model.EnableGithub
	}

	c.String(http.StatusOK, "Configured sucessfully")
}

func (configurationController *Controller) Get(c *gin.Context) {

	auth_data := map[string]bool{
		types.DISABLE_GITHUBAUTH: controller.Server.Config.DisableGithubAuth,
		types.DISABLE_LOCALAUTH:  controller.Server.Config.DisableLocalAuth,
	}

	c.JSON(http.StatusOK, auth_data)
}

// Register will rsgister this controller to the specified router
func (configurationController *Controller) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, configurationController, configurationController.routePath)
}
