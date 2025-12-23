// cmd/lazitex-cli/internal/ollama.go
// Ollama 管理相关业务命令解析

package internal

import (
	// 外部包
	"errors"
	"fmt"
	"strings"

	// 内部包
	lang "github.com/SJRnhqh/lazitex/lang"
)

// OllamaAction 表示 ollama 命令支持的动作类型
type OllamaAction int

const (
	OllamaUnknown   OllamaAction = iota
	OllamaCheck                  // 检查 Ollama 安装状态
	OllamaInstall                // 安装 Ollama
	OllamaUninstall              // 卸载 Ollama
	OllamaStatus                 // 查看 Ollama 状态
)

// String 返回动作的字符串表示，用于调试和日志
func (a OllamaAction) String() string {
	switch a {
	case OllamaCheck:
		return "check"
	case OllamaInstall:
		return "install"
	case OllamaUninstall:
		return "uninstall"
	case OllamaStatus:
		return "status"
	default:
		return "unknown"
	}
}

var (
	// 哨兵错误：错误消息作为类型标识符，用于映射到国际化键
	ErrOllamaNoAction     = errors.New("err_no_action")
	ErrOllamaMultiActions = errors.New("err_multi_actions")
	ErrOllamaUnknownFlag  = errors.New("err_unknown_flag")
)

// OllamaUnknownFlagError 包含未知标志的具体信息
type OllamaUnknownFlagError struct {
	Flag string
}

func (e *OllamaUnknownFlagError) Error() string {
	return "err_unknown_flag"
}

func (e *OllamaUnknownFlagError) Unwrap() error {
	return ErrOllamaUnknownFlag
}

// getOllamaErrorType 获取错误类型标识符
// 支持哨兵错误和自定义错误类型
func getOllamaErrorType(err error) string {
	if err == nil {
		return ""
	}

	// 先检查是否是自定义错误类型
	if flagErr, ok := err.(*OllamaUnknownFlagError); ok {
		// 保留标志信息，但返回统一的错误类型
		_ = flagErr.Flag // 可以在国际化键中使用，如果需要的话
		return "err_unknown_flag"
	}

	// 检查哨兵错误（通过 errors.Is 支持错误包装）
	if errors.Is(err, ErrOllamaNoAction) {
		return "err_no_action"
	}
	if errors.Is(err, ErrOllamaMultiActions) {
		return "err_multi_actions"
	}
	if errors.Is(err, ErrOllamaUnknownFlag) {
		return "err_unknown_flag"
	}

	return ""
}

// OllamaActionFlagMap 将命令行标志映射到对应的动作
// 使用 map 可以简化解析逻辑，便于扩展新的 action
var OllamaActionFlagMap = map[string]OllamaAction{
	"-c":          OllamaCheck,
	"--check":     OllamaCheck,
	"-i":          OllamaInstall,
	"--install":   OllamaInstall,
	"-u":          OllamaUninstall,
	"--uninstall": OllamaUninstall,
	"-s":          OllamaStatus,
	"--status":    OllamaStatus,
}

// OllamaActionToFlagsMap 预先构建的动作到标志的映射（按需初始化）
var OllamaActionToFlagsMap map[OllamaAction][]string

func init() {
	OllamaActionToFlagsMap = make(map[OllamaAction][]string, 4) // 预分配容量
	for flag, action := range OllamaActionFlagMap {
		OllamaActionToFlagsMap[action] = append(OllamaActionToFlagsMap[action], flag)
	}
}

// GetOllamaErrorI18nKey 根据错误消息返回对应的国际化键
// 如果错误不是 ollama 相关错误，返回空字符串
func GetOllamaErrorI18nKey(err error) string {
	errorType := getOllamaErrorType(err)
	if errorType == "" {
		return ""
	}

	// 构建国际化键：ollama.{errorType}
	return "ollama." + errorType
}

