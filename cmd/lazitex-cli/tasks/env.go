// cmd/lazitex-cli/tasks/env.go
// LaTex环境相关业务管理

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

// 检查 LaTex 环境
func CheckEnvironment() {
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
	env := core.CheckLaTexEnvironment(checker)
	env.PrintEnvironment()
}

// InstallLaTexEnvironment 安装 LaTex 环境
func InstallLaTexEnvironment() {
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
	if err := core.InstallLaTexEnvironment(installer); err != nil {
		fmt.Printf(lang.T("msg.install_failed")+": %v\n", err)
		return
	}
	// 注意：成功消息由安装器内部输出，这里不需要再输出
}

// UninstallLaTexEnvironment 卸载 LaTex 环境
func UninstallLaTexEnvironment() {
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
	if err := core.UninstallLaTexEnvironment(installer); err != nil {
		fmt.Printf(lang.T("msg.uninstall_failed")+": %v\n", err)
		return
	}
	// 注意：成功消息由卸载器内部输出，这里不需要再输出
}
