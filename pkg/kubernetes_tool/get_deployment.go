package kubernetes_tool

import (
	"context"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetDeploymentArgs struct {
	Name      string `json:"name" jsonschema:"required,description=Name of the deployment"`
	Namespace string `json:"namespace" jsonschema:"required,description=Namespace that the deployment should be deployed into in"`
}

func (k *KubernetesTool) GetDeployment(ctx context.Context, args GetDeploymentArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("get_deployment")

	var (
		opts metav1.GetOptions

		deployment, err = k.client.
				AppsV1().
				Deployments(args.Namespace).
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
					Text: deployment.String(),
				},
			},
		},
	}, nil
}
