// cmd/lazitex-cli/internal/latex.go
// LaTeX 管理相关业务命令解析

package internal

// LaTeXAction 表示 latex 命令支持的动作类型
type LaTeXAction int

const (
	LaTeXUnknown LaTeXAction = iota
	LaTeXCheck     // 检查 LaTeX 安装状态
	LaTeXInstall   // 安装 LaTeX
	LaTeXUninstall // 卸载 LaTeX
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