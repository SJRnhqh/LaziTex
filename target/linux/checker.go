// target/linux/checker.go

package linux

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/SJRnhqh/lazitex/core"
)

// Checker Linux 平台检查器
type Checker struct{}

// NewChecker 创建 Linux 检查器
func NewChecker() *Checker {
	return &Checker{}
}

// GetPlatformName 获取平台名称
func (c *Checker) GetPlatformName() string {
	// 尝试识别具体的 Linux 发行版
	distro := detectLinuxDistro()
	if distro != "" {
		return "Linux (" + distro + ")"
	}
	return "Linux"
}

// detectLinuxDistro 检测 Linux 发行版
func detectLinuxDistro() string {
	// 读取 /etc/os-release
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				name := strings.TrimPrefix(line, "PRETTY_NAME=")
				name = strings.Trim(name, "\"")
				return name
			}
		}
	}

	// 尝试其他方法
	if data, err := os.ReadFile("/etc/lsb-release"); err == nil {
		if strings.Contains(string(data), "Ubuntu") {
			return "Ubuntu"
		}
	}

	return ""
}

// GetSearchPaths 获取 Linux 特定的搜索路径
func (c *Checker) GetSearchPaths() []string {
	paths := []string{}

	// 标准系统路径
	standardPaths := []string{
		"/usr/bin",
		"/usr/local/bin",
		"/bin",
	}

	for _, path := range standardPaths {
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
		}
	}

	// TeX Live 特定路径
	texlivePaths := []string{
		"/usr/local/texlive/2025/bin/x86_64-linux",
		"/usr/local/texlive/2024/bin/x86_64-linux",
		"/usr/local/texlive/2023/bin/x86_64-linux",
		"/usr/local/texlive/2022/bin/x86_64-linux",
		"/usr/local/texlive/2025/bin/aarch64-linux",
		"/usr/local/texlive/2024/bin/aarch64-linux",
	}

	for _, path := range texlivePaths {
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
		}
	}

	// 用户本地安装路径
	if homeDir, err := os.UserHomeDir(); err == nil {
		userPaths := []string{
			filepath.Join(homeDir, ".local", "bin"),
			filepath.Join(homeDir, "texlive", "2025", "bin", "x86_64-linux"),
			filepath.Join(homeDir, "texlive", "2024", "bin", "x86_64-linux"),
		}

		for _, path := range userPaths {
			if _, err := os.Stat(path); err == nil {
				paths = append(paths, path)
			}
		}
	}

	// Flatpak 路径（某些发行版）
	if _, err := os.Stat("/var/lib/flatpak/exports/bin"); err == nil {
		paths = append(paths, "/var/lib/flatpak/exports/bin")
	}

	return paths
}

// DetectDistribution 检测 LaTeX 发行版
func (c *Checker) DetectDistribution(tools map[string]core.CompilerInfo) string {
	// 优先检查 pdflatex 的版本信息
	if pdflatex, exists := tools["pdflatex"]; exists && pdflatex.Installed {
		version := strings.ToLower(pdflatex.Version)

		// 检测 TeX Live
		if strings.Contains(version, "tex live") {
			// 提取年份
			for _, year := range []string{"2025", "2024", "2023", "2022", "2021", "2020"} {
				if strings.Contains(version, year) {
					return "TeX Live " + year
				}
			}
			return "TeX Live"
		}

		// 通过路径判断
		path := strings.ToLower(pdflatex.Path)
		if strings.Contains(path, "texlive") {
			for _, year := range []string{"2025", "2024", "2023", "2022", "2021"} {
				if strings.Contains(path, year) {
					return "TeX Live " + year
				}
			}
			return "TeX Live"
		}

		// Debian/Ubuntu 包管理安装的版本
		if strings.HasPrefix(path, "/usr/bin/") {
			return "TeX Live (发行版软件包)"
		}
	}

	// 通过 tlmgr 判断
	if tlmgr, exists := tools["tlmgr"]; exists && tlmgr.Installed {
		return "TeX Live"
	}

	return "未知"
}

