// config/llm.go
// 负责管理LLM Provider配置

package config

import (
	"fmt"
)

// LLMProvider LLM Provider 配置
type LLMProvider struct {
	ID         string            `json:"id"`         // 唯一标识
	Provider   string            `json:"provider"`   // 提供商类型: "ollama", "openai" 等
	Name       string            `json:"name"`       // 用户自定义名称
	Model      string            `json:"model"`      // 模型名称
	Enabled    bool              `json:"enabled"`    // 是否启用
	Verified   bool              `json:"verified"`   // 是否已验证
	VerifiedAt string            `json:"verifiedAt"` // 验证时间
	Config     map[string]string `json:"config"`     // 提供商特定配置（baseURL, timeout 等）
}

// AddLLMProvider 添加新的 LLM Provider（如果已存在则返回错误）
func AddLLMProvider(provider *LLMProvider) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// 检查是否已存在相同 ID 或名称的 Provider
	for _, existing := range config.LLMProviders {
		if existing.ID == provider.ID {
			return fmt.Errorf("LLM Provider ID '%s' 已存在", provider.ID)
		}
		if existing.Name == provider.Name {
			return fmt.Errorf("LLM Provider 名称 '%s' 已存在", provider.Name)
		}
	}

	// 添加新 Provider
	config.LLMProviders = append(config.LLMProviders, *provider)
	return SaveConfig(config)
}

// UpdateLLMProvider 更新已存在的 LLM Provider
func UpdateLLMProvider(idOrName string, provider *LLMProvider) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// 查找要更新的 Provider
	found := false
	for i := range config.LLMProviders {
		if config.LLMProviders[i].ID == idOrName || config.LLMProviders[i].Name == idOrName {
			// 检查新名称是否与其他 Provider 冲突（除了自己）
			for j, other := range config.LLMProviders {
				if i != j && other.Name == provider.Name {
					return fmt.Errorf("LLM Provider 名称 '%s' 已存在", provider.Name)
				}
			}

			// 更新 Provider（保持原有 ID）
			provider.ID = config.LLMProviders[i].ID
			config.LLMProviders[i] = *provider
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("LLM Provider '%s' 不存在", idOrName)
	}

	return SaveConfig(config)
}

// RemoveLLMProvider 删除 LLM Provider
func RemoveLLMProvider(idOrName string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// 查找并删除
	newProviders := []LLMProvider{}
	found := false
	for _, provider := range config.LLMProviders {
		if provider.ID == idOrName || provider.Name == idOrName {
			found = true
			// 如果删除的是当前激活的 LLM，从激活列表中移除
			newActiveLLMs := []string{}
			for _, activeID := range config.ActiveLLMs {
				if activeID != provider.ID {
					newActiveLLMs = append(newActiveLLMs, activeID)
				}
			}
			config.ActiveLLMs = newActiveLLMs
		} else {
			newProviders = append(newProviders, provider)
		}
	}

	if !found {
		return fmt.Errorf("LLM Provider '%s' 不存在", idOrName)
	}

	config.LLMProviders = newProviders
	return SaveConfig(config)
}

// FindLLMProvider 查找 LLM Provider
func FindLLMProvider(idOrName string) (*LLMProvider, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	for i := range config.LLMProviders {
		if config.LLMProviders[i].ID == idOrName || config.LLMProviders[i].Name == idOrName {
			return &config.LLMProviders[i], nil
		}
	}

	return nil, nil // 未找到，返回 nil
}

// SetActiveLLM 添加 LLM Provider 到激活列表（如果已存在则不重复添加）
func SetActiveLLM(idOrName string) error {
	// 验证 LLM Provider 是否存在
	provider, err := FindLLMProvider(idOrName)
	if err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("LLM Provider '%s' 不存在", idOrName)
	}

	// 加载配置
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// 检查是否已在激活列表中
	for _, activeID := range config.ActiveLLMs {
		if activeID == provider.ID {
			// 已在列表中，直接返回成功
			return nil
		}
	}

	// 添加到激活列表
	config.ActiveLLMs = append(config.ActiveLLMs, provider.ID)
	return SaveConfig(config)
}

// UnsetActiveLLM 从激活列表中移除指定的 LLM Provider（需要验证）
func UnsetActiveLLM(idOrName string) error {
	// 验证 LLM Provider 是否存在
	provider, err := FindLLMProvider(idOrName)
	if err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("LLM Provider '%s' 不存在", idOrName)
	}

	// 加载配置检查当前激活状态
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// 验证指定的 LLM 是否在激活列表中
	found := false
	for _, activeID := range config.ActiveLLMs {
		if activeID == provider.ID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("LLM Provider '%s' 不在激活列表中", idOrName)
	}

	// 从激活列表中移除
	newActiveLLMs := []string{}
	for _, activeID := range config.ActiveLLMs {
		if activeID != provider.ID {
			newActiveLLMs = append(newActiveLLMs, activeID)
		}
	}
	config.ActiveLLMs = newActiveLLMs
	return SaveConfig(config)
}
