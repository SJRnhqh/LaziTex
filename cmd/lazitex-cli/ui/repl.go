// cmd/lazitex-cli/ui/repl.go
// REPL模式界面管理

package ui

import (
	// 外部包
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	// 内部包
	tasks "github.com/SJRnhqh/lazitex/cmd/lazitex-cli/tasks"
	lang "github.com/SJRnhqh/lazitex/lang"
	lipgloss "github.com/charmbracelet/lipgloss"
	readline "github.com/chzyer/readline"
	figure "github.com/common-nighthawk/go-figure"
)

// 创建补全器
var replCompleter = readline.NewPrefixCompleter(
	// 1. 核心命令补全
	readline.PcItem("help"),
	readline.PcItem("version"),
	readline.PcItem("check"),
	readline.PcItem("install"),
	readline.PcItem("uninstall"),
	readline.PcItem("exit"),
	readline.PcItem("quit"),

	// 2. 语言补全
	readline.PcItem("lang",
		readline.PcItem("zh"),
		readline.PcItem("en"),
	),

	// 3. 带路径补全的命令
	readline.PcItem("build"), // 后面我们会动态处理
	readline.PcItem("preview"),
	readline.PcItem("cd"),
	readline.PcItem("ls"),
	readline.PcItem("pwd"),
	readline.PcItem("clear"),
	readline.PcItem("cat"),
)

type LaTeXCompleter struct{}

// Do 实现 readline.AutoCompleter 接口，支持跨目录的文件补全
func (c *LaTeXCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	strLine := string(line[:pos])
	parts := strings.Fields(strLine)

	// 如果还没有输入完第一个单词，或者是在输入第一个单词，使用基础命令补全
	if len(parts) == 0 || (len(parts) == 1 && !strings.HasSuffix(strLine, " ")) {
		return replCompleter.Do(line, pos)
	}

	// 解析当前命令
	cmd := strings.ToLower(parts[0])

	// 仅对需要路径参数的命令进行增强补全
	if cmd == "build" || cmd == "cd" || cmd == "ls" || cmd == "cat" || cmd == "preview" {
		var inputPath string
		if strings.HasSuffix(strLine, " ") {
			inputPath = ""
		} else {
			// 获取最后一个参数作为输入路径
			inputPath = parts[len(parts)-1]
		}

		// --- 核心逻辑：处理跨目录路径 ---
		dir := "."
		filePrefix := inputPath

		// 如果包含斜杠，分离出目录部分和正在输入的文件名前缀
		if lastSlash := strings.LastIndex(inputPath, "/"); lastSlash != -1 {
			dir = inputPath[:lastSlash+1]
			filePrefix = inputPath[lastSlash+1:]
		}

		// 扫描目标目录
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, 0
		}

		var suggestions [][]rune
		for _, entry := range entries {
			name := entry.Name()

			// 排除隐藏文件
			if strings.HasPrefix(name, ".") {
				continue
			}

			// --- 针对不同命令的智能过滤 ---
			isDir := entry.IsDir()

			// 1. cd 命令：只补全目录
			if cmd == "cd" && !isDir {
				continue
			}

			// 2. build 命令：只补全目录和 .tex 文件
			if cmd == "build" || cmd == "preview" {
				if !isDir && !strings.HasSuffix(strings.ToLower(name), ".tex") {
					continue
				}
			}

			// 3. cat 命令：只补全目录和相关文本文件
			if cmd == "cat" {
				ext := strings.ToLower(filepath.Ext(name))
				if !isDir && ext != ".tex" && ext != ".log" && ext != ".aux" && ext != ".txt" {
					continue
				}
			}

			// 简单的匹配过滤
			if strings.HasPrefix(name, filePrefix) {
				// 补全部分 = 完整名字 - 已经输入的前缀
				suffix := name[len(filePrefix):]
				if isDir {
					suffix += "/" // 目录自动加斜杠，方便继续 Tab 进入
				}
				suggestions = append(suggestions, []rune(suffix))
			}
		}

		// 返回建议列表和当前正在匹配的前缀长度
		return suggestions, len(filePrefix)
	}

	// 其他情况回退到默认补全器
	return replCompleter.Do(line, pos)
}

