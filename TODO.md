# TODO

- More functionality
  x add getting a deployment
  x Listing pods
  - exposing a deployment / creating a service
  - autoscaling a deployment / creating an hpa
  - creating a job
  - running a job
  x listing nodes
  - getting a node
  - top pods
  x top nodes

- Config
  - json file
  - pkg
  - configurable ports
  - configurable model

- Experiment more with MCP prompt and resource stuff
  - maybe resource could be useful for manifests?

- Figure out streaming
  - stream the response back and print it out on the screen

- Remove mcp_golang from the kubernetes_tool file
  - Make the mcp server just a server with mappers

- Clean up the repo structure

- Maybe look at implementing other cases of media later on

- Give the LLM a way to retrieve and modify the YAML to debug a pod issue