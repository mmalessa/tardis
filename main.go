package main

import (
	"tardis/pkg/app"
	"tardis/pkg/k8portforwarder"

	"github.com/gookit/config/v2"
	log "github.com/sirupsen/logrus"
)

// It can be overwriten by "go build"!
var env = "prod"
var forwarders map[string]*k8portforwarder.K8PortForwarder

func main() {

	if err := app.InitLogs(env); err != nil {
		log.Fatal(err)
	}

	log.Infof("Start application (env=\"%s\")\n", env)

	if err := app.InitConfig(env); err != nil {
		log.Fatal(err)
	}

	var forwards map[string]k8portforwarder.K8PortForwarderConfig
	config.BindStruct("forwards", &forwards)

	forwarders = make(map[string]*k8portforwarder.K8PortForwarder)
	for key, forwardConfig := range forwards {
		var err error
		forwarders[key], err = k8portforwarder.NewK8PortForwarder(forwardConfig)
		if err != nil {
			panic(err)
		}
	}

	for _, forwarder := range forwarders {
		forwarder.Test()
	}

	return

	// kpf, err := k8portforwarder.NewK8PortForwarder(kubeContext, namespace, localPort, remotePort, remoteType, remoteName)
	// if err != nil {
	// 	log.Fatal("Error setting up port forwarder: ", err)
	// }

	// ctx := context.TODO()
	// if err := kpf.Start(ctx); err != nil {
	// 	fmt.Println(err)
	// 	panic("")
	// }

	// log.Printf("Started tunnel on %d\n", kpf.LocalPort)
	// time.Sleep(60 * time.Second)
	// log.Printf("Time is up!\n")

}
