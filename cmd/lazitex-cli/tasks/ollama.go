// cmd/lazitex-cli/tasks/ollama.go
// Ollama 管理编排任务

package tasks

import (
	// 外部包
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	// 内部包
	core "github.com/SJRnhqh/lazitex/core"
	lang "github.com/SJRnhqh/lazitex/lang"
	linux "github.com/SJRnhqh/lazitex/target/linux"
	mac "github.com/SJRnhqh/lazitex/target/mac"
	win "github.com/SJRnhqh/lazitex/target/win"
)

// getOllamaManager 根据当前平台创建对应的 Ollama Manager
func getOllamaManager() (core.OllamaManager, error) {
	switch runtime.GOOS {
	case "linux":
		return linux.NewOllamaManager(), nil
	case "darwin":
		return mac.NewOllamaManager(), nil
	case "windows":
		return win.NewOllamaManager(), nil
	default:
		return nil, fmt.Errorf(lang.T("msg.unsupported_os"), runtime.GOOS)
	}
}

// formatVersionMessage 格式化版本信息输出
func formatVersionMessage(version string) {
	if version != "" && version != "installed" {
		fmt.Printf("🦙 %s: %s\n", lang.T("msg.ollama.check_installed"), version)
	} else {
		fmt.Println("🦙 " + lang.T("msg.ollama.check_installed_no_version"))
	}
}

// askForConfirmation 请求用户确认
func askForConfirmation(prompt string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// CheckOllama 检查 Ollama 是否安装
func CheckOllama() {
	manager, err := getOllamaManager()
	if err != nil {
		fmt.Println(err)
		return
	}

	installed, version, err := core.CheckOllama(manager)
	if err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.ollama.check_failed"), err)
		return
	}

	if installed {
		formatVersionMessage(version)
	} else {
		fmt.Println(lang.T("msg.ollama.not_installed"))
	}
}

// InstallOllama 安装 Ollama
func InstallOllama() {
	manager, err := getOllamaManager()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 如果已安装则提示并退出
	installed, version, err := core.CheckOllama(manager)
	if err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.ollama.check_failed"), err)
		return
	}
	if installed {
		formatVersionMessage(version)
		return
	}

	// 安装前确认（默认否）
	if !askForConfirmation(lang.T("msg.ollama.confirm_install")) {
		fmt.Println(lang.T("msg.ollama.install_cancelled"))
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
	manager, err := getOllamaManager()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 先检查是否已安装，未安装则直接提示
	installed, _, err := core.CheckOllama(manager)
	if err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.ollama.check_failed"), err)
		return
	}
	if !installed {
		fmt.Println(lang.T("msg.ollama.not_installed"))
		return
	}

	// 卸载前确认（默认否）
	if !askForConfirmation(lang.T("msg.ollama.confirm_uninstall")) {
		fmt.Println(lang.T("msg.ollama.uninstall_cancelled"))
		return
	}

	if err := core.UninstallOllama(manager); err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.ollama.uninstall_failed"), err)
		return
	}
	// 注意：成功消息由卸载器内部输出，这里不需要再输出
}

// StatusOllama 检查 Ollama 服务状态
func StatusOllama() {
	manager, err := getOllamaManager()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = core.StatusOllama(manager)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}
	// 注意：状态信息由 Status 方法内部输出，这里不需要再输出
}
