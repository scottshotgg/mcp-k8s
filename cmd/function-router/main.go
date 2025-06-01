package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	mcp_golang_http "github.com/metoro-io/mcp-golang/transport/http"
)

var (
	// TODO: try to get streaming working
	stream bool
)

type (
	Router struct {
		kubeMCPClient *mcp_golang.Client
		tools         []*Tool
	}
)

func (r *Router) fetchTools() {
	// List available tools
	var res, err = r.kubeMCPClient.ListTools(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	for _, tool := range res.Tools {
		r.tools = append(r.tools, &Tool{
			Type: "function",
			Function: &Function{
				Name:        tool.Name,
				Description: description(tool.Description),
				Parameters:  tool.InputSchema,
			},
		})
	}
}

func description(desc *string) string {
	if desc != nil {
		return *desc
	}

	return ""
}

func main() {
	var ollamaURI = os.Getenv("OLLAMA_URI")
	if ollamaURI == "" {
		ollamaURI = "localhost"
	}

	var kubeMCPURI = os.Getenv("KUBE_MCP_URI")
	if kubeMCPURI == "" {
		kubeMCPURI = "localhost"
	}

	// Create an HTTP transport that connects to the server
	var transport = mcp_golang_http.NewHTTPClientTransport("/mcp")
	transport.WithBaseURL(fmt.Sprintf("http://%s:8080", kubeMCPURI))

	// Create a new client with the transport
	var k8sMCPClient = mcp_golang.NewClient(transport)

	// Initialize the client
	_, err := k8sMCPClient.Initialize(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// TODO: make a new func
	var router = Router{
		kubeMCPClient: k8sMCPClient,
	}

	router.fetchTools()

	var (
		ollamaChatURL = fmt.Sprintf("http://%s:11434/api/chat", ollamaURI)
		scanner       = bufio.NewScanner(os.Stdin)
	)

	fmt.Println("Ask me anything:")
	for {
		fmt.Print("> ")

		var ok = scanner.Scan()
		if !ok {
			panic("not ok to scan?")
		}

		var (
			text = scanner.Text()
			err  = scanner.Err()
		)

		if err != nil {
			panic(err)
		}

		var output string

		switch {
		case text == "":
			continue

		case strings.HasPrefix(text, "cmd::"):
			output = router.command(text)

		default:
			output, err = router.loop(text, ollamaChatURL, k8sMCPClient)
			if err != nil {
				fmt.Printf("!!! %s !!!", output)
				continue
			}

			// TODO: try turing on 'no think'
			// TODO: or better yet make a struct that unmarshals this into a 'think'/'thought' field
			// and allow that to be printable with a command
			var split = strings.Split(output, "</think>")
			switch len(split) {
			case 2:
				output = split[1]

			case 1:

			default:
				fmt.Println("something is wrong here ... :", output)
			}

			output = strings.TrimSpace(output)
		}

		fmt.Printf("\n%s\n", output)
	}
}

func (r *Router) command(text string) string {
	text = strings.TrimPrefix(text, "cmd::")
	switch text {
	case "exit":
		fmt.Println("# exiting!")
		os.Exit(0)

	case "tools":
		r.tools = []*Tool{}
		r.fetchTools()
		return "# reloaded tools"

	case "stream":
		stream = !stream
		return "# streaming enabled"

	default:
		return fmt.Sprintf("# unrecognized command: `%s`\n", text)
	}

	return "# im not supposed to be here"
}

func (r *Router) loop(text, u string, k8sMCPClient *mcp_golang.Client) (string, error) {
	var (
		prompt = Message{
			Role:    "system",
			Content: "You are a helpful assistant. Never call tools unless absolutely necessary. Respond in plain language when possible.",
		}

		initMsg = Message{
			Role:    "user",
			Content: text,
		}

		req = LLMRequest{
			Model:  "qwen3:14b",
			Stream: stream,
			Tools:  r.tools,
			Messages: []*Message{
				&prompt,
				&initMsg,
			},
		}

		b, err = json.Marshal(req)
	)

	if err != nil {
		return "", err
	}

	// fmt.Println("b:", string(b))

	resp, err := http.Post(u, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// bb, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("bb:", string(bb))

	var res LLMResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	// err = json.Unmarshal(bb, &res)
	if err != nil {
		return "", err
	}

	// fmt.Printf("res: %+v\n", res)

	if res.Message == nil {
		return "", errors.New("res was nil")
	}

	if len(res.Message.ToolCalls) == 0 {
		return res.Message.Content, nil
	}

	var messages = []*Message{
		&initMsg,
	}

	// TODO: tie this into the request and stuff as well
	var ctx = context.Background()

	for k, toolCall := range res.Message.ToolCalls {
		// fmt.Println("Calling function:")
		// fmt.Println("Name:", toolCall.Function.Name)
		// fmt.Println("Arguments:", toolCall.Function.Arguments)

		// var fn, ok = toolFns[toolCall.Function.Name]
		// if !ok {
		// 	return "", errors.New("fn not found")
		// }

		// if fn == nil {
		// 	return "", errors.New("fn was nil")
		// }

		res, err := k8sMCPClient.CallTool(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
		if err != nil {
			return "", err
		}

		if len(res.Content) == 0 {
			fmt.Println("res.Content == 0")
			continue
		}

		var output string
		// TODO: handle multiple messages coming back later
		switch res.Content[0].Type {
		case mcp_golang.ContentTypeText:
			output = res.Content[0].TextContent.Text

			// TODO: implement other cases later

		default:
			fmt.Println("default content type somehow?")
			continue
		}

		toolCall.ID = strconv.Itoa(k)

		messages = append(messages, []*Message{
			{
				Role: "assistant",
				// ToolCalls: []*ToolCalls{
				// 	{
				// 		ID:       strconv.Itoa(k),
				// 		Type:     "function",
				// 		Function: toolCall.Function,
				// 	},
				// },
				ToolCalls: []*ToolCalls{
					toolCall,
				},
			},
			{
				Role:       "tool",
				ToolCallID: toolCall.ID,
				Content:    output,
			},
		}...)
	}

	var llmResp = LLMRequest{
		Model:    "qwen3:14b",
		Stream:   stream,
		Messages: messages,
	}

	bb, err := json.Marshal(llmResp)
	if err != nil {
		return "", err
	}

	resp, err = http.Post(u, "application/json", bytes.NewBuffer(bb))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var res2 LLMResponse
	err = json.NewDecoder(resp.Body).Decode(&res2)
	if err != nil {
		return "", err
	}

	return res2.Message.Content, nil
}
