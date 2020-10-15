package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/mayadata-io/kubera-auth/pkg/server"
)

// Server represents the server with default config
var Server = server.NewServer(server.NewConfig())

// UserController will do user operations
type UserController interface {
	Logout(c *gin.Context)
	UpdateUserDetails(c *gin.Context)
	GetAllUsers(c *gin.Context)
	Create(c *gin.Context)
}

// PasswordController will do password operations
type PasswordController interface {
	Update(c *gin.Context)
	Reset(c *gin.Context)
}

//LoginController will do llogin operations
type LoginController interface {
	Login(c *gin.Context)
	SocialLogin(c *gin.Context)
	CallbackRequest(c *gin.Context)
}
