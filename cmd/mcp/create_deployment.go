package main

import (
	"context"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateDeploymentArgs struct {
	Name      string `json:"name" jsonschema:"required,description=Name of the deployment"`
	Image     string `json:"image" jsonschema:"required,description=Container image that the deployment should run"`
	Namespace string `json:"namespace" jsonschema:"required,description=Namespace that the deployment should be deployed into in"`
}

func (k *KubernetesTool) CreateDeployment(ctx context.Context, args CreateDeploymentArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("create_deployment")

	var (
		opts metav1.CreateOptions

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
}
