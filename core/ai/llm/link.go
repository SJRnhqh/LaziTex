// core/ai/llm/link.go
// 对于已经注册LLMProvider的link核心业务逻辑

package llm

import (
	// 外部包
	"fmt"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

func LinkLLMProvider(llmProviderIDOrNameOrModel []string) {
	if len(llmProviderIDOrNameOrModel) == 1 {
		linkSingle(llmProviderIDOrNameOrModel[0])
		return
	}

	// 批量模式
	for _, llmProvider := range llmProviderIDOrNameOrModel {
		// 先检查是否已激活
		matches, _ := cfg.FindAllMatchingProviders(llmProvider)
		if len(matches) > 0 {
			config := cfg.LoadConfig()
			isActive := false
			for _, activeID := range config.ActiveLLMs {
				if activeID == matches[0].ID {
					isActive = true
					break
				}
			}
			if isActive {
				fmt.Printf(lang.T("msg.llm.link.batch_already_active")+"\n", llmProvider)
				continue
			}
		}

		if err := cfg.SetActiveLLM(llmProvider); err != nil {
			if err.Error() == "MULTIPLE_MATCHES" {
				fmt.Printf(lang.T("msg.llm.link.batch_multiple_skip")+"\n", llmProvider)
				continue
			}
			fmt.Printf(lang.T("msg.llm.link.batch_failed")+"\n", llmProvider, err)
			continue
		}
		fmt.Printf(lang.T("msg.llm.link.success")+"\n", llmProvider)
	}
}

func UnlinkLLMProvider(llmProviderIDOrNameOrModel []string) {
	// 如果没有参数，列出所有已注册的 Provider 供选择 unlink
	if len(llmProviderIDOrNameOrModel) == 0 {
		// 优先尝试 unlink 当前前台的 Provider
		currentLLM, err := cfg.GetCurrentLLM()
		if err == nil && currentLLM != nil {
			// 有当前前台，直接 unlink
			if err := cfg.UnsetActiveLLMByID(currentLLM.ID); err != nil {
				fmt.Printf(lang.T("msg.llm.unlink.failed")+"\n", err)
				return
			}
			fmt.Printf(lang.T("msg.llm.unlink.success")+"\n", currentLLM.Name)
			return
		}

		// 没有当前前台，列出所有已激活的 Provider 供选择
		listAndUnlinkProviders()
		return
	}

	if len(llmProviderIDOrNameOrModel) == 1 {
		unlinkSingle(llmProviderIDOrNameOrModel[0])
		return
	}

	// 批量模式
	for _, llmProvider := range llmProviderIDOrNameOrModel {
		if err := cfg.UnsetActiveLLM(llmProvider); err != nil {
			// 检查是否是多重匹配错误
			matches, matchType := cfg.FindAllMatchingProviders(llmProvider)
			if matchType == "model" && len(matches) > 1 {
				fmt.Printf(lang.T("msg.llm.unlink.batch_multiple_skip")+"\n", llmProvider)
				continue
			}
			fmt.Printf(lang.T("msg.llm.unlink.batch_failed")+"\n", llmProvider, err)
			continue
		}
		fmt.Printf(lang.T("msg.llm.unlink.success")+"\n", llmProvider)
	}
}

func linkSingle(idOrNameOrModel string) {
	// 查找所有匹配的 providers
	matches, matchType := cfg.FindAllMatchingProviders(idOrNameOrModel)

	if len(matches) == 0 {
		fmt.Printf(lang.T("msg.llm.link.not_found")+"\n", idOrNameOrModel)
		return
	}

	// 如果是唯一匹配，直接处理
	if len(matches) == 1 {
		provider := matches[0]
		config := cfg.LoadConfig()

		// 检查是否已激活
		isActive := false
		for _, activeID := range config.ActiveLLMs {
			if activeID == provider.ID {
				isActive = true
				break
			}
		}

		// 如果未激活，先激活
		if !isActive {
			if err := cfg.SetActiveLLM(idOrNameOrModel); err != nil {
				fmt.Printf(lang.T("msg.llm.link.set_active_failed")+"\n", err)
				return
			}
			fmt.Printf(lang.T("msg.llm.link.activated")+"\n", idOrNameOrModel)
		} else {
			fmt.Printf(lang.T("msg.llm.link.already_active")+"\n", idOrNameOrModel)
		}

		// 检查是否已经是当前前台
		if config.CurrentLLM == provider.ID {
			fmt.Printf(lang.T("msg.llm.link.already_current")+"\n", idOrNameOrModel)
			return
		}

		// 设置为前台
		if err := cfg.SetCurrentLLM(idOrNameOrModel); err != nil {
			fmt.Printf(lang.T("msg.llm.link.set_current_failed")+"\n", err)
			return
		}
		fmt.Println(lang.T("msg.llm.link.set_current_success"))
		return
	}

	// 多个匹配（Model 匹配），显示交互式选择
	if matchType == "model" {
		handleMultipleMatches(matches, idOrNameOrModel)
	}
}

func handleMultipleMatches(matches []*cfg.LLMProvider, modelName string) {
	title := fmt.Sprintf(lang.T("msg.llm.link.model_multiple_title"), modelName)
	SelectAndLinkProviders(matches, title, formatProviderSimple)
}

func unlinkSingle(idOrNameOrModel string) {
	// 查找所有匹配的 providers
	matches, matchType := cfg.FindAllMatchingProviders(idOrNameOrModel)

	if len(matches) == 0 {
		fmt.Printf(lang.T("msg.llm.unlink.not_found")+"\n", idOrNameOrModel)
		return
	}

	// 如果是唯一匹配，直接处理
	if len(matches) == 1 {
		if err := cfg.UnsetActiveLLM(idOrNameOrModel); err != nil {
			fmt.Printf(lang.T("msg.llm.unlink.failed")+"\n", err)
			return
		}
		fmt.Printf(lang.T("msg.llm.unlink.success")+"\n", idOrNameOrModel)
		return
	}

	// 多个匹配（Model 匹配），显示交互式选择
	if matchType == "model" {
		handleMultipleUnlinks(matches, idOrNameOrModel)
	}
}

func handleMultipleUnlinks(matches []*cfg.LLMProvider, modelName string) {
	title := fmt.Sprintf(lang.T("msg.llm.unlink.model_multiple_title"), modelName)
	SelectAndUnlinkProviders(matches, title, formatProviderSimple)
}

// listAndUnlinkProviders 列出所有已激活的 Provider 并支持选择 unlink
func listAndUnlinkProviders() {
	activeProviders := cfg.GetActiveLLMProviders()

	if len(activeProviders) == 0 {
		fmt.Println(lang.T("msg.llm.unlink.no_active_providers"))
		return
	}

	SelectAndUnlinkProviders(activeProviders, lang.T("msg.llm.unlink.active_title"), formatProviderWithStatus)
}

// SwitchLLMProvider 切换当前前台的 LLM Provider
func SwitchLLMProvider(llmProviderIDOrNameOrModel []string) {
	// 如果没有参数，列出所有已激活的 Provider 供选择
	if len(llmProviderIDOrNameOrModel) == 0 {
		listAndSwitchProviders()
		return
	}

	// 只支持单个参数（不支持批量切换）
	if len(llmProviderIDOrNameOrModel) == 1 {
		switchSingle(llmProviderIDOrNameOrModel[0])
		return
	}

	// 批量切换没有意义，提示错误
	fmt.Println(lang.T("msg.llm.switch.batch_not_supported"))
}

// switchSingle 切换到单个指定的 LLM Provider
func switchSingle(idOrNameOrModel string) {
	// 查找所有匹配的 providers
	matches, matchType := cfg.FindAllMatchingProviders(idOrNameOrModel)

	if len(matches) == 0 {
		fmt.Printf(lang.T("msg.llm.switch.not_found")+"\n", idOrNameOrModel)
		return
	}

	// 过滤出已激活的 providers
	activeMatches := filterActiveProviders(matches)
	if len(activeMatches) == 0 {
		fmt.Printf(lang.T("msg.llm.switch.not_active")+"\n", idOrNameOrModel)
		return
	}

	// 如果是唯一匹配，直接切换
	if len(activeMatches) == 1 {
		provider := activeMatches[0]
		config := cfg.LoadConfig()
		// 检查是否已经是当前前台
		if config.CurrentLLM == provider.ID {
			fmt.Printf(lang.T("msg.llm.switch.already_current")+"\n", provider.Name)
			return
		}
		if err := cfg.SetCurrentLLM(provider.ID); err != nil {
			fmt.Printf(lang.T("msg.llm.switch.failed")+"\n", err)
			return
		}
		fmt.Printf(lang.T("msg.llm.switch.success")+"\n", provider.Name)
		return
	}

	// 多个匹配（Model 匹配），显示交互式选择
	if matchType == "model" {
		handleMultipleSwitchMatches(activeMatches, idOrNameOrModel)
	}
}

// handleMultipleSwitchMatches 处理多个匹配的 Model 切换
func handleMultipleSwitchMatches(matches []*cfg.LLMProvider, modelName string) {
	title := fmt.Sprintf(lang.T("msg.llm.switch.model_multiple_title"), modelName)
	SelectAndSwitchProvider(matches, title, formatProviderSimple)
}

// listAndSwitchProviders 列出所有已激活的 Provider 并支持选择切换
func listAndSwitchProviders() {
	activeProviders := cfg.GetActiveLLMProviders()

	if len(activeProviders) == 0 {
		fmt.Println(lang.T("msg.llm.switch.no_active"))
		return
	}

	SelectAndSwitchProvider(activeProviders, lang.T("msg.llm.switch.active_title"), formatProviderWithStatus)
}

// filterActiveProviders 过滤出已激活的 providers
func filterActiveProviders(providers []*cfg.LLMProvider) []*cfg.LLMProvider {
	config := cfg.LoadConfig()
	var activeProviders []*cfg.LLMProvider

	for _, provider := range providers {
		for _, activeID := range config.ActiveLLMs {
			if provider.ID == activeID {
				activeProviders = append(activeProviders, provider)
				break
			}
		}
	}

	return activeProviders
}
