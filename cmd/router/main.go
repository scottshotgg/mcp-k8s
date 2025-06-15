package main

import (
	"os"

	"git.6740.io/scottshotgg/mcp-k8s/pkg/router"
)

func main() {
	// TODO: turn this into a config file with a config pkg later on
	var ollamaURI = os.Getenv("OLLAMA_URI")
	if ollamaURI == "" {
		ollamaURI = "127.0.0.1"
	}

	var kubeMCPURI = os.Getenv("KUBE_MCP_URI")
	if kubeMCPURI == "" {
		kubeMCPURI = "127.0.0.1"
	}

	var model = os.Getenv("MODEL")
	if model == "" {
		panic("MODEL not provided")
	}

	var r, err = router.New(model, ollamaURI, kubeMCPURI)
	if err != nil {
		panic(err)
	}

	go func() {
		var err = server(r)
		if err != nil {
			panic(err)
		}
	}()

	r.Start()
}
