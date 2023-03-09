package xds

import (
	// "encoding/base64"
	// "io/ioutil"
	// "net/http"
	// "strings"
	// "time"

	"os"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/utils"
	"github.com/sirupsen/logrus"

	// yaml "gopkg.in/yaml.v2"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// type BootstrapConfig struct {
// 	Name   string
// 	Server string
// 	CA     string // CA is expected to be base64 encoded PEM file
// }

// var bootstrapConfigs []BootstrapConfig

// func CreateBootstrapClients() ([]kubernetes.Interface, error) {
// 	var bootstrapClients []kubernetes.Interface

// 	bootstrapYaml, err := ioutil.ReadFile("./bootstrap.yaml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = yaml.Unmarshal(bootstrapYaml, &bootstrapConfigs)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, cluster := range bootstrapConfigs {
// 		// CA is base64 encoded, so we decode
// 		caReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(cluster.CA))
// 		caBytes, _ := ioutil.ReadAll(caReader)

// 		var restConfig *rest.Config
// 		restConfig = &rest.Config{
// 			Host: cluster.Server,
// 			TLSClientConfig: rest.TLSClientConfig{
// 				CAData: caBytes,
// 			},
// 		}

// 		previousWrappedTransport := restConfig.WrapTransport
// 		restConfig.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
// 			if previousWrappedTransport != nil {
// 				rt = previousWrappedTransport(rt)
// 			}
// 			return &TokenRoundtripper{
// 				TokenProvider: &NullTokenProvider{},
// 				RoundTripper:  rt,
// 			}
// 		}

// 		restConfig.Timeout = time.Second * 5

// 		client, err := kubernetes.NewForConfig(restConfig)
// 		if err != nil {
// 			return nil, err
// 		}
// 		bootstrapClients = append(bootstrapClients, client)
// 	}
// 	return bootstrapClients, nil
// }

func CreateClusterClient() (kubernetes.Interface, error) {
	homeDie, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	kubeConfigPath := homeDie + "/.kube/config"
	var config *rest.Config
	if utils.FileExists(kubeConfigPath) {
		logrus.Info("kube config file exist")
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
	// dynamic.NewForConfig()
	return clientset, nil
}
