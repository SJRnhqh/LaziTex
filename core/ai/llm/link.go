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
		if err := cfg.SetActiveLLM(idOrNameOrModel); err != nil {
			fmt.Printf(lang.T("msg.llm.link.set_active_failed")+"\n", err)
			return
		}
		fmt.Printf(lang.T("msg.llm.link.success")+"\n", idOrNameOrModel)

		if err := cfg.SetCurrentLLM(idOrNameOrModel); err != nil {
			fmt.Printf(lang.T("msg.llm.link.set_current_failed")+"\n", err)
			return
		}
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

// listAndUnlinkProviders 列出所有已注册的 Provider 并支持选择 unlink
func listAndUnlinkProviders() {
	config := cfg.LoadConfig()
	allProviders := make([]*cfg.LLMProvider, len(config.LLMProviders))
	for i := range config.LLMProviders {
		allProviders[i] = &config.LLMProviders[i]
	}

	if len(allProviders) == 0 {
		fmt.Println(lang.T("msg.llm.unlink.no_providers"))
		return
	}

	SelectAndUnlinkProviders(allProviders, lang.T("msg.llm.unlink.registered_title"), formatProviderWithStatus)
}