// StartREPL 启动 REPL 模式
func StartREPL() {
	printLogo()
	printREPLWelcome()

	// 1. 配置历史记录文件的路径
	homeDir, _ := os.UserHomeDir()
	historyFile := filepath.Join(homeDir, ".lazitex_history")

	// 2. 初始化 readline 实例
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          lang.T("repl.prompt"),
		HistoryFile:     historyFile,
		AutoComplete:    &LaTeXCompleter{},
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Printf(lang.T("repl.err_init")+"\n", err)
		return
	}
	defer rl.Close()

	for {
		// 3. 使用 readline 读取输入
		line, err := rl.Readline()
		if err != nil { // 处理 Ctrl+C 或 Ctrl+D
			break
		}

		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}

		// 解析并执行命令
		if shouldExit := handleREPLCommand(input); shouldExit {
			break
		}
	}
	fmt.Println(lang.T("repl.goodbye"))
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

// StartREPL 启动 REPL 模式
func handleREPLCommand(input string) bool {
	// 解析命令和参数
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false
	}
	command := strings.ToLower(parts[0])

	// 先尝试处理 REPL 内置的 shell 类命令
	if handled := handleShellLikeCommands(parts); handled {
		return false
	}

	switch command {
	case "quit", "exit":
		return true // 返回 true 表示退出 REPL

	case "help":
		tasks.ShowREPLHelp()

	case "version":
		fmt.Println("LaziTex v0.0.1")

	case "check":
		tasks.CheckEnvironment()

	case "install":
		tasks.InstallLaTeXEnvironment()

	case "uninstall":
		tasks.UninstallLaTeXEnvironment()

	case "ollama":
		if len(parts) < 2 {
			fmt.Println(lang.T("repl.ollama_usage"))
			return false
		}
		action := parts[1]
		switch action {
		case "check", "-c", "--check":
			tasks.CheckOllama()
		case "install", "-i", "--install":
			tasks.InstallOllama()
		case "uninstall", "-u", "--uninstall":
			tasks.UninstallOllama()
		default:
			fmt.Println(lang.T("repl.ollama_usage"))
		}

	case "build":
		if len(parts) < 2 {
			fmt.Println(lang.T("repl.build_usage"))
			return false
		}

		filePath := ""
		outputPath := ""
		show := false
		pendingOutput := false
		quiet := false
		tidy := false
		// 解析参数：build <file> [-s] [-o output_path]
		for i := 1; i < len(parts); i++ {
			arg := parts[i]
			if arg == "-s" || arg == "--show" {
				show = true
			} else if arg == "-o" || arg == "--output" {
				pendingOutput = true
			} else if arg == "-q" || arg == "--quiet" {
				quiet = true
			} else if arg == "-t" || arg == "--tidy" {
				tidy = true
			} else if strings.HasPrefix(arg, "-") {
				// 严谨处理：未知标志位报错
				fmt.Printf(lang.T("repl.unknown_command")+"\n", arg)
				return false
			} else {
				if pendingOutput && outputPath == "" {
					outputPath = arg
					pendingOutput = false
				} else if filePath == "" {
					filePath = arg
				}
			}
		}

		// 检查：如果开启了 -o 但没拿到路径，或者没提供输入文件
		if (pendingOutput && outputPath == "") || filePath == "" {
			fmt.Println(lang.T("repl.build_usage"))
			return false
		}

		// 调用统一的构建入口
		tasks.BuildLaTeX(filePath, outputPath, show, quiet, tidy)
	case "preview":
		if len(parts) < 2 {
			fmt.Println(lang.T("repl.preview_usage"))
			return false
		}

		var filePath string
		port := 8080
		quiet := false
		tidy := false
		devMode := false

		// 解析参数：preview <file> [-q] [-t]
		for i := 1; i < len(parts); i++ {
			arg := parts[i]

			// 处理标志位
			if arg == "-q" || arg == "--quiet" {
				quiet = true
				continue
			}
			if arg == "-t" || arg == "--tidy" {
				tidy = true
				continue
			}
			if arg == "-d" || arg == "--dev" {
				devMode = true
				continue
			}

			// 处理端口号参数（独立的:port格式）
			if strings.HasPrefix(arg, ":") {
				portStr := arg[1:]
				parsedPort, err := strconv.Atoi(portStr)
				if err != nil || parsedPort < 1 || parsedPort > 65535 {
					fmt.Println(lang.T("msg.invalid_port"))
					return false
				}
				port = parsedPort
				continue
			}

			// 处理文件路径（如果 filePath 还没设置）
			if filePath == "" {
				if strings.Contains(arg, ":") {
					lastColonIndex := strings.LastIndex(arg, ":")
					if lastColonIndex > 0 && lastColonIndex < len(arg)-1 {
						portStr := arg[lastColonIndex+1:]
						parsedPort, err := strconv.Atoi(portStr)
						if err == nil && parsedPort >= 1 && parsedPort <= 65535 {
							filePath = arg[:lastColonIndex]
							port = parsedPort
							continue
						}
					}
				}
				filePath = arg
			}
		}

		// 检查是否提供了文件路径
		if filePath == "" {
			fmt.Println(lang.T("repl.preview_usage"))
			return false
		}

		tasks.StartLivePreview(filePath, port, quiet, tidy, devMode)
	case "lang", "language":
		// 语言切换命令
		if len(parts) < 2 {
			fmt.Println(lang.T("repl.lang_usage"))
			fmt.Println(lang.T("repl.lang_current") + lang.GetCurrentLanguageName())
		} else {
			handleLanguageSwitch(parts[1])
		}

	default:
		fmt.Printf(lang.T("repl.unknown_command")+"\n", command)
		fmt.Println(lang.T("repl.type_help"))
	}

	return false
}

