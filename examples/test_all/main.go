// test_all 通过环境变量 ZAI_API_KEY 测试所有已实现的 API 接口。
// 在项目根目录运行: ZAI_API_KEY=your-key go run ./examples/test_all
// 注: Embeddings/图像 可能因账号余额或资源包返回 429(code 1113)，属业务限制非 SDK 错误。
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/yonwoo9/zai-go-sdk"
)

func main() {
	apiKey := os.Getenv("ZAI_API_KEY")
	if apiKey == "" {
		log.Fatal("ZAI_API_KEY environment variable is required")
	}

	client, err := zai.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ok, fail := 0, 0

	// 1. Chat 同步
	fmt.Println("=== 1. Chat 同步 (CreateChatCompletion) ===")
	chatResp, err := client.Chat.CreateChatCompletion(ctx, &zai.ChatCompletionRequest{
		Model: "glm-4-flash",
		Messages: []zai.Message{
			zai.NewUserMessage("用一句话介绍你自己。"),
		},
	})
	if err != nil {
		fmt.Printf("FAIL: %v\n\n", err)
		fail++
	} else {
		fmt.Printf("OK: %s\n\n", *chatResp.Choices[0].Message.Content)
		ok++
	}

	// 2. Chat 流式
	fmt.Println("=== 2. Chat 流式 (CreateChatCompletionStream) ===")
	stream, err := client.Chat.CreateChatCompletionStream(ctx, &zai.ChatCompletionRequest{
		Model: "glm-4-flash",
		Messages: []zai.Message{
			zai.NewUserMessage("说三个字：你好世界"),
		},
	})
	if err != nil {
		fmt.Printf("FAIL: %v\n\n", err)
		fail++
	} else {
		defer stream.Close()
		var text string
		for {
			chunk, err := stream.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("FAIL (stream): %v\n\n", err)
				fail++
				break
			}
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != nil {
				text += *chunk.Choices[0].Delta.Content
			}
		}
		if text != "" {
			fmt.Printf("OK: %s\n\n", text)
			ok++
		}
	}

	// 3. Embeddings
	fmt.Println("=== 3. Embeddings (CreateEmbeddings) ===")
	embResp, err := client.Embeddings.CreateEmbeddings(ctx, zai.NewEmbeddingsRequest("embedding-3", "测试文本"))
	if err != nil {
		fmt.Printf("FAIL: %v\n\n", err)
		fail++
	} else {
		fmt.Printf("OK: model=%s, dim=%d\n\n", embResp.Model, len(embResp.Data[0].Embedding))
		ok++
	}

	// 4. 图像生成（同步）
	fmt.Println("=== 4. 图像生成 (Images.Generations) ===")
	imgResp, err := client.Images.Generations(ctx, zai.NewImageGenerationRequest("一只可爱的猫", "glm-image"))
	if err != nil {
		fmt.Printf("FAIL: %v\n\n", err)
		fail++
	} else {
		if len(imgResp.Data) > 0 && imgResp.Data[0].URL != nil {
			fmt.Printf("OK: url=%s\n\n", *imgResp.Data[0].URL)
		} else {
			fmt.Printf("OK: created=%d\n\n", imgResp.Created)
		}
		ok++
	}

	// 5. 视频生成（提交任务 + 轮询一次，不长时间等待）
	fmt.Println("=== 5. 视频生成 (Videos.Generations + RetrieveVideosResult) ===")
	quality := "speed"
	vidResp, err := client.Videos.Generations(ctx, &zai.VideoGenerationRequest{
		Model:   "cogvideox-3",
		Prompt:  zai.String("一只猫在草地上奔跑"),
		Quality: &quality,
	})
	if err != nil {
		fmt.Printf("FAIL: %v\n\n", err)
		fail++
	} else {
		fmt.Printf("OK: task submitted, id=%s\n", *vidResp.ID)
		ok++
		// 轮询一次看状态
		time.Sleep(3 * time.Second)
		result, err := client.Videos.RetrieveVideosResult(ctx, *vidResp.ID)
		if err != nil {
			fmt.Printf("  RetrieveVideosResult: %v\n", err)
		} else {
			fmt.Printf("  status=%s\n", result.TaskStatus)
		}
		fmt.Println()
	}

	fmt.Printf("=== 结果: %d 通过, %d 失败 ===\n", ok, fail)
	if fail > 0 {
		os.Exit(1)
	}
}
