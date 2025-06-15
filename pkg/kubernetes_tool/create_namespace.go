package kubernetes_tool

import (
	"context"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateNamespaceArgs struct {
	Name string `json:"name" jsonschema:"required,description=Name of the namespace"`
}

func (k *KubernetesTool) CreateNamespace(ctx context.Context, args CreateNamespaceArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("create_namespace")

	var (
		opts metav1.CreateOptions

		k8sNS = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: args.Name,
			},
		}

		_, err = k.client.
			CoreV1().
			Namespaces().
			Create(ctx, k8sNS, opts)
	)

	if err != nil {
		return nil, err
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: fmt.Sprintf("Your Kubernetes namespace (%s) has been created", args.Name),
				},
			},
		},
	}, nil
}
