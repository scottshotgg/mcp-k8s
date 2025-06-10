package main

import (
	"context"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetNamespaceArgs struct {
	Name string `json:"name" jsonschema:"required,description=Name of the namespace"`
}

func (k *KubernetesTool) GetNamespace(ctx context.Context, args GetNamespaceArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("get_namespace")

	var (
		opts metav1.GetOptions

		namespace, err = k.client.
				CoreV1().
				Namespaces().
				Get(ctx, args.Name, opts)
	)

	if err != nil {
		return nil, err
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: namespace.String(),
				},
			},
		},
	}, nil
}
