// core/ai/llm/test.go
// 测试LLMProvider的核心业务逻辑

package llm

import (
	// 外部包
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// TestLLMProvider 测试 LLM Provider 的核心业务逻辑
func TestLLMProvider(llmProviderIDOrNameOrModel []string) {
	// 如果没有参数，优先测试当前前台的 Provider
	if len(llmProviderIDOrNameOrModel) == 0 {
		// 优先尝试测试当前前台的 Provider
		currentLLM, err := cfg.GetCurrentLLM()
		if err == nil && currentLLM != nil {
			// 有当前前台，直接测试
			testProvider(currentLLM)
			return
		}

		// 没有当前前台，列出所有已注册的 Provider 供选择
		listAndTestProviders()
		return
	}

	// 单个参数
	if len(llmProviderIDOrNameOrModel) == 1 {
		testSingle(llmProviderIDOrNameOrModel[0])
		return
	}

	// 批量模式
	for _, llmProvider := range llmProviderIDOrNameOrModel {
		testSingle(llmProvider)
	}
}

// testSingle 测试单个指定的 LLM Provider
func testSingle(idOrNameOrModel string) {
	// 查找所有匹配的 providers
	matches, matchType := cfg.FindAllMatchingProviders(idOrNameOrModel)

	if len(matches) == 0 {
		fmt.Printf(lang.T("msg.llm.test.not_found")+"\n", idOrNameOrModel)
		return
	}

	// 如果是唯一匹配，直接测试
	if len(matches) == 1 {
		testProvider(matches[0])
		return
	}

	// 多个匹配（Model 匹配），显示交互式选择
	if matchType == "model" {
		handleMultipleTestMatches(matches, idOrNameOrModel)
	}
}

// handleMultipleTestMatches 处理多个匹配的 Model 测试
func handleMultipleTestMatches(matches []*cfg.LLMProvider, modelName string) {
	title := fmt.Sprintf(lang.T("msg.llm.test.model_multiple_title"), modelName)
	SelectAndTestProviders(matches, title, formatProviderSimple)
}

// listAndTestProviders 列出所有已注册的 Provider 并支持选择测试
func listAndTestProviders() {
	allProviders := cfg.GetAllLLMProviders()

	if len(allProviders) == 0 {
		fmt.Println(lang.T("msg.llm.test.no_providers"))
		return
	}

	SelectAndTestProviders(allProviders, lang.T("msg.llm.test.registered_title"), formatProviderWithStatus)
}

// testProvider 测试指定的 Provider（根据类型调用不同的测试方法）
func testProvider(provider *cfg.LLMProvider) {
	var result string
	var err error

	// 根据 provider 类型调用不同的测试方法
	switch strings.ToLower(provider.Provider) {
	case "ollama":
		result, err = TestOllamaConnection(provider)
	default:
		fmt.Printf(lang.T("msg.llm.test.unsupported_provider")+"\n", provider.Provider)
		return
	}

	if err != nil {
		// 测试失败，更新验证状态为 false
		updateProviderVerifiedStatus(provider.ID, false, "")
		fmt.Printf(lang.T("msg.llm.test.test_failed")+"\n", provider.Name, err)
		return
	}

	// 测试成功，更新验证状态
	updateProviderVerifiedStatus(provider.ID, true, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf(lang.T("msg.llm.test.test_success")+"\n", provider.Name, result)
}

// updateProviderVerifiedStatus 更新 Provider 的验证状态并持久化
func updateProviderVerifiedStatus(providerID string, verified bool, verifiedAt string) {
	config := cfg.LoadConfig()

	// 查找并更新 Provider
	for i := range config.LLMProviders {
		if config.LLMProviders[i].ID == providerID {
			config.LLMProviders[i].Verified = verified
			if verifiedAt != "" {
				config.LLMProviders[i].VerifiedAt = verifiedAt
			} else {
				config.LLMProviders[i].VerifiedAt = ""
			}
			break
		}
	}

	// 保存配置
	if err := cfg.SaveConfig(config); err != nil {
		fmt.Printf(lang.T("msg.llm.test.update_config_failed")+"\n", err)
	}
}

// SelectAndTestProviders 通用的 Provider 选择和测试函数
// providers: 要显示的 Provider 列表
// title: 显示标题
// formatter: 格式化每个 Provider 显示的函数，如果为 nil 则使用默认格式
func SelectAndTestProviders(providers []*cfg.LLMProvider, title string, formatter ProviderFormatter) {
	if len(providers) == 0 {
		fmt.Println(lang.T("msg.llm.test.no_providers"))
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
	fmt.Println("\n" + lang.T("msg.llm.test.select_prompt"))
	fmt.Printf(lang.T("msg.llm.test.select_number")+"\n", len(providers))
	fmt.Println(lang.T("msg.llm.test.select_all"))
	fmt.Println(lang.T("msg.llm.test.select_cancel"))
	fmt.Print("\n" + lang.T("msg.llm.test.select_input"))

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
		fmt.Println(lang.T("msg.llm.cancelled"))
		return
	}

	// 处理全选
	if input == "a" || input == "all" {
		testAllProviders(providers)
		return
	}

	// 处理数字选择
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(providers) {
		fmt.Printf(lang.T("msg.llm.link.invalid_choice")+"\n", input)
		return
	}

	// 测试选中的 Provider
	selectedProvider := providers[choice-1]
	testProvider(selectedProvider)
}

// testAllProviders 测试所有 Provider
func testAllProviders(providers []*cfg.LLMProvider) {
	testedCount := 0
	successCount := 0

	for _, p := range providers {
		var result string
		var err error

		// 根据 provider 类型调用不同的测试方法
		switch strings.ToLower(p.Provider) {
		case "ollama":
			result, err = TestOllamaConnection(p)
		default:
			fmt.Printf(lang.T("msg.llm.test.unsupported_provider")+"\n", p.Provider)
			continue
		}

		if err != nil {
			// 测试失败，更新验证状态为 false
			updateProviderVerifiedStatus(p.ID, false, "")
			fmt.Printf(lang.T("msg.llm.test.test_failed")+"\n", p.Name, err)
		} else {
			// 测试成功，更新验证状态
			updateProviderVerifiedStatus(p.ID, true, time.Now().Format("2006-01-02 15:04:05"))
			fmt.Printf(lang.T("msg.llm.test.test_success")+"\n", p.Name, result)
			successCount++
		}
		testedCount++
	}

	fmt.Printf("\n"+lang.T("msg.llm.test.test_all_count")+"\n", testedCount, successCount)
}
