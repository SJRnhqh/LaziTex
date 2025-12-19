// cmd/lazitex-cli/modes/repl.go

package modes

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	core "github.com/SJRnhqh/lazitex/core"
	linux "github.com/SJRnhqh/lazitex/target/linux"
	mac "github.com/SJRnhqh/lazitex/target/mac"
	win "github.com/SJRnhqh/lazitex/target/win"
	readline "github.com/chzyer/readline"
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

type LaTexCompleter struct{}

// Do 实现 readline.AutoCompleter 接口，支持跨目录的文件补全
func (c *LaTexCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
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
	printREPLWelcome()

	// 1. 配置历史记录文件的路径
	homeDir, _ := os.UserHomeDir()
	historyFile := filepath.Join(homeDir, ".lazitex_history")

	// 2. 初始化 readline 实例
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          core.T("repl.prompt"),
		HistoryFile:     historyFile,
		AutoComplete:    &LaTexCompleter{},
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Printf(core.T("repl.err_init")+"\n", err)
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
	fmt.Println(core.T("repl.goodbye"))
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
		showREPLHelp()

	case "version":
		fmt.Println("LaziTex v0.0.1")

	case "check":
		checkREPLEnvironment()

	case "install":
		InstallLaTeXEnvironment()

	case "uninstall":
		UninstallLaTeXEnvironment()

	case "build":
		if len(parts) < 2 {
			fmt.Println(core.T("repl.build_usage"))
			return false
		}

		filePath := ""
		outputPath := ""
		show := false
		pendingOutput := false

		// 解析参数：build <file> [-s] [-o output_path]
		for i := 1; i < len(parts); i++ {
			arg := parts[i]
			if arg == "-s" || arg == "--show" {
				show = true
			} else if arg == "-o" || arg == "--output" {
				pendingOutput = true
			} else if strings.HasPrefix(arg, "-") {
				// 严谨处理：未知标志位报错
				fmt.Printf(core.T("repl.unknown_command")+"\n", arg)
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
			fmt.Println(core.T("repl.build_usage"))
			return false
		}

		// 调用统一的构建入口
		BuildLaTeX(filePath, outputPath, show)
	case "preview":
		if len(parts) < 2 {
			fmt.Println(core.T("repl.preview_usage"))
			return false
		}
		filePath := parts[1]
		StartLivePreview(filePath)
	case "lang", "language":
		// 语言切换命令
		if len(parts) < 2 {
			fmt.Println(core.T("repl.lang_usage"))
			fmt.Println(core.T("repl.lang_current") + getCurrentLanguageName())
		} else {
			handleLanguageSwitch(parts[1])
		}

	default:
		fmt.Printf(core.T("repl.unknown_command")+"\n", command)
		fmt.Println(core.T("repl.type_help"))
	}

	return false
}

// printREPLWelcome 打印欢迎信息
func printREPLWelcome() {
	fmt.Println(core.T("repl.welcome"))
	fmt.Println(core.T("repl.help_hint"))
	fmt.Println()
}

// getCurrentLanguageName 获取当前语言的友好名称
func getCurrentLanguageName() string {
	switch core.GetLanguage() {
	case core.LangZH:
		return "中文 (Chinese)"
	case core.LangEN:
		return "English"
	default:
		return "English"
	}
}

// handleLanguageSwitch 处理语言切换
func handleLanguageSwitch(lang string) {
	lang = strings.ToLower(lang)

	switch lang {
	case "zh", "chinese", "中文":
		core.SetLanguage(core.LangZH)
		core.SaveLanguagePreference(core.LangZH)
		fmt.Println(core.T("repl.lang_switched"))

	case "en", "english", "英文":
		core.SetLanguage(core.LangEN)
		core.SaveLanguagePreference(core.LangEN)
		fmt.Println(core.T("repl.lang_switched"))

	default:
		fmt.Printf(core.T("repl.lang_unsupported")+"\n", lang)
		fmt.Println(core.T("repl.lang_available"))
	}
}

// showREPLHelp 显示帮助信息
func showREPLHelp() {
	fmt.Println(core.T("repl.help_title"))
	fmt.Printf("  %-30s - %s\n", "help", core.T("repl.help_desc"))
	fmt.Printf("  %-30s - %s\n", "version", core.T("repl.version_desc"))
	fmt.Printf("  %-30s - %s\n", "check", core.T("repl.check_desc"))
	fmt.Printf("  %-30s - %s\n", "install", core.T("repl.install_desc"))
	fmt.Printf("  %-30s - %s\n", "uninstall", core.T("repl.uninstall_desc"))
	fmt.Printf("  %-30s - %s\n", "build <file> [-o path] [-s]", core.T("repl.build_desc"))
	fmt.Printf("  %-30s - %s\n", "preview <file>", core.T("repl.preview_desc"))
	fmt.Printf("  %-30s - %s\n", "lang <zh|en>", core.T("repl.lang_desc"))
	fmt.Printf("  %-30s - %s\n", "lang", core.T("repl.lang_current")+getCurrentLanguageName())
	fmt.Printf("  %-30s - %s\n", "quit / exit", core.T("repl.quit_desc"))
	fmt.Println()
	fmt.Println(core.T("repl.shell_title") + ":")
	fmt.Printf("  %-30s - %s\n", "cd [path]", core.T("repl.help_cd"))
	fmt.Printf("  %-30s - %s\n", "ls [path]", core.T("repl.help_ls"))
	fmt.Printf("  %-30s - %s\n", "pwd", core.T("repl.help_pwd"))
	fmt.Printf("  %-30s - %s\n", "clear", core.T("repl.help_clear"))
	fmt.Printf("  %-30s - %s\n", "cat <file>", core.T("repl.help_cat"))
}

