// cmd/lazitex-cli/tasks/llm.go
// LLM 管理相关业务

package tasks

import (
	// 外部包
	"bufio"
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

// 智能：不存在则注册+设为当前，存在则直接设为当前
func SmartLinkOrAddLLM(alias string) {
	existing, err := cfg.FindLLMProvider(alias)
	if err != nil {
		fmt.Printf("❌ 读取配置失败: %v\n", err)
		return
	}

	if existing != nil {
		if err := cfg.SetActiveLLM(existing.ID); err != nil {
			fmt.Printf("❌ 连接失败: %v\n", err)
			return
		}
		fmt.Printf("🔗 已连接到现有 LLM: %s (%s)\n", existing.Name, existing.Model)
		return
	}

	llmProvider := cfg.NewLLMProvider(alias, "ollama", alias)
	if err := cfg.AddLLMProvider(llmProvider); err != nil {
		fmt.Printf("❌ 注册失败: %v\n", err)
		return
	}
	if err := cfg.SetActiveLLM(llmProvider.ID); err != nil {
		fmt.Printf("❌ 注册成功但连接失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 已注册并连接 LLM: %s (%s)\n", llmProvider.Name, llmProvider.Model)
}

// 显式注册（不自动连接）
func AddLLM(provider, name, model string) {
	llmProvider := cfg.NewLLMProvider(provider, name, model)
	verified, err := llm.LLMConnectivity(llmProvider)
	if err != nil {
		fmt.Printf(lang.T("msg.llm.test.test_failed")+"\n", err)
		return
	}
	if verified {
		fmt.Printf("连接测试验证成功，正在注册 LLM: %s (%s)\n", llmProvider.Name, llmProvider.Model)
	} else {
		fmt.Printf("❌ 连接测试失败，取消注册: %v\n", err)
	}
	if err := cfg.AddLLMProvider(llmProvider); err != nil {
		fmt.Printf("❌ 注册失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 已注册 LLM: %s (%s)\n", llmProvider.Name, llmProvider.Model)
}

func LinkLLM(name string) {
	p, err := cfg.FindLLMProvider(name)
	if err != nil {
		fmt.Printf(lang.T("msg.llm.link.read_failed")+"\n", err)
		return
	}
	if p == nil {
		fmt.Printf(lang.T("msg.llm.link.not_found")+"\n", name)
		return
	}

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

// UnlinkLLM 取消连接指定的 LLM
func UnlinkLLM(model string) {
	// 先查找 provider 以便显示名称
	provider, err := cfg.FindLLMProvider(model)
	if err != nil {
		fmt.Printf(lang.T("msg.llm.unlink.read_failed")+"\n", err)
		return
	}
	if provider == nil {
		fmt.Printf(lang.T("msg.llm.unlink.not_found")+"\n", model)
		return
	}

	// 调用 config 层的 UnsetActiveLLM，它会验证激活状态
	if err := cfg.UnsetActiveLLM(model); err != nil {
		// 根据错误类型显示不同的消息
		if err.Error() == fmt.Sprintf("LLM Provider '%s' 不在激活列表中", model) {
			fmt.Printf(lang.T("msg.llm.unlink.not_active")+"\n", provider.Name, provider.Model)
		} else {
			fmt.Printf(lang.T("msg.llm.unlink.failed")+"\n", err)
		}
		return
	}

	// 成功消息
	fmt.Printf(lang.T("msg.llm.unlink.success")+"\n", provider.Name, provider.Model)
}

// 测试：检查配置并测试真实连通性（支持已注册和未注册两种模式）
func TestLLM(idOrNameOrModel string) {
	// 尝试作为已注册的 Provider 查找
	p, err := cfg.FindLLMProvider(idOrNameOrModel)
	// 未找到对应的 LLM Provider
	if err != nil {
		fmt.Printf(lang.T("msg.llm.test.read_failed")+"\n", err)
		return
	}
	// 已注册的模型，执行完整测试和验证状态更新
	testRegisteredLLM(p)
}

// 删除LLM Provider（支持 ID、Name、Model）
func RemoveLLM(identifier string) {
	config := cfg.LoadConfig()

	// 查找匹配的providers（调用核心层函数）
	matches := llm.FindMatchingLLMProviders(config.LLMProviders, identifier)

	if len(matches) == 0 {
		fmt.Printf(lang.T("msg.llm.remove.not_found")+"\n", identifier)
		return
	}

	// 处理匹配结果
	var targetID string

	if len(matches) == 1 {
		// 单个匹配：显示信息并确认删除
		fmt.Println(llm.FormatProviderInfo(matches[0]))
		if !confirmDeletion() {
			fmt.Println(lang.T("msg.llm.remove.cancelled"))
			return
		}
		targetID = matches[0].ID
	} else {
		// 多个匹配：显示列表让用户选择
		fmt.Printf(lang.T("msg.llm.remove.multiple_matches")+"\n", len(matches))
		formattedList := llm.FormatProviderList(matches)
		for _, item := range formattedList {
			fmt.Println("  " + item)
		}

		selectedIndex := selectProvider(len(matches))
		if selectedIndex < 0 {
			fmt.Println(lang.T("msg.llm.remove.cancelled"))
			return
		}

		// selectedIndex 已经通过 selectProvider 验证，一定是有效索引
		targetID = matches[selectedIndex].ID
	}

	// 执行删除（调用 config 层函数，只通过 ID 删除）
	found, err := cfg.RemoveLLMProvider(targetID)
	if err != nil {
		fmt.Printf(lang.T("msg.llm.remove.save_failed")+"\n", err)
		return
	}

	if found {
		fmt.Printf(lang.T("msg.llm.remove.success")+"\n", identifier)
	} else {
		fmt.Printf(lang.T("msg.llm.remove.not_found")+"\n", targetID)
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

	// 表头：使用国际化字符串
	printLLMRow(
		lang.T("msg.llm.list_header_id"),
		lang.T("msg.llm.list_header_provider"),
		lang.T("msg.llm.list_header_model"),
		lang.T("msg.llm.list_header_enabled"),
		lang.T("msg.llm.list_header_verified"),
		lang.T("msg.llm.list_header_verified_at"),
		lang.T("msg.llm.list_header_active"),
	)

	for _, item := range items {
		enabledText := lang.T("msg.common.no")
		if item.Enabled {
			enabledText = lang.T("msg.common.yes")
		}

		printLLMRow(
			item.Name,
			item.Provider,
			item.Model,
			enabledText,
			item.VerifiedMark,
			item.VerifiedAt,
			item.ActiveMark,
		)
	}
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

// selectProvider 让用户从多个匹配项中选择（用户交互逻辑，保留在 tasks 层）
func selectProvider(count int) int {
	fmt.Print(lang.T("msg.llm.remove.select_prompt"))

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

// --- CLI 内部辅助：对齐表格输出（避免侵入核心层） ---
const (
	llmColIDWidth       = 20
	llmColProviderWidth = 15
	llmColModelWidth    = 20
	llmColEnabledWidth  = 10
	llmColVerifiedWidth = 8
	llmColVerifiedAt    = 12
)

func padRunewidth(s string, width int) string {
	w := runewidth.StringWidth(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

func printLLMRow(id, provider, model, enabled, verified, verifiedAt, active string) {
	fmt.Printf("  %s %s %s %s %s %s %s\n",
		padRunewidth(id, llmColIDWidth),
		padRunewidth(provider, llmColProviderWidth),
		padRunewidth(model, llmColModelWidth),
		padRunewidth(enabled, llmColEnabledWidth),
		padRunewidth(verified, llmColVerifiedWidth),
		padRunewidth(verifiedAt, llmColVerifiedAt),
		active)
}
