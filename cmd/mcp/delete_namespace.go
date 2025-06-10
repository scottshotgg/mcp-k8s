package main

import (
	"context"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeleteNamespaceArgs struct {
	Name string `json:"name" jsonschema:"required,description=Name of the namespace"`
}

func (k *KubernetesTool) DeleteNamespace(ctx context.Context, args DeleteNamespaceArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("delete_namespace")

	var (
		opts metav1.DeleteOptions

		err = k.client.
			CoreV1().
			Namespaces().
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
					Text: fmt.Sprintf("Your Kubernetes namespace (%s) has been deleted", args.Name),
				},
			},
		},
	}, nil
}
