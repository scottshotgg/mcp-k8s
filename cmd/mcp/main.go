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

	// Register create_namespace
	var err = server.RegisterTool("create_namespace", "Create a namespace in Kubernetes", k.CreateNamespace)
	if err != nil {
		panic(err)
	}

	// Register create_deployment
	err = server.RegisterTool("create_deployment", "Create a deployment in Kubernetes", k.CreateDeployment)
	if err != nil {
		panic(err)
	}

	// server.RegisterResource()

	fmt.Println("started!")

	err = server.Serve()
	if err != nil {
		panic(err)
	}
}
