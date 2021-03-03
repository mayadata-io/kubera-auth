package configuration

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mayadata-io/kubera-auth/pkg/errors"
	"github.com/mayadata-io/kubera-auth/pkg/k8s"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
)

type Controller struct {
	controller.GenericController
	routePath string
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
	}
}

// Put updates the password of the concerned user
// nolint: cyclop
func (configurationController *Controller) Put(c *gin.Context) {
	configurationModel := &Model{}
	tokenString, err := getTokenFromHeader(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	userInfo, err := controller.Server.GetUserFromToken(tokenString)
	if err != nil || userInfo == nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	} else if userInfo.GetRole() != models.RoleAdmin {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	err = c.BindJSON(configurationModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	if configurationModel.ClientID != nil && configurationModel.EnableGithub != nil {
		credentialType := c.Query("type")
		var data map[string]string

		switch credentialType {
		case string(models.GithubAuth):
			{
				data = map[string]string{
					types.GITHUB_CLIENT_ID:     *configurationModel.ClientID,
					types.GITHUB_CLIENT_SECRET: *configurationModel.ClientSecret,
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
		_, err = k8s.ClientSet.CoreV1().Secrets(types.DefaultNamespace).Update(c.Request.Context(), secret, metav1.UpdateOptions{})
		if err != nil {
			c.String(http.StatusInternalServerError, "Error storing credentials : %v", err)
			log.Errorln("Error updating secret ", err)
			return
		}

		controller.Server.GithubConfig.ClientID = *configurationModel.ClientID
		controller.Server.GithubConfig.ClientSecret = *configurationModel.ClientSecret
	}
	if configurationModel.EnableGithub != nil {
		cm, err := k8s.ClientSet.CoreV1().ConfigMaps(types.DefaultNamespace).Get(c.Request.Context(), types.DefaultConfigMap, metav1.GetOptions{})
		if err != nil {
			c.String(http.StatusInternalServerError, "Error enabling github login : %v", err)
			log.Errorln("Error getting configmap ", err)
			return
		}

		cm.Data[types.DISABLE_GITHUBAUTH] = strconv.FormatBool(!*configurationModel.EnableGithub)
		_, err = k8s.ClientSet.CoreV1().ConfigMaps(types.DefaultNamespace).Update(c.Request.Context(), cm, metav1.UpdateOptions{})
		if err != nil {
			c.String(http.StatusInternalServerError, "Error enabling github login : %v", err)
			log.Errorln("Error updating configmap ", err)
			return
		}

		controller.Server.Config.DisableGithubAuth = !*configurationModel.EnableGithub
	}

	c.String(http.StatusOK, "Configured successfully")
}

func (configurationController *Controller) Get(c *gin.Context) {
	authData := map[string]interface{}{
		types.DISABLE_GITHUBAUTH: controller.Server.Config.DisableGithubAuth,
		types.DISABLE_LOCALAUTH:  controller.Server.Config.DisableLocalAuth,
	}

	tokenString, err := getTokenFromHeader(c.Request)
	if err != nil {
		log.Errorln("Invalid Token: Unable to parse jwt")
	}

	userInfo, err := controller.Server.GetUserFromToken(tokenString)
	if err == nil && userInfo.GetRole() == models.RoleAdmin {
		authData[types.GITHUB_CLIENT_ID] = controller.Server.GithubConfig.ClientID
		authData[types.GITHUB_CLIENT_SECRET] = controller.Server.GithubConfig.ClientSecret
	}

	c.JSON(http.StatusOK, authData)
}

func getTokenFromHeader(r *http.Request) (string, error) {
	auth := r.Header.Get(types.AuthHeaderKey)
	prefix := types.AuthHeaderPrefix
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}

	if token == "" {
		return token, errors.ErrInvalidAccessToken
	}

	return token, nil
}

// Register will rsgister this controller to the specified router
func (configurationController *Controller) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, configurationController, configurationController.routePath)
}
