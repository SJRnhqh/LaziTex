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
	fmt.Printf("  %-38s %s\n", "-h, --help", lang.T("help.show_help"))
	fmt.Printf("  %-38s %s\n", "-v, --version", lang.T("help.show_version"))
	fmt.Printf("  %-38s %s\n", "-c, --check", lang.T("help.check_env"))
	fmt.Printf("  %-38s %s\n", "-r, --repl", lang.T("help.start_repl"))
	fmt.Printf("  %-38s %s\n", "-t, --tui", lang.T("help.start_tui"))
	fmt.Printf("  %-38s %s\n", "-i, --install", lang.T("help.install_latex"))
	fmt.Printf("  %-38s %s\n", "-u, --uninstall", lang.T("help.uninstall_latex"))
	fmt.Printf("  %-38s %s\n", "-b, --build <file> [-o path] [-s]", lang.T("help.build_latex"))
	fmt.Printf("  %-38s %s\n", "-p, --preview <file>", lang.T("help.live_preview"))
	fmt.Printf("  %-38s %s\n", "-l, --lang <lang>", lang.T("help.set_language"))
}

// showREPLHelp 显示帮助信息
func ShowREPLHelp() {
	fmt.Println(lang.T("repl.help_title"))
	fmt.Printf("  %-30s - %s\n", "help", lang.T("repl.help_desc"))
	fmt.Printf("  %-30s - %s\n", "version", lang.T("repl.version_desc"))
	fmt.Printf("  %-30s - %s\n", "check", lang.T("repl.check_desc"))
	fmt.Printf("  %-30s - %s\n", "install", lang.T("repl.install_desc"))
	fmt.Printf("  %-30s - %s\n", "uninstall", lang.T("repl.uninstall_desc"))
	fmt.Printf("  %-30s - %s\n", "build <file> [-o path] [-s]", lang.T("repl.build_desc"))
	fmt.Printf("  %-30s - %s\n", "preview <file>", lang.T("repl.preview_desc"))
	fmt.Printf("  %-30s - %s\n", "lang <zh|en>", lang.T("repl.lang_desc"))
	fmt.Printf("  %-30s - %s\n", "lang", lang.T("repl.lang_current")+lang.GetCurrentLanguageName())
	fmt.Printf("  %-30s - %s\n", "quit / exit", lang.T("repl.quit_desc"))
	fmt.Println()
	fmt.Println(lang.T("repl.shell_title") + ":")
	fmt.Printf("  %-30s - %s\n", "cd [path]", lang.T("repl.help_cd"))
	fmt.Printf("  %-30s - %s\n", "ls [path]", lang.T("repl.help_ls"))
	fmt.Printf("  %-30s - %s\n", "pwd", lang.T("repl.help_pwd"))
	fmt.Printf("  %-30s - %s\n", "clear", lang.T("repl.help_clear"))
	fmt.Printf("  %-30s - %s\n", "cat <file>", lang.T("repl.help_cat"))
}
