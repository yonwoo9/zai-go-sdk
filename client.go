package zai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// DefaultMaxRetries is the default maximum number of retries
	DefaultMaxRetries = 2
	// DefaultTimeout is the default request timeout
	DefaultTimeout = 300 * time.Second
	// ZaiBaseURL is the default base URL for overseas regions
	ZaiBaseURL = "https://api.z.ai/api/paas/v4"
	// ZhipuAiBaseURL is the default base URL for mainland China regions
	ZhipuAiBaseURL = "https://open.bigmodel.cn/api/paas/v4"
)

// ClientConfig holds the configuration for the ZAI client
type ClientConfig struct {
	APIKey             string
	BaseURL            string
	HTTPClient         *http.Client
	MaxRetries         int
	DisableTokenCache  bool
	SourceChannel      string
	CustomHeaders      map[string]string
}

// BaseClient is the base client for ZAI API
type BaseClient struct {
	apiKey            string
	baseURL           string
	httpClient        *http.Client
	maxRetries        int
	disableTokenCache bool
	sourceChannel     string
	customHeaders     map[string]string
}

// Client is the main client for ZAI API (overseas regions)
type Client struct {
	*BaseClient
	Chat       *ChatService
	Embeddings *EmbeddingsService
	Images     *ImagesService
	Videos     *VideosService
	Audio      *AudioService
	Files      *FilesService
}

// ZhipuClient is the client for Zhipu AI (mainland China regions)
type ZhipuClient struct {
	*BaseClient
	Chat       *ChatService
	Embeddings *EmbeddingsService
	Images     *ImagesService
	Videos     *VideosService
	Audio      *AudioService
	Files      *FilesService
}

// NewClient creates a new ZAI client for overseas regions
func NewClient(apiKey string, config ...*ClientConfig) (*Client, error) {
	cfg := getConfig(apiKey, ZaiBaseURL, config...)
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("ZAI_API_KEY")
	}
	if cfg.APIKey == "" {
		return nil, &Error{Message: "api_key not provided, please provide it through parameters or environment variables"}
	}

	baseClient := newBaseClient(cfg)
	client := &Client{
		BaseClient: baseClient,
	}

	// Initialize services
	client.Chat = &ChatService{client: baseClient}
	client.Embeddings = &EmbeddingsService{client: baseClient}
	client.Images = &ImagesService{client: baseClient}
	client.Videos = &VideosService{client: baseClient}
	client.Audio = &AudioService{client: baseClient}
	client.Files = &FilesService{client: baseClient}

	return client, nil
}

// NewZhipuClient creates a new Zhipu AI client for mainland China regions
func NewZhipuClient(apiKey string, config ...*ClientConfig) (*ZhipuClient, error) {
	cfg := getConfig(apiKey, ZhipuAiBaseURL, config...)
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("ZAI_API_KEY")
	}
	if cfg.APIKey == "" {
		return nil, &Error{Message: "api_key not provided, please provide it through parameters or environment variables"}
	}

	baseClient := newBaseClient(cfg)
	client := &ZhipuClient{
		BaseClient: baseClient,
	}

	// Initialize services
	client.Chat = &ChatService{client: baseClient}
	client.Embeddings = &EmbeddingsService{client: baseClient}
	client.Images = &ImagesService{client: baseClient}
	client.Videos = &VideosService{client: baseClient}
	client.Audio = &AudioService{client: baseClient}
	client.Files = &FilesService{client: baseClient}

	return client, nil
}

func getConfig(apiKey, defaultBaseURL string, config ...*ClientConfig) *ClientConfig {
	if len(config) > 0 && config[0] != nil {
		cfg := config[0]
		if cfg.APIKey == "" {
			cfg.APIKey = apiKey
		}
		if cfg.BaseURL == "" {
			cfg.BaseURL = defaultBaseURL
		}
		if cfg.HTTPClient == nil {
			cfg.HTTPClient = &http.Client{Timeout: DefaultTimeout}
		}
		if cfg.MaxRetries == 0 {
			cfg.MaxRetries = DefaultMaxRetries
		}
		if cfg.SourceChannel == "" {
			cfg.SourceChannel = "go-sdk"
		}
		return cfg
	}

	return &ClientConfig{
		APIKey:        apiKey,
		BaseURL:       defaultBaseURL,
		HTTPClient:    &http.Client{Timeout: DefaultTimeout},
		MaxRetries:    DefaultMaxRetries,
		SourceChannel: "go-sdk",
	}
}

func newBaseClient(cfg *ClientConfig) *BaseClient {
	if cfg.BaseURL == "" {
		if envURL := os.Getenv("ZAI_BASE_URL"); envURL != "" {
			cfg.BaseURL = envURL
		}
	}

	return &BaseClient{
		apiKey:            cfg.APIKey,
		baseURL:           cfg.BaseURL,
		httpClient:        cfg.HTTPClient,
		maxRetries:        cfg.MaxRetries,
		disableTokenCache: cfg.DisableTokenCache,
		sourceChannel:     cfg.SourceChannel,
		customHeaders:     cfg.CustomHeaders,
	}
}

// doRequest performs an HTTP request with retry logic
func (c *BaseClient) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt) * time.Second
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		err := c.doRequestOnce(ctx, method, path, body, result)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on certain errors
		switch err.(type) {
		case *APIAuthenticationError, *APIRequestFailedError:
			return err
		}
	}

	return lastErr
}

// doRequestOnce performs a single HTTP request
func (c *BaseClient) doRequestOnce(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return &Error{Message: fmt.Sprintf("failed to marshal request body: %v", err)}
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return &Error{Message: fmt.Sprintf("failed to create request: %v", err)}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("x-source-channel", c.sourceChannel)

	// Add custom headers
	for key, value := range c.customHeaders {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &APITimeoutError{Err: &Error{Message: fmt.Sprintf("request failed: %v", err)}}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Error{Message: fmt.Sprintf("failed to read response body: %v", err)}
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		var errResp struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error.Message != "" {
			return NewError(resp.StatusCode, errResp.Error.Message, errResp.Error.Type, errResp.Error.Code)
		}
		return NewError(resp.StatusCode, string(respBody), "", "")
	}

	// Parse response
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return &Error{Message: fmt.Sprintf("failed to unmarshal response: %v", err)}
		}
	}

	return nil
}
