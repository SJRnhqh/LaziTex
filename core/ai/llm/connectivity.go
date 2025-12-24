// core/ai/llm/connectivity.go
// LLM 连通性测试的统一入口

package llm

import (
	"fmt"

	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// LLMConnectivity 测试 LLM 是否可用（统一入口）
func LLMConnectivity(llmprovider *cfg.LLMProvider) (bool, error) {
	// 如果已禁用，直接返回 false（不进行连通性测试）
	if !llmprovider.Enabled {
		return false, nil
	}

	// 从注册表获取 Provider 实现
	provider := GetProvider(llmprovider.Provider)
	if provider == nil {
		return false, fmt.Errorf(lang.T("msg.llm.provider.unsupported"), llmprovider.Provider)
	}

	// 调用 provider 的测试方法
	return provider.Connectivity(llmprovider.Model, llmprovider.Config)
}
