// cmd/lazitex-cli/modes/repl.go

package modes

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	core "github.com/SJRnhqh/lazitex/core"
	linux "github.com/SJRnhqh/lazitex/target/linux"
	mac "github.com/SJRnhqh/lazitex/target/mac"
	win "github.com/SJRnhqh/lazitex/target/win"
)

// StartREPL 启动 REPL 模式
func StartREPL() {
	printREPLWelcome()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// 显示提示符
		fmt.Print(core.T("repl.prompt"))

		// 读取用户输入
		if !scanner.Scan() {
			// 扫描结束
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// 处理空输入
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
	fmt.Println("  help            - " + core.T("repl.help_desc"))
	fmt.Println("  version         - " + core.T("repl.version_desc"))
	fmt.Println("  check           - " + core.T("repl.check_desc"))
	fmt.Println("  install         - " + core.T("repl.install_desc"))
	fmt.Println("  uninstall       - " + core.T("repl.uninstall_desc"))
	fmt.Println("  lang <zh|en>    - " + core.T("repl.lang_desc"))
	fmt.Println("  quit/exit       - " + core.T("repl.quit_desc"))
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
