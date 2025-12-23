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


func SetAllActionHandlers() {
	SetOllamaActionHandlers(map[OllamaAction]func(){
		OllamaCheck:     tasks.CheckOllama,
		OllamaInstall:   tasks.InstallOllama,
		OllamaUninstall: tasks.UninstallOllama,
		OllamaStatus:    tasks.StatusOllama,
	})
}