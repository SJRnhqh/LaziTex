// target/win/checker.go
// Windows 平台特定 LaTeX 检查器实现

package win

import (
	// 外部包
	"os"
	"path/filepath"
	"strings"

	// 内部包
	core "github.com/SJRnhqh/lazitex/core"
)

// Checker Windows 平台检查器
type Checker struct{}

// NewChecker 创建 Windows 检查器
func NewChecker() *Checker {
	return &Checker{}
}

// GetPlatformName 获取平台名称
func (c *Checker) GetPlatformName() string {
	return "Windows"
}

// GetSearchPaths 获取 Windows 特定的搜索路径
func (c *Checker) GetSearchPaths() []string {
	paths := []string{}

	// TeX Live 常见安装路径
	texLivePaths := []string{
		"C:\\texlive\\2025\\bin\\windows",
		"C:\\texlive\\2024\\bin\\windows",
		"C:\\texlive\\2023\\bin\\windows",
		"C:\\texlive\\2022\\bin\\windows",
	}

	for _, path := range texLivePaths {
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
		}
	}

	// MiKTeX 常见安装路径
	miktexPaths := []string{
		"C:\\Program Files\\MiKTeX\\miktex\\bin\\x64",
		"C:\\Program Files (x86)\\MiKTeX\\miktex\\bin",
		"C:\\MiKTeX\\miktex\\bin\\x64",
	}

	for _, path := range miktexPaths {
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
		}
	}

	// 检查用户目录下的安装
	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		userPaths := []string{
			filepath.Join(userProfile, "AppData", "Local", "Programs", "MiKTeX", "miktex", "bin", "x64"),
			filepath.Join(userProfile, "texlive", "2025", "bin", "windows"),
			filepath.Join(userProfile, "texlive", "2024", "bin", "windows"),
		}

		for _, path := range userPaths {
			if _, err := os.Stat(path); err == nil {
				paths = append(paths, path)
			}
		}
	}

	return paths
}

// DetectDistribution 检测 LaTeX 发行版
func (c *Checker) DetectDistribution(tools map[string]core.CompilerInfo) string {
	if pdflatex, exists := tools["pdflatex"]; exists && pdflatex.Installed {
		version := strings.ToLower(pdflatex.Version)

		if strings.Contains(version, "miktex") {
			return "MiKTeX"
		}

		if strings.Contains(version, "tex live") {
			for _, year := range []string{"2025", "2024", "2023", "2022", "2021"} {
				if strings.Contains(version, year) {
					return "TeX Live " + year
				}
			}
			return "TeX Live"
		}

		path := strings.ToLower(pdflatex.Path)
		if strings.Contains(path, "miktex") {
			return "MiKTeX"
		}
		if strings.Contains(path, "texlive") {
			for _, year := range []string{"2025", "2024", "2023", "2022", "2021"} {
				if strings.Contains(path, year) {
					return "TeX Live " + year
				}
			}
			return "TeX Live"
		}
	}

	if tlmgr, exists := tools["tlmgr"]; exists && tlmgr.Installed {
		return "TeX Live"
	}

	if mpm, exists := tools["mpm"]; exists && mpm.Installed {
		return "MiKTeX"
	}

	return "未知"
}

// GetInstallGuide 获取 Windows 安装指南
func (c *Checker) GetInstallGuide() string {
	guide := "📖 Windows 安装指南:\n\n"
	guide += "  推荐选项 1: TeX Live (完整功能)\n"
	guide += "  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	guide += "  1. 访问: https://www.tug.org/texlive/acquire-netinstall.html\n"
	guide += "  2. 下载: install-tl-windows.exe\n"
	guide += "  3. 运行安装程序，选择完整安装（约 7GB）\n"
	guide += "  4. 安装后重启命令行\n"
	guide += "\n"
	guide += "  推荐选项 2: MiKTeX (轻量按需)\n"
	guide += "  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	guide += "  1. 访问: https://miktex.org/download\n"
	guide += "  2. 下载 Windows 安装包\n"
	guide += "  3. 安装时选择\"自动安装缺失的包\"\n"
	guide += "  4. 安装后重启命令行\n"
	guide += "\n"
	guide += "  💡 提示: TeX Live 适合完整功能，MiKTeX 适合快速开始\n"

	return guide
}

// PostCheck Windows 特定的后处理检查
func (c *Checker) PostCheck(env *core.LaTeXEnvironment) {
	// Windows 特定检查可以在这里添加
}
