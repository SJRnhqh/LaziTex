// cmd/lazitex-cli/tasks/llm.go
// LLM 管理相关业务

package tasks

import (
	"fmt"
	"strings"
	"time"

	cfg "github.com/SJRnhqh/lazitex/config"
	llm "github.com/SJRnhqh/lazitex/core/ai/llm"
	lang "github.com/SJRnhqh/lazitex/lang"
	"github.com/mattn/go-runewidth"
)

// 默认用 Ollama，把模型名和 ID/Name 统一成传入的别名，方便「智能注册+链接」
func newDefaultProvider(alias string) *cfg.LLMProvider {
	return &cfg.LLMProvider{
		ID:       alias,
		Name:     alias,
		Provider: "ollama",
		Model:    alias,
		Enabled:  true,
		Config:   map[string]string{},
	}
}

// 智能：不存在则注册+设为当前，存在则直接设为当前
func SmartLinkOrAddLLM(alias string) {
	existing, err := cfg.FindLLMProvider(alias)
	if err != nil {
		fmt.Printf("❌ 读取配置失败: %v\n", err)
		return
	}

	if existing != nil {
		if err := cfg.SetActiveLLM(existing.ID); err != nil {
			fmt.Printf("❌ 链接失败: %v\n", err)
			return
		}
		fmt.Printf("🔗 已链接到现有 LLM: %s (%s)\n", existing.Name, existing.Model)
		return
	}

	provider := newDefaultProvider(alias)
	if err := cfg.AddLLMProvider(provider); err != nil {
		fmt.Printf("❌ 注册失败: %v\n", err)
		return
	}
	if err := cfg.SetActiveLLM(provider.ID); err != nil {
		fmt.Printf("❌ 注册成功但链接失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 已注册并链接 LLM: %s (%s)\n", provider.Name, provider.Model)
}

// 显式注册（不自动链接）
func AddLLM(alias string) {
	provider := newDefaultProvider(alias)
	if err := cfg.AddLLMProvider(provider); err != nil {
		fmt.Printf("❌ 注册失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 已注册 LLM: %s (%s)\n", provider.Name, provider.Model)
}

// 显式链接已存在的
func LinkLLM(alias string) {
	if err := cfg.SetActiveLLM(alias); err != nil {
		fmt.Printf("❌ 链接失败: %v\n", err)
		return
	}
	fmt.Printf("🔗 已切换到 LLM: %s\n", alias)
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

// 测试：检查配置并测试真实连通性（支持已注册和未注册两种模式）
func TestLLM(model string) {
	// 模式A：尝试作为已注册的 Provider 查找
	p, err := cfg.FindLLMProvider(model)
	if err != nil {
		fmt.Printf(lang.T("msg.llm.test.read_failed")+"\n", err)
		return
	}

	if p == nil {
		testUnregisteredLLM(model)
		return
	}

	// 模式A：已注册的模型，执行完整测试和验证状态更新
	testRegisteredLLM(p)
}

// 删除
func RemoveLLM(alias string) {
	if err := cfg.RemoveLLMProvider(alias); err != nil {
		fmt.Printf(lang.T("msg.llm.remove.failed")+"\n", err)
		return
	}
	fmt.Printf(lang.T("msg.llm.remove.success")+"\n", alias)
}

// 未注册模型的连通性测试（默认使用 ollama）
func testUnregisteredLLM(model string) {
	tempProvider := &cfg.LLMProvider{
		Provider: "ollama",
		Model:    model,
		Enabled:  true,
		Config:   map[string]string{},
	}

	available, err := llm.TestLLMConnectivity(tempProvider)
	if err != nil {
		fmt.Printf(lang.T("msg.llm.test.test_failed")+"\n", err)
		return
	}

	if available {
		fmt.Printf(lang.T("msg.llm.test.available_unregistered")+"\n", model, model)
	} else {
		fmt.Printf(lang.T("msg.llm.test.unavailable_unregistered")+"\n", model, model)
	}
}

// 已注册模型的连通性测试及状态更新
func testRegisteredLLM(p *cfg.LLMProvider) {
	available, err := llm.TestLLMConnectivity(p)
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
