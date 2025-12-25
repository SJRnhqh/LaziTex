package llm

import (
	"errors"
	"fmt"

	cfg "github.com/SJRnhqh/lazitex/config"
)

// Call 调用当前 LLM（单轮非流式）
func Call(llmprovider *cfg.LLMProvider, prompt string) (string, error) {
	if llmprovider == nil {
		return "", errors.New("LLM 配置为空")
	}
	if llmprovider.Model == "" {
		return "", errors.New("未选择任何 LLM 模型")
	}
	if prompt == "" {
		return "", errors.New("prompt 不能为空")
	}

	provider := GetProvider(llmprovider.Provider)
	if provider == nil {
		return "", fmt.Errorf("不支持的 LLM Provider: %s", llmprovider.Provider)
	}
	return provider.Generate(llmprovider.Model, prompt, llmprovider.Config)
}
