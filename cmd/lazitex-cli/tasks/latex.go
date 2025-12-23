// cmd/lazitex-cli/tasks/latex.go
// LaTeX 管理编排任务

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

// getLaTeXManager 根据当前平台创建对应的 LaTeX Manager
func getLaTeXManager() (core.LaTeXManager, error) {
	switch runtime.GOOS {
	case "linux":
		return linux.NewLaTeXManager(), nil
	case "darwin":
		return mac.NewLaTeXManager(), nil
	case "windows":
		return win.NewLaTeXManager(), nil
	default:
		return nil, fmt.Errorf(lang.T("msg.unsupported_os"), runtime.GOOS)
	}
}

// CheckLaTeX 检查 LaTeX 环境
func CheckLaTeX() {
	manager, err := getLaTeXManager()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 执行检测
	env := core.CheckLaTeXEnvironment(manager)
	env.PrintEnvironment()
}

// InstallLaTeX 安装 LaTeX 环境
func InstallLaTeX() {
	manager, err := getLaTeXManager()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 执行安装
	if err := core.InstallLaTeXEnvironment(manager); err != nil {
		fmt.Printf(lang.T("msg.install_failed")+": %v\n", err)
		return
	}
	// 注意：成功消息由安装器内部输出，这里不需要再输出
}

// UninstallLaTeX 卸载 LaTeX 环境
func UninstallLaTeX() {
	manager, err := getLaTeXManager()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 执行卸载
	if err := core.UninstallLaTeXEnvironment(manager); err != nil {
		fmt.Printf(lang.T("msg.uninstall_failed")+": %v\n", err)
		return
	}
	// 注意：成功消息由卸载器内部输出，这里不需要再输出
}
