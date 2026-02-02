package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	// Create client
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	temperature := 0.5
	maxTokens := 2000

	// Create chat completion with web search tool
	response, err := client.Chat.CreateChatCompletion(context.Background(), &zai.ChatCompletionRequest{
		Model: "glm-4.7",
		Messages: []zai.Message{
			zai.NewSystemMessage("You are a helpful assistant."),
			zai.NewUserMessage("What are the latest developments in artificial intelligence in 2024?"),
		},
		Tools: []zai.Tool{
			zai.NewWebSearchTool("latest AI developments 2024", true),
		},
		Temperature: &temperature,
		MaxTokens:   &maxTokens,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response with web search:")
	fmt.Println(*response.Choices[0].Message.Content)
	fmt.Printf("\nUsage: %d tokens\n", response.Usage.TotalTokens)
}
