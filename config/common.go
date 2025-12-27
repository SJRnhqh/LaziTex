// config/common.go
// 负责管理应用配置

package config

import (
	// 外部包
	"encoding/json"
	"os"
	"path/filepath"
)

// Config 应用配置
type Config struct {
	Language    string `json:"language"`    // 用户偏好语言: "zh" 或 "en"
	DisplayMode string `json:"displayMode"` // 显示模式: "virtual" 或 "pagination"

	// LLM 配置
	LLMProviders []LLMProvider `json:"llmProviders,omitempty"` // LLMProvider列表
	ActiveLLMs   []string      `json:"activeLLMs,omitempty"`   // 当前连接的 LLM Provider ID 列表（支持多激活）
	CurrentLLM   string        `json:"currentLLM,omitempty"`   // 当前激活的 LLM Provider ID（用于 ask/chat）
}

// GetConfigPath 获取配置文件路径
//
// 返回 LaziTex 配置文件的完整路径，位置因操作系统而异：
//   - Linux/Unix: ~/.config/lazitex/config.json
//   - macOS: ~/Library/Application Support/lazitex/config.json
//   - Windows: %APPDATA%\lazitex\config.json
//
// 可能返回的错误：
//   - os.UserConfigDir() 失败：无法获取用户配置目录
//   - 用户主目录不存在
//   - 系统环境变量问题（Windows: %APPDATA%, Unix: $XDG_CONFIG_HOME 或 ~/.config）
//   - os.MkdirAll() 失败：无法创建配置目录
//   - 父目录权限不足
//   - 磁盘空间不足
//   - 文件系统只读
//
// 返回值：
//   - string: 配置文件路径
//   - error: 错误信息，如果成功则为nil
func GetConfigPath() (string, error) {
	// 获取用户配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// 创建 LaziTex 配置目录
	lazitexDir := filepath.Join(configDir, "lazitex")
	if err := os.MkdirAll(lazitexDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(lazitexDir, "config.json"), nil
}

// LoadConfig 加载配置
//
// 加载用户配置文件，如果文件不存在或损坏则返回默认配置。
// 此函数保证总是返回有效的配置对象，不会出错。
//
// 配置文件位置因操作系统而异：
//   - Linux/Unix: ~/.config/lazitex/config.json
//   - macOS: ~/Library/Application Support/lazitex/config.json
//   - Windows: %APPDATA%\lazitex\config.json
//
// 返回值：
//   - *Config: 有效的配置对象，失败时返回默认配置
func LoadConfig() *Config {
	configPath, err := GetConfigPath()
	if err != nil {
		return &Config{
			Language:     "en",
			DisplayMode:  "virtual",
			LLMProviders: []LLMProvider{},
			ActiveLLMs:   []string{},
			CurrentLLM:   "",
		}
	}

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Language:     "en",
			DisplayMode:  "virtual",
			LLMProviders: []LLMProvider{},
			ActiveLLMs:   []string{},
		}
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return &Config{
			Language:     "en",
			DisplayMode:  "virtual",
			LLMProviders: []LLMProvider{},
			ActiveLLMs:   []string{},
		}
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return &Config{
			Language:     "en",
			DisplayMode:  "virtual",
			LLMProviders: []LLMProvider{},
			ActiveLLMs:   []string{},
		}
	}

	// 设置默认值
	if config.DisplayMode == "" {
		config.DisplayMode = "virtual"
	}
	if config.Language == "" {
		config.Language = "en"
	}
	if config.LLMProviders == nil {
		config.LLMProviders = []LLMProvider{}
	}
	if config.ActiveLLMs == nil {
		config.ActiveLLMs = []string{}
	}
	if config.CurrentLLM == "" {
	}

	return &config
}

// SaveConfig 保存配置
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// UpdateConfig 部分更新配置（只更新提供的字段）
func UpdateConfig(updates map[string]interface{}) error {
	// 加载现有配置
	config := LoadConfig()

	// 更新提供的字段
	if displayMode, ok := updates["displayMode"].(string); ok {
		if displayMode == "virtual" || displayMode == "pagination" {
			config.DisplayMode = displayMode
		}
	}
	if language, ok := updates["language"].(string); ok {
		if language == "zh" || language == "en" {
			config.Language = language
		}
	}
	// 支持更新 activeLLMs（列表）
	if activeLLMs, ok := updates["activeLLMs"].([]string); ok {
		config.ActiveLLMs = activeLLMs
	}
	// 向后兼容：如果传入了 activeLLM（单值），转换为列表
	if activeLLM, ok := updates["activeLLM"].(string); ok {
		if activeLLM != "" {
			// 如果不在列表中，添加
			found := false
			for _, id := range config.ActiveLLMs {
				if id == activeLLM {
					found = true
					break
				}
			}
			if !found {
				config.ActiveLLMs = append(config.ActiveLLMs, activeLLM)
			}
		} else {
			// 如果为空字符串，清空列表（向后兼容）
			config.ActiveLLMs = []string{}
		}
	}

	// 保存更新后的配置
	return SaveConfig(config)
}
