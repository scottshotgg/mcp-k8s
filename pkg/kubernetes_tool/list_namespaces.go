package kubernetes_tool

import (
	"context"
	"fmt"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListNamespaceArgs struct{}

func (k *KubernetesTool) ListNamespaces(ctx context.Context, _ CreateNamespaceArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("list_namespaces")

	var (
		opts metav1.ListOptions

		namespacesList, err = k.client.
					CoreV1().
					Namespaces().
					List(ctx, opts)
	)

	if err != nil {
		return nil, err
	}

	var namespaces []string
	for _, namespace := range namespacesList.Items {
		namespaces = append(namespaces, namespace.Name)
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: fmt.Sprintf("The currently accessible Kubernetes namespaces are: %s", strings.Join(namespaces, ", ")),
				},
			},
		},
	}, nil
}
