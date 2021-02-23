package types

import (
	"time"
)

// define the type of authorization request
var (
	DefaultName                                      = "ADMIN"
	DefaultEmail                                     = "admin@kubera.com"
	DefaultAuthDB                                    = "auth"
	DefaultLocalAuthCollection                       = "usercredentials"
	GithubState                                      = "github"
	UserInfoKey                                      = "userInfo"
	TemplatePath                                     = "./templates"
	AuthHeaderKey                                    = "Authorization"
	AuthHeaderPrefix                                 = "Bearer "
	TimeFormat                                       = time.RFC1123Z
	TrueValue                                        = true
	FalseValue                                       = false
	VerificationLinkExpirationTimeUnit time.Duration = 10
	PasswordEncryptionCost                           = 15
)
