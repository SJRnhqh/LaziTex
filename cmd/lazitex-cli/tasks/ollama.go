// cmd/lazitex-cli/tasks/ollama.go
// Ollama 管理相关业务

package tasks

import (
	// 外部包
	"fmt"
	"runtime"

	// 内部包
	core "github.com/SJRnhqh/lazitex/core"
	lang "github.com/SJRnhqh/lazitex/lang"
	win "github.com/SJRnhqh/lazitex/target/win"
)

// CheckOllama 检查 Ollama 是否安装
func CheckOllama() {
	var manager core.OllamaManager

	switch runtime.GOOS {
	case "windows":
		manager = win.NewOllamaManager()
	default:
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	installed, version, err := core.CheckOllama(manager)
	if err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.ollama.check_failed"), err)
		return
	}

	if installed {
		if version != "" && version != "installed" {
			fmt.Printf("🦙 %s: %s\n", lang.T("msg.ollama.check_installed"), version)
		} else {
			fmt.Println("🦙 " + lang.T("msg.ollama.check_installed_no_version"))
		}
	} else {
		fmt.Println("🦙 " + lang.T("msg.ollama.not_installed"))
	}
}

// InstallOllama 安装 Ollama
func InstallOllama() {
	var manager core.OllamaManager

	switch runtime.GOOS {
	case "windows":
		manager = win.NewOllamaManager()
	default:
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	if err := core.InstallOllama(manager); err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.ollama.install_failed"), err)
		return
	}
	// 注意：成功消息由安装器内部输出，这里不需要再输出
}

// UninstallOllama 卸载 Ollama
func UninstallOllama() {
	var manager core.OllamaManager

	switch runtime.GOOS {
	case "windows":
		manager = win.NewOllamaManager()
	default:
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	if err := core.UninstallOllama(manager); err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.ollama.uninstall_failed"), err)
		return
	}
	// 注意：成功消息由卸载器内部输出，这里不需要再输出
}
