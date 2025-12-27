// core/ai/llm/list.go
// 列出所有注册的LLMProvider的核心业务逻辑

package llm

import (
	// 外部包
	"fmt"
	"strings"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// ListLLMProviders 列出 LLM Provider 的核心业务逻辑
func ListLLMProviders(llmProviderIDOrNameOrModel []string) {
	var providers []*cfg.LLMProvider

	// 如果没有参数，列出所有已注册的 Provider
	if len(llmProviderIDOrNameOrModel) == 0 {
		providers = cfg.GetAllLLMProviders()
	} else {
		// 有参数，查找匹配的 Provider
		providers = findMatchingProviders(llmProviderIDOrNameOrModel)
	}

	if len(providers) == 0 {
		fmt.Println(lang.T("msg.llm.list.no_providers"))
		return
	}

	// 以表格形式显示
	printProvidersTable(providers)
}

// findMatchingProviders 查找匹配的 Provider（支持多个参数）
func findMatchingProviders(idsOrNamesOrModels []string) []*cfg.LLMProvider {
	var allMatches []*cfg.LLMProvider
	matchedIDs := make(map[string]bool)

	for _, idOrNameOrModel := range idsOrNamesOrModels {
		matches, _ := cfg.FindAllMatchingProviders(idOrNameOrModel)
		for _, match := range matches {
			// 避免重复
			if !matchedIDs[match.ID] {
				allMatches = append(allMatches, match)
				matchedIDs[match.ID] = true
			}
		}
	}

	return allMatches
}

// printProvidersTable 以表格形式打印 Provider 列表（类似 docker ps）
func printProvidersTable(providers []*cfg.LLMProvider) {
	config := cfg.LoadConfig()

	// 计算每列的最大宽度
	widths := calculateColumnWidths(providers, config)

	// 打印表头
	printTableHeader(widths)

	// 打印分隔线
	printTableSeparator(widths)

	// 打印数据行
	for _, p := range providers {
		printTableRow(p, widths, config)
	}
}

// getDisplayWidth 计算字符串的显示宽度（中文字符占2个宽度）
func getDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		// 中文字符、全角字符等占2个宽度
		if r >= 0x1100 && (r <= 0x115F || r >= 0x2E80 && r <= 0x9FFF ||
			r >= 0xAC00 && r <= 0xD7AF || r >= 0xF900 && r <= 0xFAFF ||
			r >= 0xFE30 && r <= 0xFE4F || r >= 0xFF00 && r <= 0xFFEF) {
			width += 2
		} else {
			width += 1
		}
	}
	return width
}

// calculateColumnWidths 计算每列的最大宽度（使用显示宽度）
func calculateColumnWidths(providers []*cfg.LLMProvider, config *cfg.Config) map[string]int {
	widths := map[string]int{
		"Marker":     1, // 标记列固定宽度为 1（用于显示 *）
		"ID":         getDisplayWidth(lang.T("msg.llm.list_header_id")),
		"Name":       getDisplayWidth(lang.T("msg.llm.list_header_name")),
		"Provider":   getDisplayWidth(lang.T("msg.llm.list_header_provider")),
		"Model":      getDisplayWidth(lang.T("msg.llm.list_header_model")),
		"Enabled":    getDisplayWidth(lang.T("msg.llm.list_header_enabled")),
		"Verified":   getDisplayWidth(lang.T("msg.llm.list_header_verified")),
		"VerifiedAt": getDisplayWidth(lang.T("msg.llm.list_header_verified_at")),
		"Active":     getDisplayWidth(lang.T("msg.llm.list_header_active")),
	}

	// 遍历所有 Provider，找出每列的最大宽度
	for _, p := range providers {
		// ID（显示前8位）
		idDisplay := p.ID[:8]
		if w := getDisplayWidth(idDisplay); w > widths["ID"] {
			widths["ID"] = w
		}

		// Name
		if w := getDisplayWidth(p.Name); w > widths["Name"] {
			widths["Name"] = w
		}

		// Provider
		if w := getDisplayWidth(p.Provider); w > widths["Provider"] {
			widths["Provider"] = w
		}

		// Model
		if w := getDisplayWidth(p.Model); w > widths["Model"] {
			widths["Model"] = w
		}

		// Enabled
		enabledStr := formatBool(p.Enabled)
		if w := getDisplayWidth(enabledStr); w > widths["Enabled"] {
			widths["Enabled"] = w
		}

		// Verified
		verifiedStr := formatBool(p.Verified)
		if w := getDisplayWidth(verifiedStr); w > widths["Verified"] {
			widths["Verified"] = w
		}

		// VerifiedAt
		verifiedAtStr := formatVerifiedAt(p.VerifiedAt)
		if w := getDisplayWidth(verifiedAtStr); w > widths["VerifiedAt"] {
			widths["VerifiedAt"] = w
		}

		// Active（状态字符串）
		activeStr := getActiveStatus(p.ID, config)
		if w := getDisplayWidth(activeStr); w > widths["Active"] {
			widths["Active"] = w
		}
	}

	// 设置最小宽度（使用显示宽度）
	minWidths := map[string]int{
		"Marker":     1,
		"ID":         8,
		"Name":       4,
		"Provider":   8,
		"Model":      5,
		"Enabled":    6,
		"Verified":   8,
		"VerifiedAt": 10,
		"Active":     6,
	}

	for key, minWidth := range minWidths {
		if widths[key] < minWidth {
			widths[key] = minWidth
		}
	}

	return widths
}