// printREPLWelcome 打印欢迎信息
func printREPLWelcome() {
	fmt.Println(lang.T("repl.welcome"))
	fmt.Println(lang.T("repl.help_hint"))
	fmt.Println()
}

// handleLanguageSwitch 处理语言切换
func handleLanguageSwitch(langStr string) {
	langStr = strings.ToLower(langStr)

	switch langStr {
	case "zh", "chinese", "中文":
		lang.SetLanguage(lang.LangZH)
		lang.SaveLanguagePreference(lang.LangZH)
		fmt.Println(lang.T("repl.lang_switched"))

	case "en", "english", "英文":
		lang.SetLanguage(lang.LangEN)
		lang.SaveLanguagePreference(lang.LangEN)
		fmt.Println(lang.T("repl.lang_switched"))

	default:
		fmt.Printf(lang.T("repl.lang_unsupported")+"\n", langStr)
		fmt.Println(lang.T("repl.lang_available"))
	}
}

// handleShellLikeCommands 处理 REPL 内置的简单 shell 类命令
func handleShellLikeCommands(parts []string) bool {
	if len(parts) == 0 {
		return false
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "cd":
		target := ""
		if len(parts) < 2 {
			if home, err := os.UserHomeDir(); err == nil {
				target = home
			} else {
				fmt.Printf(lang.T("repl.err_cd")+"\n", err)
				return true
			}
		} else {
			target = parts[1]
		}

		if err := os.Chdir(target); err != nil {
			fmt.Printf(lang.T("repl.err_cd")+"\n", err)
		}
		return true

	case "pwd":
		if cwd, err := os.Getwd(); err == nil {
			fmt.Println(cwd)
		} else {
			fmt.Printf(lang.T("repl.err_pwd")+"\n", err)
		}
		return true

	case "ls":
		dir := "."
		if len(parts) > 1 {
			dir = parts[1]
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Printf(lang.T("repl.err_ls")+"\n", err)
			return true
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() {
				name += "/"
			}
			fmt.Println(name)
		}
		return true

	case "clear":
		// ANSI 清屏：移动到左上并清空屏幕
		fmt.Print("\033[H\033[2J")
		return true

	case "cat":
		if len(parts) < 2 {
			fmt.Println(lang.T("repl.cat_usage"))
			return true
		}
		content, err := os.ReadFile(parts[1])
		if err != nil {
			fmt.Printf(lang.T("repl.err_cat")+"\n", err)
		} else {
			fmt.Print(string(content))
			if !strings.HasSuffix(string(content), "\n") {
				fmt.Println()
			}
		}
		return true
	}

	return false
}
