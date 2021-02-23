package types

import (
	"os"
)

var (
	DefaultNamespace    = os.Getenv("POD_NAMESPACE")
	DefaultConfigMap    = os.Getenv("CONFIGMAP_NAME")
	DefaultUserName     = os.Getenv("ADMIN_USERNAME")
	DefaultUserPassword = os.Getenv("ADMIN_PASSWORD")
	DefaultDBServerURL  = os.Getenv("DB_SERVER")
	PortalURL           = os.Getenv("PORTAL_URL")
	DBUser              = os.Getenv("DB_USER")
	DBPassword          = os.Getenv("DB_PASSWORD")
	Credentials         = os.Getenv("SECRET_NAME")
	EmailUsername       = os.Getenv("EMAIL_USERNAME")
	EmailPassword       = os.Getenv("EMAIL_PASSWORD")
)
