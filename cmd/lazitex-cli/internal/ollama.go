// cmd/lazitex-cli/internal/ollama.go
// Ollama 管理相关业务命令解析

package internal

import (
	// 外部包
	"errors"
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

// 错误定义
var (
	ErrOllamaNoAction      = errors.New("err_no_action")
	ErrOllamaMultiActions  = errors.New("err_multi_actions")
	ErrOllamaUnknownFlag   = errors.New("err_unknown_flag")
	ErrOllamaExtraArgument = errors.New("err_extra_argument")
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

// OllamaExtraArgumentError 包含额外参数的具体信息
type OllamaExtraArgumentError struct {
	Arg string
}

func (e *OllamaExtraArgumentError) Error() string {
	return "err_extra_argument"
}

func (e *OllamaExtraArgumentError) Unwrap() error {
	return ErrOllamaExtraArgument
}

// GetArg 返回额外参数的值
func (e *OllamaExtraArgumentError) GetArg() string {
	return e.Arg
}

// GetFlag 返回标志的值（用于 UnknownFlagError）
func (e *OllamaUnknownFlagError) GetFlag() string {
	return e.Flag
}

// OllamaActionFlagMap 将命令行标志映射到对应的动作
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

// OllamaActionToFlagsMap 预先构建的动作到标志的映射
var OllamaActionToFlagsMap map[OllamaAction][]string

// ollamaErrorConfig 错误配置
var ollamaErrorConfig = ActionErrorConfig{
	NoAction:              ErrOllamaNoAction,
	MultiActions:          ErrOllamaMultiActions,
	UnknownFlag:           ErrOllamaUnknownFlag,
	ExtraArgument:         ErrOllamaExtraArgument,
	NewUnknownFlagError:   func(flag string) error { return &OllamaUnknownFlagError{Flag: flag} },
	NewExtraArgumentError: func(arg string) error { return &OllamaExtraArgumentError{Arg: arg} },
	IsUnknownFlagError: func(err error) (string, bool) {
		if flagErr, ok := err.(*OllamaUnknownFlagError); ok {
			return flagErr.GetFlag(), true
		}
		return "", false
	},
	IsExtraArgumentError: func(err error) (string, bool) {
		if argErr, ok := err.(*OllamaExtraArgumentError); ok {
			return argErr.GetArg(), true
		}
		return "", false
	},
	I18nPrefix: "ollama",
}

func init() {
	OllamaActionToFlagsMap = BuildActionToFlagsMap(OllamaActionFlagMap)
}

// GetOllamaErrorI18nKey 根据错误消息返回对应的国际化键
func GetOllamaErrorI18nKey(err error) string {
	return GetActionErrorI18nKey(err, ollamaErrorConfig)
}

// GetSupportedOllamaFlags 返回所有支持的 ollama 命令标志
func GetSupportedOllamaFlags() []string {
	return GetSupportedFlags(OllamaActionFlagMap)
}

// GetOllamaActionForFlag 返回指定标志对应的动作
func GetOllamaActionForFlag(flag string) (OllamaAction, bool) {
	return GetActionForFlag(OllamaActionFlagMap, flag)
}

// GetOllamaFlagsForAction 返回指定动作对应的所有标志
func GetOllamaFlagsForAction(action OllamaAction) []string {
	return GetFlagsForAction(OllamaActionToFlagsMap, action)
}

// IsValidOllamaFlag 检查给定的字符串是否是有效的 ollama 命令标志
func IsValidOllamaFlag(flag string) bool {
	return IsValidFlag(OllamaActionFlagMap, flag)
}

// ParseOllamaArgs 解析 -o/--ollama 后的参数，确保且仅有一个动作
func ParseOllamaArgs(args []string) (OllamaAction, error) {
	return ParseActionArgs(args, OllamaActionFlagMap, OllamaUnknown, ollamaErrorConfig)
}

// OllamaActionHandlers 预定义的动作处理函数映射
var OllamaActionHandlers map[OllamaAction]func()

// SetOllamaActionHandlers 设置动作处理函数（初始化时调用一次）
func SetOllamaActionHandlers(handlers map[OllamaAction]func()) {
	OllamaActionHandlers = handlers
}

// HandleOllamaCommand 统一的 ollama 命令处理入口
func HandleOllamaCommand(args []string, mode, usageKey string) bool {
	return HandleActionCommand(
		args,
		usageKey,
		ParseOllamaArgs,
		GetOllamaErrorI18nKey,
		OllamaActionHandlers,
	)
}
