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

	// Create embeddings for a single text
	response, err := client.Embeddings.CreateEmbeddings(context.Background(),
		zai.NewEmbeddingsRequest("embedding-3", "Hello, world!"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Model: %s\n", response.Model)
	fmt.Printf("Embedding dimension: %d\n", len(response.Data[0].Embedding))
	fmt.Printf("First 5 values: %v\n", response.Data[0].Embedding[:5])
	fmt.Printf("Usage: %d tokens\n", response.Usage.TotalTokens)

	// Create embeddings for multiple texts
	batchResponse, err := client.Embeddings.CreateEmbeddings(context.Background(),
		zai.NewBatchEmbeddingsRequest("embedding-3", []string{
			"Hello, world!",
			"Artificial intelligence is fascinating.",
			"Go is a great programming language.",
		}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nBatch embeddings created: %d\n", len(batchResponse.Data))
	for i, emb := range batchResponse.Data {
		fmt.Printf("Embedding %d dimension: %d\n", i+1, len(emb.Embedding))
	}
}
