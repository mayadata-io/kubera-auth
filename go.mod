module github.com/mayadata-io/kubera-auth

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/evanphx/json-patch v4.9.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	go.mongodb.org/mongo-driver v1.4.2
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.2.0 // indirect
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6 // indirect
	k8s.io/utils v0.0.0-20201015054608-420da100c033 // indirect
	sigs.k8s.io/structured-merge-diff v0.0.0-20190525122527-15d366b2352e // indirect
)

replace (
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.1.0
	k8s.io/api => k8s.io/api v0.18.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.0
)
