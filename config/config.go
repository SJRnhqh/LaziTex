// config/config.go
// 负责管理应用配置

package config

import (
	// 外部包
	"encoding/json"
	"os"
	"path/filepath"

	// 内部包
	"github.com/SJRnhqh/lazitex/model"
)

// Config 应用配置
type Config struct {
	Language string `json:"language"` // 用户偏好语言: "zh" 或 "en"
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
		return &Config{Language: "en"}, nil // 默认英文
	}

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{Language: "en"}, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return &Config{Language: "en"}, nil
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return &Config{Language: "en"}, nil
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

// SaveLanguagePreference 保存语言偏好
func SaveLanguagePreference(lang model.Language) error {
	config := &Config{
		Language: string(lang),
	}
	return SaveConfig(config)
}

// LoadLanguagePreference 加载语言偏好
func LoadLanguagePreference() model.Language {
	config, err := LoadConfig()
	if err != nil {
		return model.LangEN // 默认英文
	}

	if config.Language == "zh" {
		return model.LangZH
	}
	return model.LangEN
}
