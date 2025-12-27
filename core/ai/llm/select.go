// core/ai/llm/select.go
// 统一的 Provider 选择交互逻辑

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

// ProviderFormatter 用于格式化显示每个 Provider 的函数类型
type ProviderFormatter func(index int, provider *cfg.LLMProvider) string

// SelectAndLinkProviders 通用的 Provider 选择和激活函数
// providers: 要显示的 Provider 列表
// title: 显示标题（例如："模型 'xxx' 匹配到多个 LLM" 或 "已注册的 LLM Provider"）
// formatter: 格式化每个 Provider 显示的函数，如果为 nil 则使用默认格式
func SelectAndLinkProviders(providers []*cfg.LLMProvider, title string, formatter ProviderFormatter) {
	if len(providers) == 0 {
		fmt.Println(lang.T("msg.llm.link.no_available"))
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
	fmt.Println("\n" + lang.T("msg.llm.link.select_prompt"))
	fmt.Printf(lang.T("msg.llm.link.select_number")+"\n", len(providers))
	fmt.Println(lang.T("msg.llm.link.select_all"))
	fmt.Println(lang.T("msg.llm.link.select_cancel"))
	fmt.Print("\n" + lang.T("msg.llm.link.select_input"))

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
		activateAllProviders(providers)
		return
	}

	// 处理数字选择
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(providers) {
		fmt.Printf(lang.T("msg.llm.link.invalid_choice")+"\n", input)
		return
	}

	// 激活选中的 Provider
	selectedProvider := providers[choice-1]
	activateSingleProvider(selectedProvider)
}

// activateAllProviders 激活所有 Provider（不设置前台）
func activateAllProviders(providers []*cfg.LLMProvider) {
	activatedCount := 0

	for _, p := range providers {
		if err := cfg.SetActiveLLMByID(p.ID); err != nil {
			fmt.Printf(lang.T("msg.llm.link.activate_failed")+"\n", p.Name, err)
		} else {
			fmt.Printf(lang.T("msg.llm.link.activate_success")+"\n", p.Name)
			activatedCount++
		}
	}

	fmt.Printf("\n"+lang.T("msg.llm.link.activate_all_count")+"\n", activatedCount)
}

// activateSingleProvider 激活单个 Provider 并设为前台
func activateSingleProvider(provider *cfg.LLMProvider) {
	// 激活选中的 Provider
	if err := cfg.SetActiveLLMByID(provider.ID); err != nil {
		fmt.Printf(lang.T("msg.llm.link.activate_failed")+"\n", err)
		return
	}
	fmt.Printf(lang.T("msg.llm.link.activate_success")+"\n", provider.Name)

	// 设置为前台
	if err := cfg.SetCurrentLLM(provider.ID); err != nil {
		fmt.Printf(lang.T("msg.llm.link.set_current_failed")+"\n", err)
		return
	}
	fmt.Println(lang.T("msg.llm.link.set_current_success"))
}

// formatProviderWithStatus 格式化显示 Provider（带状态信息）
func formatProviderWithStatus(index int, provider *cfg.LLMProvider) string {
	config := cfg.LoadConfig()

	// 检查状态
	isActive := false
	for _, activeID := range config.ActiveLLMs {
		if activeID == provider.ID {
			isActive = true
			break
		}
	}
	isCurrent := config.CurrentLLM == provider.ID

	status := ""
	if isCurrent {
		status = lang.T("msg.llm.link.status_current")
	} else if isActive {
		status = lang.T("msg.llm.link.status_active")
	} else {
		status = lang.T("msg.llm.link.status_inactive")
	}

	return fmt.Sprintf(lang.T("msg.llm.link.provider_format")+"%s",
		index, provider.Name, provider.Provider, provider.Model, status)
}

// formatProviderSimple 简单的格式化显示（用于 model 匹配场景）
func formatProviderSimple(index int, provider *cfg.LLMProvider) string {
	return fmt.Sprintf(lang.T("msg.llm.link.provider_format_simple"),
		index, provider.Name, provider.ID[:8]+"...", provider.Provider)
}

