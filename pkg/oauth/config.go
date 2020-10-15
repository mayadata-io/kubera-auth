package oauth

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

//SocialAuthConfig is the type for github and gmail config
type SocialAuthConfig struct {
	oauth2.Config
}

//NewGithubConfig returns the github config
func NewGithubConfig() SocialAuthConfig {
	return SocialAuthConfig{
		Config: oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}
