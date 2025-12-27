// core/ai/llm/remove.go
// 注销LLMProvider的核心业务逻辑

package llm

import (
	// 外部包
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// RemoveLLMProvider 删除 LLM Provider 的核心业务逻辑
func RemoveLLMProvider(llmProviderIDOrNameOrModel []string) {
	// 如果没有参数，列出所有已注册的 Provider 供选择删除
	if len(llmProviderIDOrNameOrModel) == 0 {
		listAndRemoveProviders()
		return
	}

	// 单个参数
	if len(llmProviderIDOrNameOrModel) == 1 {
		removeSingle(llmProviderIDOrNameOrModel[0])
		return
	}

	// 批量模式：统计后询问是否全部删除
	removeBatch(llmProviderIDOrNameOrModel)
}

// removeSingle 删除单个指定的 LLM Provider
func removeSingle(idOrNameOrModel string) {
	// 查找所有匹配的 providers
	matches, matchType := cfg.FindAllMatchingProviders(idOrNameOrModel)

	if len(matches) == 0 {
		fmt.Printf(lang.T("msg.llm.remove.not_found")+"\n", idOrNameOrModel)
		return
	}

	// 如果是唯一匹配，直接询问确认
	if len(matches) == 1 {
		if confirmRemove(matches[0]) {
			removeProviderByID(matches[0].ID)
		}
		return
	}

	// 多个匹配（Model 匹配），显示交互式选择
	if matchType == "model" {
		handleMultipleRemoveMatches(matches, idOrNameOrModel)
	}
}

// handleMultipleRemoveMatches 处理多个匹配的 Model 删除
func handleMultipleRemoveMatches(matches []*cfg.LLMProvider, modelName string) {
	title := fmt.Sprintf(lang.T("msg.llm.remove.model_multiple_title"), modelName)
	SelectAndRemoveProviders(matches, title, formatProviderSimple)
}

// listAndRemoveProviders 列出所有已注册的 Provider 并支持选择删除
func listAndRemoveProviders() {
	allProviders := cfg.GetAllLLMProviders()

	if len(allProviders) == 0 {
		fmt.Println(lang.T("msg.llm.remove.no_providers"))
		return
	}

	SelectAndRemoveProviders(allProviders, lang.T("msg.llm.remove.registered_title"), formatProviderWithStatus)
}

// removeBatch 批量删除多个 Provider
func removeBatch(providers []string) {
	// 统计所有匹配的 Provider
	var allMatches []*cfg.LLMProvider
	matchedNames := make(map[string]bool)

	for _, provider := range providers {
		matches, _ := cfg.FindAllMatchingProviders(provider)
		for _, match := range matches {
			// 避免重复（同一个 Provider 可能被多个参数匹配到）
			if !matchedNames[match.ID] {
				allMatches = append(allMatches, match)
				matchedNames[match.ID] = true
			}
		}
	}

	if len(allMatches) == 0 {
		fmt.Println(lang.T("msg.llm.remove.batch_no_matches"))
		return
	}

	// 显示将要删除的 Provider 列表
	fmt.Printf("\n"+lang.T("msg.llm.remove.batch_list")+"\n", len(allMatches))
	for i, p := range allMatches {
		fmt.Printf("  %d. %s (%s:%s)\n", i+1, p.Name, p.Provider, p.Model)
	}

	// 询问是否全部删除
	if confirmBatchRemove(len(allMatches)) {
		removedCount := 0
		for _, p := range allMatches {
			if removeProviderByID(p.ID) {
				removedCount++
			}
		}
		fmt.Printf("\n"+lang.T("msg.llm.remove.batch_success")+"\n", removedCount, len(allMatches))
	} else {
		fmt.Println(lang.T("msg.llm.remove.cancelled"))
	}
}

// removeProviderByID 通过 ID 删除 Provider（内部函数）
func removeProviderByID(providerID string) bool {
	removed, err := cfg.RemoveLLMProvider(providerID)
	if err != nil {
		fmt.Printf(lang.T("msg.llm.remove.failed")+"\n", providerID, err)
		return false
	}
	if !removed {
		fmt.Printf(lang.T("msg.llm.remove.not_found")+"\n", providerID)
		return false
	}
	return true
}

// confirmRemove 确认是否删除单个 Provider
func confirmRemove(provider *cfg.LLMProvider) bool {
	// 显示 Provider 信息
	fmt.Printf("\n"+lang.T("msg.llm.remove.provider_info")+"\n",
		provider.Name, provider.Provider, provider.Model)

	// 检查状态并提示
	config := cfg.LoadConfig()
	isActive := false
	for _, activeID := range config.ActiveLLMs {
		if activeID == provider.ID {
			isActive = true
			break
		}
	}
	isCurrent := config.CurrentLLM == provider.ID

	if isCurrent {
		fmt.Println(lang.T("msg.llm.remove.warning_current"))
	} else if isActive {
		fmt.Println(lang.T("msg.llm.remove.warning_active"))
	}

	// 询问确认
	return confirmYesNo(lang.T("msg.llm.remove.confirm_prompt"))
}

// confirmBatchRemove 确认是否批量删除
func confirmBatchRemove(count int) bool {
	return confirmYesNo(fmt.Sprintf(lang.T("msg.llm.remove.batch_confirm"), count))
}

// SelectAndRemoveProviders 通用的 Provider 选择和删除函数
// providers: 要显示的 Provider 列表
// title: 显示标题
// formatter: 格式化每个 Provider 显示的函数，如果为 nil 则使用默认格式
func SelectAndRemoveProviders(providers []*cfg.LLMProvider, title string, formatter ProviderFormatter) {
	if len(providers) == 0 {
		fmt.Println(lang.T("msg.llm.remove.no_providers"))
		return
	}

	// 显示标题
	fmt.Printf("\n%s：\n\n", title)

	// 显示列表
	for i, p := range providers {
		if formatter != nil {
			fmt.Println(formatter(i+1, p))
		} else {
			// 默认格式
			fmt.Printf(lang.T("msg.llm.link.provider_format")+"\n",
				i+1, p.Name, p.Provider, p.Model)
		}
	}

	// 显示选择提示
	fmt.Println("\n" + lang.T("msg.llm.remove.select_prompt"))
	fmt.Printf(lang.T("msg.llm.remove.select_number")+"\n", len(providers))
	fmt.Println(lang.T("msg.llm.remove.select_all"))
	fmt.Println(lang.T("msg.llm.remove.select_cancel"))
	fmt.Print("\n" + lang.T("msg.llm.remove.select_input"))

	// 读取用户输入
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf(lang.T("msg.llm.link.read_input_failed")+"\n", err)
		return
	}

	input = strings.TrimSpace(strings.ToLower(input))

	// 处理取消
	if input == "q" || input == "quit" {
		fmt.Println(lang.T("msg.llm.remove.cancelled"))
		return
	}

	// 处理全选
	if input == "a" || input == "all" {
		removeAllProviders(providers)
		return
	}

	// 处理数字选择
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(providers) {
		fmt.Printf(lang.T("msg.llm.link.invalid_choice")+"\n", input)
		return
	}

	// 删除选中的 Provider（需要二次确认）
	selectedProvider := providers[choice-1]
	if confirmRemove(selectedProvider) {
		if removeProviderByID(selectedProvider.ID) {
			fmt.Printf(lang.T("msg.llm.remove.success")+"\n", selectedProvider.Name)
		}
	} else {
		fmt.Println(lang.T("msg.llm.remove.cancelled"))
	}
}

// removeAllProviders 删除所有选中的 Provider
func removeAllProviders(providers []*cfg.LLMProvider) {
	// 显示将要删除的 Provider 列表
	fmt.Printf("\n"+lang.T("msg.llm.remove.batch_list")+"\n", len(providers))
	for i, p := range providers {
		fmt.Printf("  %d. %s (%s:%s)\n", i+1, p.Name, p.Provider, p.Model)
	}

	// 询问是否全部删除
	if confirmBatchRemove(len(providers)) {
		removedCount := 0
		for _, p := range providers {
			if removeProviderByID(p.ID) {
				removedCount++
				fmt.Printf(lang.T("msg.llm.remove.success")+"\n", p.Name)
			}
		}
		fmt.Printf("\n"+lang.T("msg.llm.remove.batch_success")+"\n", removedCount, len(providers))
	} else {
		fmt.Println(lang.T("msg.llm.remove.cancelled"))
	}
}
