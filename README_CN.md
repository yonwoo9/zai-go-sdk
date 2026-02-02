# Z.ai Open Platform Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/yonwoo9/zai-go-sdk.svg)](https://pkg.go.dev/github.com/yonwoo9/zai-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/yonwoo9/zai-go-sdk)](https://goreportcard.com/report/github.com/yonwoo9/zai-go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[English](README.md) | 中文文档

[Z.ai 开放平台](https://docs.z.ai/) 官方 Go SDK，让 Go 开发者更容易调用 Z.ai 的开放 API。

## ✨ 核心功能

### 🤖 **对话补全**

- **标准对话**: 使用包括 `glm-4.7` 在内的多种模型创建对话补全
- **流式支持**: 实时流式响应，适用于交互式应用
- **工具调用**: 函数调用能力，增强 AI 交互
- **多模态对话**: 支持视觉模型的图像理解能力

### 🧠 **向量嵌入**

- **文本嵌入**: 为文本生成高质量的向量嵌入
- **可配置维度**: 可自定义嵌入维度
- **批量处理**: 支持单次请求处理多个输入

### 🎥 **视频生成**

- **文本生成视频**: 从文本提示生成视频
- **图像生成视频**: 从图像输入创建视频
- **可自定义参数**: 控制质量、时长、帧率和尺寸
- **音频支持**: 可选的视频音频生成

### 🎨 **图像生成**

- **文本生成图像**: 从文本提示生成图像
- **异步支持**: 支持异步图像生成和轮询
- **可自定义参数**: 控制质量、尺寸和风格

## 📦 安装

### 环境要求

- **Go**: 1.21+

### 通过 go get 安装

```bash
go get github.com/yonwoo9/zai-go-sdk
```

## 🚀 快速开始

### 创建 API Key

#### 获取 API Key

- **海外地区**: 访问 [Z.ai 开放平台](https://docs.z.ai/) 获取您的 API key
- **中国大陆地区**: 访问 [智谱 AI 开放平台](https://www.bigmodel.cn/) 获取您的 API key

#### API BASE URL

- **中国大陆地区**: `https://open.bigmodel.cn/api/paas/v4/`
- **海外地区**: `https://api.z.ai/api/paas/v4/`

### 基本使用

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	// 海外用户，创建 Client
	client, err := zai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	// 中国用户，创建 ZhipuClient
	// client, err := zai.NewZhipuClient("your-api-key")

	// 创建对话补全
	response, err := client.Chat.CreateChatCompletion(context.Background(), &zai.ChatCompletionRequest{
		Model: "glm-4.7",
		Messages: []zai.Message{
			zai.NewUserMessage("你好，Z.ai！"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*response.Choices[0].Message.Content)
}
```

### 客户端配置

SDK 支持多种方式配置 API key：

#### 环境变量

```bash
export ZAI_API_KEY="your-api-key"
export ZAI_BASE_URL="https://api.z.ai/api/paas/v4/"  # 可选
```

#### 代码配置

```go
import (
	"time"
	"net/http"
	"github.com/yonwoo9/zai-go-sdk"
)

// 基本配置
client, err := zai.NewClient("your-api-key")

// 高级配置
client, err := zai.NewClient("your-api-key", &zai.ClientConfig{
	BaseURL: "https://api.z.ai/api/paas/v4/",
	HTTPClient: &http.Client{
		Timeout: 300 * time.Second,
	},
	MaxRetries: 3,
	SourceChannel: "my-app",
})

// 使用智谱域名服务
zhipuClient, err := zai.NewZhipuClient("your-api-key")
```

## 📖 使用示例

### 流式对话

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
			zai.NewSystemMessage("你是一个有帮助的助手。"),
			zai.NewUserMessage("给我讲一个关于人工智能的故事。"),
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

### 带工具调用的对话

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
			zai.NewSystemMessage("你是一个有帮助的助手。"),
			zai.NewUserMessage("什么是人工智能？"),
		},
		Tools: []zai.Tool{
			zai.NewWebSearchTool("什么是人工智能？", true),
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

### 向量嵌入

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
		zai.NewEmbeddingsRequest("embedding-3", "你好，世界！"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("嵌入维度: %d\n", len(response.Data[0].Embedding))
	fmt.Printf("前 5 个值: %v\n", response.Data[0].Embedding[:5])
}
```

### 图像生成

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
		zai.NewImageGenerationRequest("一幅美丽的山间日落景色", "cogview-3-plus"))
	if err != nil {
		log.Fatal(err)
	}

	if len(response.Data) > 0 && response.Data[0].URL != nil {
		fmt.Printf("图像 URL: %s\n", *response.Data[0].URL)
	}
}
```

### 视频生成

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

	// 生成视频
	response, err := client.Videos.Generations(context.Background(), &zai.VideoGenerationRequest{
		Model:     "cogvideox-3",
		Prompt:    zai.String("一只猫在玩球。"),
		Quality:   &quality,
		WithAudio: &withAudio,
		Size:      &size,
		FPS:       &fps,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("视频生成已启动。任务 ID: %s\n", *response.ID)

	// 轮询结果
	for {
		result, err := client.Videos.RetrieveVideosResult(context.Background(), *response.ID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("任务状态: %s\n", result.TaskStatus)

		if result.TaskStatus == "SUCCESS" {
			if len(result.VideoResult) > 0 {
				fmt.Printf("视频 URL: %s\n", result.VideoResult[0].URL)
			}
			break
		} else if result.TaskStatus == "FAIL" {
			fmt.Println("视频生成失败")
			break
		}

		time.Sleep(5 * time.Second)
	}
}
```

## 🚨 错误处理

SDK 提供了全面的错误处理：

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
			zai.NewUserMessage("你好，Z.ai！"),
		},
	})

	if err != nil {
		switch e := err.(type) {
		case *zai.APIAuthenticationError:
			fmt.Printf("认证失败: %v\n", e)
		case *zai.APIReachLimitError:
			fmt.Printf("超出速率限制: %v\n", e)
		case *zai.APITimeoutError:
			fmt.Printf("请求超时: %v\n", e)
		default:
			fmt.Printf("意外错误: %v\n", e)
		}
		return
	}

	fmt.Println(*response.Choices[0].Message.Content)
}
```

### 错误类型

| 错误类型                   | 描述                 |
| -------------------------- | -------------------- |
| `APIRequestFailedError`    | 无效的请求参数 (400) |
| `APIAuthenticationError`   | 认证失败 (401)       |
| `APIReachLimitError`       | 超出速率限制 (429)   |
| `APIInternalError`         | 内部服务器错误 (500) |
| `APIServerFlowExceedError` | 服务器过载 (503)     |
| `APITimeoutError`          | 请求超时             |
| `APIStatusError`           | 通用 API 错误        |

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。

## 📞 支持

如有问题和技术支持，请访问 [Z.ai 开放平台](https://docs.z.ai/) 或查看文档。
