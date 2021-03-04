package oauth

import (
	"os"

	"github.com/mayadata-io/kubera-auth/pkg/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

//SocialAuthConfig is the type for github and gmail config
type SocialAuthConfig struct {
	oauth2.Config
}

//NewGithubConfig returns the github config
func NewGithubConfig() SocialAuthConfig {
	return SocialAuthConfig{
		Config: oauth2.Config{
			ClientID:     os.Getenv(types.GITHUB_CLIENT_ID),
			ClientSecret: os.Getenv(types.GITHUB_CLIENT_SECRET),
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

func NewGoogleConfig() SocialAuthConfig {
	return SocialAuthConfig{
		Config: oauth2.Config{
			ClientID:     os.Getenv(types.GOOGLE_CLIENT_ID),
			ClientSecret: os.Getenv(types.GOOGLE_CLIENT_SECRET),
			Scopes:       []string{},
			Endpoint:     google.Endpoint,
		},
	}
}
