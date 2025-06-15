package main

import (
	"fmt"

	kt "git.6740.io/scottshotgg/mcp-k8s/pkg/kubernetes_tool"
	mcp_golang "github.com/metoro-io/mcp-golang"
	mcp_golang_http "github.com/metoro-io/mcp-golang/transport/http"
)

func main() {
	var k, err = kt.NewKubernetesTool()
	if err != nil {
		panic(err)
	}

	var transport = mcp_golang_http.NewHTTPTransport("/mcp")
	transport.WithAddr(":8080")

	var server = mcp_golang.NewServer(transport)

	// --- Kubectl tools ---

	// Register get_manifest
	err = server.RegisterTool("run_kubectl_command", "Run a Kubectl command against the Kubernetes cluster", k.RunKubectlCommand)
	if err != nil {
		panic(err)
	}

	// // --- Initial State command ---

	// // Register query_initial_state
	// err = server.RegisterTool("query_initial_state", "Query the world view of the Kubernetes cluster", k.QueryInitialState)
	// if err != nil {
	// 	panic(err)
	// }

	// // --- Namespace tools ---

	// // Register create_namespace
	// err = server.RegisterTool("create_namespace", "Create a namespace in Kubernetes", k.CreateNamespace)
	// if err != nil {
	// 	panic(err)
	// }

	// // Register delete_namespace
	// err = server.RegisterTool("delete_namespace", "Remove or delete an existing namespace in Kubernetes", k.DeleteNamespace)
	// if err != nil {
	// 	panic(err)
	// }

	// // Register list_namespaces
	// err = server.RegisterTool("list_namespaces", "List all accessible existing namespaces in Kubernetes", k.ListNamespaces)
	// if err != nil {
	// 	panic(err)
	// }

	// // Register get_namespace
	// err = server.RegisterTool("get_namespace", "Get details for a particular namespace in Kubernetes", k.GetNamespace)
	// if err != nil {
	// 	panic(err)
	// }

	// // --- Pod tools ---

	// // Register list_pods
	// err = server.RegisterTool("list_pods", "List pods for a particular namespace in Kubernetes", k.ListPods)
	// if err != nil {
	// 	panic(err)
	// }

	// // --- Deployment tools ---

	// // Register create_deployment
	// err = server.RegisterTool("create_deployment", "Create a deployment in Kubernetes", k.CreateDeployment)
	// if err != nil {
	// 	panic(err)
	// }

	// // Register delete_deployment
	// err = server.RegisterTool("delete_deployment", "Remove or delete an existing deployment in Kubernetes", k.DeleteDeployment)
	// if err != nil {
	// 	panic(err)
	// }

	// // Register get_deployment
	// err = server.RegisterTool("get_deployment", "Get details for a particular deployment in a particular namespace in Kubernetes", k.GetDeployment)
	// if err != nil {
	// 	panic(err)
	// }

	// // Register list_deployments
	// err = server.RegisterTool("list_deployments", "List all deployments for a particular namespace in Kubernetes", k.ListDeployments)
	// if err != nil {
	// 	panic(err)
	// }

	// // --- Nodes tools ---

	// // Register list_nodes
	// err = server.RegisterTool("list_nodes", "List all nodes in Kubernetes", k.ListNodes)
	// if err != nil {
	// 	panic(err)
	// }

	// // Register top_nodes
	// err = server.RegisterTool("top_nodes", "Get the resource usage of all nodes in Kubernetes", k.TopNodes)
	// if err != nil {
	// 	panic(err)
	// }

	// // --- Debugging tools ---

	// // Register debug_workload
	// err = server.RegisterTool("debug_workload", "Debug a failing or crashing workload (pod, container) in Kubernetes", k.DebugWorkload)
	// if err != nil {
	// 	panic(err)
	// }

	// // --- Manifest tools ---

	// // Register get_manifest
	// err = server.RegisterTool("get_manifest", "Retrieve the manifest for a deployment in Kubernetes", k.GetManifest)
	// if err != nil {
	// 	panic(err)
	// }

	// // Register apply_manifest
	// err = server.RegisterTool("apply_manifest", "Apply/Update the manifest for a deployment in Kubernetes", k.ApplyManifest)
	// if err != nil {
	// 	panic(err)
	// }

	// // --- Resources ---

	// // Register a resource
	// err = server.RegisterResource("test://resource", "Resource Name", "This is a test resource", "application/json", resourceTest)

	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("started!")

	err = server.Serve()
	if err != nil {
		panic(err)
	}
}

func resourceTest() (*mcp_golang.ResourceResponse, error) {
	// Define the resource content
	var content = mcp_golang.NewTextEmbeddedResource("test://resource", "This is a test resource", "application/json")

	return mcp_golang.NewResourceResponse(content), nil
}

// func registerTool(server *mcp_golang.Server, name string, description string, handler any) error {
// 	var handler = func()

// 	var err = server.RegisterTool(name, description, handler)
// 	if err != nil {
// 		return err
// 	}

// }
