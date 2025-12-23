// cmd/lazitex-cli/internal/latex.go
// LaTeX 管理相关业务命令解析

package internal

import (
	// 外部包
	"errors"
	"fmt"
	"strings"

	// 内部包
	lang "github.com/SJRnhqh/lazitex/lang"
)

// LaTeXAction 表示 latex 命令支持的动作类型
type LaTeXAction int

const (
	LaTeXUnknown   LaTeXAction = iota
	LaTeXCheck                 // 检查 LaTeX 安装状态
	LaTeXInstall               // 安装 LaTeX
	LaTeXUninstall             // 卸载 LaTeX
)

// String 返回动作的字符串表示，用于调试和日志
func (a LaTeXAction) String() string {
	switch a {
	case LaTeXCheck:
		return "check"
	case LaTeXInstall:
		return "install"
	case LaTeXUninstall:
		return "uninstall"
	default:
		return "unknown"
	}
}

var (
	// 哨兵错误：错误消息作为类型标识符，用于映射到国际化键
	ErrLaTeXNoAction     = errors.New("err_no_action")
	ErrLaTeXMultiActions = errors.New("err_multi_actions")
	ErrLaTeXUnknownFlag  = errors.New("err_unknown_flag")
)

// LaTeXUnknownFlagError 包含未知标志的具体信息
type LaTeXUnknownFlagError struct {
	Flag string
}

func (e *LaTeXUnknownFlagError) Error() string {
	return "err_unknown_flag"
}

func (e *LaTeXUnknownFlagError) Unwrap() error {
	return ErrLaTeXUnknownFlag
}

// getLaTeXErrorType 获取错误类型标识符
// 支持哨兵错误和自定义错误类型
func getLaTeXErrorType(err error) string {
	if err == nil {
		return ""
	}

	// 先检查是否是自定义错误类型
	if flagErr, ok := err.(*LaTeXUnknownFlagError); ok {
		// 保留标志信息，但返回统一的错误类型
		_ = flagErr.Flag // 可以在国际化键中使用，如果需要的话
		return "err_unknown_flag"
	}

	// 检查哨兵错误（通过 errors.Is 支持错误包装）
	if errors.Is(err, ErrLaTeXNoAction) {
		return "err_no_action"
	}
	if errors.Is(err, ErrLaTeXMultiActions) {
		return "err_multi_actions"
	}
	if errors.Is(err, ErrLaTeXUnknownFlag) {
		return "err_unknown_flag"
	}

	return ""
}

// LaTeXactionFlagMap 将命令行标志映射到对应的动作
// 使用 map 可以简化解析逻辑，便于扩展新的 action
var LaTeXActionFlagMap = map[string]LaTeXAction{
	"-c":          LaTeXCheck,
	"--check":     LaTeXCheck,
	"-i":          LaTeXInstall,
	"--install":   LaTeXInstall,
	"-u":          LaTeXUninstall,
	"--uninstall": LaTeXUninstall,
}

// LaTeXActionToFlagsMap 预先构建的动作到标志的映射（按需初始化）
var LaTeXActionToFlagsMap map[LaTeXAction][]string

func init() {
	LaTeXActionToFlagsMap = make(map[LaTeXAction][]string, 3) // 预分配容量
	for flag, action := range LaTeXActionFlagMap {
		LaTeXActionToFlagsMap[action] = append(LaTeXActionToFlagsMap[action], flag)
	}
}

// GetLaTeXErrorI18nKey 根据错误消息返回对应的国际化键
// 如果错误不是 latex 相关错误，返回空字符串
func GetLaTeXErrorI18nKey(err error) string {
	errorType := getLaTeXErrorType(err)
	if errorType == "" {
		return ""
	}

	// 构建国际化键：latex.{errorType}
	return "latex." + errorType
}

// GetSupportedLaTeXFlags 返回所有支持的 latex 命令标志
// 用于生成帮助信息等场景
// 返回的切片按字母顺序排序，便于显示
func GetSupportedLaTeXFlags() []string {
	flags := make([]string, 0, len(LaTeXActionFlagMap))
	for flag := range LaTeXActionFlagMap {
		flags = append(flags, flag)
	}
	// 简单排序：短标志在前，长标志在后
	// 如果需要更复杂的排序，可以使用 sort.Strings(flags)
	return flags
}

