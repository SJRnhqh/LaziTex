// target/mac/checker.go

package mac

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/SJRnhqh/lazitex/core"
)

// Checker macOS 平台检查器
type Checker struct{}

// NewChecker 创建 macOS 检查器
func NewChecker() *Checker {
	return &Checker{}
}

// GetPlatformName 获取平台名称
func (c *Checker) GetPlatformName() string {
	// 尝试获取 macOS 版本
	cmd := exec.Command("sw_vers", "-productVersion")
	if output, err := cmd.Output(); err == nil {
		version := strings.TrimSpace(string(output))
		return "macOS " + version
	}
	return "macOS"
}

// GetSearchPaths 获取 macOS 特定的搜索路径
func (c *Checker) GetSearchPaths() []string {
	paths := []string{}

	// MacTeX/TeX Live 标准安装路径
	texlivePaths := []string{
		"/Library/TeX/texbin",
		"/usr/local/texlive/2025/bin/universal-darwin",
		"/usr/local/texlive/2024/bin/universal-darwin",
		"/usr/local/texlive/2023/bin/universal-darwin",
		"/usr/local/texlive/2022/bin/universal-darwin",
		"/usr/local/texlive/2025/bin/x86_64-darwin",
		"/usr/local/texlive/2024/bin/x86_64-darwin",
		"/usr/local/texlive/2023/bin/x86_64-darwin",
		"/usr/local/texlive/2025/bin/aarch64-darwin",
		"/usr/local/texlive/2024/bin/aarch64-darwin",
	}

	for _, path := range texlivePaths {
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
		}
	}

	// Homebrew 安装路径
	homebrewPaths := []string{
		"/opt/homebrew/bin",             // Apple Silicon
		"/usr/local/bin",                // Intel Mac
		"/opt/homebrew/opt/texlive/bin", // Homebrew TeX Live
	}

	for _, path := range homebrewPaths {
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
		}
	}

	// MacPorts 安装路径
	macportsPaths := []string{
		"/opt/local/bin",
	}

	for _, path := range macportsPaths {
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
		}
	}

	// 用户本地安装
	if homeDir, err := os.UserHomeDir(); err == nil {
		userPaths := []string{
			filepath.Join(homeDir, ".local", "bin"),
			filepath.Join(homeDir, "Library", "TeX", "texbin"),
			filepath.Join(homeDir, "texlive", "2025", "bin", "universal-darwin"),
			filepath.Join(homeDir, "texlive", "2024", "bin", "universal-darwin"),
		}

		for _, path := range userPaths {
			if _, err := os.Stat(path); err == nil {
				paths = append(paths, path)
			}
		}
	}

	return paths
}

// DetectDistribution 检测 LaTex 发行版
func (c *Checker) DetectDistribution(tools map[string]core.CompilerInfo) string {
	// 优先检查 pdflatex 的版本信息
	if pdflatex, exists := tools["pdflatex"]; exists && pdflatex.Installed {
		version := strings.ToLower(pdflatex.Version)
		path := strings.ToLower(pdflatex.Path)

		// 检测 MacTeX
		if strings.Contains(path, "/library/tex/") {
			// 这是 MacTeX 的标准安装路径
			for _, year := range []string{"2025", "2024", "2023", "2022", "2021"} {
				if strings.Contains(version, year) {
					return "MacTeX " + year
				}
			}
			return "MacTeX (TeX Live)"
		}

		// 检测 TeX Live
		if strings.Contains(version, "tex live") {
			for _, year := range []string{"2025", "2024", "2023", "2022", "2021"} {
				if strings.Contains(version, year) {
					return "TeX Live " + year
				}
			}
			return "TeX Live"
		}

		// 通过路径判断
		if strings.Contains(path, "texlive") {
			for _, year := range []string{"2025", "2024", "2023", "2022", "2021"} {
				if strings.Contains(path, year) {
					return "TeX Live " + year
				}
			}
			return "TeX Live"
		}

		// Homebrew 安装
		if strings.Contains(path, "homebrew") || strings.Contains(path, "/opt/homebrew/") {
			return "TeX Live (Homebrew)"
		}

		// MacPorts 安装
		if strings.Contains(path, "macports") || strings.Contains(path, "/opt/local/") {
			return "TeX Live (MacPorts)"
		}
	}

	// 通过 tlmgr 判断
	if tlmgr, exists := tools["tlmgr"]; exists && tlmgr.Installed {
		path := strings.ToLower(tlmgr.Path)
		if strings.Contains(path, "/library/tex/") {
			return "MacTeX (TeX Live)"
		}
		return "TeX Live"
	}

	return "未知"
}

