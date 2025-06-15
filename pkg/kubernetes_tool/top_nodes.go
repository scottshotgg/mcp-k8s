package kubernetes_tool

import (
	"context"
	"encoding/json"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TopNodesArgs struct{}

// kubectl get --raw /apis/metrics.k8s.io/v1beta1/nodes
func (k *KubernetesTool) TopNodes(ctx context.Context, args TopNodesArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("top_nodes")

	var (
		opts metav1.ListOptions

		metricsList, err = k.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, opts)
	)

	if err != nil {
		return nil, err
	}

	var metrics = map[string]corev1.ResourceList{}
	for _, m := range metricsList.Items {
		metrics[m.Name] = m.Usage
	}

	metricsStr, err := json.Marshal(metrics)

	// TODO: MAYBE we should make multiple text responses instead of concatting everything together?
	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: string(metricsStr),
				},
			},
		},
	}, nil
}
