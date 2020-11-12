package types

import (
	"os"
	"time"
)

// define the type of authorization request
var (
	JWTSecretString            string = "JWT_SECRET"
	GITHUB_CLIENT_ID           string = "GITHUB_CLIENT_ID"
	GITHUB_CLIENT_SECRET       string = "GITHUB_CLIENT_SECRET"
	DISABLE_LOCALAUTH          string = "DISABLE_LOCALAUTH"
	DISABLE_GITHUBAUTH         string = "DISABLE_GITHUBAUTH"
	DefaultNamespace           string = os.Getenv("POD_NAMESPACE")
	DefaultConfigMap           string = os.Getenv("CONFIGMAP_NAME")
	DefaultUserName            string = os.Getenv("ADMIN_USERNAME")
	DefaultName                string = "ADMIN"
	DefaultEmail               string = "admin@kubera.com"
	DefaultUserPassword        string = os.Getenv("ADMIN_PASSWORD")
	DefaultDBServerURL         string = os.Getenv("DB_SERVER")
	PortalURL                  string = os.Getenv("PORTAL_URL")
	DBUser                     string = os.Getenv("DB_USER")
	DBPassword                 string = os.Getenv("DB_PASSWORD")
	Credentials                string = os.Getenv("SECRET_NAME")
	DefaultAuthDB              string = "auth"
	DefaultLocalAuthCollection string = "usercredentials"
	PasswordEncryptionCost     int    = 15
	TimeFormat                 string = time.RFC1123Z
	GithubState                string = "github"
)