// printTableHeader 打印表头（使用显示宽度对齐）
func printTableHeader(widths map[string]int) {
	headers := []string{
		"", // 标记列（空）
		lang.T("msg.llm.list_header_id"),
		lang.T("msg.llm.list_header_name"),
		lang.T("msg.llm.list_header_provider"),
		lang.T("msg.llm.list_header_model"),
		lang.T("msg.llm.list_header_enabled"),
		lang.T("msg.llm.list_header_verified"),
		lang.T("msg.llm.list_header_verified_at"),
		lang.T("msg.llm.list_header_active"),
	}
	keys := []string{"Marker", "ID", "Name", "Provider", "Model", "Enabled", "Verified", "VerifiedAt", "Active"}

	for i, header := range headers {
		key := keys[i]
		width := widths[key]
		displayWidth := getDisplayWidth(header)
		padding := width - displayWidth
		if padding < 0 {
			padding = 0
		}
		fmt.Print(header + strings.Repeat(" ", padding))
		if i < len(headers)-1 {
			fmt.Print("  ")
		}
	}
	fmt.Println()
}

// printTableSeparator 打印分隔线
func printTableSeparator(widths map[string]int) {
	totalWidth := widths["Marker"] + widths["ID"] + widths["Name"] + widths["Provider"] + widths["Model"] +
		widths["Enabled"] + widths["Verified"] + widths["VerifiedAt"] + widths["Active"] + 16 // 16 = 8个空格分隔符
	separator := strings.Repeat("-", totalWidth)
	fmt.Println(separator)
}

// getCurrentMarker 获取当前前台标记（* 或空格）
func getCurrentMarker(providerID string, config *cfg.Config) string {
	if config.CurrentLLM == providerID {
		return "*"
	}
	return " "
}

// printTableRow 打印数据行（使用显示宽度对齐）
func printTableRow(p *cfg.LLMProvider, widths map[string]int, config *cfg.Config) {
	values := []string{
		getCurrentMarker(p.ID, config), // 标记列
		p.ID[:8],
		p.Name,
		p.Provider,
		p.Model,
		formatBool(p.Enabled),
		formatBool(p.Verified),
		formatVerifiedAt(p.VerifiedAt),
		getActiveStatus(p.ID, config),
	}
	keys := []string{"Marker", "ID", "Name", "Provider", "Model", "Enabled", "Verified", "VerifiedAt", "Active"}

	for i, value := range values {
		key := keys[i]
		width := widths[key]
		displayWidth := getDisplayWidth(value)
		padding := width - displayWidth
		if padding < 0 {
			padding = 0
		}
		fmt.Print(value + strings.Repeat(" ", padding))
		if i < len(values)-1 {
			fmt.Print("  ")
		}
	}
	fmt.Println()
}

// formatBool 格式化布尔值
func formatBool(b bool) string {
	if b {
		return lang.T("msg.llm.list.yes")
	}
	return lang.T("msg.llm.list.no")
}

// formatVerifiedAt 格式化验证时间
func formatVerifiedAt(verifiedAt string) string {
	if verifiedAt == "" {
		return lang.T("msg.llm.list.never")
	}
	return verifiedAt
}

// getActiveStatus 获取激活状态字符串
func getActiveStatus(providerID string, config *cfg.Config) string {
	// 检查是否激活
	isActive := false
	for _, activeID := range config.ActiveLLMs {
		if activeID == providerID {
			isActive = true
			break
		}
	}

	// 检查是否当前前台
	isCurrent := config.CurrentLLM == providerID

	if isCurrent {
		return lang.T("msg.llm.list.status_current")
	} else if isActive {
		return lang.T("msg.llm.list.status_active")
	}
	return lang.T("msg.llm.list.status_inactive")
}
