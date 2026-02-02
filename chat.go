package zai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ChatService handles chat completion operations
type ChatService struct {
	client *BaseClient
}

// Message represents a chat message
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // Can be string or array of content parts
}

// ContentPart represents a part of multimodal content
type ContentPart struct {
	Type     string    `json:"type"`      // "text" or "image_url"
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// ImageURL represents an image URL in multimodal content
type ImageURL struct {
	URL string `json:"url"`
}

// Function represents a function call
type Function struct {
	Arguments string `json:"arguments"`
	Name      string `json:"name"`
}

// ToolCall represents a tool call in a message
type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// CompletionMessage represents a completion message
type CompletionMessage struct {
	Content          *string    `json:"content,omitempty"`
	Role             string     `json:"role"`
	ReasoningContent *string    `json:"reasoning_content,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
}

// PromptTokensDetails represents detailed token usage for prompts
type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

// CompletionTokensDetails represents detailed token usage for completions
type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

// CompletionUsage represents token usage information
type CompletionUsage struct {
	PromptTokens            int                      `json:"prompt_tokens"`
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`
	CompletionTokens        int                      `json:"completion_tokens"`
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"`
	TotalTokens             int                      `json:"total_tokens"`
}

// CompletionChoice represents a completion choice
type CompletionChoice struct {
	Index        int               `json:"index"`
	FinishReason string            `json:"finish_reason"`
	Message      CompletionMessage `json:"message"`
}

// ChatCompletion represents a chat completion response
type ChatCompletion struct {
	ID      string             `json:"id"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   CompletionUsage    `json:"usage"`
}

// ChatCompletionChunk represents a streaming chunk
type ChatCompletionChunk struct {
	ID      string                    `json:"id"`
	Created int64                     `json:"created"`
	Model   string                    `json:"model"`
	Choices []ChatCompletionChunkChoice `json:"choices"`
}

// ChatCompletionChunkChoice represents a streaming choice
type ChatCompletionChunkChoice struct {
	Index        int                         `json:"index"`
	Delta        ChatCompletionChunkDelta    `json:"delta"`
	FinishReason *string                     `json:"finish_reason,omitempty"`
}

// ChatCompletionChunkDelta represents the delta in a streaming chunk
type ChatCompletionChunkDelta struct {
	Content          *string    `json:"content,omitempty"`
	Role             *string    `json:"role,omitempty"`
	ReasoningContent *string    `json:"reasoning_content,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
}

// Tool represents a tool that can be called
type Tool struct {
	Type      string                 `json:"type"`
	Function  *FunctionDefinition    `json:"function,omitempty"`
	WebSearch *WebSearchTool         `json:"web_search,omitempty"`
}

// FunctionDefinition represents a function definition
type FunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// WebSearchTool represents a web search tool
type WebSearchTool struct {
	SearchQuery  string `json:"search_query,omitempty"`
	SearchResult bool   `json:"search_result,omitempty"`
}

// SensitiveWordCheck represents sensitive word check configuration
type SensitiveWordCheck struct {
	Type string `json:"type"`
}

// ChatCompletionRequest represents a chat completion request
type ChatCompletionRequest struct {
	Model              string              `json:"model"`
	Messages           []Message           `json:"messages"`
	RequestID          *string             `json:"request_id,omitempty"`
	UserID             *string             `json:"user_id,omitempty"`
	DoSample           *bool               `json:"do_sample,omitempty"`
	Stream             *bool               `json:"stream,omitempty"`
	Temperature        *float64            `json:"temperature,omitempty"`
	TopP               *float64            `json:"top_p,omitempty"`
	MaxTokens          *int                `json:"max_tokens,omitempty"`
	Seed               *int                `json:"seed,omitempty"`
	Stop               interface{}         `json:"stop,omitempty"` // string or []string
	SensitiveWordCheck *SensitiveWordCheck `json:"sensitive_word_check,omitempty"`
	Tools              []Tool              `json:"tools,omitempty"`
	ToolChoice         *string             `json:"tool_choice,omitempty"`
	Meta               map[string]string   `json:"meta,omitempty"`
	ResponseFormat     interface{}         `json:"response_format,omitempty"`
	Thinking           interface{}         `json:"thinking,omitempty"`
	WatermarkEnabled   *bool               `json:"watermark_enabled,omitempty"`
	ToolStream         *bool               `json:"tool_stream,omitempty"`
}

// ChatCompletionStream represents a streaming response
type ChatCompletionStream struct {
	reader  *bufio.Reader
	response *http.Response
}

// Next reads the next chunk from the stream
func (s *ChatCompletionStream) Next() (*ChatCompletionChunk, error) {
	for {
		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			return nil, fmt.Errorf("failed to read stream: %w", err)
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// SSE format: "data: {...}"
		if !bytes.HasPrefix(line, []byte("data: ")) {
			continue
		}

		data := bytes.TrimPrefix(line, []byte("data: "))

		// Check for stream end
		if bytes.Equal(data, []byte("[DONE]")) {
			return nil, io.EOF
		}

		var chunk ChatCompletionChunk
		if err := json.Unmarshal(data, &chunk); err != nil {
			return nil, fmt.Errorf("failed to unmarshal chunk: %w", err)
		}

		return &chunk, nil
	}
}

// Close closes the stream
func (s *ChatCompletionStream) Close() error {
	if s.response != nil && s.response.Body != nil {
		return s.response.Body.Close()
	}
	return nil
}

// CreateChatCompletion creates a chat completion
func (s *ChatService) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletion, error) {
	// Validate and adjust temperature and top_p
	if req.Temperature != nil {
		temp := *req.Temperature
		if temp <= 0 {
			doSample := false
			req.DoSample = &doSample
			temp = 0.01
			req.Temperature = &temp
		} else if temp >= 1 {
			temp = 0.99
			req.Temperature = &temp
		}
	}

	if req.TopP != nil {
		topP := *req.TopP
		if topP <= 0 {
			topP = 0.01
			req.TopP = &topP
		} else if topP >= 1 {
			topP = 0.99
			req.TopP = &topP
		}
	}

	var result ChatCompletion
	err := s.client.doRequest(ctx, http.MethodPost, "/chat/completions", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateChatCompletionStream creates a streaming chat completion
func (s *ChatService) CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionStream, error) {
	// Set stream to true
	stream := true
	req.Stream = &stream

	// Validate and adjust temperature and top_p
	if req.Temperature != nil {
		temp := *req.Temperature
		if temp <= 0 {
			doSample := false
			req.DoSample = &doSample
			temp = 0.01
			req.Temperature = &temp
		} else if temp >= 1 {
			temp = 0.99
			req.Temperature = &temp
		}
	}

	if req.TopP != nil {
		topP := *req.TopP
		if topP <= 0 {
			topP = 0.01
			req.TopP = &topP
		} else if topP >= 1 {
			topP = 0.99
			req.TopP = &topP
		}
	}

	url := s.client.baseURL + "/chat/completions"

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("failed to marshal request body: %v", err)}
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("failed to create request: %v", err)}
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.client.apiKey)
	httpReq.Header.Set("x-source-channel", s.client.sourceChannel)
	httpReq.Header.Set("Accept", "text/event-stream")

	// Add custom headers
	for key, value := range s.client.customHeaders {
		httpReq.Header.Set(key, value)
	}

	resp, err := s.client.httpClient.Do(httpReq)
	if err != nil {
		return nil, &APITimeoutError{Err: &Error{Message: fmt.Sprintf("request failed: %v", err)}}
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)

		var errResp struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error.Message != "" {
			return nil, NewError(resp.StatusCode, errResp.Error.Message, errResp.Error.Type, errResp.Error.Code)
		}
		return nil, NewError(resp.StatusCode, string(respBody), "", "")
	}

	return &ChatCompletionStream{
		reader:   bufio.NewReader(resp.Body),
		response: resp,
	}, nil
}

