package main

import (
	"time"
)

type LLMResponse struct {
	Model              string    `json:"model,omitempty"`
	Stream             bool      `json:"stream,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	Message            *Message  `json:"message,omitempty"`
	DoneReason         string    `json:"done_reason,omitempty"`
	Done               bool      `json:"done,omitempty"`
	TotalDuration      int64     `json:"total_duration,omitempty"`
	LoadDuration       int       `json:"load_duration,omitempty"`
	PromptEvalCount    int       `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int       `json:"prompt_eval_duration,omitempty"`
	EvalCount          int       `json:"eval_count,omitempty"`
	EvalDuration       int64     `json:"eval_duration,omitempty"`
}

type Function struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	// Parameters  *Parameters       `json:"parameters,omitempty"`
	Parameters interface{}       `json:"parameters,omitempty"`
	Arguments  map[string]string `json:"arguments,omitempty"`
}

type ToolCalls struct {
	ID       string    `json:"id,omitempty"`
	Type     string    `json:"type,omitempty"`
	Function *Function `json:"function,omitempty"`
}

type Message struct {
	Role       string       `json:"role,omitempty"`
	Content    string       `json:"content,omitempty"`
	ToolCalls  []*ToolCalls `json:"tool_calls,omitempty"`
	ToolCallID string       `json:"tool_call_id,omitempty"`
}

type LLMRequest struct {
	// TODO: not sure how that works
	// EnableThinking bool       `json:"enable_thinking,omitempty"`
	Model    string     `json:"model,omitempty"`
	Stream   bool       `json:"stream"`
	Messages []*Message `json:"messages,omitempty"`
	Tools    []*Tool    `json:"tools,omitempty"`
}

type Property struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

type Parameters struct {
	Type       string              `json:"type,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

type Tool struct {
	Type     string    `json:"type,omitempty"`
	Function *Function `json:"function,omitempty"`
}

type WeatherRequest struct {
	Location string `json:"location"`
}

type WeatherResponse struct {
	Location string `json:"location"`
	TempC    string `json:"temp_c"`
	Summary  string `json:"summary"`
}
