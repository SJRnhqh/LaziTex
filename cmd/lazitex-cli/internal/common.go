// cmd/lazitex-cli/internal/common.go
// 通用业务命令解析

package internal

import (
	// 外部包
	"errors"
	"fmt"
	"strings"

	// 内部包
	tasks "github.com/SJRnhqh/lazitex/cmd/lazitex-cli/tasks"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// I18n 模式常量
const (
	I18nModeMsg  = "msg"  // 普通命令行模式
	I18nModeRepl = "repl" // REPL 模式
)

// ============================================================================
// 泛型 Action 命令解析框架
// ============================================================================

// ActionType 约束：Action 类型必须是整数类型，且 Unknown 值必须为 0
type ActionType interface {
	~int
}

// ActionErrorConfig 错误配置结构，用于配置不同命令的错误处理
type ActionErrorConfig struct {
	// 哨兵错误
	NoAction      error
	MultiActions  error
	UnknownFlag   error
	ExtraArgument error

	// 创建 UnknownFlagError 的工厂函数
	// 参数是未知标志字符串，返回对应的错误
	NewUnknownFlagError func(flag string) error

	// 创建 ExtraArgumentError 的工厂函数
	// 参数是额外参数字符串，返回对应的错误
	NewExtraArgumentError func(arg string) error

	// 错误类型检查函数：检查错误是否是 UnknownFlagError 类型
	IsUnknownFlagError func(err error) (string, bool)

	// 错误类型检查函数：检查错误是否是 ExtraArgumentError 类型
	IsExtraArgumentError func(err error) (string, bool)

	// 国际化键前缀（如 "ollama" 或 "latex"）
	I18nPrefix string
}

// getActionErrorType 获取错误类型标识符（泛型版本）
func getActionErrorType(err error, config ActionErrorConfig) string {
	if err == nil {
		return ""
	}

	// 先检查是否是自定义 ExtraArgumentError 类型
	if config.IsExtraArgumentError != nil {
		if arg, ok := config.IsExtraArgumentError(err); ok {
			_ = arg // 参数信息可用于日志等
			return "err_extra_argument"
		}
	}

	// 先检查是否是自定义 UnknownFlagError 类型
	if config.IsUnknownFlagError != nil {
		if flag, ok := config.IsUnknownFlagError(err); ok {
			_ = flag // 标志信息可用于日志等
			return "err_unknown_flag"
		}
	}

	// 检查哨兵错误
	if errors.Is(err, config.NoAction) {
		return "err_no_action"
	}
	if errors.Is(err, config.MultiActions) {
		return "err_multi_actions"
	}
	if errors.Is(err, config.UnknownFlag) {
		return "err_unknown_flag"
	}
	if errors.Is(err, config.ExtraArgument) {
		return "err_extra_argument"
	}

	return ""
}

// GetActionErrorI18nKey 根据错误配置返回对应的国际化键（泛型版本）
func GetActionErrorI18nKey(err error, config ActionErrorConfig) string {
	errorType := getActionErrorType(err, config)
	if errorType == "" {
		return ""
	}
	return config.I18nPrefix + "." + errorType
}

// GetSupportedFlags 返回所有支持的标志（泛型版本）
func GetSupportedFlags[Action ActionType](flagMap map[string]Action) []string {
	flags := make([]string, 0, len(flagMap))
	for flag := range flagMap {
		flags = append(flags, flag)
	}
	// 简单排序：短标志在前，长标志在后
	// 如果需要更复杂的排序，可以使用 sort.Strings(flags)
	return flags
}

// GetActionForFlag 返回指定标志对应的动作（泛型版本）
func GetActionForFlag[Action ActionType](flagMap map[string]Action, flag string) (Action, bool) {
	action, ok := flagMap[flag]
	return action, ok
}

// BuildActionToFlagsMap 构建动作到标志的反向映射（泛型版本）
func BuildActionToFlagsMap[Action ActionType](flagMap map[string]Action) map[Action][]string {
	result := make(map[Action][]string, len(flagMap))
	for flag, action := range flagMap {
		result[action] = append(result[action], flag)
	}
	return result
}

// GetFlagsForAction 返回指定动作对应的所有标志（泛型版本）
func GetFlagsForAction[Action ActionType](actionToFlagsMap map[Action][]string, action Action) []string {
	if flags, ok := actionToFlagsMap[action]; ok {
		return flags
	}
	return nil
}

// IsValidFlag 检查给定的字符串是否是有效的标志（泛型版本）
func IsValidFlag[Action ActionType](flagMap map[string]Action, flag string) bool {
	_, ok := flagMap[flag]
	return ok
}

// ParseActionArgs 解析 Action 命令参数（泛型版本）
// 确保且仅有一个动作
//
// 参数：
//   - args: 要解析的参数列表（不包含命令本身）
//   - flagMap: 标志到动作的映射
//   - unknownAction: Unknown 动作值（通常为 0）
//   - config: 错误配置
//
// 返回：
//   - Action: 解析得到的动作，如果解析失败返回 unknownAction
//   - error: 解析错误
func ParseActionArgs[Action ActionType](
	args []string,
	flagMap map[string]Action,
	unknownAction Action,
	config ActionErrorConfig,
) (Action, error) {
	var foundAction Action
	actionCount := 0

	for _, arg := range args {
		if action, ok := flagMap[arg]; ok {
			// 找到有效的 action 标志
			if actionCount == 0 {
				foundAction = action
			} else if foundAction != action {
				// 检测到不同的 action，说明用户指定了多个不同的 action
				return unknownAction, config.MultiActions
			}
			actionCount++
		} else if strings.HasPrefix(arg, "-") {
			// 未知的标志，返回包含具体标志信息的错误
			if config.NewUnknownFlagError != nil {
				return unknownAction, config.NewUnknownFlagError(arg)
			}
			return unknownAction, config.UnknownFlag
		} else {
			// 非 flag 参数：这些命令不接受额外的非标志参数
			if config.NewExtraArgumentError != nil {
				return unknownAction, config.NewExtraArgumentError(arg)
			}
			return unknownAction, config.ExtraArgument
		}
	}

	// 验证：必须指定且只能指定一个 action
	if actionCount == 0 {
		return unknownAction, config.NoAction
	}

	return foundAction, nil
}

// HandleActionCommand 统一的 Action 命令处理入口（泛型版本）
func HandleActionCommand[Action ActionType](
	args []string,
	usageKey string,
	parseFunc func([]string) (Action, error),
	errorMapper func(error) string,
	handlers map[Action]func(),
) bool {
	action, err := parseFunc(args)
	if !HandleCommandError(err, errorMapper, usageKey) {
		return false
	}
	if handlers == nil {
		fmt.Println(lang.T(usageKey))
		return false
	}
	if handler, ok := handlers[action]; ok {
		handler()
		return true
	}
	return false
}

// HandleCommandError 通用的命令错误处理函数
// 用于统一处理命令解析错误，显示错误消息和使用说明
//
// 参数：
//   - err: 解析错误，如果为 nil 则不执行任何操作
//   - mapper: 错误映射器函数，用于将错误转换为国际化键（可以为 nil）
//   - usageKey: 使用说明的国际化键（如 "msg.build_usage"）
//
// 返回：
//   - bool: 如果没有错误返回 true，有错误则返回 false（便于在 Handler 中使用 return）
//
// 示例：
//
//	if err != nil {
//	    return HandleCommandError(err, GetOllamaErrorI18nKey, "msg.ollama_usage")
//	}
func HandleCommandError(err error, mapper func(error) string, usageKey string) bool {
	if err == nil {
		return true // 没有错误，成功
	}

	// 如果提供了映射器，尝试获取错误消息
	if mapper != nil {
		if i18nKey := mapper(err); i18nKey != "" {
			// 检查错误消息是否包含格式字符串（%s, %v 等）
			msg := lang.T(i18nKey)
			if strings.Contains(msg, "%") {
				// 尝试从错误中提取参数值
				// 检查是否是 ExtraArgumentError 或 UnknownFlagError
				if argErr, ok := err.(interface{ GetArg() string }); ok {
					fmt.Printf(msg+"\n", argErr.GetArg())
				} else if flagErr, ok := err.(interface{ GetFlag() string }); ok {
					fmt.Printf(msg+"\n", flagErr.GetFlag())
				} else {
					// 如果无法提取参数，直接显示消息
					fmt.Println(msg)
				}
			} else {
				fmt.Println(msg)
			}
		}
	}

	// 显示使用说明
	fmt.Println(lang.T(usageKey))
	return false
}

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
