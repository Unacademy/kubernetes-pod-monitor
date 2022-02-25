package sessions

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clientset *kubernetes.Clientset
)

func newClientset() *kubernetes.Clientset {
	env := viper.GetString("DEPLOY_ENV")
	kubeconfig := ""
	if env == "local" {
		kubeconfig = filepath.Join(
			"config", viper.GetString("CLUSTER_NAME")+".yml",
		)
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create an rest client not targeting specific API version
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}

// InitClientset initializes all common components
func InitClientset() {
	if clientset == nil {
		clientset = newClientset()
	}
}

// GetClientset returns the instance of Clientset that have
// already been initialized through InitClientset
func GetClientset() *kubernetes.Clientset {
	InitClientset()
	return clientset
}
