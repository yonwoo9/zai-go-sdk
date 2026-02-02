package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	// Create client
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	// Create a streaming chat completion
	stream, err := client.Chat.CreateChatCompletionStream(context.Background(), &zai.ChatCompletionRequest{
		Model: "glm-4.7",
		Messages: []zai.Message{
			zai.NewSystemMessage("You are a helpful assistant."),
			zai.NewUserMessage("Tell me a short story about artificial intelligence."),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	fmt.Println("Streaming response:")
	for {
		chunk, err := stream.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != nil {
			fmt.Print(*chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
}
