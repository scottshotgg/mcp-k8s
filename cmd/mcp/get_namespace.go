package main

import (
	"context"

	mcp_golang "github.com/metoro-io/mcp-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetNamespaceArgs struct {
	Name string `json:"name" jsonschema:"required,description=Name of the namespace"`
}

func (k *KubernetesTool) GetNamespace(ctx context.Context, args GetNamespaceArgs) (*mcp_golang.ToolResponse, error) {
	var (
		opts metav1.GetOptions

		ns, err = k.client.
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
					Text: ns.String(),
				},
			},
		},
	}, nil

	// return fmt.Sprintf("Your Kubernetes namespace (%s) has been created", name), nil
}
