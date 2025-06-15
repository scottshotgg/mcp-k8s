package kubernetes_tool

import (
	"context"
	"fmt"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListNodesArgs struct{}

func (k *KubernetesTool) ListNodes(ctx context.Context, args ListNodesArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("list_nodes")

	var (
		opts metav1.ListOptions

		nodeList, err = k.client.
				CoreV1().
				Nodes().
				List(ctx, opts)
	)

	if err != nil {
		return nil, err
	}

	var nodes []string
	for _, n := range nodeList.Items {
		nodes = append(nodes, n.Name)
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: fmt.Sprintf("The Kubernetes nodes are: %s", strings.Join(nodes, ", ")),
				},
			},
		},
	}, nil
}
