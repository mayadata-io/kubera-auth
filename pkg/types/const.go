package types

import (
	"os"
	"time"
)

// define the type of authorization request
var (
	JWTSecretString            string = "JWT_SECRET"
	DefaultNamespace           string = "kubera"
	DefaultConfigMap           string = os.Getenv("CONFIGMAP_NAME")
	DefaultUserName            string = os.Getenv("ADMIN_USERNAME")
	DefaultUserPassword        string = os.Getenv("ADMIN_PASSWORD")
	DefaultDBServerURL         string = os.Getenv("DB_SERVER")
	PortalURL                  string = os.Getenv("PORTAL_URL")
	DBUser                     string = os.Getenv("DB_USER")
	DBPassword                 string = os.Getenv("DB_PASSWORD")
	DefaultAuthDB              string = "auth"
	DefaultLocalAuthCollection string = "usercredentials"
	PasswordEncryptionCost     int    = 15
	TimeFormat                 string = time.RFC1123Z
	GithubState                string = "github"
)
