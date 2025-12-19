// cmd/lazitex-cli/main.go
// 程序入口

package main

import (
	// 外部包
	"fmt"
	"os"
	"strings"

	// 内部包
	tasks "github.com/SJRnhqh/lazitex/cmd/lazitex-cli/tasks"
	ui "github.com/SJRnhqh/lazitex/cmd/lazitex-cli/ui"
	lang "github.com/SJRnhqh/lazitex/lang"
)

func main() {
	// 初始化语言设置并移除 --lang 参数
	args := processLanguageFlag()

	// 如果没有参数，显示帮助
	if len(args) == 0 {
		tasks.ShowHelp()
		return
	}

	// 简单处理第一个参数
	switch args[0] {
	case "-h", "--help":
		tasks.ShowHelp()
	case "-v", "--version":
		fmt.Println("LaziTex v0.0.1")
	case "-c", "--check":
		tasks.CheckEnvironment()
	case "-t", "--tui":
		ui.StartTUI()
	case "-r", "--repl":
		ui.StartREPL()
	case "-i", "--install":
		tasks.InstallLaTexEnvironment()
	case "-u", "--uninstall":
		tasks.UninstallLaTexEnvironment()
	case "-b", "--build":
		if len(args) < 2 {
			fmt.Println(lang.T("msg.build_usage"))
			return
		}

		var filePath string
		var outputPath string
		show := false
		pendingOutput := false // 标记是否正在等待输出路径值

		// 灵活解析：遍历 -b 之后的所有参数
		for i := 1; i < len(args); i++ {
			arg := args[i]
			if arg == "-s" || arg == "--show" {
				show = true
			} else if arg == "-o" || arg == "--output" {
				pendingOutput = true
			} else if strings.HasPrefix(arg, "-") {
				// 严谨处理：未知的标志位直接报错，防止静默错误
				fmt.Printf(lang.T("msg.unknown_command")+"\n", arg)
				return
			} else {
				// 遇到不以 - 开头的参数
				if pendingOutput && outputPath == "" {
					outputPath = arg
					pendingOutput = false
				} else if filePath == "" {
					filePath = arg
				}
			}
		}

		// 检查：如果开启了 -o 但没拿到路径，或者没提供输入文件
		if pendingOutput && outputPath == "" {
			fmt.Println(lang.T("msg.build_usage"))
			return
		}

		if filePath == "" {
			fmt.Println(lang.T("msg.build_usage"))
			return
		}

		tasks.BuildLaTex(filePath, outputPath, show)
	case "-p", "--preview":
		if len(args) < 2 {
			// 如果没传文件名，显示用法
			fmt.Println(lang.T("msg.preview_usage"))
			return
		}
		filePath := args[1]
		tasks.StartLivePreview(filePath)
	default:
		fmt.Printf(lang.T("msg.unknown_command")+"\n", args[0])
		tasks.ShowHelp()
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
			langStr := args[i+1]
			switch langStr {
			case "en", "english":
				lang.SetLanguage(lang.LangEN)
				langSet = true
				// 保存用户偏好
				lang.SaveLanguagePreference(lang.LangEN)
			case "zh", "chinese":
				lang.SetLanguage(lang.LangZH)
				langSet = true
				// 保存用户偏好
				lang.SaveLanguagePreference(lang.LangZH)
			}
			i++ // 跳过语言值
		} else {
			newArgs = append(newArgs, args[i])
		}
	}

	// 如果没有通过 --lang 设置，使用环境检测（会读取配置文件）
	if !langSet {
		lang.SetLanguage(lang.DetectSystemLanguage())
	}

	return newArgs
}
