package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	mcp_golang_http "github.com/metoro-io/mcp-golang/transport/http"
)

const (
	weatherTemplate = "The weather in %s is %d°C and partly cloudy."
)

var (
	stream   bool
	stdTools = []*Tools{
		{
			Type: "function",
			Function: &Function{
				Name:        "create_namespace",
				Description: "Create a namespace in Kubernetes",
				Parameters: &Parameters{
					Type: "object",
					Required: []string{
						"name",
					},
					Properties: map[string]Property{
						"name": {
							Type:        "string",
							Description: "Name of the namespace",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: &Function{
				Name:        "create_deployment",
				Description: "Create a deployment on Kubernetes",
				Parameters: &Parameters{
					Type: "object",
					Required: []string{
						"name",
						"image",
						"namespace",
					},
					Properties: map[string]Property{
						"name": {
							Type:        "string",
							Description: "Name of the deployment",
						},
						"image": {
							Type:        "string",
							Description: "Image to run",
						},
						"namespace": {
							Type:        "string",
							Description: "Namespace that the deployment should be in",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: &Function{
				Name:        "get_weather",
				Description: "Get the current weather for a city",
				Parameters: &Parameters{
					Type: "object",
					Required: []string{
						"city_name",
					},
					Properties: map[string]Property{
						"city_name": {
							Type:        "string",
							Description: "Name of the city",
						},
					},
				},
			},
		},
	}
)

func main() {
	// Create an HTTP transport that connects to the server
	transport := mcp_golang_http.NewHTTPClientTransport("/mcp")
	transport.WithBaseURL("http://localhost:8080")

	// Create a new client with the transport
	client := mcp_golang.NewClient(transport)

	// Initialize the client
	_, err := client.Initialize(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// List available tools
	tools, err := client.ListTools(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	log.Println("Available Tools:")
	for _, tool := range tools.Tools {
		desc := ""
		if tool.Description != nil {
			desc = *tool.Description
		}

		log.Printf("Tool: %s. Description: %s", tool.Name, desc)
	}

	var (
		k = NewKubernetesTool()

		toolFns = map[string]ToolFn{
			"create_deployment": k.CreateDeployment,
			"create_namespace":  k.CreateNamespace,
			"get_weather":       getWeather,
		}

		baseURI = os.Getenv("BASE_URI")
	)

	if baseURI == "" {
		baseURI = "localhost"
	}

	var (
		u       = fmt.Sprintf("http://%s:11434/api/chat", baseURI)
		scanner = bufio.NewScanner(os.Stdin)
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
			output = command(text)

		default:
			output, err = loop(text, u, toolFns)
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

func command(text string) string {
	text = strings.TrimPrefix(text, "cmd::")
	switch text {
	case "exit":
		fmt.Println("# exiting!")
		os.Exit(0)

	case "stream":
		stream = !stream
		return "# streaming enabled"

	default:
		return fmt.Sprintf("# unrecognized command: `%s`\n", text)
	}

	return "# im not supposed to be here"
}

func loop(text, u string, toolFns map[string]ToolFn) (string, error) {
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
			Model:  "qwen3:4b",
			Stream: stream,
			Tools:  stdTools,
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

		var fn, ok = toolFns[toolCall.Function.Name]
		if !ok {
			return "", errors.New("fn not found")
		}

		if fn == nil {
			return "", errors.New("fn was nil")
		}

		toolCall.ID = strconv.Itoa(k)
		output, err := fn(ctx, toolCall.Function.Arguments)
		if err != nil {
			return "", err
		}

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
		Model:    "qwen3:4b",
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

func getWeather(ctx context.Context, args map[string]string) (string, error) {
	// fmt.Printf("args: %+v\n", args)

	var cityName, ok = args["city_name"]
	if !ok {
		return "", errors.New("city_name was not found")
	}

	return fmt.Sprintf(weatherTemplate, cityName, rand.Intn(15)+10), nil
}