// GetInstallGuide 获取 macOS 安装指南
func (c *Checker) GetInstallGuide() string {
	guide := "📖 macOS 安装指南:\n\n"

	guide += "  推荐选项 1: MacTeX (官方完整版)\n"
	guide += "  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	guide += "  这是 macOS 上最推荐的方式，包含完整的 TeX Live 和 GUI 工具。\n"
	guide += "\n"
	guide += "  安装方法:\n"
	guide += "    1. 访问: https://www.tug.org/mactex/\n"
	guide += "    2. 下载 MacTeX.pkg (约 4.5 GB)\n"
	guide += "    3. 双击安装包，按照提示安装\n"
	guide += "    4. 安装后重启终端\n"
	guide += "\n"
	guide += "  包含工具:\n"
	guide += "    - 完整的 TeX Live 发行版\n"
	guide += "    - TeXShop 编辑器\n"
	guide += "    - BibDesk 文献管理\n"
	guide += "    - LaTexiT 公式编辑器\n"
	guide += "\n"

	guide += "  推荐选项 2: Homebrew (轻量快捷)\n"
	guide += "  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	guide += "  适合已经使用 Homebrew 的用户，安装更快。\n"
	guide += "\n"
	guide += "  安装方法:\n"
	guide += "    # 完整安装 MacTeX\n"
	guide += "    brew install --cask mactex\n"
	guide += "\n"
	guide += "    # 或者只安装基础版（不含 GUI 工具，约 1.5 GB）\n"
	guide += "    brew install --cask basictex\n"
	guide += "\n"
	guide += "  安装后更新:\n"
	guide += "    sudo tlmgr update --self\n"
	guide += "    sudo tlmgr update --all\n"
	guide += "\n"

	guide += "  选项 3: MacPorts\n"
	guide += "  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	guide += "  如果你使用 MacPorts:\n"
	guide += "    sudo port install texlive\n"
	guide += "\n"

	// 检测当前是否安装了 Homebrew
	if _, err := exec.LookPath("brew"); err == nil {
		guide += "  💡 提示: 检测到你已安装 Homebrew，推荐使用选项 2\n"
	} else {
		guide += "  💡 提示: 推荐直接安装 MacTeX (选项 1) 获得最佳体验\n"
	}

	guide += "\n"
	guide += "  安装后验证:\n"
	guide += "    lazitex --check\n"

	return guide
}

// PostCheck macOS 特定的后处理检查
func (c *Checker) PostCheck(env *core.LaTexEnvironment) {
	// 检查是否是 MacTeX 完整安装
	isMacTeX := false
	if _, err := os.Stat("/Library/TeX/texbin"); err == nil {
		isMacTeX = true
	}

	// 检查 GUI 工具（MacTeX 特有）
	if isMacTeX {
		guiApps := map[string]string{
			"TeXShop": "/Applications/TeX/TeXShop.app",
			"BibDesk": "/Applications/TeX/BibDesk.app",
			"LaTexiT": "/Applications/TeX/LaTexiT.app",
		}

		// 这里可以检查 GUI 应用是否存在
		// 但由于我们主要关注命令行工具，暂不在检测结果中显示
		_ = guiApps
	}

	// 检查 XeLaTex 的字体访问
	if xelatex, exists := env.Tools["xelatex"]; exists && xelatex.Installed {
		// macOS 上 XeLaTex 应该能访问系统字体
		// 可以通过 fc-list 检查，但 macOS 默认使用 fontconfig
	}

	// 检查是否通过 Homebrew 安装
	if pdflatex, exists := env.Tools["pdflatex"]; exists && pdflatex.Installed {
		if strings.Contains(pdflatex.Path, "homebrew") {
			// Homebrew 安装的版本
			// 可能需要额外配置
		}
	}

	// 检查 Ghostscript (ps2pdf 依赖)
	if ps2pdf, exists := env.Tools["ps2pdf"]; exists && ps2pdf.Installed {
		if _, err := exec.LookPath("gs"); err != nil {
			// Ghostscript 未安装，ps2pdf 可能无法正常工作
		}
	}
}
