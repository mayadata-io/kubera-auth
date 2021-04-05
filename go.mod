module github.com/mayadata-io/kubera-auth

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/imdario/mergo v0.3.12
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/oauth2 v0.0.0-20210323180902-22b0adad7558
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v0.20.5
)
