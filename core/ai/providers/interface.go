// core/ai/providers/interface.go
// LLM Provider 接口定义

package providers

// Provider 接口定义
type Provider interface {
	Connectivity(model string, config map[string]string) (bool, error)
}

