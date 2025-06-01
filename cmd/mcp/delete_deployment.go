package main

import (
	"context"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeleteDeploymentArgs struct {
	Name      string `json:"name" jsonschema:"required,description=Name of the namespace"`
	Namespace string `json:"namespace" jsonschema:"required,description=Namespace that the deployment should be deployed into in"`
}

func (k *KubernetesTool) DeleteDeployment(ctx context.Context, args DeleteDeploymentArgs) (*mcp_golang.ToolResponse, error) {
	var (
		opts metav1.DeleteOptions

		err = k.client.
			AppsV1().
			Deployments(args.Namespace).
			Delete(ctx, args.Name, opts)
	)

	if err != nil {
		return nil, err
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: fmt.Sprintf("Your Kubernetes deployment (%s) has been deleted", args.Name),
				},
			},
		},
	}, nil

	// return fmt.Sprintf("Your Kubernetes deployment (%s) has been created", name), nil
}
