// config/config.go
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
	Language string `json:"language"` // 用户偏好语言: "zh" 或 "en"
	DisplayMode string `json:"displayMode"` // 显示模式: "virtual" 或 "pagination"
}

// getConfigPath 获取配置文件路径
func getConfigPath() (string, error) {
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
func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return &Config{Language: "en", DisplayMode: "virtual"}, nil // 默认英文
	}

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{Language: "en", DisplayMode: "virtual"}, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return &Config{Language: "en", DisplayMode: "virtual"}, nil
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return &Config{Language: "en", DisplayMode: "virtual"}, nil
	}

	if config.DisplayMode == "" {
		config.DisplayMode = "virtual"
	}

	if config.Language == "" {
		config.Language = "en"
	}

	return &config, nil
}

// SaveConfig 保存配置
func SaveConfig(config *Config) error {
	configPath, err := getConfigPath()
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
	config, err := LoadConfig()
	if err != nil {
		return err
	}

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

	// 保存更新后的配置
	return SaveConfig(config)
}