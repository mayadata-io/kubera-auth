package oauth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

//GetToken gets the temp code for oauth and exchanges this code with github in order to get auth token
func (config SocialAuthConfig) GetToken(c *gin.Context) (*oauth2.Token, error) {
	code := c.Query("code")
	if code == "" {
		return nil, errors.New("Code not found")
	}
	return config.Exchange(c.Request.Context(), code)
}
