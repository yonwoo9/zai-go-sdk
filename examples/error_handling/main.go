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

	// Example of handling different error types
	response, err := client.Chat.CreateChatCompletion(context.Background(), &zai.ChatCompletionRequest{
		Model: "glm-4.7",
		Messages: []zai.Message{
			zai.NewUserMessage("Hello, Z.ai!"),
		},
	})

	if err != nil {
		switch e := err.(type) {
		case *zai.APIAuthenticationError:
			fmt.Printf("Authentication failed: %v\n", e)
			fmt.Println("Please check your API key.")
		case *zai.APIReachLimitError:
			fmt.Printf("Rate limit exceeded: %v\n", e)
			fmt.Println("Please wait before making more requests.")
		case *zai.APITimeoutError:
			fmt.Printf("Request timeout: %v\n", e)
			fmt.Println("Please try again later.")
		case *zai.APIRequestFailedError:
			fmt.Printf("Invalid request: %v\n", e)
			fmt.Println("Please check your request parameters.")
		case *zai.APIInternalError:
			fmt.Printf("Internal server error: %v\n", e)
			fmt.Println("Please try again later.")
		default:
			fmt.Printf("Unexpected error: %v\n", e)
		}
		return
	}

	fmt.Println("Success!")
	fmt.Println(*response.Choices[0].Message.Content)
}
