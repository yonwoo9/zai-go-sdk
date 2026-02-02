package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	// Create client (will use ZAI_API_KEY environment variable if not provided)
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	// Create a simple chat completion
	response, err := client.Chat.CreateChatCompletion(context.Background(), &zai.ChatCompletionRequest{
		Model: "glm-4.7",
		Messages: []zai.Message{
			zai.NewUserMessage("Hello, Z.ai! Please introduce yourself."),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:")
	fmt.Println(*response.Choices[0].Message.Content)
	fmt.Printf("\nUsage: %d tokens\n", response.Usage.TotalTokens)
}
