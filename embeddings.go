package zai

import (
	"context"
	"net/http"
)

// EmbeddingsService handles embeddings operations
type EmbeddingsService struct {
	client *BaseClient
}

// Embedding represents a single embedding vector
type Embedding struct {
	Object    string    `json:"object"`
	Index     *int      `json:"index,omitempty"`
	Embedding []float64 `json:"embedding"`
}

// EmbeddingsResponse represents the embeddings API response
type EmbeddingsResponse struct {
	Object string             `json:"object"`
	Data   []Embedding        `json:"data"`
	Model  string             `json:"model"`
	Usage  CompletionUsage    `json:"usage"`
}

// EmbeddingsRequest represents an embeddings request
type EmbeddingsRequest struct {
	Input              interface{}         `json:"input"` // string, []string, []int, or [][]int
	Model              string              `json:"model"`
	Dimensions         *int                `json:"dimensions,omitempty"`
	EncodingFormat     *string             `json:"encoding_format,omitempty"`
	User               *string             `json:"user,omitempty"`
	RequestID          *string             `json:"request_id,omitempty"`
	SensitiveWordCheck *SensitiveWordCheck `json:"sensitive_word_check,omitempty"`
}

// CreateEmbeddings creates embeddings for the given input
func (s *EmbeddingsService) CreateEmbeddings(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error) {
	var result EmbeddingsResponse
	err := s.client.doRequest(ctx, http.MethodPost, "/embeddings", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Helper function to create an embeddings request with a single text input
func NewEmbeddingsRequest(model, text string) *EmbeddingsRequest {
	return &EmbeddingsRequest{
		Model: model,
		Input: text,
	}
}

// Helper function to create an embeddings request with multiple text inputs
func NewBatchEmbeddingsRequest(model string, texts []string) *EmbeddingsRequest {
	return &EmbeddingsRequest{
		Model: model,
		Input: texts,
	}
}
