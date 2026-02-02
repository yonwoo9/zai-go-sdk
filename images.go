package zai

import (
	"context"
	"fmt"
	"net/http"
)

// ImagesService handles image generation operations
type ImagesService struct {
	client *BaseClient
}

// GeneratedImage represents a generated image
type GeneratedImage struct {
	B64JSON       *string `json:"b64_json,omitempty"`
	URL           *string `json:"url,omitempty"`
	RevisedPrompt *string `json:"revised_prompt,omitempty"`
}

// ImagesResponse represents the image generation response
type ImagesResponse struct {
	Created int64            `json:"created"`
	Data    []GeneratedImage `json:"data"`
}

// AsyncImagesResponse represents the async image generation response
type AsyncImagesResponse struct {
	ID          *string          `json:"id,omitempty"`
	Model       string           `json:"model"`
	RequestID   string           `json:"request_id"`
	TaskStatus  string           `json:"task_status"` // PROCESSING, SUCCESS, FAIL
	ImageResult []GeneratedImage `json:"image_result,omitempty"`
}

// ImageGenerationRequest represents an image generation request
type ImageGenerationRequest struct {
	Prompt             string              `json:"prompt"`
	Model              *string             `json:"model,omitempty"`
	N                  *int                `json:"n,omitempty"`
	Quality            *string             `json:"quality,omitempty"`
	ResponseFormat     *string             `json:"response_format,omitempty"`
	Size               *string             `json:"size,omitempty"`
	Style              *string             `json:"style,omitempty"`
	SensitiveWordCheck *SensitiveWordCheck `json:"sensitive_word_check,omitempty"`
	User               *string             `json:"user,omitempty"`
	UserID             *string             `json:"user_id,omitempty"`
	RequestID          *string             `json:"request_id,omitempty"`
	WatermarkEnabled   *bool               `json:"watermark_enabled,omitempty"`
}

// AsyncImageGenerationRequest represents an async image generation request
type AsyncImageGenerationRequest struct {
	Prompt           string  `json:"prompt"`
	Model            *string `json:"model,omitempty"`
	Quality          *string `json:"quality,omitempty"`
	Size             *string `json:"size,omitempty"`
	UserID           *string `json:"user_id,omitempty"`
	RequestID        *string `json:"request_id,omitempty"`
	WatermarkEnabled *bool   `json:"watermark_enabled,omitempty"`
}

// Generations generates images from text prompts
func (s *ImagesService) Generations(ctx context.Context, req *ImageGenerationRequest) (*ImagesResponse, error) {
	var result ImagesResponse
	err := s.client.doRequest(ctx, http.MethodPost, "/images/generations", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// AsyncGenerations asynchronously generates images from text prompts
// Only supports glm-image model. Use RetrieveImagesResult() to poll for the result.
func (s *ImagesService) AsyncGenerations(ctx context.Context, req *AsyncImageGenerationRequest) (*AsyncImagesResponse, error) {
	var result AsyncImagesResponse
	err := s.client.doRequest(ctx, http.MethodPost, "/async/images/generations", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// RetrieveImagesResult retrieves the result of an async image generation operation
func (s *ImagesService) RetrieveImagesResult(ctx context.Context, id string) (*AsyncImagesResponse, error) {
	if id == "" {
		return nil, &Error{Message: "id must be provided"}
	}

	var result AsyncImagesResponse
	err := s.client.doRequest(ctx, http.MethodGet, fmt.Sprintf("/async-result/%s", id), nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Helper function to create a simple image generation request
func NewImageGenerationRequest(prompt, model string) *ImageGenerationRequest {
	return &ImageGenerationRequest{
		Prompt: prompt,
		Model:  &model,
	}
}

// Helper function to create an async image generation request
func NewAsyncImageGenerationRequest(prompt, model string) *AsyncImageGenerationRequest {
	return &AsyncImageGenerationRequest{
		Prompt: prompt,
		Model:  &model,
	}
}
