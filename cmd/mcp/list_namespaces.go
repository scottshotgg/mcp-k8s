package main

import (
	"context"
	"fmt"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListNamespaceArgs struct{}

func (k *KubernetesTool) ListNamespaces(ctx context.Context, _ CreateNamespaceArgs) (*mcp_golang.ToolResponse, error) {
	var (
		opts metav1.ListOptions

		namespaces, err = k.client.
				CoreV1().
				Namespaces().
				List(ctx, opts)
	)

	if err != nil {
		return nil, err
	}

	var nss []string
	for _, ns := range namespaces.Items {
		nss = append(nss, ns.Name)
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: fmt.Sprintf("The currently accessible Kubernetes namespaces are: %s", strings.Join(nss, ", ")),
				},
			},
		},
	}, nil

	// return fmt.Sprintf("Your Kubernetes namespace (%s) has been created", name), nil
}
