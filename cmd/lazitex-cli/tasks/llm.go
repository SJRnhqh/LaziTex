// tasks/llm.go
// 处理LLM管理业务的编排

package tasks

import (
	// 内部包
	llm "github.com/SJRnhqh/lazitex/core/ai/llm"
)

func LinkLLM(sub []string) {
	switch len(sub) {
	case 1:
		llm.RegisterLLMProvider()
	default:
		llm.LinkLLMProvider(sub[1:])
	}
}

func UnlinkLLM() {
	// Mock 断掉LLM连接的业务逻辑
}

func SwitchLLM() {
	// Mock 切换LLM的业务逻辑
}

func TestLLM() {
	// Mock 测试LLM是否正常工作的业务逻辑
}

func ListLLM() {
	// Mock 列出所有注册的LLM的业务逻辑
}

func RemoveLLM() {
	// Mock 移除LLM注册的业务逻辑
}

func AskLLM() {
	// Mock 与LLM进行交互的业务逻辑
}

func ChatLLM() {
	// Mock 与LLM进行持续性对话交互的业务逻辑
}
