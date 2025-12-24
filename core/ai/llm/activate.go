// core/ai/llm/active.go
// LLM 激活状态管理
// 负责管理当前激活的 LLM，并提供状态显示功能
// 为 ask/chat 功能提供基础支持

package llm

import (
	"fmt"

	cfg "github.com/SJRnhqh/lazitex/config"
)

// SetCurrentLLM 设置当前激活的 LLM（仅限已 link 的 LLM）
//
// 参数：
//   - providerID: LLM Provider 的 ID
//
// 返回：
//   - error: 如果 LLM 未 link 或不存在，返回错误
func SetCurrentLLM(providerID string) error {
	config := cfg.LoadConfig()

	// 验证是否在激活列表中（已 link）
	found := false
	for _, activeID := range config.ActiveLLMs {
		if activeID == providerID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("LLM Provider '%s' is not linked", providerID)
	}

	// 查找 provider 信息，验证是否存在
	var provider *cfg.LLMProvider
	for i := range config.LLMProviders {
		if config.LLMProviders[i].ID == providerID {
			provider = &config.LLMProviders[i]
			break
		}
	}

	if provider == nil {
		return fmt.Errorf("LLM Provider '%s' not found", providerID)
	}

	// 更新当前激活的 LLM
	config.CurrentLLM = providerID
	return cfg.SaveConfig(config)
}

// GetCurrentLLM 获取当前激活的 LLM 信息
//
// 返回：
//   - *cfg.LLMProvider: 当前激活的 LLM Provider，如果未设置则返回 nil
//   - string: 当前激活的 LLM ID，如果未设置则返回空字符串
func GetCurrentLLM() (*cfg.LLMProvider, string) {
	config := cfg.LoadConfig()

	if config.CurrentLLM == "" {
		return nil, ""
	}

	// 查找 provider
	for i := range config.LLMProviders {
		if config.LLMProviders[i].ID == config.CurrentLLM {
			return &config.LLMProviders[i], config.CurrentLLM
		}
	}

	// 如果配置中的 CurrentLLM 指向的 Provider 不存在，返回 nil
	return nil, ""
}

// FormatActivePrompt 格式化激活状态提示（用于显示在 prompt 中）
//
// 参数：
//   - provider: LLM Provider 信息，如果为 nil 则返回空字符串
//   - mode: 模式标识，"m" 表示 model（LLM），"a" 表示 agent（未来扩展）
//
// 返回：
//   - 格式化的字符串，如 "(my-llm:m)" 或 "(my-agent:a)"，如果 provider 为 nil 则返回空字符串
func FormatActivePrompt(provider *cfg.LLMProvider, mode string) string {
	if provider == nil {
		return ""
	}

	// mode: "m" 表示 model, "a" 表示 agent
	return fmt.Sprintf("(%s:%s)", provider.Name, mode)
}

// GetLinkedLLMs 获取所有已 link 的 LLM 列表
//
// 返回：
//   - []cfg.LLMProvider: 所有已 link 的 LLM Provider 列表
func GetLinkedLLMs() []cfg.LLMProvider {
	config := cfg.LoadConfig()

	var linked []cfg.LLMProvider
	activeMap := make(map[string]bool)
	for _, id := range config.ActiveLLMs {
		activeMap[id] = true
	}

	// 遍历所有 Provider，找出已 link 的
	for _, provider := range config.LLMProviders {
		if activeMap[provider.ID] {
			linked = append(linked, provider)
		}
	}

	return linked
}

// IsCurrentLLMLinked 检查当前激活的 LLM 是否仍在 link 列表中
//
// 返回：
//   - bool: 如果当前 LLM 已 link 则返回 true，否则返回 false
func IsCurrentLLMLinked() bool {
	config := cfg.LoadConfig()

	if config.CurrentLLM == "" {
		return false
	}

	// 检查是否在激活列表中
	for _, activeID := range config.ActiveLLMs {
		if activeID == config.CurrentLLM {
			return true
		}
	}

	return false
}
