// core/ai/llm/list.go
// LLM 列表构建的业务逻辑

package llm

import (
	"time"

	cfg "github.com/SJRnhqh/lazitex/config"
)

// ProviderListItem 用于展示的 LLM Provider 列表项
type ProviderListItem struct {
	Name         string
	Provider     string
	Model        string
	Enabled      bool
	VerifiedMark string
	VerifiedAt   string
	ActiveMark   string
}

// BuildProviderList 负责构建展示层需要的 LLM 列表数据
func BuildProviderList(dateFormat string) ([]ProviderListItem, error) {
	cfgData, err := cfg.LoadConfig()
	if err != nil {
		return nil, err
	}

	items := make([]ProviderListItem, 0, len(cfgData.LLMProviders))
	for _, p := range cfgData.LLMProviders {
		item := ProviderListItem{
			Name:         p.Name,
			Provider:     p.Provider,
			Model:        p.Model,
			Enabled:      p.Enabled,
			VerifiedMark: formatVerifiedMark(p.Verified),
			VerifiedAt:   formatVerifiedAt(p.VerifiedAt, dateFormat),
		}

		if cfgData.ActiveLLM == p.ID {
			item.ActiveMark = "*"
		}

		items = append(items, item)
	}

	return items, nil
}

func formatVerifiedMark(verified bool) string {
	if verified {
		return "✓"
	}
	return "✗"
}

func formatVerifiedAt(raw string, dateFormat string) string {
	if raw == "" {
		return ""
	}

	if dateFormat == "" {
		dateFormat = "2006-01-02"
	}

	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t.Format(dateFormat)
	}

	if len(raw) >= 10 {
		return raw[:10]
	}

	return raw
}
