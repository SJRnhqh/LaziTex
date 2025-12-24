// core/ai/llm/remove.go
// 删除指定LLM Provider的核心业务逻辑

package llm

import (
	// 外部包
	"fmt"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
)

// FindMatchingLLMProviders 查找匹配的LLM providers
// 支持通过 ID、Name 或 Model 进行匹配
func FindMatchingLLMProviders(providers []cfg.LLMProvider, identifier string) []cfg.LLMProvider {
	var matches []cfg.LLMProvider
	for _, provider := range providers {
		if provider.ID == identifier || provider.Name == identifier || provider.Model == identifier {
			matches = append(matches, provider)
		}
	}
	return matches
}

// FormatProviderInfo 格式化 Provider 信息为字符串（用于显示）
// 返回格式化的字符串，不直接打印，保持核心层无副作用
func FormatProviderInfo(provider cfg.LLMProvider) string {
	return fmt.Sprintf("ID: %s, Name: %s, Provider: %s, Model: %s",
		provider.ID, provider.Name, provider.Provider, provider.Model)
}

// FormatProviderList 格式化多个 Provider 为列表字符串（用于选择菜单）
// 返回格式化的字符串列表，每个元素包含索引和 Provider 信息
func FormatProviderList(providers []cfg.LLMProvider) []string {
	formatted := make([]string, len(providers))
	for i, p := range providers {
		formatted[i] = fmt.Sprintf("[%d] %s (ID: %s, Provider: %s, Model: %s)",
			i+1, p.Name, p.ID, p.Provider, p.Model)
	}
	return formatted
}

// ExtractProviderID 从格式化的列表项中提取 Provider ID
// 用于用户选择后获取对应的 ID
func ExtractProviderID(providers []cfg.LLMProvider, index int) string {
	if index < 0 || index >= len(providers) {
		return ""
	}
	return providers[index].ID
}
