// cmd/lazitex-cli/main.go
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	modes "github.com/SJRnhqh/lazitex/cmd/lazitex-cli/modes"
	core "github.com/SJRnhqh/lazitex/core"
	linux "github.com/SJRnhqh/lazitex/target/linux"
	mac "github.com/SJRnhqh/lazitex/target/mac"
	win "github.com/SJRnhqh/lazitex/target/win"
	lipgloss "github.com/charmbracelet/lipgloss"
	figure "github.com/common-nighthawk/go-figure"
)

func main() {
	// 初始化语言设置并移除 --lang 参数
	args := processLanguageFlag()

	// 如果没有参数，显示帮助
	if len(args) == 0 {
		showHelp()
		return
	}

	// 简单处理第一个参数
	switch args[0] {
	case "-h", "--help":
		showHelp()
	case "-v", "--version":
		fmt.Println("LaziTex v0.0.1")
	case "-c", "--check":
		checkEnvironment()
	case "-t", "--tui":
		modes.StartTUI()
	case "-r", "--repl":
		printLogo()
		modes.StartREPL()
	case "-i", "--install":
		installLaTeXEnvironment()
	case "-u", "--uninstall":
		uninstallLaTeXEnvironment()
	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		showHelp()
	}
}

// processLanguageFlag 处理语言参数并返回剩余参数
func processLanguageFlag() []string {
	args := os.Args[1:] // 排除程序名
	newArgs := []string{}
	langSet := false

	for i := 0; i < len(args); i++ {
		// 支持 --lang 和 -l 两种形式
		if (args[i] == "--lang" || args[i] == "-l") && i+1 < len(args) {
			// 设置语言
			lang := args[i+1]
			switch lang {
			case "en", "english":
				core.SetLanguage(core.LangEN)
				langSet = true
				// 保存用户偏好
				core.SaveLanguagePreference(core.LangEN)
			case "zh", "chinese":
				core.SetLanguage(core.LangZH)
				langSet = true
				// 保存用户偏好
				core.SaveLanguagePreference(core.LangZH)
			}
			i++ // 跳过语言值
		} else {
			newArgs = append(newArgs, args[i])
		}
	}

	// 如果没有通过 --lang 设置，使用环境检测（会读取配置文件）
	if !langSet {
		core.SetLanguage(core.DetectSystemLanguage())
	}

	return newArgs
}

// 显示命令帮助
func showHelp() {
	fmt.Println(core.T("help.description"))
	fmt.Println()
	fmt.Println(core.T("help.usage"))
	fmt.Println("  lazitex [command]")
	fmt.Println()
	fmt.Println(core.T("help.commands"))
	fmt.Println("  -h, --help       " + core.T("help.show_help"))
	fmt.Println("  -v, --version    " + core.T("help.show_version"))
	fmt.Println("  -c, --check      " + core.T("help.check_env"))
	fmt.Println("  -r, --repl       " + core.T("help.start_repl"))
	fmt.Println("  -t, --tui        " + core.T("help.start_tui"))
	fmt.Println("  -i, --install    " + core.T("help.install_latex"))
	fmt.Println("  -u, --uninstall  " + core.T("help.uninstall_latex"))
	fmt.Println("  -l, --lang <lang> " + core.T("help.set_language"))
}

// 打印Logo
func printLogo() {
	// 定义科技感颜色样式
	orangeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF8C00")).
		Background(lipgloss.Color("#1a1a1a")).
		Bold(true) // 橙色，深色背景，粗体

	blueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF")).
		Background(lipgloss.Color("#1a1a1a")).
		Bold(true) // 亮蓝色，深色背景，粗体

	greenStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF7F")).
		Background(lipgloss.Color("#1a1a1a")).
		Bold(true) // 亮绿色，深色背景，粗体

	// 生成每个字母的ASCII艺术字
	letters := []string{"L", "a", "z", "i", "t", "e", "x"}
	styles := []lipgloss.Style{orangeStyle, blueStyle, orangeStyle, blueStyle, greenStyle, greenStyle, greenStyle}

	// 存储每个字母的艺术字行
	var letterLines [][]string
	maxHeight := 0
	maxWidth := 0

	// 生成每个字母的艺术字
	for _, letter := range letters {
		fig := figure.NewFigure(letter, "", true)
		lines := strings.Split(fig.String(), "\n")
		// 移除空行并计算最大宽度
		var cleanLines []string
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				cleanLines = append(cleanLines, line)
				if len(line) > maxWidth {
					maxWidth = len(line)
				}
			}
		}
		letterLines = append(letterLines, cleanLines)
		if len(cleanLines) > maxHeight {
			maxHeight = len(cleanLines)
		}
	}

	// 打印带白色横条背景的Logo
	// fmt.Println()

	// 计算总宽度
	// totalWidth := (maxWidth+1)*len(letters) - 1

	// 白色横条样式
	// whiteBarStyle := lipgloss.NewStyle().
	// 	Background(lipgloss.Color("#FFFFFF")).
	// 	Foreground(lipgloss.Color("#000000"))

	// 顶部白色横条
	// topBar := strings.Repeat(" ", totalWidth)
	// fmt.Println(whiteBarStyle.Render(topBar))

	// 打印彩色Logo
	for i := 0; i < maxHeight; i++ {
		var line string
		for j, letterLine := range letterLines {
			if i < len(letterLine) {
				// 右对齐填充到固定宽度
				paddedLine := fmt.Sprintf("%-*s", maxWidth, letterLine[i])
				line += styles[j].Render(paddedLine)
			} else {
				// 如果当前字母的行数不够，用空格填充
				line += strings.Repeat(" ", maxWidth)
			}
			// 字母之间添加一个空格
			if j < len(letterLines)-1 {
				line += " "
			}
		}
		fmt.Println(line)
	}

	// 底部白色横条
	// bottomBar := strings.Repeat(" ", totalWidth)
	// fmt.Println(whiteBarStyle.Render(bottomBar))
	// fmt.Println()
}

// 检查 LaTeX 环境
func checkEnvironment() {
	// 根据平台创建对应的检查器
	var checker core.EnvironmentChecker

	switch runtime.GOOS {
	case "windows":
		checker = win.NewChecker()
	case "linux":
		checker = linux.NewChecker()
	case "darwin":
		checker = mac.NewChecker()
	default:
		// 不支持的平台，使用一个简单的错误提示
		fmt.Printf(core.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	// 执行检测
	env := core.CheckLaTeXEnvironment(checker)
	env.PrintEnvironment()
}

// 安装 LaTeX 环境
func installLaTeXEnvironment() {
	modes.InstallLaTeXEnvironment()
}

// 卸载 LaTeX 环境
func uninstallLaTeXEnvironment() {
	modes.UninstallLaTeXEnvironment()
}
