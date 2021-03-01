package types

import (
	"time"
)

// define the admin & login details(true, false, nil)
var (
	DefaultName  = "ADMIN"
	DefaultEmail = "admin@kubera.com"
	TrueValue    = true
	FalseValue   = false
)

// Authentication related constants
const (
	DefaultAuthDB                      string        = "auth"
	DefaultLocalAuthCollection         string        = "usercredentials"
	GithubState                        string        = "github"
	JWTUserCredentialsKey              string        = "userCredentials"
	TemplatePath                       string        = "./templates"
	AuthHeaderKey                      string        = "Authorization"
	AuthHeaderPrefix                   string        = "Bearer "
	TimeFormat                         string        = time.RFC1123Z
	VerificationLinkExpirationTimeUnit time.Duration = 10
	PasswordEncryptionCost             int           = 15
)
