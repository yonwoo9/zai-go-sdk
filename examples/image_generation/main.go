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

	// Generate image
	response, err := client.Images.Generations(context.Background(),
		zai.NewImageGenerationRequest("A beautiful sunset over mountains with a lake in the foreground", "glm-image"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Image generated successfully!")
	if len(response.Data) > 0 && response.Data[0].URL != nil {
		fmt.Printf("Image URL: %s\n", *response.Data[0].URL)
	}
}
