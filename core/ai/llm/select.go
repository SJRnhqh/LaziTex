// core/ai/llm/select.go
// LLM Provider 智能匹配和选择逻辑
// 支持通过 ID、Name 或 Model 进行匹配，并提供交互式选择功能
// 主要用于已注册的 LLM Provider 的查找和选择

package llm

import (
	"errors"
	"fmt"

	cfg "github.com/SJRnhqh/lazitex/config"
)

// FindMatchingProviders 查找匹配的 LLM providers
// 支持通过 ID、Name 或 Model 进行匹配（针对已注册的 Provider）
func FindMatchingProviders(providers []cfg.LLMProvider, identifier string) []cfg.LLMProvider {
	var matches []cfg.LLMProvider
	for _, provider := range providers {
		if provider.ID == identifier ||
			provider.Name == identifier ||
			provider.Model == identifier {
			matches = append(matches, provider)
		}
	}
	return matches
}

// SelectProviderResult 选择结果
type SelectProviderResult struct {
	Provider *cfg.LLMProvider
	Error    error
}

// SelectProviderOptions 选择选项
type SelectProviderOptions struct {
	// MultipleMatchesPrompt 多个匹配时的提示信息（支持 %d 占位符）
	MultipleMatchesPrompt string
	// SelectPrompt 选择提示信息
	SelectPrompt string
	// OnSelect 选择函数，接收匹配数量和提示信息，返回选中的索引（0-based），-1 表示取消
	OnSelect func(count int, selectPrompt string) int
}

// FindOrSelectProvider 查找或选择 Provider（统一入口）
//
// 逻辑：
// - ID 和 Name 是唯一的，如果匹配到直接返回
// - Model 可能匹配多个，需要调用 OnSelect 让用户选择
// - 如果未找到，返回错误
//
// 参数：
//   - providers: 所有已注册的 Provider 列表
//   - identifier: 标识符（ID、Name 或 Model）
//   - options: 选择选项（仅在多个匹配时使用）
//
// 返回：
//   - 选中的 Provider 指针，如果未找到或取消则返回 nil
//   - 错误信息
func FindOrSelectProvider(
	providers []cfg.LLMProvider,
	identifier string,
	options *SelectProviderOptions,
) (*cfg.LLMProvider, error) {
	// 查找所有匹配项
	matches := FindMatchingProviders(providers, identifier)

	// 未找到
	if len(matches) == 0 {
		return nil, fmt.Errorf("no matching LLM Provider found: %s", identifier)
	}

	// 单个匹配：直接返回（ID 和 Name 是唯一的，Model 如果只有一个也直接返回）
	if len(matches) == 1 {
		return &matches[0], nil
	}

	// 多个匹配：需要用户选择（通常是 Model 匹配到多个）
	if options == nil || options.OnSelect == nil {
		return nil, fmt.Errorf("found %d matching LLM Providers, but no selection function provided", len(matches))
	}

	// 显示匹配列表
	if options.MultipleMatchesPrompt != "" {
		fmt.Printf(options.MultipleMatchesPrompt+"\n", len(matches))
	} else {
		fmt.Printf("Found %d matching LLM Providers, please select:\n", len(matches))
	}

	formattedList := FormatProviderList(matches)
	for _, item := range formattedList {
		fmt.Println("  " + item)
	}

	// 调用选择函数
	selectPrompt := options.SelectPrompt
	if selectPrompt == "" {
		selectPrompt = "Please select (enter number, or 'q' to cancel): "
	}

	selectedIndex := options.OnSelect(len(matches), selectPrompt)
	if selectedIndex < 0 {
		return nil, errors.New("cancelled")
	}

	// 验证索引有效性
	if selectedIndex >= len(matches) {
		return nil, fmt.Errorf("invalid selection index: %d", selectedIndex)
	}

	return &matches[selectedIndex], nil
}

// FormatProviderInfo 格式化单个 Provider 信息（用于显示）
// 返回格式化的字符串，不直接打印，保持核心层无副作用
func FormatProviderInfo(provider cfg.LLMProvider) string {
	return fmt.Sprintf("ID: %s, Name: %s, Provider: %s, Model: %s",
		provider.ID, provider.Name, provider.Provider, provider.Model)
}

// FormatProviderList 格式化多个 Provider 为列表（用于选择菜单）
// 返回格式化的字符串列表，每个元素包含索引和 Provider 信息
func FormatProviderList(providers []cfg.LLMProvider) []string {
	formatted := make([]string, len(providers))
	for i, p := range providers {
		formatted[i] = fmt.Sprintf("[%d] %s (ID: %s, Provider: %s, Model: %s)",
			i+1, p.Name, p.ID, p.Provider, p.Model)
	}
	return formatted
}

// ExtractProviderID 从匹配列表中提取指定索引的 Provider ID
// 用于用户选择后获取对应的 ID
func ExtractProviderID(providers []cfg.LLMProvider, index int) string {
	if index < 0 || index >= len(providers) {
		return ""
	}
	return providers[index].ID
}