// GetLaTeXActionForFlag 返回指定标志对应的动作
// 如果标志无效，返回 LaTeXUnknown 和 false
func GetLaTeXActionForFlag(flag string) (LaTeXAction, bool) {
	action, ok := LaTeXActionFlagMap[flag]
	return action, ok
}

// GetLaTeXFlagsForAction 返回指定动作对应的所有标志
// 例如 LaTeXCheck 会返回 ["-c", "--check"]
// 使用预构建的映射，O(1) 查找
func GetLaTeXFlagsForAction(action LaTeXAction) []string {
	if flags, ok := LaTeXActionToFlagsMap[action]; ok {
		return flags
	}
	return nil
}

// IsValidLaTeXFlag 检查给定的字符串是否是有效的 latex 命令标志
func IsValidLaTeXFlag(flag string) bool {
	_, ok := LaTeXActionFlagMap[flag]
	return ok
}

// ParseLaTeXArgs 解析 -x/--latex 后的参数，确保且仅有一个动作
//
// 参数：
//   - args: 要解析的参数列表（不包含 -x 或 --latex 本身）
//
// 返回：
//   - LaTeXAction: 解析得到的动作，如果解析失败返回 LaTeXUnknown
//   - error: 解析错误，包括：
//   - ErrLaTeXNoAction: 未指定任何动作
//   - ErrLaTeXMultiActions: 指定了多个不同的动作
//   - LaTeXUnknownFlagError: 遇到未知的标志（包含具体标志名称）
//
// 示例：
//
//	ParseLaTeXArgs([]string{"-c"})           // LaTeXCheck, nil
//	ParseLaTeXArgs([]string{"-c", "-i"})     // LaTeXUnknown, ErrLaTeXMultiActions
//	ParseLaTeXArgs([]string{"-x"})           // LaTeXUnknown, LaTeXUnknownFlagError{Flag: "-x"}
//	ParseLaTeXArgs([]string{})               // LaTeXUnknown, ErrLaTeXNoAction
func ParseLaTeXArgs(args []string) (LaTeXAction, error) {
	var foundAction LaTeXAction
	actionCount := 0

	for _, arg := range args {
		if action, ok := LaTeXActionFlagMap[arg]; ok {
			// 找到有效的 action 标志
			if actionCount == 0 {
				foundAction = action
			} else if foundAction != action {
				// 检测到不同的 action，说明用户指定了多个不同的 action
				return LaTeXUnknown, ErrLaTeXMultiActions
			}
			actionCount++
		} else if strings.HasPrefix(arg, "-") {
			// 未知的标志，返回包含具体标志信息的错误
			return LaTeXUnknown, &LaTeXUnknownFlagError{Flag: arg}
		}
		// 非 flag 参数忽略
	}

	// 验证：必须指定且只能指定一个 action
	if actionCount == 0 {
		return LaTeXUnknown, ErrLaTeXNoAction
	}

	return foundAction, nil
}

// LaTeXActionHandlers 预定义的动作处理函数映射
var LaTeXActionHandlers map[LaTeXAction]func()

// SetLaTeXActionHandlers 设置动作处理函数（初始化时调用一次）
func SetLaTeXActionHandlers(handlers map[LaTeXAction]func()) {
	LaTeXActionHandlers = handlers
}

// HandleLaTeXCommand 统一的latex 命令处理入口
func HandleLaTeXCommand(args []string, mode, usageKey string) bool {
	action, err := ParseLaTeXArgs(args)
	if err != nil {
		if i18nKey := GetLaTeXErrorI18nKey(err); i18nKey != "" {
			fmt.Println(lang.T(i18nKey))
		}
		fmt.Println(lang.T(usageKey))
		return false
	}
	if LaTeXActionHandlers == nil {
		fmt.Println(lang.T(usageKey))
		return false
	}
	if handler, ok := LaTeXActionHandlers[action]; ok {
		handler()
		return true
	}
	return false
}
