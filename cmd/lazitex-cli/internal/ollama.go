// cmd/lazitex-cli/internal/ollama.go
// Ollama 管理相关业务命令解析及执行

package internal

import (
	"errors"
	"strings"
)

type OllamaAction int

const (
	OllamaUnknown OllamaAction = iota
	OllamaCheck
	OllamaInstall
	OllamaUninstall
	OllamaStatus
)

var (
	ErrOllamaNoAction     = errors.New("no action")
	ErrOllamaMultiActions = errors.New("multiple actions")
	ErrOllamaUnknownFlag  = errors.New("unknown flag")
)

// ParseOllamaArgs 解析 -o/--ollama 后的参数，确保且仅有一个动作
func ParseOllamaArgs(args []string) (OllamaAction, error) {
	checkMode, installMode, uninstallMode, statusMode := false, false, false, false
	actionCount := 0

	for _, arg := range args {
		switch arg {
		case "-c", "--check":
			checkMode = true
			actionCount++
		case "-i", "--install":
			installMode = true
			actionCount++
		case "-u", "--uninstall":
			uninstallMode = true
			actionCount++
		case "-s", "--status":
			statusMode = true
			actionCount++
		default:
			if strings.HasPrefix(arg, "-") {
				return OllamaUnknown, ErrOllamaUnknownFlag
			}
			// 非 flag 参数忽略
		}
	}

	if actionCount == 0 {
		return OllamaUnknown, ErrOllamaNoAction
	}
	if actionCount > 1 {
		return OllamaUnknown, ErrOllamaMultiActions
	}

	switch {
	case checkMode:
		return OllamaCheck, nil
	case installMode:
		return OllamaInstall, nil
	case uninstallMode:
		return OllamaUninstall, nil
	case statusMode:
		return OllamaStatus, nil
	default:
		return OllamaUnknown, ErrOllamaUnknownFlag
	}
}