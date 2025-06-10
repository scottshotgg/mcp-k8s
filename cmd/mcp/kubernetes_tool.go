package main

import (
	"flag"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type KubernetesTool struct {
	client          *kubernetes.Clientset
	restClient      *rest.RESTClient
	metricsClient   *versioned.Clientset
	dynClient       *dynamic.DynamicClient
	discoveryClient *discovery.DiscoveryClient
}

func NewKubernetesTool() (*KubernetesTool, error) {
	var (
		kubeconfig *string
		home       = homedir.HomeDir()
	)

	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	// use the current context in kubeconfig
	var config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	config.QPS = 1000
	config.Burst = 2000

	var k KubernetesTool

	// Create the clientset
	k.client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Create a metrics client
	k.metricsClient, err = versioned.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// copy the original
	var configClone = *config
	configClone.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{
		CodecFactory: scheme.Codecs,
	}

	configClone.GroupVersion = &schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}

	// Create a REST client
	k.restClient, err = rest.RESTClientFor(&configClone)
	if err != nil {
		return nil, err
	}

	// Set up dynamic client and discovery client
	k.dynClient, err = dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Create a discovery client
	k.discoveryClient, err = discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	return &k, nil
}
