package configuration

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"
	"github.com/imdario/mergo"
	"github.com/mayadata-io/kubera-auth/pkg/errors"
	"github.com/mayadata-io/kubera-auth/pkg/k8s"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Controller struct {
	controller.GenericController
	routePath string
}

// Model is the configModel object for this route
type Model struct {
	// TODO: See the impact of unexported fields in this struct
	GithubClientID     *string `json:"GITHUB_CLIENT_ID,omitempty"`
	GithubClientSecret *string `json:"GITHUB_CLIENT_SECRET,omitempty"`
	EnableGithub       *bool   `json:"ENABLE_GITHUB,omitempty"`
	GoogleClientID     *string `json:"GOOGLE_CLIENT_ID,omitempty"`
	GoogleClientSecret *string `json:"GOOGLE_CLIENT_SECRET,omitempty"`
	EnableGoogle       *bool   `json:"ENABLE_GOOGLE,omitempty"`
}

// New creates a new controller for configs endpoint
func New() *Controller {
	return &Controller{
		routePath: controller.ConfigurationRoute,
	}
}

// Put updates the password of the concerned user
// nolint: cyclop
func (configurationController *Controller) Put(c *gin.Context) {
	// 1. Verify authentication
	tokenString, err := getTokenFromHeader(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	// 2. Verify authorization
	jwtUserCredentials, err := controller.Server.GetUserFromToken(tokenString)
	if err != nil || jwtUserCredentials == nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	} else if jwtUserCredentials.Role != models.RoleAdmin {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	// 3. Validate request body
	configModel := &Model{}
	err = c.BindJSON(configModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}
	// 4. Persist the configuration for further usage
	configurationController.updateToK8s(c, configModel)
}

// updateToK8s updates the value of Oauth secrets in secrets & configmaps
// TODO: Need to lock this operation, so only one request can update it at a time
// TODO(vharsh): Use pflags to configure all of these
func (configurationController *Controller) updateToK8s(c *gin.Context, requestModel *Model) {
	// 1. Get current configuration from configmap
	// 2. See if any update is required
	// 3. Respond back with updated *Model
	cm, err := k8s.ClientSet.CoreV1().ConfigMaps(types.DefaultNamespace).
		Get(c.Request.Context(), types.DefaultConfigMap, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, gin.H{
			"message": "Unable to persist config change",
		})
	}
	githubDisable, _ := strconv.ParseBool(cm.Data[types.DISABLE_GITHUBAUTH])
	githubEnable := !githubDisable
	googDisable, _ := strconv.ParseBool(cm.Data[types.DISABLE_GOOGLEAUTH])
	googEnable := !googDisable
	githubClientID := cm.Data[types.GITHUB_CLIENT_ID]
	githubClientSecret := cm.Data[types.GITHUB_CLIENT_SECRET]
	googClientID := cm.Data[types.GOOGLE_CLIENT_ID]
	googClientSecret := cm.Data[types.GOOGLE_CLIENT_SECRET]
	cfgMapModel := Model{
		GithubClientID:     &githubClientID,
		GithubClientSecret: &githubClientSecret,
		EnableGithub:       &githubEnable,
		GoogleClientID:     &googClientID,
		GoogleClientSecret: &googClientSecret,
		EnableGoogle:       &googEnable,
	}
	// update the configmap model with data from the request-model
	if err := mergo.Merge(&cfgMapModel, &requestModel); err != nil {
		log.Error("Error merging to cfgMapModelg", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update configs",
		})
	}

	// merge the request data
	if err := mergo.Merge(&cm.Data, &cfgMapModel); err != nil {
		log.Error("Error merging to cfgMap", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update config",
		})
	}
	_, err = k8s.ClientSet.CoreV1().ConfigMaps(types.DefaultNamespace).Update(c.Request.Context(), cm, metav1.UpdateOptions{})
	if err != nil {
		log.Errorln("Error updating configmap ", err)
		c.String(http.StatusInternalServerError, "Error enabling OAuth config")
	}
	// TODO: Add config for localAuth
	controller.Server.Config.DisableGoogleAuth = !*cfgMapModel.EnableGoogle
	controller.Server.Config.DisableGithubAuth = !*cfgMapModel.EnableGithub
	c.JSON(http.StatusOK, cfgMapModel)
	// Set a nice success response with the Model
}

func (configurationController *Controller) Get(c *gin.Context) {
	authData := map[string]interface{}{
		types.DISABLE_GITHUBAUTH: controller.Server.Config.DisableGithubAuth,
		types.DISABLE_LOCALAUTH:  controller.Server.Config.DisableLocalAuth,
		types.DISABLE_GOOGLEAUTH: controller.Server.Config.DisableGithubAuth,
	}

	tokenString, err := getTokenFromHeader(c.Request)
	if err != nil {
		log.Errorln("Invalid Token: Unable to parse jwt")
	}

	jwtUserCredentials, err := controller.Server.GetUserFromToken(tokenString)
	if err == nil && jwtUserCredentials.Role == models.RoleAdmin {
		authData[types.GITHUB_CLIENT_ID] = controller.Server.GithubConfig.ClientID
		authData[types.GITHUB_CLIENT_SECRET] = controller.Server.GithubConfig.ClientSecret
		authData[types.GOOGLE_CLIENT_ID] = controller.Server.GoogleConfig.ClientID
		authData[types.GOOGLE_CLIENT_SECRET] = controller.Server.GoogleConfig.ClientSecret
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
