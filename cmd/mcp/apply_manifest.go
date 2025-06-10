package main

import (
	"bytes"
	"context"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

// TODO: Take deployment as a type
type ApplyManifestArgs struct {
	Name      string `json:"name" jsonschema:"required,description=Name of the deployment"`
	Namespace string `json:"namespace" jsonschema:"required,description=Namespace that the deployment should be deployed into in"`
	Manifest  string `json:"manifest" jsonschema:"required,description=Manifest of the deployment"`
}

func (k *KubernetesTool) ApplyManifest(ctx context.Context, args ApplyManifestArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("apply_manifest")

	// Use UniversalDeserializer from scheme
	var deserializer = scheme.Codecs.UniversalDeserializer()

	// Decode the YAML into a runtime.Object
	var obj, gvk, err = deserializer.Decode([]byte(args.Manifest), nil, nil)
	if err != nil {
		return nil, err
	}

	// Type assert to *appsv1.Deployment
	var deployment, ok = obj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("decoded object is not a Deployment: got %v", gvk)
	}

	var (
		opts metav1.UpdateOptions
	)

	deployment, err = k.client.
		AppsV1().
		Deployments(args.Namespace).
		Update(ctx, deployment, opts)

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
