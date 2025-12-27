// core/ai/llm/registry.go
// LLMProvider注册核心业务逻辑

package llm

import (
	// 外部包
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// RegisterLLMProvider 注册LLM供应商的入口函数
func RegisterLLMProvider() {
	fmt.Println(lang.T("msg.llm.register.title"))
	var baseURL string

	// 先确认是否要注册
	if !confirmYesNo(lang.T("msg.llm.register.confirm_new")) {
		listAndSelectProviders()
		return
	}

	reader := bufio.NewReader(os.Stdin)

	// 1. 输入供应商类型（循环直到输入有效）
	provider := readNonEmpty(reader, lang.T("msg.llm.register.input_provider"))
	switch provider {
	case "ollama":
		// 5. 输入基础URL（循环直到输入有效）
		baseURL = readNonEmpty(reader, lang.T("msg.llm.register.input_baseurl"))
	default:
		fmt.Println(lang.T("msg.llm.register.unsupported_provider"))
		// 询问是否要连接现有的
		if confirmYesNo(lang.T("msg.llm.register.confirm_link_existing")) {
			listAndSelectProviders()
		}
		return
	}

	// 2. 输入名称（循环直到输入有效且不重复）
	name := readUniqueName(reader)

	// 3. 输入模型名称（循环直到输入有效）
	model := readNonEmpty(reader, lang.T("msg.llm.register.input_model"))

	// 4. 保存配置
	newLLMProvider := cfg.NewLLMProvider(provider, name, model, baseURL)
	switch provider {
	case "ollama":
		// 通过Eino框架去Test连接状况，若失败则无法注册
		resp, err := TestOllamaConnection(newLLMProvider)
		if err != nil {
			fmt.Printf(lang.T("msg.llm.register.test_failed")+"\n", err)
			// 询问是否要连接现有的
			if confirmYesNo(lang.T("msg.llm.register.confirm_link_existing")) {
				listAndSelectProviders()
			}
			return
		}
		fmt.Printf(lang.T("msg.llm.register.test_success")+"\n", resp)
	default:
		fmt.Println(lang.T("msg.llm.register.unsupported_provider"))
		// 询问是否要连接现有的
		if confirmYesNo(lang.T("msg.llm.register.confirm_link_existing")) {
			listAndSelectProviders()
		}
		return
	}

	if err := cfg.AddLLMProvider(newLLMProvider); err != nil {
		fmt.Printf(lang.T("msg.llm.register.save_failed")+"\n", err)
		// 询问是否要连接现有的
		if confirmYesNo(lang.T("msg.llm.register.confirm_link_existing")) {
			listAndSelectProviders()
		}
		return
	}

	if err := cfg.SetActiveLLM(newLLMProvider.ID); err != nil {
		fmt.Printf(lang.T("msg.llm.register.set_active_failed")+"\n", err)
		// 即使激活失败，注册也算成功，只是提示
	}

	// 设置为当前前台 LLM
	if err := cfg.SetCurrentLLM(newLLMProvider.ID); err != nil {
		fmt.Printf(lang.T("msg.llm.register.set_current_failed")+"\n", err)
		// 这个错误不影响注册成功，只是提示
	}

	fmt.Println(lang.T("msg.llm.register.success"))
}

// confirmYesNo 询问用户 Yes/No，只接受有效输入
func confirmYesNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s (Y/n): ", prompt)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		// 空输入默认为 Yes
		if input == "" || input == "y" || input == "yes" {
			return true
		}

		// 明确的 No
		if input == "n" || input == "no" {
			return false
		}

		// 其他输入：提示并重新询问
		fmt.Println(lang.T("msg.llm.register.confirm_invalid"))
	}
}

// readNonEmpty 读取非空输入，一直循环直到输入有效
func readNonEmpty(reader *bufio.Reader, prompt string) string {
	for {
		fmt.Printf("%s: ", prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input != "" {
			return input
		}

		fmt.Println(lang.T("msg.llm.register.input_empty"))
	}
}

// readUniqueName 读取唯一的名称，一直循环直到输入有效且不重复
func readUniqueName(reader *bufio.Reader) string {
	config := cfg.LoadConfig()

	for {
		fmt.Printf("%s", lang.T("msg.llm.register.input_name")+": ")
		input, _ := reader.ReadString('\n')
		name := strings.TrimSpace(input)

		// 检查是否为空
		if name == "" {
			fmt.Println(lang.T("msg.llm.register.name_empty"))
			continue
		}

		// 检查是否重复
		exists := false
		for _, provider := range config.LLMProviders {
			if provider.Name == name {
				exists = true
				break
			}
		}

		if exists {
			fmt.Printf(lang.T("msg.llm.register.name_exists")+"\n", name)
			continue
		}

		return name
	}
}

// TestOllamaConnection 测试Ollama连接是否成功
// 测试成功会自动更新 provider 的验证状态
func TestOllamaConnection(llmprovider *cfg.LLMProvider) (string, error) {
	// 构造请求 URL
	tagsURL := strings.TrimRight(llmprovider.Config.BaseURL, "/") + "/api/tags"

	// 创建带超时的 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 发送 GET 请求
	resp, err := client.Get(tagsURL)
	if err != nil {
		return "", fmt.Errorf(lang.T("msg.llm.register.ollama.connect_failed"), err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(lang.T("msg.llm.register.ollama.status_error"), resp.StatusCode)
	}

	// 解析响应，检查模型是否存在
	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf(lang.T("msg.llm.register.ollama.parse_failed"), err)
	}

	// 检查目标模型是否存在
	modelExists := false
	for _, model := range result.Models {
		if model.Name == llmprovider.Model {
			modelExists = true
			break
		}
	}

	if !modelExists {
		return "", fmt.Errorf(lang.T("msg.llm.register.ollama.model_not_found"),
			llmprovider.Model, llmprovider.Model)
	}

	// ✨ 测试成功，更新验证状态
	llmprovider.Verified = true
	llmprovider.VerifiedAt = time.Now().Format("2006-01-02 15:04:05")

	return fmt.Sprintf(lang.T("msg.llm.register.ollama.test_success"), llmprovider.Model), nil
}

// listAndSelectProviders 列出所有已注册的 Provider 并支持选择 link
func listAndSelectProviders() {
	config := cfg.LoadConfig()
	allProviders := make([]*cfg.LLMProvider, len(config.LLMProviders))
	for i := range config.LLMProviders {
		allProviders[i] = &config.LLMProviders[i]
	}

	if len(allProviders) == 0 {
		fmt.Println(lang.T("msg.llm.link.no_providers"))
		return
	}

	SelectAndLinkProviders(allProviders, lang.T("msg.llm.link.registered_title"), formatProviderWithStatus)
}
