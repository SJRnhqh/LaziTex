// core/ai/llm/client.go
// LLM ChatModel 客户端接口抽象

package llm

import (
	// 外部包
	"context"
	"fmt"

	ollama "github.com/cloudwego/eino-ext/components/model/ollama"
	schema "github.com/cloudwego/eino/schema"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// StreamResult 流式响应结果接口
type StreamResult interface {
	Recv() (*schema.Message, error)
}

// ChatModelClient LLM ChatModel 客户端接口
// 用于抽象不同 provider 的 ChatModel 实现
type ChatModelClient interface {
	// Stream 流式调用模型
	// ctx: 上下文（包含超时控制）
	// messages: 消息列表
	// 返回: 流式响应结果和错误
	Stream(ctx context.Context, messages []*schema.Message) (StreamResult, error)
}

// ollamaChatModelWrapper Ollama ChatModel 包装器
// 实现 ChatModelClient 接口
type ollamaChatModelWrapper struct {
	model *ollama.ChatModel
}

// Stream 实现 ChatModelClient 接口
func (w *ollamaChatModelWrapper) Stream(ctx context.Context, messages []*schema.Message) (StreamResult, error) {
	return w.model.Stream(ctx, messages)
}

// NewChatModelClient 创建 ChatModel 客户端（工厂函数）
// 根据 provider 类型返回相应的 ChatModelClient 实现
func NewChatModelClient(ctx context.Context, provider *cfg.LLMProvider) (ChatModelClient, error) {
	switch provider.Provider {
	case "ollama":
		cm, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
			BaseURL: provider.Config.BaseURL,
			Model:   provider.Model,
		})
		if err != nil {
			return nil, fmt.Errorf(lang.T("msg.llm.client.create_failed")+": %w", err)
		}
		return &ollamaChatModelWrapper{model: cm}, nil
	default:
		return nil, fmt.Errorf(lang.T("msg.llm.client.unsupported_provider"), provider.Provider)
	}
}
