package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mayadata-io/kubera-auth/pkg/server"
)

// Server is the global server
var Server *server.Server

//InitializeServer initializes the server
func InitializeServer() {
	// Server represents the server with default config
	Server = server.NewServer(server.NewConfig())
}

type Controller interface {
	Get(c *gin.Context)
	Post(c *gin.Context)
	Put(c *gin.Context)
	Delete(c *gin.Context)
	Patch(c *gin.Context)
	Register(router *gin.RouterGroup)
}

func RegisterController(router *gin.RouterGroup, controller Controller, routePath string) {
	router.GET(routePath, controller.Get)
	router.POST(routePath, controller.Post)
	router.PUT(routePath, controller.Put)
	router.DELETE(routePath, controller.Delete)
	router.PATCH(routePath, controller.Patch)
}

type GenericController struct {
	Controller
}

func (genericController *GenericController) Get(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotFound)
}

func (genericController *GenericController) Post(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotFound)
}

func (genericController *GenericController) Put(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotFound)
}

func (genericController *GenericController) Delete(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotFound)
}

func (genericController *GenericController) Patch(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotFound)
}
