package main

import (
	"fmt"
	"time"

	mcp_golang "github.com/metoro-io/mcp-golang"
	mcp_golang_http "github.com/metoro-io/mcp-golang/transport/http"
)

// tools/list
// tools/call

// TimeArgs represents the arguments for the current time tool
type TimeArgs struct {
	Format string `json:"format,omitempty" jsonschema:"description=Optional time format (default: RFC3339)"`
}

func main() {
	// var addr = fmt.Sprintf("%s:%d", net.IPv4allrouter.String(), 8080)
	// var err = http.ListenAndServe(addr, nil)
	// if err != nil {
	// 	panic(err)
	// }

	transport := mcp_golang_http.NewHTTPTransport("/mcp")
	transport.WithAddr(":8080")
	server := mcp_golang.NewServer(transport)

	// Register current time tool
	var err = server.RegisterTool("time", "Returns the current time", func(args TimeArgs) (*mcp_golang.ToolResponse, error) {
		format := time.RFC3339
		if args.Format != "" {
			format = args.Format
		}

		message := time.Now().Format(format)
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(message)), nil
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

// type ListToolsRes struct {
// 	Tool []*Tool
// }

// type Tool struct{}

// func call(w http.ResponseWriter, r *http.Request) {

// }

// func list(w http.ResponseWriter, r *http.Request) {}
