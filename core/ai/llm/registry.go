// core/ai/llm/registry.go
// LLM 注册表管理

package llm

import (
	"sync"

	"github.com/SJRnhqh/lazitex/core/ai/providers"
)

var (
	// providerRegistry Provider 注册表（存储工厂函数）
	providerRegistry = make(map[string]func() providers.Provider)
	registryMu       sync.RWMutex // 保护注册表的互斥锁
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