// GetSupportedOllamaFlags 返回所有支持的 ollama 命令标志
// 用于生成帮助信息等场景
// 返回的切片按字母顺序排序，便于显示
func GetSupportedOllamaFlags() []string {
	flags := make([]string, 0, len(OllamaActionFlagMap))
	for flag := range OllamaActionFlagMap {
		flags = append(flags, flag)
	}
	// 简单排序：短标志在前，长标志在后
	// 如果需要更复杂的排序，可以使用 sort.Strings(flags)
	return flags
}

// GetOllamaActionForFlag 返回指定标志对应的动作
// 如果标志无效，返回 OllamaUnknown 和 false
func GetOllamaActionForFlag(flag string) (OllamaAction, bool) {
	action, ok := OllamaActionFlagMap[flag]
	return action, ok
}

// GetOllamaFlagsForAction 返回指定动作对应的所有标志
// 例如 OllamaCheck 会返回 ["-c", "--check"]
// 使用预构建的映射，O(1) 查找
func GetOllamaFlagsForAction(action OllamaAction) []string {
	if flags, ok := OllamaActionToFlagsMap[action]; ok {
		return flags
	}
	return nil
}

// IsValidOllamaFlag 检查给定的字符串是否是有效的 ollama 命令标志
func IsValidOllamaFlag(flag string) bool {
	_, ok := OllamaActionFlagMap[flag]
	return ok
}

// ParseOllamaArgs 解析 -o/--ollama 后的参数，确保且仅有一个动作
//
// 参数：
//   - args: 要解析的参数列表（不包含 -o 或 --ollama 本身）
//
// 返回：
//   - OllamaAction: 解析得到的动作，如果解析失败返回 OllamaUnknown
//   - error: 解析错误，包括：
//   - ErrOllamaNoAction: 未指定任何动作
//   - ErrOllamaMultiActions: 指定了多个不同的动作
//   - OllamaUnknownFlagError: 遇到未知的标志（包含具体标志名称）
//
// 示例：
//
//	ParseOllamaArgs([]string{"-c"})           // OllamaCheck, nil
//	ParseOllamaArgs([]string{"-c", "-i"})     // OllamaUnknown, ErrOllamaMultiActions
//	ParseOllamaArgs([]string{"-x"})           // OllamaUnknown, OllamaUnknownFlagError{Flag: "-x"}
//	ParseOllamaArgs([]string{})               // OllamaUnknown, ErrOllamaNoAction
func ParseOllamaArgs(args []string) (OllamaAction, error) {
	var foundAction OllamaAction
	actionCount := 0

	for _, arg := range args {
		if action, ok := OllamaActionFlagMap[arg]; ok {
			// 找到有效的 action 标志
			if actionCount == 0 {
				foundAction = action
			} else if foundAction != action {
				// 检测到不同的 action，说明用户指定了多个不同的 action
				return OllamaUnknown, ErrOllamaMultiActions
			}
			actionCount++
		} else if strings.HasPrefix(arg, "-") {
			// 未知的标志，返回包含具体标志信息的错误
			return OllamaUnknown, &OllamaUnknownFlagError{Flag: arg}
		}
		// 非 flag 参数忽略
	}

	// 验证：必须指定且只能指定一个 action
	if actionCount == 0 {
		return OllamaUnknown, ErrOllamaNoAction
	}

	return foundAction, nil
}

// OllamaActionHandlers 预定义的动作处理函数映射
var OllamaActionHandlers map[OllamaAction]func()

// SetOllamaActionHandlers 设置动作处理函数（初始化时调用一次）
func SetOllamaActionHandlers(handlers map[OllamaAction]func()) {
	OllamaActionHandlers = handlers
}

// HandleOllamaCommand 统一的 ollama 命令处理入口
func HandleOllamaCommand(args []string, mode, usageKey string) bool {
	action, err := ParseOllamaArgs(args)
	if err != nil {
		if i18nKey := GetOllamaErrorI18nKey(err); i18nKey != "" {
			fmt.Println(lang.T(i18nKey))
		}
		fmt.Println(lang.T(usageKey))
		return false
	}
	if OllamaActionHandlers == nil {
		fmt.Println(lang.T(usageKey))
		return false
	}
	if handler, ok := OllamaActionHandlers[action]; ok {
		handler()
		return true
	}
	return false
}
