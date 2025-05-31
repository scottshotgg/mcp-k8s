package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

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

type KubernetesTool struct {
	client *kubernetes.Clientset
}

type ToolFn func(ctx context.Context, args map[string]string) (string, error)

func (k *KubernetesTool) CreateNamespace(ctx context.Context, args map[string]string) (string, error) {
	// fmt.Printf("args: %+v\n", args)

	var name, ok = args["name"]
	if !ok {
		return "", errors.New("name was not found")
	}

	var (
		k8sNS = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}

		opts = metav1.CreateOptions{}

		_, err = k.client.
			CoreV1().
			Namespaces().
			Create(ctx, k8sNS, opts)
	)

	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("Your Kubernetes namespace (%s) has been created", name), nil
}

func (k *KubernetesTool) CreateDeployment(ctx context.Context, args map[string]string) (string, error) {
	// fmt.Printf("args: %+v\n", args)

	// TODO: more validation on these in the future
	var name, ok = args["name"]
	if !ok {
		// TODO: make better errors, define these
		return "", errors.New("name was not found")
	}

	image, ok := args["image"]
	if !ok {
		return "", errors.New("image was not found")
	}

	namespace, ok := args["namespace"]
	if !ok {
		return "", errors.New("namespace was not found")
	}

	// spec.template.metadata.labels

	var (
		labels = map[string]string{
			"app": name,
		}

		deploy = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name:   name,
						Labels: labels,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  name,
								Image: image,
								// TODO: save that for later
								// Command: []string{},
							},
						},
					},
				},
			},
		}

		opts = metav1.CreateOptions{}

		_, err = k.client.
			AppsV1().
			Deployments(namespace).
			Create(ctx, deploy, opts)
	)

	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("Your Kubernetes deployment (%s) has been created", name), nil
}

// Expose a deployment
// Create a service
