// core/ai/llm/chat.go
// 处理LLM的持续性对话交互

package llm

import (
	// 外部包
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	prompt "github.com/cloudwego/eino/components/prompt"
	schema "github.com/cloudwego/eino/schema"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

func ChatLLMProvider() {
	// 1. 从配置中获取当前LLM Provider
	currentProvider, err := cfg.GetCurrentLLM()
	if err != nil {
		fmt.Printf(lang.T("msg.llm.chat.get_current_failed")+"\n", err)
		return
	}
	if currentProvider == nil {
		fmt.Println(lang.T("msg.llm.chat.no_current_llm"))
		return
	}

	// 2. 创建基础 context（用于模型初始化）
	baseCtx := context.Background()

	// 3. 创建 ChatModel 客户端（使用工厂函数）
	client, err := NewChatModelClient(baseCtx, currentProvider)
	if err != nil {
		fmt.Printf(lang.T("msg.llm.chat.create_model_failed")+"\n", err)
		return
	}

	// 4. 创建对话模板
	tmpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(lang.T("msg.llm.chat.system_message")),
		schema.MessagesPlaceholder("history", false),
		schema.UserMessage("{input}"),
	)

	var history []*schema.Message

	// 5. 显示欢迎信息
	fmt.Println(lang.T("msg.llm.chat.welcome"))
	fmt.Println(lang.T("msg.llm.chat.exit_hint"))

	scanner := bufio.NewScanner(os.Stdin)

	// 6. 对话循环
	for {
		fmt.Print(lang.T("msg.llm.chat.user_prompt"))
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// 检查退出命令
		if isExitCommand(input) {
			fmt.Println(lang.T("msg.llm.chat.goodbye"))
			break
		}

		// 验证输入是否为空
		if input == "" {
			continue
		}

		// 为每次请求创建新的超时 context（30秒）
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		// 使用 template 将 input + history 格式化为完整消息列表
		messages, err := tmpl.Format(ctx, map[string]any{
			"input":   input,
			"history": history,
		})
		if err != nil {
			cancel()
			fmt.Printf(lang.T("msg.llm.chat.format_failed")+"\n", err)
			continue
		}

		// 调用模型（流式）
		streamResult, err := client.Stream(ctx, messages)
		if err != nil {
			cancel()
			// 检查是否是超时错误
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Println(lang.T("msg.llm.chat.timeout"))
			} else {
				fmt.Printf(lang.T("msg.llm.chat.stream_failed")+"\n", err)
			}
			continue
		}

		// 显示 AI 响应前缀
		fmt.Print(lang.T("msg.llm.chat.ai_prompt"))

		// 流式输出响应
		var fullResponse string
		for {
			chunk, err := streamResult.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				// 检查是否是超时错误
				if ctx.Err() == context.DeadlineExceeded {
					fmt.Println("\n" + lang.T("msg.llm.chat.timeout"))
				} else {
					fmt.Printf("\n"+lang.T("msg.llm.chat.receive_error")+"\n", err)
				}
				break
			}
			fmt.Print(chunk.Content)
			fullResponse += chunk.Content
		}
		fmt.Println() // 换行

		cancel()

		// 将本轮对话加入历史（用于下一轮）
		history = append(history,
			&schema.Message{Role: schema.User, Content: input},
			&schema.Message{Role: schema.Assistant, Content: fullResponse},
		)
	}
}

// isExitCommand 检查是否是退出命令
func isExitCommand(input string) bool {
	lowerInput := strings.ToLower(strings.TrimSpace(input))
	return lowerInput == "quit" || lowerInput == "exit" || lowerInput == "q"
}
