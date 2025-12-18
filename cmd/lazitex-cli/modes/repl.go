// cmd/lazitex-cli/modes/repl.go

package modes

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
	readline.PcItem("cd"),
	readline.PcItem("ls"),
	readline.PcItem("pwd"),
	readline.PcItem("clear"),
	readline.PcItem("cat"),
)

type LaTexCompleter struct{}

// Do 实现 readline.AutoCompleter 接口
func (c *LaTexCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	strLine := string(line[:pos])
	parts := strings.Fields(strLine)

	// 如果还没有输入完第一个单词，或者是在输入第一个单词
	if len(parts) == 0 || (len(parts) == 1 && !strings.HasSuffix(strLine, " ")) {
		return replCompleter.Do(line, pos)
	}

	// 如果第一个单词是 build, cd, ls 或 cat，则进行文件/目录补全
	cmd := strings.ToLower(parts[0])
	if cmd == "build" || cmd == "cd" || cmd == "ls" || cmd == "cat" {
		// 拿到用户正在输入的参数部分
		var prefix string
		if strings.HasSuffix(strLine, " ") {
			prefix = ""
		} else {
			prefix = parts[len(parts)-1]
		}

		// 扫描当前目录
		files, _ := os.ReadDir(".")
		var suggestions [][]rune
		for _, f := range files {
			name := f.Name()
			// 1. 对于 cd 命令，只补全文件夹
			if cmd == "cd" {
				if !f.IsDir() {
					continue
				}
			}
			// 2. 对于 build 命令，只补全 .tex 文件和文件夹
			if cmd == "build" {
				if !f.IsDir() && !strings.HasSuffix(strings.ToLower(name), ".tex") {
					continue
				}
			}
			// 3. 对于 cat 命令，补全 .tex, .log, .aux 文件和文件夹
			if cmd == "cat" {
				ext := strings.ToLower(filepath.Ext(name))
				if !f.IsDir() && ext != ".tex" && ext != ".log" && ext != ".aux" {
					continue
				}
			}

			// 简单的匹配过滤
			if strings.HasPrefix(name, prefix) {
				suffix := name[len(prefix):]
				if f.IsDir() {
					suffix += "/" // 文件夹加斜杠
				}
				suggestions = append(suggestions, []rune(suffix))
			}
		}
		return suggestions, len(prefix)
	}

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
		preview := false
		// 解析参数：build <file> [-p]
		for i := 1; i < len(parts); i++ {
			arg := parts[i]
			if arg == "-p" || arg == "--preview" {
				preview = true
			} else if filePath == "" {
				filePath = arg
			}
		}

		if filePath == "" {
			fmt.Println(core.T("repl.build_usage"))
			return false
		}

		// 执行构建
		BuildLaTeX(filePath, preview)

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
	fmt.Println("  help                - " + core.T("repl.help_desc"))
	fmt.Println("  version             - " + core.T("repl.version_desc"))
	fmt.Println("  check               - " + core.T("repl.check_desc"))
	fmt.Println("  install             - " + core.T("repl.install_desc"))
	fmt.Println("  uninstall           - " + core.T("repl.uninstall_desc"))
	fmt.Println("  build <file> [-p]   - " + core.T("repl.build_desc"))
	fmt.Println("  lang <zh|en>        - " + core.T("repl.lang_desc"))
	fmt.Println("  quit/exit           - " + core.T("repl.quit_desc"))
	fmt.Println()
	fmt.Println("  // " + core.T("repl.shell_title"))
	fmt.Println("  cd [path]           - " + core.T("repl.help_cd"))
	fmt.Println("  ls [path]           - " + core.T("repl.help_ls"))
	fmt.Println("  pwd                 - " + core.T("repl.help_pwd"))
	fmt.Println("  clear               - " + core.T("repl.help_clear"))
	fmt.Println("  cat [file]          - " + core.T("repl.help_cat"))
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
func BuildLaTeX(filePath string, preview bool) {
	opts := core.BuildOptions{
		InputPath: filePath,
		Preview:   preview,
	}

	pdfPath, err := core.Build(opts)
	if err != nil {
		fmt.Printf(core.T("msg.build_failed")+": %v\n", err)
		os.Exit(1)
	}

	if preview && pdfPath != "" {
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
		fmt.Printf(core.T("msg.unsupported_os")+"\n", runtime.GOOS)
	case "linux":
		fmt.Printf(core.T("msg.unsupported_os")+"\n", runtime.GOOS)
	default:
		// 如果不支持，就打印个提示，不强求
		fmt.Printf(core.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	if err != nil {
		fmt.Printf("Warning: Failed to open preview: %v\n", err)
	}
}
