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
}

// NewConfig create to configuration instance
func NewConfig() *Config {

	config := &Config{
		TokenType: "Bearer",
	}

	var err error

	disableGithubAuth := os.Getenv(types.DISABLE_GITHUBAUTH)
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
