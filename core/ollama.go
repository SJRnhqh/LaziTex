// core/ollama.go
// 核心业务：Ollama 管理器接口定义和包装函数

package core

import (
	"fmt"

	lang "github.com/SJRnhqh/lazitex/lang"
)

// OllamaManager 定义 Ollama 管理器接口
// 各平台需要实现此接口以提供平台特定的 Ollama 管理功能
type OllamaManager interface {
	// Check 检查 Ollama 是否安装
	// 返回: (是否安装, 版本信息, 错误)
	Check() (bool, string, error)

	// Install 安装 Ollama
	// 返回: 错误信息（如果安装失败）
	Install() error

	// Uninstall 卸载 Ollama
	// 返回: 错误信息（如果卸载失败）
	Uninstall() error

	// Status 获取 Ollama 状态
	// 返回: 状态信息
	Status() (bool, error)
}

// CheckOllama 检查 Ollama 是否安装
// 使用平台特定的管理器实现进行检查
func CheckOllama(manager OllamaManager) (bool, string, error) {
	if manager == nil {
		return false, "", fmt.Errorf("%s", lang.T("msg.ollama.manager_nil"))
	}
	return manager.Check()
}

// InstallOllama 安装 Ollama
// 使用平台特定的管理器实现进行安装
func InstallOllama(manager OllamaManager) error {
	if manager == nil {
		return fmt.Errorf("%s", lang.T("msg.ollama.manager_nil"))
	}
	return manager.Install()
}

// UninstallOllama 卸载 Ollama
// 使用平台特定的管理器实现进行卸载
func UninstallOllama(manager OllamaManager) error {
	if manager == nil {
		return fmt.Errorf("%s", lang.T("msg.ollama.manager_nil"))
	}
	return manager.Uninstall()
}

// StatusOllama 检查 Ollama 服务状态
// 使用平台特定的管理器实现进行检查
func StatusOllama(manager OllamaManager) (bool, error) {
	if manager == nil {
		return false, fmt.Errorf("%s", lang.T("msg.ollama.manager_nil"))
	}
	return manager.Status()
}