package kubernetes_tool

import (
	"context"
	"fmt"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListPodsArgs struct {
	Namespace string `json:"namespace" jsonschema:"required,description=Namespace that the deployment should be deployed into in"`
}

func (k *KubernetesTool) ListPods(ctx context.Context, args ListPodsArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("list_pods")

	var (
		opts metav1.ListOptions

		podsList, err = k.client.
				CoreV1().
				Pods(args.Namespace).
				List(ctx, opts)
	)

	if err != nil {
		return nil, err
	}

	var pods []string
	for _, pod := range podsList.Items {
		pods = append(pods, pod.Name)
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: fmt.Sprintf("The Kubernetes pod in the %s namespace are: %s", args.Namespace, strings.Join(pods, ", ")),
				},
			},
		},
	}, nil
}
