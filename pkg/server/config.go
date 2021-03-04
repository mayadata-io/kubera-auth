package server

import (
	"log"
	"os"
	"strconv"

	"github.com/mayadata-io/kubera-auth/pkg/types"
)

// Config configuration parameters
type Config struct {
	TokenType         string // token type
	DisableLocalAuth  bool
	DisableGithubAuth bool
	DisableGoogleAuth bool
}

// NewConfig create to configuration instance
func NewConfig() *Config {
	config := &Config{
		TokenType: types.BEARER,
	}
	var err error
	// TODO: Think of something to do away of repetitive code
	disableGithubAuth := os.Getenv(types.DISABLE_GITHUBAUTH)
	disableGoogleAuth := os.Getenv(types.DISABLE_GOOGLEAUTH)
	disableLocalAuth := os.Getenv(types.DISABLE_LOCALAUTH)
	if disableLocalAuth == "" {
		// Will be enabled by default
		config.DisableLocalAuth = false
	} else {
		config.DisableLocalAuth, err = strconv.ParseBool(disableLocalAuth)
		if err != nil {
			log.Fatal("Error parsing ", types.DISABLE_LOCALAUTH, err)
		}
	}

	if disableGoogleAuth == "" {
		// Will be enabled by default
		config.DisableGoogleAuth = false
	} else {
		config.DisableGoogleAuth, err = strconv.ParseBool(disableGoogleAuth)
		if err != nil {
			log.Fatal("Error parsing ", types.DISABLE_GOOGLEAUTH, err)
		}
	}

	if disableGithubAuth == "" {
		// Will be disabled by default
		config.DisableGithubAuth = true
	} else {
		config.DisableGithubAuth, err = strconv.ParseBool(disableGithubAuth)
		if err != nil {
			log.Fatal("Error parsing ", types.DISABLE_GITHUBAUTH, err)
		}
	}
	return config
}