// Helper function to create a simple text message
func NewUserMessage(content string) Message {
	return Message{
		Role:    "user",
		Content: content,
	}
}

// Helper function to create a system message
func NewSystemMessage(content string) Message {
	return Message{
		Role:    "system",
		Content: content,
	}
}

// Helper function to create an assistant message
func NewAssistantMessage(content string) Message {
	return Message{
		Role:    "assistant",
		Content: content,
	}
}

// Helper function to create a multimodal message with text and image
func NewMultimodalMessage(role, text, imageURL string) Message {
	content := []ContentPart{
		{Type: "text", Text: text},
		{Type: "image_url", ImageURL: &ImageURL{URL: imageURL}},
	}
	return Message{
		Role:    role,
		Content: content,
	}
}

// Helper function to create a web search tool
func NewWebSearchTool(query string, includeResult bool) Tool {
	return Tool{
		Type: "web_search",
		WebSearch: &WebSearchTool{
			SearchQuery:  query,
			SearchResult: includeResult,
		},
	}
}

// Helper function to create a function tool
func NewFunctionTool(name, description string, parameters map[string]interface{}) Tool {
	return Tool{
		Type: "function",
		Function: &FunctionDefinition{
			Name:        name,
			Description: description,
			Parameters:  parameters,
		},
	}
}

// String returns a string representation of the completion
func (c *ChatCompletion) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ChatCompletion(id=%s, model=%s, choices=%d)\n", c.ID, c.Model, len(c.Choices)))
	for i, choice := range c.Choices {
		content := ""
		if choice.Message.Content != nil {
			content = *choice.Message.Content
		}
		sb.WriteString(fmt.Sprintf("  Choice %d: %s (finish_reason=%s)\n", i, content, choice.FinishReason))
	}
	sb.WriteString(fmt.Sprintf("  Usage: prompt=%d, completion=%d, total=%d\n",
		c.Usage.PromptTokens, c.Usage.CompletionTokens, c.Usage.TotalTokens))
	return sb.String()
}
