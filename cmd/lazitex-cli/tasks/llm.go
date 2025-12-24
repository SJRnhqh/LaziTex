// cmd/lazitex-cli/tasks/llm.go
// LLM 管理相关业务

package tasks

import (
	// 外部包
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
	llm "github.com/SJRnhqh/lazitex/core/ai/llm"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// 显式注册（不自动连接）
func AddLLM(provider, name, model string) {
	// 创建 Provider 配置
	llmProvider := cfg.NewLLMProvider(provider, name, model)

	// 调用核心层注册函数
	result, err := llm.RegisterLLMProvider(llmProvider, nil)
	if err != nil {
		// 根据错误类型提供更友好的错误消息
		if errors.Is(err, llm.ErrConnectivityTest) || errors.Is(err, llm.ErrLLMUnavailable) {
			fmt.Printf(lang.T("msg.llm.add.test_failed")+"\n", err)
		} else if errors.Is(err, cfg.ErrProviderIDExists) || errors.Is(err, cfg.ErrProviderNameExists) {
			fmt.Printf(lang.T("msg.llm.add.register_failed")+"\n", err)
		} else {
			fmt.Printf(lang.T("msg.llm.add.register_failed")+"\n", err)
		}
		return
	}

	// 注册成功，显示成功消息
	if result.Verified {
		fmt.Printf(lang.T("msg.llm.add.test_success")+"\n", result.Provider.Name, result.Provider.Model)
	}
	fmt.Printf(lang.T("msg.llm.add.registered")+"\n", result.Provider.Name, result.Provider.Model)
}

// LinkLLM 使用统一函数
func LinkLLM(identifier string) {
	config := cfg.LoadConfig()

	p, err := llm.FindOrSelectProvider(
		config.LLMProviders,
		identifier,
		&llm.SelectProviderOptions{
			MultipleMatchesPrompt: lang.T("msg.llm.link.multiple_matches"),
			SelectPrompt:          lang.T("msg.llm.select_prompt"),
			OnSelect:              selectProviderIndex,
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "cancelled") || strings.Contains(err.Error(), "已取消") {
			fmt.Println(lang.T("msg.llm.cancelled"))
		} else {
			fmt.Printf(lang.T("msg.llm.link.not_found")+"\n", identifier)
		}
		return
	}

	// 后续逻辑保持不变
	available, testErr := llm.LLMConnectivity(p)
	if testErr != nil {
		fmt.Printf(lang.T("msg.llm.link.test_failed")+"\n", testErr)
		return
	}
	if !available {
		fmt.Printf(lang.T("msg.llm.link.unavailable")+"\n", p.Name, p.Model)
		return
	}

	if err := cfg.SetActiveLLM(p.ID); err != nil {
		fmt.Printf(lang.T("msg.llm.link.set_active_failed")+"\n", err)
		return
	}
	fmt.Printf(lang.T("msg.llm.link.success")+"\n", p.Name, p.Model)
}

// UnlinkLLM 使用统一函数
func UnlinkLLM(identifier string) {
	config := cfg.LoadConfig()

	provider, err := llm.FindOrSelectProvider(
		config.LLMProviders,
		identifier,
		&llm.SelectProviderOptions{
			MultipleMatchesPrompt: lang.T("msg.llm.unlink.multiple_matches"),
			SelectPrompt:          lang.T("msg.llm.select_prompt"),
			OnSelect:              selectProviderIndex,
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "cancelled") || strings.Contains(err.Error(), "已取消") {
			fmt.Println(lang.T("msg.llm.cancelled"))
		} else {
			fmt.Printf(lang.T("msg.llm.unlink.not_found")+"\n", identifier)
		}
		return
	}

	// 后续逻辑保持不变
	if err := cfg.UnsetActiveLLM(provider.ID); err != nil {
		if err.Error() == fmt.Sprintf("LLM Provider '%s' 不在激活列表中", provider.ID) {
			fmt.Printf(lang.T("msg.llm.unlink.not_active")+"\n", provider.Name, provider.Model)
		} else {
			fmt.Printf(lang.T("msg.llm.unlink.failed")+"\n", err)
		}
		return
	}

	fmt.Printf(lang.T("msg.llm.unlink.success")+"\n", provider.Name, provider.Model)
}

// TestLLM 使用统一函数
func TestLLM(identifier string) {
	config := cfg.LoadConfig()

	p, err := llm.FindOrSelectProvider(
		config.LLMProviders,
		identifier,
		&llm.SelectProviderOptions{
			MultipleMatchesPrompt: lang.T("msg.llm.test.multiple_matches"),
			SelectPrompt:          lang.T("msg.llm.select_prompt"),
			OnSelect:              selectProviderIndex,
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "cancelled") || strings.Contains(err.Error(), "已取消") {
			fmt.Println(lang.T("msg.llm.cancelled"))
		} else {
			fmt.Printf(lang.T("msg.llm.test.not_found")+"\n", identifier)
		}
		return
	}

	// 后续逻辑保持不变
	testRegisteredLLM(p)
}

// RemoveLLM 使用统一函数（需要额外的确认步骤）
func RemoveLLM(identifier string) {
	config := cfg.LoadConfig()

	p, err := llm.FindOrSelectProvider(
		config.LLMProviders,
		identifier,
		&llm.SelectProviderOptions{
			MultipleMatchesPrompt: lang.T("msg.llm.remove.multiple_matches"),
			SelectPrompt:          lang.T("msg.llm.remove.select_prompt"),
			OnSelect:              selectProviderIndex,
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "cancelled") || strings.Contains(err.Error(), "已取消") {
			fmt.Println(lang.T("msg.llm.remove.cancelled"))
		} else {
			fmt.Printf(lang.T("msg.llm.remove.not_found")+"\n", identifier)
		}
		return
	}

	// remove 特有的：删除前确认
	fmt.Println(llm.FormatProviderInfo(*p))
	if !confirmDeletion() {
		fmt.Println(lang.T("msg.llm.remove.cancelled"))
		return
	}

	// 执行删除
	found, err := cfg.RemoveLLMProvider(p.ID)
	if err != nil {
		fmt.Printf(lang.T("msg.llm.remove.save_failed")+"\n", err)
		return
	}

	if found {
		fmt.Printf(lang.T("msg.llm.remove.success")+"\n", identifier)
	} else {
		fmt.Printf(lang.T("msg.llm.remove.not_found")+"\n", p.ID)
	}
}

// 列出所有纯粹读取本地已经注册配置的LLM
func ListLLM() {
	items, err := llm.BuildProviderList(lang.T("msg.llm.date_format"))
	if err != nil {
		fmt.Printf(lang.T("msg.llm.load_config_failed")+"\n", err)
		return
	}
	if len(items) == 0 {
		fmt.Println(lang.T("msg.llm.no_providers"))
		return
	}

	printLLMTableHeader()
	for _, item := range items {
		printLLMTableRow(item)
	}
}

// 通用的选择函数（在 tasks 层）
func selectProviderIndex(count int, prompt string) int {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return -1
	}

	input = strings.TrimSpace(input)
	if input == "" || strings.ToLower(input) == "q" || strings.ToLower(input) == "quit" {
		return -1
	}

	var index int
	if _, err := fmt.Sscanf(input, "%d", &index); err != nil {
		return -1
	}

	// 转换为 0-based 索引
	index--
	if index < 0 || index >= count {
		return -1
	}

	return index
}

