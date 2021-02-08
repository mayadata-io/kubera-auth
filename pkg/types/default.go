package types

import (
	"time"
)

// define the type of authorization request
var (
	DefaultName                string = "ADMIN"
	DefaultEmail               string = "admin@kubera.com"
	DefaultAuthDB              string = "auth"
	DefaultLocalAuthCollection string = "usercredentials"
	PasswordEncryptionCost     int    = 15
	TimeFormat                 string = time.RFC1123Z
	GithubState                string = "github"
	UserInfoKey                string = "userInfo"
	TemplatePath               string = "./templates"
	TrueValue                  bool   = true
	FalseValue                 bool   = false
)
