#!/usr/bin/bash

echo "Building MCP server"
CGO_ENABLED=0 $(which go) build -o "./bin/mcp" "./cmd/mcp"