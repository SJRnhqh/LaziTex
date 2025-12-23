// cmd/lazitex-cli/main.go
// 程序入口

package main

import (
	// 外部包
	"fmt"
	"os"
	"strconv"
	"strings"

	// 内部包
	internal "github.com/SJRnhqh/lazitex/cmd/lazitex-cli/internal"
	tasks "github.com/SJRnhqh/lazitex/cmd/lazitex-cli/tasks"
	ui "github.com/SJRnhqh/lazitex/cmd/lazitex-cli/ui"
	lang "github.com/SJRnhqh/lazitex/lang"
)

func main() {
	// 初始化 ollama 动作处理函数（只需一次）
	internal.SetOllamaActionHandlers(map[internal.OllamaAction]func(){
		internal.OllamaCheck:     tasks.CheckOllama,
		internal.OllamaInstall:   tasks.InstallOllama,
		internal.OllamaUninstall: tasks.UninstallOllama,
		internal.OllamaStatus:    tasks.StatusOllama,
	})

	// 初始化语言设置并移除 --lang 参数
	args := processLanguageFlag()

	// 如果没有其他参数，输出当前语言即可
	if len(args) == 0 {
		fmt.Println(lang.T("repl.lang_current") + lang.GetCurrentLanguageName())
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
		tasks.InstallLaTeXEnvironment()

	case "-u", "--uninstall":
		tasks.UninstallLaTeXEnvironment()

	case "config":
		tasks.OpenConfigFile()

	case "-o", "--ollama":// Ollama 管理 TODO: 测试覆盖
		internal.HandleOllamaCommand(args[1:], internal.I18nModeMsg, "msg.ollama_usage")
		return

	case "-m", "--llm":
		// LLM 管理：智能 / add / link / list / test / remove(unlink)
		if len(args) < 2 {
			fmt.Println(lang.T("msg.llm.cli_usage"))
			return
		}

		sub := args[1:]
		action := sub[0]

		// 先处理不需要额外参数的命令（如 list）
		switch action {
		case "list":
			tasks.ListLLM()
			return
		case "add":
			if len(sub) < 2 {
				fmt.Println(lang.T("msg.llm.cli_usage"))
				return
			}
			tasks.AddLLM(sub[1])
			return
		case "link":
			if len(sub) < 2 {
				fmt.Println(lang.T("msg.llm.cli_usage"))
				return
			}
			tasks.LinkLLM(sub[1])
			return
		case "test":
			if len(sub) < 2 {
				fmt.Println(lang.T("msg.llm.cli_usage"))
				return
			}
			tasks.TestLLM(sub[1])
			return
		case "remove", "unlink":
			if len(sub) < 2 {
				fmt.Println(lang.T("msg.llm.cli_usage"))
				return
			}
			tasks.RemoveLLM(sub[1])
			return
		default:
			// 默认当作模型名（智能模式）：暂时不实现，显示提示
			if len(sub) == 1 {
				fmt.Println("⚠️  智能注册/链接功能暂未实现")
				fmt.Println("请使用: lazitex -m add <model> 注册，或 lazitex -m link <model> 链接")
				return
			}
			// 未知命令
			fmt.Printf("未知操作: %s\n", action)
			fmt.Println(lang.T("msg.llm.cli_usage"))
		}

	case "-b", "--build":
		if len(args) < 2 {
			fmt.Println(lang.T("msg.build_usage"))
			return
		}

		var filePath string
		var outputPath string
		show := false
		quiet := false
		tidy := false
		pendingOutput := false // 标记是否正在等待输出路径值

		// 灵活解析：遍历 -b 之后的所有参数
		for i := 1; i < len(args); i++ {
			arg := args[i]
			if arg == "-s" || arg == "--show" {
				show = true
			} else if arg == "-q" || arg == "--quiet" {
				quiet = true
			} else if arg == "-t" || arg == "--tidy" {
				tidy = true
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

		tasks.BuildLaTeX(filePath, outputPath, show, quiet, tidy)

	case "-p", "--preview":
		if len(args) < 2 {
			// 如果没传文件名，显示用法
			fmt.Println(lang.T("msg.preview_usage"))
			return
		}

		var filePath string
		quiet := false
		tidy := false
		devMode := false
		port := 8080

		// 灵活解析：遍历 -p 之后的所有参数
		for i := 1; i < len(args); i++ {
			arg := args[i]

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

			// 处理未知标志位
			if strings.HasPrefix(arg, "-") {
				fmt.Printf(lang.T("msg.unknown_command")+"\n", arg)
				return
			}

			// 处理端口号参数（独立的:port格式）
			if strings.HasPrefix(arg, ":") {
				portStr := arg[1:]
				parsedPort, err := strconv.Atoi(portStr)
				if err != nil || parsedPort < 1 || parsedPort > 65535 {
					fmt.Println(lang.T("msg.invalid_port"))
					return
				}
				port = parsedPort
				continue
			}

			if filePath == "" {
				// 检查是否包含端口号
				if strings.Contains(arg, ":") {
					lastColonIndex := strings.LastIndex(arg, ":")
					if lastColonIndex > 0 && lastColonIndex < len(arg)-1 {
						portStr := arg[lastColonIndex+1:]
						parsedPort, err := strconv.Atoi(portStr)
						// 如果解析成功，说明是 file.tex:8080 格式
						if err != nil && parsedPort >= 1 && parsedPort <= 65535 {
							// 提取文件路径部分（去掉端口号）
							filePath = arg[:lastColonIndex]
							port = parsedPort
							continue
						}
						// 如果解析失败，说明冒号后面不是端口号（比如 Windows 路径 C:\）
						// 当作普通文件路径处理
					}
				}
				// 普通文件路径（不包含端口号，或者包含冒号但不是端口号格式）
				filePath = arg
			}
		}
		// 检查是否提供了文件路径
		if filePath == "" {
			fmt.Println(lang.T("msg.preview_usage"))
			return
		}

		tasks.StartLivePreview(filePath, port, quiet, tidy, devMode)
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
