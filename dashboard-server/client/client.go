package client

import (
	"dashboard-server/config"
	"log"

	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

var DefaultClient v1.CoreV1Interface
var DefaultConfig rest.Config

func init() {
	projectConfig := config.GetConfig()
	DefaultConfig = rest.Config{
		Host:            projectConfig.K8s.Host,
		BearerToken:     projectConfig.K8s.Token,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}
	cli, err := kubernetes.NewForConfig(&DefaultConfig)
	if err != nil {
		log.Panic("Error: ", err)
	}
	DefaultClient = cli.CoreV1()
}
