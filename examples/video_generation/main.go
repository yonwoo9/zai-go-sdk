package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	// Create client
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	quality := "quality"
	withAudio := true
	size := "1920x1080"
	fps := 30

	// Generate video from text
	response, err := client.Videos.Generations(context.Background(), &zai.VideoGenerationRequest{
		Model:     "cogvideox-3",
		Prompt:    zai.String("A cat is playing with a ball in a sunny garden."),
		Quality:   &quality,
		WithAudio: &withAudio,
		Size:      &size,
		FPS:       &fps,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Video generation started. Task ID: %s\n", *response.ID)
	fmt.Println("Polling for result...")

	// Poll for result
	for {
		result, err := client.Videos.RetrieveVideosResult(context.Background(), *response.ID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Task status: %s\n", result.TaskStatus)

		if result.TaskStatus == "SUCCESS" {
			fmt.Println("\nVideo generation completed!")
			if len(result.VideoResult) > 0 {
				fmt.Printf("Video URL: %s\n", result.VideoResult[0].URL)
				fmt.Printf("Cover Image URL: %s\n", result.VideoResult[0].CoverImageURL)
			}
			break
		} else if result.TaskStatus == "FAIL" {
			fmt.Println("Video generation failed")
			break
		}

		time.Sleep(5 * time.Second)
	}
}
