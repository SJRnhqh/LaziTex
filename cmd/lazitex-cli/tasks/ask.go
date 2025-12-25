// cmd/lazitex-cli/tasks/ask.go
// 与LLM单次交互的任务编排

package tasks

import (
	// 外部包
	"fmt"

	// 内部包
	llm "github.com/SJRnhqh/lazitex/core/ai/llm"
)

// 单次询问当前激活的 LLM
func AskLLM(prompt string) {
	currentLLMProvider, _ := llm.GetCurrentLLM()
	if currentLLMProvider == nil {
		fmt.Println("当前未调用任何LLM到前台")
		return
	}

	response, err := llm.Call(currentLLMProvider, prompt)
	if err != nil {
		fmt.Printf("LLM 调用失败: %v\n", err)
		return
	}
	fmt.Println(response)
}
