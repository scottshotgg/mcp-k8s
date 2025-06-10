package main

import (
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	mcp_golang_http "github.com/metoro-io/mcp-golang/transport/http"
)

func main() {
	var (
		k         = NewKubernetesTool()
		transport = mcp_golang_http.NewHTTPTransport("/mcp")
	)

	transport.WithAddr(":8080")

	var server = mcp_golang.NewServer(transport)

	// --- Namespace commands ---

	// Register create_namespace
	var err = server.RegisterTool("create_namespace", "Create a namespace in Kubernetes", k.CreateNamespace)
	if err != nil {
		panic(err)
	}

	// Register delete_namespace
	err = server.RegisterTool("delete_namespace", "Remove or delete an existing namespace in Kubernetes", k.DeleteNamespace)
	if err != nil {
		panic(err)
	}

	// Register list_namespaces
	err = server.RegisterTool("list_namespaces", "List all accessible existing namespaces in Kubernetes", k.ListNamespaces)
	if err != nil {
		panic(err)
	}

	// Register get_namespace
	err = server.RegisterTool("get_namespace", "Get details for a particular namespace in Kubernetes", k.GetNamespace)
	if err != nil {
		panic(err)
	}

	// --- Deployment commands ---

	// Register create_deployment
	err = server.RegisterTool("create_deployment", "Create a deployment in Kubernetes", k.CreateDeployment)
	if err != nil {
		panic(err)
	}

	// Register delete_deployment
	err = server.RegisterTool("delete_deployment", "Remove or delete an existing deployment in Kubernetes", k.DeleteDeployment)
	if err != nil {
		panic(err)
	}

	// Register get_deployment
	err = server.RegisterTool("get_deployment", "Get details for a particular deployment in a particular namespace in Kubernetes", k.GetDeployment)
	if err != nil {
		panic(err)
	}

	// Register list_deployments
	err = server.RegisterTool("list_deployments", "List all deployments for a particular namespace in Kubernetes", k.ListDeployments)
	if err != nil {
		panic(err)
	}

	// --- Nodes commands ---

	// Register list_deployments
	err = server.RegisterTool("top_nodes", "Get the resource usage of all nodes in Kubernetes", k.TopNodes)
	if err != nil {
		panic(err)
	}

	// --- Debugging commands ---

	// Register list_deployments
	err = server.RegisterTool("debug_workload", "Debug a failing or crashing workload (pod, container) in Kubernetes", k.DebugWorkload)
	if err != nil {
		panic(err)
	}

	// --- Resources ---

	// Register a resource
	err = server.RegisterResource("test://resource", "Resource Name", "This is a test resource", "application/json", func() (*mcp_golang.ResourceResponse, error) {
		// Define the resource content
		content := mcp_golang.NewTextEmbeddedResource("test://resource", "This is a test resource", "application/json")
		return mcp_golang.NewResourceResponse(content), nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("started!")

	err = server.Serve()
	if err != nil {
		panic(err)
	}
}
