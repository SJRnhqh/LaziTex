// core/ai/providers/interface.go
// LLM Provider 接口定义

package providers

// Provider 接口定义
type Provider interface {
	Connectivity(model string, config map[string]string) (bool, error)
	Generate(model string, prompt string, config map[string]string) (string, error)
}

