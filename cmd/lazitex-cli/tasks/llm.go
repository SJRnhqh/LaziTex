// cmd/lazitex-cli/tasks/llm.go
// LLM 管理相关业务

package tasks

import (
	"fmt"
	"time"

	config "github.com/SJRnhqh/lazitex/config"
)

// 默认用 Ollama，把模型名和 ID/Name 统一成传入的别名，方便「智能注册+链接」
func newDefaultProvider(alias string) *config.LLMProvider {
	return &config.LLMProvider{
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
	existing, err := config.FindLLMProvider(alias)
	if err != nil {
		fmt.Printf("❌ 读取配置失败: %v\n", err)
		return
	}

	if existing != nil {
		if err := config.SetActiveLLM(existing.ID); err != nil {
			fmt.Printf("❌ 链接失败: %v\n", err)
			return
		}
		fmt.Printf("🔗 已链接到现有 LLM: %s (%s)\n", existing.Name, existing.Model)
		return
	}

	provider := newDefaultProvider(alias)
	if err := config.AddLLMProvider(provider); err != nil {
		fmt.Printf("❌ 注册失败: %v\n", err)
		return
	}
	if err := config.SetActiveLLM(provider.ID); err != nil {
		fmt.Printf("❌ 注册成功但链接失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 已注册并链接 LLM: %s (%s)\n", provider.Name, provider.Model)
}

// 显式注册（不自动链接）
func AddLLM(alias string) {
	provider := newDefaultProvider(alias)
	if err := config.AddLLMProvider(provider); err != nil {
		fmt.Printf("❌ 注册失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 已注册 LLM: %s (%s)\n", provider.Name, provider.Model)
}

// 显式链接已存在的
func LinkLLM(alias string) {
	if err := config.SetActiveLLM(alias); err != nil {
		fmt.Printf("❌ 链接失败: %v\n", err)
		return
	}
	fmt.Printf("🔗 已切换到 LLM: %s\n", alias)
}

// 列出所有
func ListLLM() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ 读取配置失败: %v\n", err)
		return
	}
	if len(config.LLMProviders) == 0 {
		fmt.Println("（暂无 LLM 配置）")
		return
	}
	fmt.Println("ID/Name\tProvider\tModel\tEnabled\tActive")
	for _, p := range config.LLMProviders {
		active := ""
		if config.ActiveLLM == p.ID {
			active = "*"
		}
		fmt.Printf("%s\t%s\t%s\t%v\t%s\n", p.Name, p.Provider, p.Model, p.Enabled, active)
	}
}

// 测试：目前只是检查配置存在与否（后续可扩展真实连通性）
func TestLLM(alias string) {
	p, err := config.FindLLMProvider(alias)
	if err != nil {
		fmt.Printf("❌ 读取失败: %v\n", err)
		return
	}
	if p == nil {
		fmt.Printf("❌ 未找到 LLM '%s'\n", alias)
		return
	}
	status := "✅ 可用"
	if !p.Enabled {
		status = "⚠️ 已禁用"
	}
	fmt.Printf("%s: %s (%s:%s)\n", status, p.Name, p.Provider, p.Model)
}

// 删除
func RemoveLLM(alias string) {
	if err := config.RemoveLLMProvider(alias); err != nil {
		fmt.Printf("❌ 删除失败: %v\n", err)
		return
	}
	fmt.Printf("🗑️ 已删除 LLM: %s\n", alias)
}

// 可选：标记验证时间的工具函数，供以后真实验证时使用
func markVerified(p *config.LLMProvider) {
	p.Verified = true
	p.VerifiedAt = time.Now().Format(time.RFC3339)
}