package kubernetes_tool

import (
	"bytes"
	"context"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type GetManifestArgs struct {
	Name      string `json:"name" jsonschema:"required,description=Name of the deployment"`
	Namespace string `json:"namespace" jsonschema:"required,description=Namespace that the deployment should be deployed into in"`
}

func (k *KubernetesTool) GetManifest(ctx context.Context, args GetManifestArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("get_manifest")

	var (
		opts metav1.GetOptions

		deployment, err = k.client.
				AppsV1().
			// TODO: Take deployment as a type
			Deployments(args.Namespace).
			Get(ctx, args.Name, opts)
	)

	if err != nil {
		return nil, err
	}

	var (
		// Strip metadata/status fields
		cleanDeployment = stripDeploymentForExport(deployment)

		scheme = runtime.NewScheme()
		codec  = json.NewYAMLSerializer(json.DefaultMetaFactory, scheme, scheme)

		b bytes.Buffer
	)

	err = codec.Encode(cleanDeployment, &b)
	if err != nil {
		return nil, err
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: b.String(),
				},
			},
		},
	}, nil
}

func stripDeploymentForExport(d *appsv1.Deployment) *appsv1.Deployment {
	// DeepCopy to avoid mutating original
	var copy = d.DeepCopy()

	// Preserve TypeMeta (Kind and APIVersion)
	copy.TypeMeta = metav1.TypeMeta{
		Kind:       "Deployment",
		APIVersion: "apps/v1",
	}

	// Clean metadata
	copy.ObjectMeta = metav1.ObjectMeta{
		Name:        copy.Name,
		Namespace:   copy.Namespace,
		Labels:      copy.Labels,
		Annotations: copy.Annotations,
	}

	// Clean spec.template.metadata
	copy.Spec.Template.ObjectMeta = metav1.ObjectMeta{
		Labels:      copy.Spec.Template.Labels,
		Annotations: copy.Spec.Template.Annotations,
	}

	// Remove status
	copy.Status = appsv1.DeploymentStatus{}

	return copy
}
