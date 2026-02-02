package zai

import (
	"context"
	"fmt"
	"net/http"
)

// VideosService handles video generation operations
type VideosService struct {
	client *BaseClient
}

// VideoResult represents a video generation result
type VideoResult struct {
	URL           string `json:"url"`
	CoverImageURL string `json:"cover_image_url"`
}

// VideoObject represents a video generation object
type VideoObject struct {
	ID          *string       `json:"id,omitempty"`
	Model       string        `json:"model"`
	VideoResult []VideoResult `json:"video_result"`
	TaskStatus  string        `json:"task_status"` // PROCESSING, SUCCESS, FAIL
	RequestID   string        `json:"request_id"`
}

// VideoGenerationRequest represents a video generation request
type VideoGenerationRequest struct {
	Model              string              `json:"model"`
	Prompt             *string             `json:"prompt,omitempty"`
	ImageURL           interface{}         `json:"image_url,omitempty"` // string, []string, or map
	Quality            *string             `json:"quality,omitempty"`   // "quality" or "speed"
	WithAudio          *bool               `json:"with_audio,omitempty"`
	Size               *string             `json:"size,omitempty"`
	Duration           *int                `json:"duration,omitempty"`
	FPS                *int                `json:"fps,omitempty"`
	Style              *string             `json:"style,omitempty"`
	AspectRatio        *string             `json:"aspect_ratio,omitempty"`
	OffPeak            *bool               `json:"off_peak,omitempty"`
	MovementAmplitude  *string             `json:"movement_amplitude,omitempty"`
	SensitiveWordCheck *SensitiveWordCheck `json:"sensitive_word_check,omitempty"`
	RequestID          *string             `json:"request_id,omitempty"`
	UserID             *string             `json:"user_id,omitempty"`
	WatermarkEnabled   *bool               `json:"watermark_enabled,omitempty"`
}

// Generations generates videos from text prompts or images
func (s *VideosService) Generations(ctx context.Context, req *VideoGenerationRequest) (*VideoObject, error) {
	if req.Model == "" {
		return nil, &Error{Message: "model must be provided"}
	}

	var result VideoObject
	err := s.client.doRequest(ctx, http.MethodPost, "/videos/generations", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// RetrieveVideosResult retrieves the result of a video generation operation
func (s *VideosService) RetrieveVideosResult(ctx context.Context, id string) (*VideoObject, error) {
	if id == "" {
		return nil, &Error{Message: "id must be provided"}
	}

	var result VideoObject
	err := s.client.doRequest(ctx, http.MethodGet, fmt.Sprintf("/async-result/%s", id), nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Helper function to create a text-to-video request
func NewTextToVideoRequest(model, prompt string) *VideoGenerationRequest {
	return &VideoGenerationRequest{
		Model:  model,
		Prompt: &prompt,
	}
}

// Helper function to create an image-to-video request
func NewImageToVideoRequest(model, imageURL string) *VideoGenerationRequest {
	return &VideoGenerationRequest{
		Model:    model,
		ImageURL: imageURL,
	}
}
