// core/ai/llm/ask.go
// 单次调用LLM的核心业务逻辑

package llm

import (
	// 外部包
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	ollama "github.com/cloudwego/eino-ext/components/model/ollama"
	schema "github.com/cloudwego/eino/schema"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

func AskLLMProvider(query string) {
	// 1. 验证查询是否为空
	query = strings.TrimSpace(query)
	if query == "" {
		fmt.Println(lang.T("msg.llm.ask.empty_query"))
		return
	}

	// 2. 从配置中获取当前LLM Provider
	currentProvider, err := cfg.GetCurrentLLM()
	if err != nil {
		fmt.Printf(lang.T("msg.llm.ask.get_current_failed")+"\n", err)
		return
	}
	if currentProvider == nil {
		fmt.Println(lang.T("msg.llm.ask.no_current_llm"))
		return
	}

	// 3. 创建带超时的 context（30秒超时）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	provider := currentProvider.Provider

	switch provider {
	case "ollama":
		// 调用 Ollama 模型
		cm, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
			BaseURL: currentProvider.Config.BaseURL,
			Model:   currentProvider.Model,
		})
		if err != nil {
			fmt.Printf(lang.T("msg.llm.ask.create_model_failed")+"\n", err)
			return
		}

		resp, err := cm.Stream(ctx, []*schema.Message{
			{
				Role:    schema.User,
				Content: query,
			},
		})
		if err != nil {
			// 检查是否是超时错误
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Println(lang.T("msg.llm.ask.timeout"))
			} else {
				fmt.Printf(lang.T("msg.llm.ask.stream_failed")+"\n", err)
			}
			return
		}

		// 显示 Provider 名称
		fmt.Printf("%s: ", currentProvider.Name)

		// 流式输出响应
		for {
			chunk, err := resp.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				// 检查是否是超时错误
				if ctx.Err() == context.DeadlineExceeded {
					fmt.Println("\n" + lang.T("msg.llm.ask.timeout"))
				} else {
					fmt.Printf("\n"+lang.T("msg.llm.ask.receive_error")+"\n", err)
				}
				break
			}
			fmt.Print(chunk.Content)
		}
		fmt.Println() // 换行结束

	default:
		fmt.Println(lang.T("msg.llm.ask.unsupported_provider") + ": " + provider)
		return
	}
}
