package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
		client        *mcp_golang.Client
		tools         []*Tool
		ollamaChatURL string
		model         string

		messages []*Message
	}
)

func (r *Router) fetchTools(ctx context.Context) error {
	// List available tools
	var res, err = r.client.ListTools(ctx, nil)
	if err != nil {
		return err
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

	return nil
}

func (r *Router) fetchResources(ctx context.Context) error {
	// List available tools
	var res, err = r.client.ListResources(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Printf("res: %+v\n", res)

	for _, resource := range res.Resources {
		readRes, err := r.client.ReadResource(ctx, resource.Uri)
		if err != nil {
			return err
		}

		fmt.Printf("resource: %+v\n", resource)
		fmt.Printf("readRes: %+v\n", readRes.Contents)
		for _, content := range readRes.Contents {
			// TODO: could we somehow offer templates as resources that can be used to apply into kubernetes?
			fmt.Printf("content: %+v\n", content.TextResourceContents.Text)
		}
	}

	return nil
}

func NewRouter(model, ollamaURI, kubeMCPURI string) (*Router, error) {
	// Create an HTTP transport that connects to the server
	var transport = mcp_golang_http.NewHTTPClientTransport("/mcp")
	transport.WithBaseURL(fmt.Sprintf("http://%s:8080", kubeMCPURI))

	// Create a new client with the transport
	var k8sMCPClient = mcp_golang.NewClient(transport)

	// Initialize the client
	_, err := k8sMCPClient.Initialize(context.Background())
	if err != nil {
		return nil, err
	}

	var (
		prompt = Message{
			Role:    "system",
			Content: "You are a helpful assistant. Never call tools unless absolutely necessary. Respond in plain language when possible. Be brief with your responses when using tools.",
		}

		router = Router{
			client:        k8sMCPClient,
			ollamaChatURL: fmt.Sprintf("http://%s:11434/api/chat", ollamaURI),
			model:         model,

			messages: []*Message{
				&prompt,
			},
		}
	)

	// TODO:
	var ctx = context.TODO()

	err = router.fetchTools(ctx)
	if err != nil {
		return nil, err
	}

	// err = router.fetchResources(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return &router, nil
}

func main() {
	// TODO: turn this into a config file with a config pkg later on
	var ollamaURI = os.Getenv("OLLAMA_URI")
	if ollamaURI == "" {
		ollamaURI = "localhost"
	}

	var kubeMCPURI = os.Getenv("KUBE_MCP_URI")
	if kubeMCPURI == "" {
		kubeMCPURI = "localhost"
	}

	var model = os.Getenv("MODEL")
	if model == "" {
		panic("MODEL not provided")
	}

	var (
		router, err = NewRouter(model, ollamaURI, kubeMCPURI)
		scanner     = bufio.NewScanner(os.Stdin)
	)

	if err != nil {
		panic(err)
	}

	go func() {
		var err = server(router)
		if err != nil {
			panic(err)
		}
	}()

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

		var output = router.handleText(text)

		fmt.Print(output)
	}
}

func (r *Router) handleText(text string) string {
	var output string

	switch {
	case text == "":
		return ""

	case strings.HasPrefix(text, "cmd::"):
		output = r.command(text)

	default:
		var err error
		output, err = r.loop(text)
		if err != nil {
			output = fmt.Sprintf("!!! %s !!!", err)
		}

		// TODO: or better yet make a struct that unmarshals this into a 'think'/'thought' field
		// and allow that to be printable with a command
	}

	return fmt.Sprintf("%s\n", output)
}

func trimOutput(output string) string {
	var split = strings.Split(output, "</think>")
	switch len(split) {
	case 2:
		output = split[1]

	case 1:

	default:
		fmt.Println("something is wrong here ... :", output)
	}

	return strings.TrimSpace(output)
}

func (r *Router) command(text string) string {
	text = strings.TrimPrefix(text, "cmd::")
	switch text {
	case "exit":
		fmt.Println("# exiting!")
		os.Exit(0)

	case "tools":
		r.tools = []*Tool{}
		// TODO:
		var err = r.fetchTools(context.TODO())
		if err != nil {
			return err.Error()
		}

		return "# reloaded tools"

	case "stream":
		stream = !stream
		return "# streaming enabled"

	case "nothink":
		if len(r.messages) == 0 {
			// TODO:
		}

		if r.messages[0].Content != "/no_think" {
			r.messages = append([]*Message{
				{
					Role:    "system",
					Content: "/no_think",
				},
			}, r.messages...)
		}

		return "# thinking disabled"

	case "think":
		if len(r.messages) == 0 {
			// TODO:
		}

		if len(r.messages) > 0 && r.messages[0].Content == "/no_think" {
			r.messages = r.messages[1:]
		}

		return "# thinking enabled"

	default:
		return fmt.Sprintf("# unrecognized command: `%s`\n", text)
	}

	return "# im not supposed to be here"
}

func (r *Router) loop(text string) (string, error) {
	var (
		initMsg = Message{
			Role:    "user",
			Content: text,
		}
	)

	// TODO: we are going to need to measure the context length at some point
	// and then we can start to either trim this or maybe summarize it in the background
	r.messages = append(r.messages, &initMsg)

	// TODO: make a command to enable debugging which will print this stuff out
	// for _, message := range r.messages {
	// 	fmt.Println("MESSAGE:", message)
	// }

	var (
		req = LLMRequest{
			Model:    r.model,
			Stream:   stream,
			Tools:    r.tools,
			Messages: r.messages,
		}

		b, err = json.Marshal(req)
	)

	if err != nil {
		return "", err
	}

	// fmt.Println("b:", string(b))

	resp, err := http.Post(r.ollamaChatURL, "application/json", bytes.NewBuffer(b))
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

	// TODO: make a command to enable debugging which will print this stuff out
	// fmt.Printf("res: %+v\n", res)

	if res.Error != "" {
		return "", errors.New(res.Error)
	}

	if res.Message == nil {
		return "", fmt.Errorf("no message sent back for some reason: %+v", res)
	}

	res.Message.Content = trimOutput(res.Message.Content)
	r.messages = append(r.messages, res.Message)

	// TODO: do we need to append the toolsCalls to the context?
	if len(res.Message.ToolCalls) == 0 {
		return res.Message.Content, nil
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

		toolRes, err := r.client.CallTool(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
		if err != nil {
			return "", err
		}

		if len(toolRes.Content) == 0 {
			fmt.Println("res.Content == 0")
			continue
		}

		var toolOutputs = []*Message{
			{
				Role: "assistant",
				ToolCalls: []*ToolCalls{
					toolCall,
				},
			},
		}

		toolCall.ID = strconv.Itoa(k)

		for _, content := range toolRes.Content {
			switch content.Type {
			case mcp_golang.ContentTypeText:
				toolOutputs = append(toolOutputs, &Message{
					Role:       "tool",
					ToolCallID: toolCall.ID,
					// TODO: need to figure out how to handle other types?
					Content: content.TextContent.Text,
				})

			// TODO: handle these other types
			case mcp_golang.ContentTypeImage:
				return "", errors.New("Images type content is not supported")

			case mcp_golang.ContentTypeEmbeddedResource:
				return "", errors.New("EmbeddedResource type content is not supported")

			default:
				return "", fmt.Errorf("default content type somehow?: %+v\n", content)
			}
		}

		// TODO: we are going to need to measure the context length at some point
		// and then we can start to either trim this or maybe summarize it in the background
		r.messages = append(r.messages, toolOutputs...)
	}

	var llmResp = LLMRequest{
		Model:    r.model,
		Stream:   stream,
		Messages: r.messages,
	}

	bb, err := json.Marshal(llmResp)
	if err != nil {
		return "", err
	}

	resp, err = http.Post(r.ollamaChatURL, "application/json", bytes.NewBuffer(bb))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var res2 LLMResponse
	err = json.NewDecoder(resp.Body).Decode(&res2)
	if err != nil {
		return "", err
	}

	return trimOutput(res2.Message.Content), nil
}
