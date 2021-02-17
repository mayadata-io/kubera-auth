package types

import (
	"time"
)

// define the type of authorization request
var (
	DefaultName                        string        = "ADMIN"
	DefaultEmail                       string        = "admin@kubera.com"
	DefaultAuthDB                      string        = "auth"
	DefaultLocalAuthCollection         string        = "usercredentials"
	GithubState                        string        = "github"
	UserInfoKey                        string        = "userInfo"
	TemplatePath                       string        = "./templates"
	AuthHeaderKey                      string        = "Authorization"
	AuthHeaderPrefix                   string        = "Bearer "
	TimeFormat                         string        = time.RFC1123Z
	TrueValue                          bool          = true
	FalseValue                         bool          = false
	VerificationLinkExpirationTimeUnit time.Duration = 10
	PasswordEncryptionCost             int           = 15
)
