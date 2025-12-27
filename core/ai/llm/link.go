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
