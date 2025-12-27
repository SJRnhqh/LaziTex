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

func UnlinkLLM(sub []string) {
	llm.UnlinkLLMProvider(sub[1:])
}

func SwitchLLM(sub []string) {
	llm.SwitchLLMProvider(sub[1:])
}

func TestLLM(sub []string) {
	llm.TestLLMProvider(sub[1:])
}

func RemoveLLM(sub []string) {
	llm.RemoveLLMProvider(sub[1:])
}

func ListLLM() {
	// Mock 列出所有注册的LLM的业务逻辑
}

func AskLLM() {
	// Mock 与LLM进行交互的业务逻辑
}

func ChatLLM() {
	// Mock 与LLM进行持续性对话交互的业务逻辑
}
