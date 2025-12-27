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
	Provider   string            `json:"provider"`   // 供应商类型: "ollama", "openai" 等
	Name       string            `json:"name"`       // 用户自定义名称
	Model      string            `json:"model"`      // 模型名称
	Enabled    bool              `json:"enabled"`    // 是否启用
	Verified   bool              `json:"verified"`   // 是否已验证
	VerifiedAt string            `json:"verifiedAt"` // 验证时间
	Config     LLMProviderConfig `json:"config"`     // 供应商特定配置（baseURL, timeout 等）
}

type LLMProviderConfig struct {
	BaseURL string `json:"baseURL"` // 服务器基础URL
}

// 创建新的 LLM Provider
func NewLLMProvider(provider, name, model, baseURL string) *LLMProvider {
	return &LLMProvider{
		ID:         generateLLMID(provider, name, model),
		Provider:   provider,
		Name:       name,
		Model:      model,
		Enabled:    true,  // 默认注册时启用的状态（主要用于未来web端的启用与禁用，cli默认启用，但是先不加命令控制启用禁用的功能）
		Verified:   false, // TODO: 未来添加验证机制
		VerifiedAt: "",    // TODO: 未来添加验证时间
		Config: LLMProviderConfig{
			BaseURL: baseURL,
		},
	}
}

// 错误定义
var (
	ErrProviderIDExists   = fmt.Errorf("LLM Provider ID already exists")
	ErrProviderNameExists = fmt.Errorf("LLM Provider name already exists")
)

// AddLLMProvider 添加新的 LLM Provider（如果已存在则返回错误）
func AddLLMProvider(llmProvider *LLMProvider) error {
	config := LoadConfig()

	// 检查是否已存在相同 ID 或名称的 Provider
	for _, existing := range config.LLMProviders {
		if existing.ID == llmProvider.ID {
			return fmt.Errorf("%w: %s", ErrProviderIDExists, llmProvider.ID)
		}
		if existing.Name == llmProvider.Name {
			return fmt.Errorf("%w: %s", ErrProviderNameExists, llmProvider.Name)
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
// 如果通过 Model 匹配到多个，返回第一个匹配项（保持向后兼容）
func FindLLMProvider(idOrNameOrModel string) (*LLMProvider, error) {
	matches, _ := FindAllMatchingProviders(idOrNameOrModel)
	if len(matches) == 0 {
		return nil, nil
	}
	return matches[0], nil
}

// SetActiveLLM 添加 LLM Provider 到激活列表（如果已存在则不重复添加）
func SetActiveLLM(idOrNameOrModel string) error {
	// 查找所有匹配的 providers
	matches, matchType := FindAllMatchingProviders(idOrNameOrModel)

	if len(matches) == 0 {
		return fmt.Errorf("LLM Provider '%s' 不存在", idOrNameOrModel)
	}

	// 如果通过 Model 匹配到多个，返回特殊错误
	if matchType == "model" && len(matches) > 1 {
		return fmt.Errorf("MULTIPLE_MATCHES") // 特殊标记，由调用方处理交互
	}

	// 唯一匹配，添加到激活列表
	provider := matches[0]
	config := LoadConfig()

	// 检查是否已在激活列表中
	for _, activeID := range config.ActiveLLMs {
		if activeID == provider.ID {
			return nil // 已在列表中
		}
	}

	// 添加到激活列表
	config.ActiveLLMs = append(config.ActiveLLMs, provider.ID)
	return SaveConfig(config)
}

// SetActiveLLMByID 通过 ID 直接添加到激活列表（内部使用，不做查找）
func SetActiveLLMByID(providerID string) error {
	config := LoadConfig()
	
	// 验证 ID 是否存在
	found := false
	for _, p := range config.LLMProviders {
		if p.ID == providerID {
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("LLM Provider ID '%s' 不存在", providerID)
	}
	
	// 检查是否已在激活列表中
	for _, activeID := range config.ActiveLLMs {
		if activeID == providerID {
			return nil // 已在列表中
		}
	}
	
	// 添加到激活列表
	config.ActiveLLMs = append(config.ActiveLLMs, providerID)
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

// SetCurrentLLM 设置当前前台使用的 LLM Provider
func SetCurrentLLM(idOrName string) error {
	// 验证 LLM Provider 是否存在
	provider, err := FindLLMProvider(idOrName)
	if err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("LLM Provider '%s' 不存在", idOrName)
	}

	// 检查是否在激活列表中
	config := LoadConfig()
	found := false
	for _, activeID := range config.ActiveLLMs {
		if activeID == provider.ID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("LLM Provider '%s' 未激活，请先使用 link 命令", idOrName)
	}

	// 设置当前 LLM
	config.CurrentLLM = provider.ID
	return SaveConfig(config)
}

// GetCurrentLLM 获取当前前台使用的 LLM Provider
func GetCurrentLLM() (*LLMProvider, error) {
	config := LoadConfig()
	if config.CurrentLLM == "" {
		return nil, nil // 没有设置当前 LLM
	}

	return FindLLMProvider(config.CurrentLLM)
}

// FindAllMatchingProviders 查找所有匹配的 LLM Providers（导出供外部使用）
// 返回匹配的 providers 和匹配类型 ("id", "name", "model")
func FindAllMatchingProviders(idOrNameOrModel string) ([]*LLMProvider, string) {
	config := LoadConfig()
	var matches []*LLMProvider

	// 优先匹配 ID（唯一）
	for i := range config.LLMProviders {
		if config.LLMProviders[i].ID == idOrNameOrModel {
			return []*LLMProvider{&config.LLMProviders[i]}, "id"
		}
	}

	// 然后匹配 Name（唯一）
	for i := range config.LLMProviders {
		if config.LLMProviders[i].Name == idOrNameOrModel {
			return []*LLMProvider{&config.LLMProviders[i]}, "name"
		}
	}

	// 最后匹配 Model（可能多个）
	for i := range config.LLMProviders {
		if config.LLMProviders[i].Model == idOrNameOrModel {
			matches = append(matches, &config.LLMProviders[i])
		}
	}

	if len(matches) > 0 {
		return matches, "model"
	}

	return nil, ""
}

// generateLLMID 根据 provider、name 和 model 生成唯一的 ID
// 返回前 16 个字符的十六进制哈希值作为 ID
func generateLLMID(provider, name, model string) string {
	var input string
	// 使用 name + provider + model 确保唯一性
	input = fmt.Sprintf("%s:%s:%s", provider, name, model)

	// 生成 SHA256 哈希
	hash := sha256.Sum256([]byte(input))
	// 返回前 16 个字符的十六进制字符串作为 ID
	return hex.EncodeToString(hash[:])[:16]
}