// 已注册模型的连通性测试及状态更新
func testRegisteredLLM(p *cfg.LLMProvider) {
	available, err := llm.LLMConnectivity(p)
	if err != nil {
		fmt.Printf(lang.T("msg.llm.test.test_failed")+"\n", err)
		return
	}

	if available {
		// 测试通过，更新 Verified 状态
		p.Verified = true
		p.VerifiedAt = time.Now().Format(time.RFC3339)

		if err := cfg.UpdateLLMProvider(p.ID, p); err != nil {
			fmt.Printf(lang.T("msg.llm.test.update_config_failed_on_success")+"\n", err)
		} else {
			fmt.Printf(lang.T("msg.llm.test.available_verified")+"\n", p.Name, p.Provider, p.Model)
		}
		return
	}

	// 测试失败，清除验证状态
	p.Verified = false
	p.VerifiedAt = ""

	if err := cfg.UpdateLLMProvider(p.ID, p); err != nil {
		fmt.Printf(lang.T("msg.llm.test.update_config_failed_on_failure")+"\n", err)
	} else {
		fmt.Printf(lang.T("msg.llm.test.unavailable")+"\n", p.Name, p.Provider, p.Model)
	}
}

// confirmDeletion 确认删除（用户交互逻辑，保留在 tasks 层）
func confirmDeletion() bool {
	return askForConfirmation(lang.T("msg.llm.remove.confirm_prompt"))
}

// --- CLI 内部辅助：对齐表格输出（避免侵入核心层） ---
const (
	llmColIDWidth       = 18 // ID 列宽（SHA256前16字符的十六进制，留2字符余量）
	llmColNameWidth     = 20
	llmColProviderWidth = 15
	llmColModelWidth    = 20
	llmColEnabledWidth  = 10
	llmColVerifiedWidth = 8
	llmColVerifiedAt    = 12
)

// padRunewidth 使用 runewidth 计算字符串宽度并填充空格
func padRunewidth(s string, width int) string {
	w := runewidth.StringWidth(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

// printLLMTableHeader 打印 LLM 列表表头
func printLLMTableHeader() {
	printLLMRow(
		lang.T("msg.llm.list_header_id"),
		lang.T("msg.llm.list_header_name"),
		lang.T("msg.llm.list_header_provider"),
		lang.T("msg.llm.list_header_model"),
		lang.T("msg.llm.list_header_enabled"),
		lang.T("msg.llm.list_header_verified"),
		lang.T("msg.llm.list_header_verified_at"),
		lang.T("msg.llm.list_header_active"),
	)
}

// printLLMTableRow 打印 LLM 列表的一行数据
func printLLMTableRow(item llm.ProviderListItem) {
	printLLMRow(
		item.ID,
		item.Name,
		item.Provider,
		item.Model,
		formatEnabledText(item.Enabled),
		item.VerifiedMark,
		item.VerifiedAt,
		item.ActiveMark,
	)
}

// formatEnabledText 格式化启用状态文本
func formatEnabledText(enabled bool) string {
	if enabled {
		return lang.T("msg.common.yes")
	}
	return lang.T("msg.common.no")
}

// printLLMRow 打印表格行（内部辅助函数）
func printLLMRow(id, name, provider, model, enabled, verified, verifiedAt, active string) {
	fmt.Printf("  %s %s %s %s %s %s %s %s\n",
		padRunewidth(id, llmColIDWidth),
		padRunewidth(name, llmColNameWidth),
		padRunewidth(provider, llmColProviderWidth),
		padRunewidth(model, llmColModelWidth),
		padRunewidth(enabled, llmColEnabledWidth),
		padRunewidth(verified, llmColVerifiedWidth),
		padRunewidth(verifiedAt, llmColVerifiedAt),
		active)
}
