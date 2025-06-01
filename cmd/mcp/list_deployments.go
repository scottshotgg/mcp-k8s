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
	var (
		opts metav1.ListOptions

		deployments, err = k.client.
					AppsV1().
					Deployments(args.Namespace).
					List(ctx, opts)
	)

	if err != nil {
		return nil, err
	}

	var ds []string
	for _, d := range deployments.Items {
		ds = append(ds, d.Name)
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: fmt.Sprintf("The Kubernetes deployments in the %s namespace are: %s", args.Namespace, strings.Join(ds, ", ")),
				},
			},
		},
	}, nil

	// return fmt.Sprintf("Your Kubernetes namespace (%s) has been created", name), nil
}
