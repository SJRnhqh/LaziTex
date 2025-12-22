package tasks

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// OpenConfigFile 打开配置文件（不存在则先写入默认配置）
func OpenConfigFile() {
	configPath, err := cfg.GetConfigPath()
	if err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.config_path_error"), err)
		return
	}

	// 确保文件存在：加载后保存一次默认配置
	conf, err := cfg.LoadConfig()
	if err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.config_load_error"), err)
		return
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := cfg.SaveConfig(conf); err != nil {
			fmt.Printf("❌ %s: %v\n", lang.T("msg.config_save_error"), err)
			return
		}
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", "", configPath) // 使用系统默认关联程序
	case "darwin":
		cmd = exec.Command("open", configPath)
	default:
		cmd = exec.Command("xdg-open", configPath)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.config_open_error"), err)
		return
	}

	fmt.Printf("✅ %s: %s\n", lang.T("msg.config_opening"), configPath)
}
