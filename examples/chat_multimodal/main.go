package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/yonwoo9/zai-go-sdk"
)

func encodeImage(imagePath string) (string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func main() {
	// Create client
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	// Read and encode image
	base64Image, err := encodeImage("test_image.png")
	if err != nil {
		log.Fatal(err)
	}

	temperature := 0.5
	maxTokens := 2000

	// Create multimodal chat completion
	response, err := client.Chat.CreateChatCompletion(context.Background(), &zai.ChatCompletionRequest{
		Model: "glm-4.6v",
		Messages: []zai.Message{
			zai.NewMultimodalMessage("user", "What's in this image? Please describe it in detail.",
				fmt.Sprintf("data:image/jpeg;base64,%s", base64Image)),
		},
		Temperature: &temperature,
		MaxTokens:   &maxTokens,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Image analysis:")
	fmt.Println(*response.Choices[0].Message.Content)
}
