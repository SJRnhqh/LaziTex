// cmd/lazitex-cli/internal/latex.go
// LaTeX 管理相关业务命令解析

package internal

import (
	// 外部包
	"errors"
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

// 错误定义
var (
	ErrLaTeXNoAction      = errors.New("err_no_action")
	ErrLaTeXMultiActions  = errors.New("err_multi_actions")
	ErrLaTeXUnknownFlag   = errors.New("err_unknown_flag")
	ErrLaTeXExtraArgument = errors.New("err_extra_argument")
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

// LaTeXExtraArgumentError 包含额外参数的具体信息
type LaTeXExtraArgumentError struct {
	Arg string
}

func (e *LaTeXExtraArgumentError) Error() string {
	return "err_extra_argument"
}

func (e *LaTeXExtraArgumentError) Unwrap() error {
	return ErrLaTeXExtraArgument
}

// GetArg 返回额外参数的值
func (e *LaTeXExtraArgumentError) GetArg() string {
	return e.Arg
}

// GetFlag 返回标志的值（用于 UnknownFlagError）
func (e *LaTeXUnknownFlagError) GetFlag() string {
	return e.Flag
}

// LaTeXActionFlagMap 将命令行标志映射到对应的动作
var LaTeXActionFlagMap = map[string]LaTeXAction{
	"-c":          LaTeXCheck,
	"--check":     LaTeXCheck,
	"-i":          LaTeXInstall,
	"--install":   LaTeXInstall,
	"-u":          LaTeXUninstall,
	"--uninstall": LaTeXUninstall,
}

// LaTeXActionToFlagsMap 预先构建的动作到标志的映射
var LaTeXActionToFlagsMap map[LaTeXAction][]string

// latexErrorConfig 错误配置
var latexErrorConfig = ActionErrorConfig{
	NoAction:              ErrLaTeXNoAction,
	MultiActions:          ErrLaTeXMultiActions,
	UnknownFlag:           ErrLaTeXUnknownFlag,
	ExtraArgument:         ErrLaTeXExtraArgument,
	NewUnknownFlagError:   func(flag string) error { return &LaTeXUnknownFlagError{Flag: flag} },
	NewExtraArgumentError: func(arg string) error { return &LaTeXExtraArgumentError{Arg: arg} },
	IsUnknownFlagError: func(err error) (string, bool) {
		if flagErr, ok := err.(*LaTeXUnknownFlagError); ok {
			return flagErr.GetFlag(), true
		}
		return "", false
	},
	IsExtraArgumentError: func(err error) (string, bool) {
		if argErr, ok := err.(*LaTeXExtraArgumentError); ok {
			return argErr.GetArg(), true
		}
		return "", false
	},
	I18nPrefix: "latex",
}

func init() {
	LaTeXActionToFlagsMap = BuildActionToFlagsMap(LaTeXActionFlagMap)
}

// GetLaTeXErrorI18nKey 根据错误消息返回对应的国际化键
func GetLaTeXErrorI18nKey(err error) string {
	return GetActionErrorI18nKey(err, latexErrorConfig)
}

// GetSupportedLaTeXFlags 返回所有支持的 latex 命令标志
func GetSupportedLaTeXFlags() []string {
	return GetSupportedFlags(LaTeXActionFlagMap)
}

// GetLaTeXActionForFlag 返回指定标志对应的动作
func GetLaTeXActionForFlag(flag string) (LaTeXAction, bool) {
	return GetActionForFlag(LaTeXActionFlagMap, flag)
}

// GetLaTeXFlagsForAction 返回指定动作对应的所有标志
func GetLaTeXFlagsForAction(action LaTeXAction) []string {
	return GetFlagsForAction(LaTeXActionToFlagsMap, action)
}

// IsValidLaTeXFlag 检查给定的字符串是否是有效的 latex 命令标志
func IsValidLaTeXFlag(flag string) bool {
	return IsValidFlag(LaTeXActionFlagMap, flag)
}

// ParseLaTeXArgs 解析 -x/--latex 后的参数，确保且仅有一个动作
func ParseLaTeXArgs(args []string) (LaTeXAction, error) {
	return ParseActionArgs(args, LaTeXActionFlagMap, LaTeXUnknown, latexErrorConfig)
}

// LaTeXActionHandlers 预定义的动作处理函数映射
var LaTeXActionHandlers map[LaTeXAction]func()

// SetLaTeXActionHandlers 设置动作处理函数（初始化时调用一次）
func SetLaTeXActionHandlers(handlers map[LaTeXAction]func()) {
	LaTeXActionHandlers = handlers
}

// HandleLaTeXCommand 统一的latex 命令处理入口
func HandleLaTeXCommand(args []string, mode, usageKey string) bool {
	return HandleActionCommand(
		args,
		usageKey,
		ParseLaTeXArgs,
		GetLaTeXErrorI18nKey,
		LaTeXActionHandlers,
	)
}