// 检查 LaTeX 环境
func checkREPLEnvironment() {
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

// InstallLaTeXEnvironment 安装 LaTeX 环境（共享函数，供 REPL 和命令行模式使用）
func InstallLaTeXEnvironment() {
	// 根据平台创建对应的安装器
	var installer core.EnvironmentInstaller

	switch runtime.GOOS {
	case "windows":
		// TODO: 后续实现 Windows 安装器
		fmt.Printf(core.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	case "linux":
		installer = linux.NewInstaller()
	case "darwin":
		installer = mac.NewInstaller()
	default:
		fmt.Printf(core.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	// 执行安装
	if err := core.InstallLaTeXEnvironment(installer); err != nil {
		fmt.Printf(core.T("msg.install_failed")+": %v\n", err)
		return
	}
	// 注意：成功消息由安装器内部输出，这里不需要再输出
}

// UninstallLaTeXEnvironment 卸载 LaTeX 环境（共享函数，供 REPL 和命令行模式使用）
func UninstallLaTeXEnvironment() {
	// 根据平台创建对应的安装器
	var installer core.EnvironmentInstaller

	switch runtime.GOOS {
	case "windows":
		// TODO: 后续实现 Windows 安装器
		fmt.Printf(core.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	case "linux":
		installer = linux.NewInstaller()
	case "darwin":
		installer = mac.NewInstaller()
	default:
		fmt.Printf(core.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	// 执行卸载
	if err := core.UninstallLaTeXEnvironment(installer); err != nil {
		fmt.Printf(core.T("msg.uninstall_failed")+": %v\n", err)
		return
	}
	// 注意：成功消息由卸载器内部输出，这里不需要再输出
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
				fmt.Printf(core.T("repl.err_cd")+"\n", err)
				return true
			}
		} else {
			target = parts[1]
		}

		if err := os.Chdir(target); err != nil {
			fmt.Printf(core.T("repl.err_cd")+"\n", err)
		}
		return true

	case "pwd":
		if cwd, err := os.Getwd(); err == nil {
			fmt.Println(cwd)
		} else {
			fmt.Printf(core.T("repl.err_pwd")+"\n", err)
		}
		return true

	case "ls":
		dir := "."
		if len(parts) > 1 {
			dir = parts[1]
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Printf(core.T("repl.err_ls")+"\n", err)
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
			fmt.Println(core.T("repl.cat_usage"))
			return true
		}
		content, err := os.ReadFile(parts[1])
		if err != nil {
			fmt.Printf(core.T("repl.err_cat")+"\n", err)
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

// 构建 LaTeX 文档
func BuildLaTeX(filePath string, outputPath string, show bool) {
	opts := core.BuildOptions{
		InputPath:  filePath,
		OutputPath: outputPath,
		Show:       show,
	}

	pdfPath, err := core.Build(opts)
	if err != nil {
		fmt.Printf(core.T("msg.build_failed")+": %v\n", err)
		return
	}

	if show && pdfPath != "" {
		OpenPDF(pdfPath)
	}
}

// OpenPDF 根据平台分发预览任务
func OpenPDF(pdfPath string) {
	var err error
	switch runtime.GOOS {
	case "darwin":
		err = mac.PreviewPDF(pdfPath) // macOS 预览 PDF
	case "windows":
		err = win.PreviewPDF(pdfPath) // Windows 预览 PDF
	case "linux":
		fmt.Printf(core.T("msg.unsupported_os")+"\n", runtime.GOOS)
	default:
		// 如果不支持，就打印个提示，不强求
		fmt.Printf(core.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	if err != nil {
		fmt.Printf(core.T("msg.preview_failed")+"\n", err)
	}
}

// StartLivePreview 启动实时预览模式
// filePath: 要预览的 .tex 文件路径
func StartLivePreview(filePath string) {
	// 1. 获取绝对路径，确保监听准确
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Printf(core.T("msg.err_abs_path")+"\n", filePath)
		return
	}

	// 2. 启动后立即执行一次“初次构建并展示”
	// 这里调用我们已有的 BuildLaTeX 函数，设置展示标志为 true
	BuildLaTeX(absPath, "", true)

	// 3. 打印监听提示（文案已在 i18n 中定义）
	fmt.Printf(core.T("msg.watching_file")+"\n", filepath.Base(absPath))

	// 4. 调用 core 层的监听引擎
	// 当文件变动时，它会回调执行我们定义的闭包函数
	err = core.WatchAndAction(absPath, func() {
		// 这里是文件变动后的动作
		currentTime := time.Now().Format("15:04:05")
		fmt.Printf("\n🔄 [%s] %s\n", currentTime, core.T("msg.building_doc"))

		// 重新执行编译和展示逻辑
		BuildLaTeX(absPath, "", true)
	})

	// 5. 错误处理（如果监听器意外崩溃）
	if err != nil {
		fmt.Printf("Watcher error: %v\n", err)
	}
}
