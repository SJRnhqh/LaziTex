// core/ai/llm/registry.go
// LLM 注册表管理和 Provider 注册的核心业务逻辑
// 包含：
// 1. Provider 实现的注册表（工厂函数管理）
// 2. LLM Provider 注册的核心业务逻辑（创建 -> 验证 -> 更新状态 -> 保存）

package llm

import (
	"errors"
	"fmt"
	"sync"
	"time"

	cfg "github.com/SJRnhqh/lazitex/config"
	providers "github.com/SJRnhqh/lazitex/core/ai/providers"
)

// 错误定义
var (
	ErrProviderNil      = errors.New("provider cannot be nil")
	ErrConnectivityTest = errors.New("connectivity test failed")
	ErrLLMUnavailable   = errors.New("LLM unavailable")
	ErrSaveConfigFailed = errors.New("save config failed")
)

var (
	// providerRegistry Provider 注册表（存储工厂函数）
	providerRegistry = make(map[string]func() providers.Provider)
	registryMu       sync.RWMutex // 保护注册表的互斥锁
)

const (
	// timeFormatRFC3339 时间格式常量
	timeFormatRFC3339 = time.RFC3339
)

// RegisterProvider 注册 Provider 实现
func RegisterProvider(name string, factory func() providers.Provider) {
	registryMu.Lock()
	defer registryMu.Unlock()
	providerRegistry[name] = factory
}

// GetProvider 根据名称获取 Provider 实例（工厂方法）
func GetProvider(name string) providers.Provider {
	registryMu.RLock()
	defer registryMu.RUnlock()

	factory, ok := providerRegistry[name]
	if !ok {
		return nil
	}
	return factory()
}

// init 初始化：注册所有内置 Provider
func init() {
	// 注册 Ollama Provider
	RegisterProvider("ollama", func() providers.Provider {
		return &providers.OllamaProvider{}
	})
	// 未来可以注册其他 Provider
}

// --- LLM Provider 注册的核心业务逻辑 ---

// RegisterProviderOptions 注册选项
type RegisterProviderOptions struct {
	// SkipVerification 是否跳过连通性验证（默认 false，即需要验证）
	SkipVerification bool
	// UpdateVerifiedStatus 是否更新验证状态（默认 true）
	UpdateVerifiedStatus bool
}

// DefaultRegisterOptions 返回默认的注册选项
func DefaultRegisterOptions() *RegisterProviderOptions {
	return &RegisterProviderOptions{
		SkipVerification:     false,
		UpdateVerifiedStatus: true,
	}
}

// RegisterProviderResult 注册结果
type RegisterProviderResult struct {
	Provider   *cfg.LLMProvider
	Verified   bool
	Registered bool
}

// RegisterLLMProvider 注册新的 LLM Provider（核心业务逻辑）
//
// 流程：
// 1. 创建 Provider 配置（由调用方提供）
// 2. 测试连通性（如果未跳过）
// 3. 更新验证状态（如果启用）
// 4. 保存到配置
//
// 参数：
//   - provider: 要注册的 Provider 配置（必须包含 Provider、Name、Model）
//   - options: 注册选项（可选，nil 时使用默认值）
//
// 返回：
//   - RegisterProviderResult: 注册结果
//   - error: 错误信息
func RegisterLLMProvider(provider *cfg.LLMProvider, options *RegisterProviderOptions) (*RegisterProviderResult, error) {
	if provider == nil {
		return nil, ErrProviderNil
	}

	// 设置默认选项
	if options == nil {
		options = DefaultRegisterOptions()
	}

	result := &RegisterProviderResult{
		Provider: provider,
	}

	// 步骤1：测试连通性（如果未跳过）
	if !options.SkipVerification {
		if err := verifyConnectivity(provider, result); err != nil {
			return nil, err
		}

		// 步骤2：更新验证状态（如果启用）
		if options.UpdateVerifiedStatus {
			updateVerifiedStatus(provider)
		}
	}

	// 步骤3：保存到配置
	if err := saveProviderConfig(provider); err != nil {
		return nil, err
	}

	result.Registered = true
	return result, nil
}

// verifyConnectivity 验证 Provider 连通性
func verifyConnectivity(provider *cfg.LLMProvider, result *RegisterProviderResult) error {
	verified, err := LLMConnectivity(provider)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrConnectivityTest, err)
	}

	result.Verified = verified

	// 如果验证失败，返回错误
	if !verified {
		return fmt.Errorf("%w: %s/%s", ErrLLMUnavailable, provider.Provider, provider.Model)
	}

	return nil
}

// updateVerifiedStatus 更新 Provider 的验证状态
func updateVerifiedStatus(provider *cfg.LLMProvider) {
	provider.Verified = true
	provider.VerifiedAt = time.Now().Format(timeFormatRFC3339)
}

// saveProviderConfig 保存 Provider 配置
func saveProviderConfig(provider *cfg.LLMProvider) error {
	if err := cfg.AddLLMProvider(provider); err != nil {
		return fmt.Errorf("%w: %w", ErrSaveConfigFailed, err)
	}
	return nil
}