// GetInstallGuide 获取 Linux 安装指南
func (c *Checker) GetInstallGuide() string {
	guide := "📖 Linux 安装指南:\n\n"

	// 根据发行版给出具体建议
	distro := detectLinuxDistro()

	if strings.Contains(distro, "Ubuntu") || strings.Contains(distro, "Debian") {
		guide += "  Ubuntu/Debian 用户:\n"
		guide += "  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
		guide += "  快速安装（基础版）:\n"
		guide += "    sudo apt-get update\n"
		guide += "    sudo apt-get install texlive\n"
		guide += "\n"
		guide += "  完整安装（推荐）:\n"
		guide += "    sudo apt-get install texlive-full\n"
		guide += "\n"
		guide += "  中文支持:\n"
		guide += "    sudo apt-get install texlive-xetex texlive-lang-chinese\n"
		guide += "\n"

	} else if strings.Contains(distro, "Fedora") || strings.Contains(distro, "Red Hat") || strings.Contains(distro, "CentOS") {
		guide += "  Fedora/RHEL/CentOS 用户:\n"
		guide += "  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
		guide += "  完整安装:\n"
		guide += "    sudo dnf install texlive-scheme-full\n"
		guide += "\n"
		guide += "  基础安装:\n"
		guide += "    sudo dnf install texlive-scheme-basic\n"
		guide += "\n"

	} else if strings.Contains(distro, "Arch") || strings.Contains(distro, "Manjaro") {
		guide += "  Arch Linux/Manjaro 用户:\n"
		guide += "  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
		guide += "  完整安装:\n"
		guide += "    sudo pacman -S texlive-most texlive-lang\n"
		guide += "\n"
		guide += "  或者只安装核心:\n"
		guide += "    sudo pacman -S texlive-core\n"
		guide += "\n"

	} else {
		guide += "  通用方法（所有发行版）:\n"
		guide += "  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
		guide += "  使用包管理器安装:\n"
		guide += "    - Ubuntu/Debian: sudo apt-get install texlive-full\n"
		guide += "    - Fedora: sudo dnf install texlive-scheme-full\n"
		guide += "    - Arch: sudo pacman -S texlive-most\n"
		guide += "    - openSUSE: sudo zypper install texlive-scheme-full\n"
		guide += "\n"
	}

	guide += "  或者安装 TeX Live 官方版本:\n"
	guide += "  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	guide += "  1. 下载安装脚本:\n"
	guide += "     wget https://mirror.ctan.org/systems/texlive/tlnet/install-tl-unx.tar.gz\n"
	guide += "  2. 解压并运行:\n"
	guide += "     tar -xzf install-tl-unx.tar.gz\n"
	guide += "     cd install-tl-*\n"
	guide += "     sudo perl install-tl\n"
	guide += "  3. 按照提示完成安装\n"
	guide += "\n"
	guide += "  💡 提示: 包管理器安装更简单，官方版本功能更完整\n"

	return guide
}

// PostCheck Linux 特定的后处理检查
func (c *Checker) PostCheck(env *core.LaTeXEnvironment) {
	// 检查是否通过包管理器安装
	isPkgInstalled := false
	if pdflatex, exists := env.Tools["pdflatex"]; exists && pdflatex.Installed {
		if strings.HasPrefix(pdflatex.Path, "/usr/bin/") {
			isPkgInstalled = true
		}
	}

	// 如果是包管理器安装，检查常见的缺失问题
	if isPkgInstalled {
		// 检查是否安装了 texlive-full 或只是 texlive-base
		cmd := exec.Command("dpkg", "-l", "texlive-full")
		if err := cmd.Run(); err != nil {
			// texlive-full 未安装，可能是最小安装
			// 这里可以记录警告
		}
	}

	// 检查字体配置
	if xelatex, exists := env.Tools["xelatex"]; exists && xelatex.Installed {
		// XeLaTeX 依赖系统字体，可以检查 fontconfig
		if _, err := exec.LookPath("fc-list"); err != nil {
			// fontconfig 未安装
		}
	}
}
