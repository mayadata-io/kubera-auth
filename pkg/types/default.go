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
	DefaultLocalAuthCollection                       = "usercredentials"
	GithubState                                      = "github"
	GoogleState                                      = "google"
	JWTUserCredentialsKey                            = "userCredentials"
	TemplatePath                                     = "./templates"
	KuberaPortalImagePath                            = "/kuberaPortal.png"
	MayadataLogoImagePath                            = "/mayadata-logo.png"
	BackgroundEmailImagePath                         = "/bg-kubera-email.png"
	VerificationEmailTemplatePath                    = "/verificationEmailTemplate.html"
	ResetPasswordEmailTemplatePath                   = "/resetPasswordEmailTemplate.html"
	AuthHeaderKey                                    = "Authorization"
	AuthHeaderPrefix                                 = "Bearer "
	TimeFormat                                       = time.RFC1123Z
	VerificationLinkExpirationTimeUnit time.Duration = 10
	PasswordEncryptionCost             int           = 15
)
