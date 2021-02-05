package types

import (
	"os"
)

var (
	DefaultNamespace    string = os.Getenv("POD_NAMESPACE")
	DefaultConfigMap    string = os.Getenv("CONFIGMAP_NAME")
	DefaultUserName     string = os.Getenv("ADMIN_USERNAME")
	DefaultUserPassword string = os.Getenv("ADMIN_PASSWORD")
	DefaultDBServerURL  string = os.Getenv("DB_SERVER")
	PortalURL           string = os.Getenv("PORTAL_URL")
	DBUser              string = os.Getenv("DB_USER")
	DBPassword          string = os.Getenv("DB_PASSWORD")
	Credentials         string = os.Getenv("SECRET_NAME")
	EmailUsername       string = os.Getenv("EMAIL_USERNAME")
	EmailPassword       string = os.Getenv("EMAIL_PASSWORD")
)
