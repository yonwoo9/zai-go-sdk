# Z.ai Open Platform Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/yonwoo9/zai-go-sdk.svg)](https://pkg.go.dev/github.com/yonwoo9/zai-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/yonwoo9/zai-go-sdk)](https://goreportcard.com/report/github.com/yonwoo9/zai-go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[中文文档](README_CN.md) | English

The official Go SDK for [Z.ai Open Platform](https://docs.z.ai/), making it easier for Go developers to call Z.ai's open APIs.

## ✨ Core Features

### 🤖 **Chat Completions**

- **Standard Chat**: Create chat completions with various models including `glm-4.7`
- **Streaming Support**: Real-time streaming responses for interactive applications
- **Tool Calling**: Function calling capabilities for enhanced AI interactions
- **Multimodal Chat**: Image understanding capabilities with vision models

### 🧠 **Embeddings**

- **Text Embeddings**: Generate high-quality vector embeddings for text
- **Configurable Dimensions**: Customizable embedding dimensions
- **Batch Processing**: Support for multiple inputs in a single request

### 🎥 **Video Generation**

- **Text-to-Video**: Generate videos from text prompts
- **Image-to-Video**: Create videos from image inputs
- **Customizable Parameters**: Control quality, duration, FPS, and size
- **Audio Support**: Optional audio generation for videos

### 🎨 **Image Generation**

- **Text-to-Image**: Generate images from text prompts
- **Async Support**: Asynchronous image generation with polling
- **Customizable Parameters**: Control quality, size, and style

## 📦 Installation

### Requirements

- **Go**: 1.21+

### Install via go get

```bash
go get github.com/yonwoo9/zai-go-sdk
```

## 🚀 Quick Start

### Create API Key

#### Get API Key

- **Overseas regions**: Visit [Z.ai Open Platform](https://docs.z.ai/) to get your API key
- **Mainland China regions**: Visit [Zhipu AI Open Platform](https://www.bigmodel.cn/) to get your API key

#### API BASE URL

- **Mainland China regions**: `https://open.bigmodel.cn/api/paas/v4/`
- **Overseas regions**: `https://api.z.ai/api/paas/v4/`

### Basic Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	// For Overseas users, create the Client
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	// For Chinese users, create the ZhipuClient
	// client, err := zai.NewZhipuClient("your-api-key")

	// Create chat completion
	response, err := client.Chat.CreateChatCompletion(context.Background(), &zai.ChatCompletionRequest{
		Model: "glm-4.7",
		Messages: []zai.Message{
			zai.NewUserMessage("Hello, Z.ai!"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*response.Choices[0].Message.Content)
}
```

### Client Configuration

The SDK supports multiple ways to configure API keys:

#### Environment Variables

```bash
export ZAI_API_KEY="your-api-key"
export ZAI_BASE_URL="https://api.z.ai/api/paas/v4/"  # Optional
```

#### Code Configuration

```go
import (
	"time"
	"net/http"
	"github.com/yonwoo9/zai-go-sdk"
)

// Basic configuration
client, err := zai.NewClient("your-api-key")

// Advanced configuration
client, err := zai.NewClient("your-api-key", &zai.ClientConfig{
	BaseURL: "https://api.z.ai/api/paas/v4/",
	HTTPClient: &http.Client{
		Timeout: 300 * time.Second,
	},
	MaxRetries: 3,
	SourceChannel: "my-app",
})

// For Zhipu's domain service
zhipuClient, err := zai.NewZhipuClient("your-api-key")
```

## 📖 Usage Examples

### Streaming Chat

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	stream, err := client.Chat.CreateChatCompletionStream(context.Background(), &zai.ChatCompletionRequest{
		Model: "glm-4.7",
		Messages: []zai.Message{
			zai.NewSystemMessage("You are a helpful assistant."),
			zai.NewUserMessage("Tell me a story about AI."),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

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
}
```

### Chat With Tool Call

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	temperature := 0.5
	maxTokens := 2000

	response, err := client.Chat.CreateChatCompletion(context.Background(), &zai.ChatCompletionRequest{
		Model: "glm-4.7",
		Messages: []zai.Message{
			zai.NewSystemMessage("You are a helpful assistant."),
			zai.NewUserMessage("What is artificial intelligence?"),
		},
		Tools: []zai.Tool{
			zai.NewWebSearchTool("What is artificial intelligence?", true),
		},
		Temperature: &temperature,
		MaxTokens:   &maxTokens,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
```

### Multimodal Chat

```go
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/yonwoo9/zai-go-sdk"
)

func encodeImage(imagePath string) (string, error) {
	data, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func main() {
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	base64Image, err := encodeImage("test_image.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	temperature := 0.5
	maxTokens := 2000

	response, err := client.Chat.CreateChatCompletion(context.Background(), &zai.ChatCompletionRequest{
		Model: "glm-4.6v",
		Messages: []zai.Message{
			zai.NewMultimodalMessage("user", "What's in this image?",
				fmt.Sprintf("data:image/jpeg;base64,%s", base64Image)),
		},
		Temperature: &temperature,
		MaxTokens:   &maxTokens,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
```

### Embeddings

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Embeddings.CreateEmbeddings(context.Background(),
		zai.NewEmbeddingsRequest("embedding-3", "Hello, world!"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Embedding dimension: %d\n", len(response.Data[0].Embedding))
	fmt.Printf("First 5 values: %v\n", response.Data[0].Embedding[:5])
}
```

### Image Generation

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Images.Generations(context.Background(),
		zai.NewImageGenerationRequest("A beautiful sunset over mountains", "cogview-3-plus"))
	if err != nil {
		log.Fatal(err)
	}

	if len(response.Data) > 0 && response.Data[0].URL != nil {
		fmt.Printf("Image URL: %s\n", *response.Data[0].URL)
	}
}
```

### Video Generation

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	quality := "quality"
	withAudio := true
	size := "1920x1080"
	fps := 30

	// Generate video
	response, err := client.Videos.Generations(context.Background(), &zai.VideoGenerationRequest{
		Model:     "cogvideox-3",
		Prompt:    zai.String("A cat is playing with a ball."),
		Quality:   &quality,
		WithAudio: &withAudio,
		Size:      &size,
		FPS:       &fps,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Video generation started. Task ID: %s\n", *response.ID)

	// Poll for result
	for {
		result, err := client.Videos.RetrieveVideosResult(context.Background(), *response.ID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Task status: %s\n", result.TaskStatus)

		if result.TaskStatus == "SUCCESS" {
			if len(result.VideoResult) > 0 {
				fmt.Printf("Video URL: %s\n", result.VideoResult[0].URL)
			}
			break
		} else if result.TaskStatus == "FAIL" {
			fmt.Println("Video generation failed")
			break
		}

		time.Sleep(5 * time.Second)
	}
}
```

## 🚨 Error Handling

The SDK provides comprehensive error handling:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

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
		case *zai.APIReachLimitError:
			fmt.Printf("Rate limit exceeded: %v\n", e)
		case *zai.APITimeoutError:
			fmt.Printf("Request timeout: %v\n", e)
		default:
			fmt.Printf("Unexpected error: %v\n", e)
		}
		return
	}

	fmt.Println(*response.Choices[0].Message.Content)
}
```

### Error Types

| Error Type                 | Description                      |
| -------------------------- | -------------------------------- |
| `APIRequestFailedError`    | Invalid request parameters (400) |
| `APIAuthenticationError`   | Authentication failed (401)      |
| `APIReachLimitError`       | Rate limit exceeded (429)        |
| `APIInternalError`         | Internal server error (500)      |
| `APIServerFlowExceedError` | Server overloaded (503)          |
| `APITimeoutError`          | Request timeout                  |
| `APIStatusError`           | General API error                |

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📞 Support

For questions and technical support, please visit [Z.ai Open Platform](https://docs.z.ai/) or check documentation.
