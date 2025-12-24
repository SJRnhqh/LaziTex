package tasks

import (
	// 外部包
	"fmt"
	"os/exec"
	"runtime"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// OpenConfigFile 打开配置文件
//
// 使用系统默认程序打开 LaziTex 配置文件进行编辑。
// 如果配置文件不存在，会先创建默认配置文件。
//
// 功能流程：
//   1. 获取配置文件路径
//   2. 确保配置文件存在（通过加载+保存默认配置）
//   3. 使用系统默认程序打开文件
//
// 支持的操作系统：
//   - Windows: 使用系统默认关联程序
//   - macOS: 使用 open 命令
//   - Linux/Unix: 使用 xdg-open 命令
//
// 可能出现的错误：
//   - 配置路径获取失败（目录权限问题等）
//   - 配置文件保存失败（磁盘空间、权限等）
//   - 系统打开命令失败（程序关联问题等）
func OpenConfigFile() {
	// 获取配置文件路径
	configPath, err := cfg.GetConfigPath()
	if err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.config_path_error"), err)
		return
	}

	// 确保配置文件存在：加载配置 + 保存（创建文件）
	config := cfg.LoadConfig()
	if err := cfg.SaveConfig(config); err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.config_save_error"), err)
		return
	}

	// 使用系统默认程序打开配置文件
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// Windows: 使用系统默认关联程序打开
		cmd = exec.Command("cmd", "/C", "start", "", configPath)
	case "darwin":
		// macOS: 使用 open 命令
		cmd = exec.Command("open", configPath)
	default:
		// Linux/Unix: 使用 xdg-open
		cmd = exec.Command("xdg-open", configPath)
	}

	// 启动打开命令
	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ %s: %v\n", lang.T("msg.config_open_error"), err)
		return
	}

	// 成功提示
	fmt.Printf("✅ %s: %s\n", lang.T("msg.config_opening"), configPath)
}