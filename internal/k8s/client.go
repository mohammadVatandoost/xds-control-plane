package k8s

import (
	"log/slog"
	"os"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateClusterClient() (kubernetes.Interface, error) {
	homeDie, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	kubeConfigPath := homeDie + "/.kube/config"
	var config *rest.Config
	if utils.FileExists(kubeConfigPath) {
		slog.Info("kube config file exist")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, err
		}
	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
