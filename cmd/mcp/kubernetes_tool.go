package main

import (
	"flag"
	"log"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type KubernetesTool struct {
	client        *kubernetes.Clientset
	restClient    *rest.RESTClient
	metricsClient *versioned.Clientset
}

func NewKubernetesTool() *KubernetesTool {
	var kubeconfig *string
	home := homedir.HomeDir()
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	configClone := *config // copy the original
	configClone.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{
		CodecFactory: scheme.Codecs,
	}

	configClone.GroupVersion = &schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}

	// 2025/06/02 16:23:49 Error creating REST client: GroupVersion is required when initializing a RESTClient
	// Create a REST client
	restClient, err := rest.RESTClientFor(&configClone)
	if err != nil {
		log.Fatalf("Error creating REST client: %v", err)
	}

	// Create a metrics client
	metricsClient, err := versioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create metrics client: %v", err)
	}

	return &KubernetesTool{
		client:        clientset,
		restClient:    restClient,
		metricsClient: metricsClient,
	}
}

// Expose a deployment
// Create a service
