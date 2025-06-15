package kubernetes_tool

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
)

type RunKubectlCommandArgs struct {
	Command string `json:"command" jsonschema:"required,description=Kubectl command to run; always starts with kubectl"`
}

func (k *KubernetesTool) RunKubectlCommand(ctx context.Context, args RunKubectlCommandArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("run_kubectl_command")

	// TODO: Run an os.Command here
	// This is very open-ended and dangerous. We should wrap it with an flag to enable/disable it willingly. We should also instruct the LLM
	// to ALWAYS ask the user before running this command

	if !strings.HasPrefix(args.Command, "kubectl") {
		args.Command = "kubectl " + args.Command
	}

	fmt.Println("command:", args.Command)

	var cmd = exec.CommandContext(ctx, "sh", "-c", args.Command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	fmt.Println("output:", string(output))

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: string(output),
				},
			},
		},
	}, nil
}
