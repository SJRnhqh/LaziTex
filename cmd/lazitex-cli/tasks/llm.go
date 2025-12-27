// tasks/llm.go
// 处理LLM管理业务的编排

package tasks

func LinkLLM(sub []string) {
	switch len(sub) {
	case 1:
		// TODO:新注册的业务逻辑
	default:
		// TODO:处理单个或者多个已经注册LLM的连接 -> 确定连接个数 -> 单个连接 -> 错误处理
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
