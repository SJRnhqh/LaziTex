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
	// 使用 100 字符宽度对齐，确保所有命令和描述都能整齐对齐，中间有足够空间
	fmt.Printf("  %-100s %s\n", "-h, --help", lang.T("help.show_help"))
	fmt.Printf("  %-100s %s\n", "-v, --version", lang.T("help.show_version"))
	fmt.Printf("  %-100s %s\n", "-r, --repl", lang.T("help.start_repl"))
	fmt.Printf("  %-100s %s\n", "-t, --tui", lang.T("help.start_tui"))
	fmt.Printf("  %-100s %s\n", "-x, --latex <-c|--check|-i|--install|-u|--uninstall>", lang.T("help.latex"))
	fmt.Printf("  %-100s %s\n", "-o, --ollama <-c|--check|-i|--install|-u|--uninstall|-s|--status>", lang.T("help.ollama"))
	// fmt.Printf("  %-100s %s\n", "-m, --llm add <model> -p <provider> -n <name> | remove/link/unlink/test/switch <id/name/model> | list", lang.T("help.llm"))
	fmt.Printf("  %-100s %s\n", "-b, --build <file> [-o|--output path] [-s|--show] [-q|--quiet] [-t|--tidy]", lang.T("help.build_latex"))
	fmt.Printf("  %-100s %s\n", "-p, --preview <file> [-q|--quiet] [-t|--tidy] [:port]", lang.T("help.live_preview"))
	fmt.Printf("  %-100s %s\n", "-l, --lang <lang>", lang.T("help.set_language"))
	fmt.Printf("  %-100s %s\n", "config", lang.T("help.config"))
}

// showREPLHelp 显示帮助信息
func ShowREPLHelp() {
	fmt.Println(lang.T("repl.help_title"))
	// 使用 100 字符宽度对齐，确保所有命令和描述都能整齐对齐，中间有足够空间
	fmt.Printf("  %-100s - %s\n", "help", lang.T("repl.help_desc"))
	fmt.Printf("  %-100s - %s\n", "version", lang.T("repl.version_desc"))
	fmt.Printf("  %-100s - %s\n", "latex <-c|--check|-i|--install|-u|--uninstall>", lang.T("repl.latex"))
	fmt.Printf("  %-100s - %s\n", "ollama <-c|--check|-i|--install|-u|--uninstall|-s|--status>", lang.T("repl.ollama"))
	fmt.Printf("  %-100s - %s\n", "llm [list|link|unlink|switch|test|remove|ask|chat] ...", lang.T("repl.llm"))
	fmt.Printf("  %-100s - %s\n", "build <file> [-o|--output path] [-s|--show] [-q|--quiet] [-t|--tidy]", lang.T("repl.build_desc"))
	fmt.Printf("  %-100s - %s\n", "preview <file> [-q|--quiet] [-t|--tidy] [:port]", lang.T("repl.preview_desc"))
	fmt.Printf("  %-100s - %s\n", "lang <zh|en>", lang.T("repl.lang_desc"))
	fmt.Printf("  %-100s - %s\n", "lang", lang.T("repl.lang_current")+lang.GetCurrentLanguageName())
	fmt.Printf("  %-100s - %s\n", "config", lang.T("repl.config_desc"))
	fmt.Printf("  %-100s - %s\n", "quit / exit", lang.T("repl.quit_desc"))
	fmt.Println()
	fmt.Println(lang.T("repl.shell_title") + ":")
	fmt.Printf("  %-100s - %s\n", "cd [path]", lang.T("repl.help_cd"))
	fmt.Printf("  %-100s - %s\n", "ls [path]", lang.T("repl.help_ls"))
	fmt.Printf("  %-100s - %s\n", "pwd", lang.T("repl.help_pwd"))
	fmt.Printf("  %-100s - %s\n", "clear", lang.T("repl.help_clear"))
	fmt.Printf("  %-100s - %s\n", "cat <file>", lang.T("repl.help_cat"))
}
