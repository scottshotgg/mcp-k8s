package main

import (
	"context"
	"fmt"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListDeploymentsArgs struct {
	Namespace string `json:"namespace" jsonschema:"required,description=Namespace that the deployment should be deployed into in"`
}

func (k *KubernetesTool) ListDeployments(ctx context.Context, args ListDeploymentsArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("list_deployments")

	var (
		opts metav1.ListOptions

		deploymentsList, err = k.client.
					AppsV1().
					Deployments(args.Namespace).
					List(ctx, opts)
	)

	if err != nil {
		return nil, err
	}

	var deployments []string
	for _, deployment := range deploymentsList.Items {
		deployments = append(deployments, deployment.Name)
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: fmt.Sprintf("The Kubernetes deployments in the %s namespace are: %s", args.Namespace, strings.Join(deployments, ", ")),
				},
			},
		},
	}, nil
}
