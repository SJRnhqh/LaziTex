// core/ai/llm/connectivity.go
// LLM 连通性测试的统一入口

package llm

import (
	"fmt"

	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// TestLLMConnectivity 测试 LLM 是否可用（统一入口）
func TestLLMConnectivity(providerCfg *cfg.LLMProvider) (bool, error) {
	// 如果已禁用，直接返回 false（不进行连通性测试）
	if !providerCfg.Enabled {
		return false, nil
	}

	// 从注册表获取 Provider 实现
	provider := GetProvider(providerCfg.Provider)
	if provider == nil {
		return false, fmt.Errorf(lang.T("msg.llm.provider.unsupported"), providerCfg.Provider)
	}

	// 调用 provider 的测试方法
	return provider.Connectivity(providerCfg.Model, providerCfg.Config)
}
