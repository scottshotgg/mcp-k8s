package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	mcp_golang "github.com/metoro-io/mcp-golang"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubernetesTool struct {
	client *kubernetes.Clientset
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

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &KubernetesTool{
		client: clientset,
	}
}

type CreateNamespaceArgs struct {
	Name string `json:"name" jsonschema:"required,description=Name of the namespace"`
}

func (k *KubernetesTool) CreateNamespace(ctx context.Context, args CreateNamespaceArgs) (*mcp_golang.ToolResponse, error) {
	var (
		opts  = metav1.CreateOptions{}
		k8sNS = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: args.Name,
			},
		}

		_, err = k.client.
			CoreV1().
			Namespaces().
			Create(ctx, k8sNS, opts)
	)

	if err != nil {
		return nil, err
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: fmt.Sprintf("Your Kubernetes namespace (%s) has been created", args.Name),
				},
			},
		},
	}, nil

	// return fmt.Sprintf("Your Kubernetes namespace (%s) has been created", name), nil
}

type CreateDeploymentArgs struct {
	Name      string `json:"name" jsonschema:"required,description=Name of the namespace"`
	Image     string `json:"image" jsonschema:"required,description=Image to run"`
	Namespace string `json:"namespace" jsonschema:"required,description=Namespace that the deployment should be deployed into in"`
}

func (k *KubernetesTool) CreateDeployment(ctx context.Context, args CreateDeploymentArgs) (*mcp_golang.ToolResponse, error) {
	var (
		opts   = metav1.CreateOptions{}
		labels = map[string]string{
			"app": args.Name,
		}

		deploy = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      args.Name,
				Namespace: args.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name:   args.Name,
						Labels: labels,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  args.Name,
								Image: args.Image,
								// TODO: save that for later
								// Command: []string{},
							},
						},
					},
				},
			},
		}

		_, err = k.client.
			AppsV1().
			Deployments(args.Namespace).
			Create(ctx, deploy, opts)
	)

	if err != nil {
		return nil, err
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: fmt.Sprintf("Your Kubernetes deployment (%s) has been created", args.Name),
				},
			},
		},
	}, nil

	// return fmt.Sprintf("Your Kubernetes deployment (%s) has been created", name), nil
}

// Expose a deployment
// Create a service
