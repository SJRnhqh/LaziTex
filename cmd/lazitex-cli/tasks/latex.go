// cmd/lazitex-cli/tasks/env.go
// LaTeX环境相关业务管理

package tasks

import (
	// 外部包
	"fmt"
	"runtime"

	// 内部包
	core "github.com/SJRnhqh/lazitex/core"
	lang "github.com/SJRnhqh/lazitex/lang"
	linux "github.com/SJRnhqh/lazitex/target/linux"
	mac "github.com/SJRnhqh/lazitex/target/mac"
	win "github.com/SJRnhqh/lazitex/target/win"
)

// CheckLaTeX 检查 LaTeX 环境
func CheckLaTeX() {
	// 根据平台创建对应的检查器
	var checker core.EnvironmentChecker

	switch runtime.GOOS {
	case "windows":
		checker = win.NewChecker()
	case "linux":
		checker = linux.NewChecker()
	case "darwin":
		checker = mac.NewChecker()
	default:
		// 不支持的平台，使用一个简单的错误提示
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	// 执行检测
	env := core.CheckLaTeXEnvironment(checker)
	env.PrintEnvironment()
}

// InstallLaTeX 安装 LaTeX 环境
func InstallLaTeX() {
	// 根据平台创建对应的安装器
	var installer core.EnvironmentInstaller

	switch runtime.GOOS {
	case "windows":
		// TODO: 后续实现 Windows 安装器
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	case "linux":
		installer = linux.NewInstaller()
	case "darwin":
		installer = mac.NewInstaller()
	default:
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	// 执行安装
	if err := core.InstallLaTeXEnvironment(installer); err != nil {
		fmt.Printf(lang.T("msg.install_failed")+": %v\n", err)
		return
	}
	// 注意：成功消息由安装器内部输出，这里不需要再输出
}

// UninstallLaTeX 卸载 LaTeX 环境
func UninstallLaTeX() {
	// 根据平台创建对应的安装器
	var installer core.EnvironmentInstaller

	switch runtime.GOOS {
	case "windows":
		// TODO: 后续实现 Windows 安装器
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	case "linux":
		installer = linux.NewInstaller()
	case "darwin":
		installer = mac.NewInstaller()
	default:
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	// 执行卸载
	if err := core.UninstallLaTeXEnvironment(installer); err != nil {
		fmt.Printf(lang.T("msg.uninstall_failed")+": %v\n", err)
		return
	}
	// 注意：成功消息由卸载器内部输出，这里不需要再输出
}
