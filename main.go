package main

import (
	"github.com/flokkr/kubernetes-launcher/src"
	"flag"
	"log"
)

var version = "dev"
func main() {
	var destination, namespace, fieldSelector, labelSelector string
	flag.StringVar(&destination, "destination", "/tmp", "Destination path")
	flag.StringVar(&namespace, "namespace", "default", "Namespace name")
	flag.StringVar(&fieldSelector, "fields", "", "Field selector for the configmap(s)")
	flag.StringVar(&labelSelector, "labels", "", "Label selector for the configmap(s)")
	flag.Parse()
	log.Println("kubernetes-launcher " + version)
	kuberneteslauncher.ListOnConfigmap(destination, namespace, fieldSelector, labelSelector, flag.Args())


}