// SelectAndUnlinkProviders 通用的 Provider 选择和取消连接函数
// providers: 要显示的 Provider 列表
// title: 显示标题（例如："模型 'xxx' 匹配到多个 LLM" 或 "已注册的 LLM Provider"）
// formatter: 格式化每个 Provider 显示的函数，如果为 nil 则使用默认格式
func SelectAndUnlinkProviders(providers []*cfg.LLMProvider, title string, formatter ProviderFormatter) {
	if len(providers) == 0 {
		fmt.Println(lang.T("msg.llm.link.no_available"))
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
	fmt.Println("\n" + lang.T("msg.llm.unlink.select_prompt"))
	fmt.Printf(lang.T("msg.llm.unlink.select_number")+"\n", len(providers))
	fmt.Println(lang.T("msg.llm.unlink.select_all"))
	fmt.Println(lang.T("msg.llm.unlink.select_cancel"))
	fmt.Print("\n" + lang.T("msg.llm.unlink.select_input"))

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
		unlinkAllProviders(providers)
		return
	}

	// 处理数字选择
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(providers) {
		fmt.Printf(lang.T("msg.llm.link.invalid_choice")+"\n", input)
		return
	}

	// 取消连接选中的 Provider
	selectedProvider := providers[choice-1]
	unlinkSingleProvider(selectedProvider)
}

// unlinkAllProviders 取消连接所有 Provider
func unlinkAllProviders(providers []*cfg.LLMProvider) {
	unlinkedCount := 0

	for _, p := range providers {
		if err := cfg.UnsetActiveLLMByID(p.ID); err != nil {
			fmt.Printf(lang.T("msg.llm.unlink.unlink_failed")+"\n", p.Name, err)
		} else {
			fmt.Printf(lang.T("msg.llm.unlink.unlink_success")+"\n", p.Name)
			unlinkedCount++
		}
	}

	fmt.Printf("\n"+lang.T("msg.llm.unlink.unlink_all_count")+"\n", unlinkedCount)
}

// unlinkSingleProvider 取消连接单个 Provider
func unlinkSingleProvider(provider *cfg.LLMProvider) {
	// 取消连接选中的 Provider
	if err := cfg.UnsetActiveLLMByID(provider.ID); err != nil {
		fmt.Printf(lang.T("msg.llm.unlink.unlink_failed")+"\n", err)
		return
	}
	fmt.Printf(lang.T("msg.llm.unlink.unlink_success")+"\n", provider.Name)
}

// SelectAndSwitchProvider 通用的 Provider 选择和切换前台函数
// providers: 要显示的 Provider 列表（只包含已激活的）
// title: 显示标题
// formatter: 格式化每个 Provider 显示的函数，如果为 nil 则使用默认格式
func SelectAndSwitchProvider(providers []*cfg.LLMProvider, title string, formatter ProviderFormatter) {
	if len(providers) == 0 {
		fmt.Println(lang.T("msg.llm.switch.no_active"))
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
	fmt.Println("\n" + lang.T("msg.llm.switch.select_prompt"))
	fmt.Printf(lang.T("msg.llm.switch.select_number")+"\n", len(providers))
	fmt.Println(lang.T("msg.llm.switch.select_cancel"))
	fmt.Print("\n" + lang.T("msg.llm.switch.select_input"))

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

	// 处理数字选择
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(providers) {
		fmt.Printf(lang.T("msg.llm.link.invalid_choice")+"\n", input)
		return
	}

	// 切换到选中的 Provider
	selectedProvider := providers[choice-1]
	switchToProvider(selectedProvider)
}

// switchToProvider 切换到指定的 Provider（设为前台）
func switchToProvider(provider *cfg.LLMProvider) {
	// 检查是否已经是当前前台
	config := cfg.LoadConfig()
	if config.CurrentLLM == provider.ID {
		fmt.Printf(lang.T("msg.llm.switch.already_current")+"\n", provider.Name)
		return
	}

	// 设置为前台
	if err := cfg.SetCurrentLLM(provider.ID); err != nil {
		fmt.Printf(lang.T("msg.llm.switch.failed")+"\n", err)
		return
	}
	fmt.Printf(lang.T("msg.llm.switch.success")+"\n", provider.Name)
}
