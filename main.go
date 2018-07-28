package main

import (
	"github.com/flokkr/kubernetes-launcher/src"
	"flag"
)

func main() {
	var destination string
	var namespace string
	var configmap string
	flag.StringVar(&destination, "destination", "/tmp", "Destination path")
	flag.StringVar(&namespace, "namespace", "default", "Namespace name")
	flag.StringVar(&configmap, "configmap", "config", "Configmap name")
	flag.Parse()
	kuberneteslauncher.ListOnConfigmap(destination, namespace, configmap, flag.Args())


}
