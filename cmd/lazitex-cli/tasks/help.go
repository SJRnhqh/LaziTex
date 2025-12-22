// cmd/lazitex-cli/tasks/help.go
// 帮助类业务管理

package tasks

import (
	// 外部包
	"fmt"

	// 内部包
	lang "github.com/SJRnhqh/lazitex/lang"
)

// 显示命令帮助
func ShowHelp() {
	fmt.Println(lang.T("help.description"))
	fmt.Println()
	fmt.Println(lang.T("help.usage"))
	fmt.Println("  lazitex [command]")
	fmt.Println()
	fmt.Println(lang.T("help.commands"))
	// 使用 75 字符宽度对齐，确保所有命令和描述都能整齐对齐
	fmt.Printf("  %-75s %s\n", "-h, --help", lang.T("help.show_help"))
	fmt.Printf("  %-75s %s\n", "-v, --version", lang.T("help.show_version"))
	fmt.Printf("  %-75s %s\n", "-c, --check", lang.T("help.check_env"))
	fmt.Printf("  %-75s %s\n", "-r, --repl", lang.T("help.start_repl"))
	fmt.Printf("  %-75s %s\n", "-t, --tui", lang.T("help.start_tui"))
	fmt.Printf("  %-75s %s\n", "-i, --install", lang.T("help.install_latex"))
	fmt.Printf("  %-75s %s\n", "-u, --uninstall", lang.T("help.uninstall_latex"))
	fmt.Printf("  %-75s %s\n", "-o, --ollama <-c|--check|-i|--install|-u|--uninstall|-s|--status>", lang.T("help.ollama"))
	fmt.Printf("  %-75s %s\n", "-m, --llm [add|link|list|test|remove|unlink|<alias>]", "LLM 管理（智能/显式注册、链接、列表、测试、删除）")
	fmt.Printf("  %-75s %s\n", "-b, --build <file> [-o|--output path] [-s|--show] [-q|--quiet] [-t|--tidy]", lang.T("help.build_latex"))
	fmt.Printf("  %-75s %s\n", "-p, --preview <file>", lang.T("help.live_preview"))
	fmt.Printf("  %-75s %s\n", "-l, --lang <lang>", lang.T("help.set_language"))
}

// showREPLHelp 显示帮助信息
func ShowREPLHelp() {
	fmt.Println(lang.T("repl.help_title"))
	// 使用 75 字符宽度对齐，确保所有命令和描述都能整齐对齐
	fmt.Printf("  %-75s - %s\n", "help", lang.T("repl.help_desc"))
	fmt.Printf("  %-75s - %s\n", "version", lang.T("repl.version_desc"))
	fmt.Printf("  %-75s - %s\n", "check", lang.T("repl.check_desc"))
	fmt.Printf("  %-75s - %s\n", "install", lang.T("repl.install_desc"))
	fmt.Printf("  %-75s - %s\n", "uninstall", lang.T("repl.uninstall_desc"))
	fmt.Printf("  %-75s - %s\n", "ollama <-c|--check|-i|--install|-u|--uninstall|-s|--status>", lang.T("repl.ollama_check_desc"))
	fmt.Printf("  %-75s - %s\n", "llm [add|link|list|test|remove|unlink|<alias>]", "LLM 管理（智能/显式注册、链接、列表、测试、删除）")
	fmt.Printf("  %-75s - %s\n", "build <file> [-o|--output path] [-s|--show] [-q|--quiet] [-t|--tidy]", lang.T("repl.build_desc"))
	fmt.Printf("  %-75s - %s\n", "preview <file>", lang.T("repl.preview_desc"))
	fmt.Printf("  %-75s - %s\n", "lang <zh|en>", lang.T("repl.lang_desc"))
	fmt.Printf("  %-75s - %s\n", "lang", lang.T("repl.lang_current")+lang.GetCurrentLanguageName())
	fmt.Printf("  %-75s - %s\n", "quit / exit", lang.T("repl.quit_desc"))
	fmt.Println()
	fmt.Println(lang.T("repl.shell_title") + ":")
	fmt.Printf("  %-75s - %s\n", "cd [path]", lang.T("repl.help_cd"))
	fmt.Printf("  %-75s - %s\n", "ls [path]", lang.T("repl.help_ls"))
	fmt.Printf("  %-75s - %s\n", "pwd", lang.T("repl.help_pwd"))
	fmt.Printf("  %-75s - %s\n", "clear", lang.T("repl.help_clear"))
	fmt.Printf("  %-75s - %s\n", "cat <file>", lang.T("repl.help_cat"))
}

// LLM 用法提示
func PrintLLMUsage() {
	fmt.Println("用法: lazitex -m [add|link|list|test|remove|unlink|<别名>] <别名>")
	fmt.Println("示例:")
	fmt.Println("  lazitex -m gemma3          # 智能：不存在则注册+链接，存在则直接链接")
	fmt.Println("  lazitex -m add gemma3      # 显式注册")
	fmt.Println("  lazitex -m link gm3        # 显式链接已存在的")
	fmt.Println("  lazitex -m list            # 列出所有")
	fmt.Println("  lazitex -m test gm3        # 测试配置存在/启用状态")
	fmt.Println("  lazitex -m remove gm3      # 删除（unlink 同效）")
}

// LLM 用法提示（REPL）
func PrintREPLLLMUsage() {
	fmt.Println("用法: llm <别名> | llm [add|link|list|test|remove|unlink] <别名>")
	fmt.Println("示例:")
	fmt.Println("  llm gemma3          # 智能：不存在则注册+链接，存在则直接链接")
	fmt.Println("  llm add gemma3      # 显式注册")
	fmt.Println("  llm link gm3        # 显式链接已存在的")
	fmt.Println("  llm list            # 列出所有")
	fmt.Println("  llm test gm3        # 测试配置存在/启用状态")
	fmt.Println("  llm remove gm3      # 删除（unlink 同效）")
}