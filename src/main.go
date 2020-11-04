package main

import (
	"flag"
	"fmt"
	"runtime"

	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/pkg/k8s"
	"github.com/mayadata-io/kubera-auth/router"
)

// Constant to define the port number
const (
	Port = ":3000"
)

func printVersion() {
	log.Infoln(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Infoln(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}

func main() {
	// send logs to stderr so we can use 'kubectl logs'
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("v", "3")

	flag.Parse()
	// Version Info
	printVersion()

	k8s.InitializeClientSet()
	route := router.New()
	log.Fatal(route.Run(Port))
}
