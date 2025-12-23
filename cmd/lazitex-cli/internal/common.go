// cmd/lazitex-cli/internal/common.go
// 通用业务命令解析

package internal

import (
	// 内部包
	tasks "github.com/SJRnhqh/lazitex/cmd/lazitex-cli/tasks"
)

// I18n 模式常量
const (
	I18nModeMsg  = "msg"  // 普通命令行模式
	I18nModeRepl = "repl" // REPL 模式
)



// SetAllActionHandlers 设置所有动作处理函数
func SetAllActionHandlers() {
	SetOllamaActionHandlers(map[OllamaAction]func(){
		OllamaCheck:     tasks.CheckOllama,
		OllamaInstall:   tasks.InstallOllama,
		OllamaUninstall: tasks.UninstallOllama,
		OllamaStatus:    tasks.StatusOllama,
	})

	SetLaTeXActionHandlers(map[LaTeXAction]func(){
		LaTeXCheck:     tasks.CheckLaTeX,
		LaTeXInstall:   tasks.InstallLaTeX,
		LaTeXUninstall: tasks.UninstallLaTeX,
	})
}
