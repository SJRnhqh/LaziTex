// config/llm.go
// 负责管理LLM Provider配置

package config

import (
	// 外部包
	"crypto/sha256"
	"encoding/hex"
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

// 创建新的 LLM Provider
func NewLLMProvider(provider, name, model string) *LLMProvider {
	return &LLMProvider{
		ID:         generateLLMID(provider, name, model),
		Provider:   provider,
		Name:       name,
		Model:      model,
		Enabled:    true,  // 默认注册时启用的状态（主要用于未来web端的启用与禁用，cli默认启用，但是先不加命令控制启用禁用的功能）
		Verified:   false, // TODO: 未来添加验证机制
		VerifiedAt: "",    // TODO: 未来添加验证时间
		Config:     map[string]string{},
	}
}

// AddLLMProvider 添加新的 LLM Provider（如果已存在则返回错误）
func AddLLMProvider(llmProvider *LLMProvider) error {
	config := LoadConfig()

	// 检查是否已存在相同 ID 或名称的 Provider
	for _, existing := range config.LLMProviders {
		if existing.ID == llmProvider.ID {
			return fmt.Errorf("LLM Provider ID '%s' 已存在", llmProvider.ID)
		}
		if existing.Name == llmProvider.Name {
			return fmt.Errorf("LLM Provider 名称 '%s' 已存在", llmProvider.Name)
		}
	}

	// 添加新的 LLM Provider
	config.LLMProviders = append(config.LLMProviders, *llmProvider)
	return SaveConfig(config)
}

// UpdateLLMProvider 更新已存在的 LLM Provider
func UpdateLLMProvider(idOrNameOrModel string, provider *LLMProvider) error {
	config := LoadConfig()

	// 查找要更新的 Provider
	found := false
	for i := range config.LLMProviders {
		if config.LLMProviders[i].ID == idOrNameOrModel || config.LLMProviders[i].Name == idOrNameOrModel {
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
		return fmt.Errorf("LLM Provider '%s' 不存在", idOrNameOrModel)
	}

	return SaveConfig(config)
}

// RemoveLLMProvider 删除指定的 LLM Provider
// 只通过 ID 进行精确删除，确保无歧义
// 返回值：(是否找到并删除成功, 错误信息)
func RemoveLLMProvider(PID string) (bool, error) {
    config := LoadConfig()

    // 查找要删除的provider
    targetIndex := -1
    for i, llmprovider := range config.LLMProviders {
        if llmprovider.ID == PID {
            targetIndex = i
            break
        }
    }

    // 未找到
    if targetIndex == -1 {
        return false, nil
    }

    // 获取要删除的provider ID（用于清理激活列表）
    removedID := config.LLMProviders[targetIndex].ID

    // 从Providers列表中删除
    config.LLMProviders = append(
        config.LLMProviders[:targetIndex],
        config.LLMProviders[targetIndex+1:]...,
    )

    // 从激活列表中移除
    newActiveLLMs := []string{}
    for _, activeID := range config.ActiveLLMs {
        if activeID != removedID {
            newActiveLLMs = append(newActiveLLMs, activeID)
        }
    }
    config.ActiveLLMs = newActiveLLMs

    // 保存配置
    return true, SaveConfig(config)
}

// FindLLMProvider 根据ID、名称或模型查找 LLM Provider
func FindLLMProvider(idOrNameOrModel string) (*LLMProvider, error) {
	config := LoadConfig()

	for i := range config.LLMProviders {
		if config.LLMProviders[i].ID == idOrNameOrModel || config.LLMProviders[i].Name == idOrNameOrModel || config.LLMProviders[i].Model == idOrNameOrModel {
			return &config.LLMProviders[i], nil
		}
	}

	return nil, nil // 未找到，返回 nil
}

// SetActiveLLM 添加 LLM Provider 到激活列表（如果已存在则不重复添加）
func SetActiveLLM(idOrNameOrModel string) error {
	// 验证 LLM Provider 是否存在
	provider, err := FindLLMProvider(idOrNameOrModel)
	if err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("LLM Provider '%s' 不存在", idOrNameOrModel)
	}

	// 加载配置
	config := LoadConfig()

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
func UnsetActiveLLM(idOrNameOrModel string) error {
	// 验证 LLM Provider 是否存在
	provider, err := FindLLMProvider(idOrNameOrModel)
	if err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("LLM Provider '%s' 不存在", idOrNameOrModel)
	}

	// 加载配置检查当前激活状态
	config := LoadConfig()

	// 验证指定的 LLM 是否在激活列表中
	found := false
	for _, activeID := range config.ActiveLLMs {
		if activeID == provider.ID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("LLM Provider '%s' 不在激活列表中", idOrNameOrModel)
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

// generateLLMID 根据 provider、model 和可选的 name 生成唯一的 ID
// 如果提供了 name，使用 name+provider+model 生成哈希
// 如果没有提供 name，使用 provider+model 生成哈希
// 返回前 16 个字符的十六进制哈希值作为 ID
func generateLLMID(provider, name, model string) string {
	var input string
	if name != "" {
		// 使用 name + provider + model 确保唯一性
		input = fmt.Sprintf("%s:%s:%s", provider, name, model)
	} else {
		// 使用 provider + model 生成 ID
		input = fmt.Sprintf("%s:%s", provider, model)
	}

	// 生成 SHA256 哈希
	hash := sha256.Sum256([]byte(input))
	// 返回前 16 个字符的十六进制字符串作为 ID
	return hex.EncodeToString(hash[:])[:16]
}
