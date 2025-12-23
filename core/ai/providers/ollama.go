// core/ai/providers/ollama.go
// Ollama 本地部署 LLM 接口实现

package providers

import (
	// 外部包
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	// 内部包
	lang "github.com/SJRnhqh/lazitex/lang"
)

// OllamaProvider Ollama 本地部署 LLM 接口实现
type OllamaProvider struct {
	Model  string            `json:"model"`  // 模型名称
	Config map[string]string `json:"config"` // 配置
}

// Connectivity 测试 Ollama 模型是否可用
func (p *OllamaProvider) Connectivity(model string, config map[string]string) (bool, error) {
	// 第一步：检查 Ollama 服务是否运行
	if !p.checkService() {
		return false, fmt.Errorf("%s", lang.T("msg.llm.provider.ollama.service_not_running"))
	}

	// 第二步：检查模型是否存在
	exists, err := p.checkModelExists(model)
	if err != nil {
		return false, err // API 调用失败，返回错误
	}
	if !exists {
		return false, fmt.Errorf("%s", lang.T("msg.llm.provider.ollama.model_not_found"))
	}

	// 第三步：测试模型可用性（发送简单请求验证模型能响应）
	available, err := p.testModelGeneration(model)
	if err != nil {
		return false, err // 生成请求失败，返回错误
	}
	if !available {
		return false, nil // 模型不可用，返回 false，不返回错误
	}

	return true, nil
}

// checkService 检查 Ollama 服务是否运行（默认端口 11434）
func (p *OllamaProvider) checkService() bool {
	conn, err := net.DialTimeout("tcp", "localhost:11434", 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// checkModelExists 检查模型是否在 Ollama 中已安装
func (p *OllamaProvider) checkModelExists(modelName string) (bool, error) {
	// 调用 Ollama API 获取已安装的模型列表
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false, fmt.Errorf(lang.T("msg.llm.provider.ollama.api_connect_failed"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf(lang.T("msg.llm.provider.ollama.api_status_error"), resp.StatusCode)
	}

	// 解析 JSON 响应
	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf(lang.T("msg.llm.provider.ollama.parse_response_failed"), err)
	}

	// 检查模型是否存在（支持完整匹配和带 :latest 后缀的匹配）
	for _, model := range result.Models {
		if model.Name == modelName || model.Name == modelName+":latest" || strings.TrimSuffix(model.Name, ":latest") == modelName {
			return true, nil
		}
	}

	return false, nil
}

// testModelGeneration 测试模型是否能够生成响应（验证模型可用性）
func (p *OllamaProvider) testModelGeneration(modelName string) (bool, error) {
	// 准备请求体（发送一个简单的测试 prompt）
	requestBody := map[string]interface{}{
		"model":  modelName,
		"prompt": "test", // 简单的测试 prompt
		"stream": false,  // 非流式响应，获取完整结果
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return false, fmt.Errorf(lang.T("msg.llm.provider.ollama.serialize_request_failed"), err)
	}

	// 发送 POST 请求到 Ollama API
	client := &http.Client{Timeout: 60 * time.Second} // 增加超时时间以支持模型加载
	resp, err := client.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		return false, fmt.Errorf(lang.T("msg.llm.provider.ollama.api_connect_failed"), err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf(lang.T("msg.llm.provider.ollama.api_status_error"), resp.StatusCode)
	}

	// 解析响应（我们只需要确认响应成功，不需要完整解析内容）
	var result struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf(lang.T("msg.llm.provider.ollama.parse_response_failed"), err)
	}

	// 如果响应成功且有内容，说明模型可用
	return result.Done, nil
}