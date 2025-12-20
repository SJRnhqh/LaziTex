// target/mac/ollama.go

package mac

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/SJRnhqh/lazitex/lang"
)

// OllamaManager macOS Ollama 管理器
type OllamaManager struct{}

// NewOllamaManager 创建 macOS Ollama 管理器
func NewOllamaManager() *OllamaManager {
	return &OllamaManager{}
}

// hasHomebrew 检查 Homebrew 是否可用
func (m *OllamaManager) hasHomebrew() bool {
	if runtime.GOOS != "darwin" {
		return false
	}

	cmd := exec.Command("brew", "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// askForConfirmation 询问用户确认 (y/n)，其他键重新询问
func (m *OllamaManager) askForConfirmation(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.TrimSpace(strings.ToLower(response))

		switch response {
		case "y", "yes", "":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			// 其他键，重新询问
			fmt.Println(lang.T("msg.ollama.invalid_input"))
		}
	}
}

// Check 检测 Ollama 是否安装
func (m *OllamaManager) Check() (bool, string, error) {
	// 1. 检查 ollama 命令是否存在（优先检查 PATH）
	_, err := exec.LookPath("ollama")
	if err == nil {
		// 在 PATH 中找到，检查版本
		cmd := exec.Command("ollama", "--version")
		output, err := cmd.Output()
		if err != nil {
			return true, "installed", nil // 已安装但无法获取版本
		}
		version := strings.TrimSpace(string(output))
		return true, version, nil
	}

	// 2. 如果 PATH 中找不到，检查常见安装路径
	commonPaths := []string{
		"/usr/local/bin/ollama",
		"/opt/homebrew/bin/ollama",
		filepath.Join(os.Getenv("HOME"), ".local/bin/ollama"),
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			// 找到安装文件，尝试运行获取版本
			cmd := exec.Command(path, "--version")
			output, err := cmd.Output()
			if err != nil {
				return true, "installed", nil // 已安装但无法获取版本
			}
			version := strings.TrimSpace(string(output))
			return true, version, nil
		}
	}

	// 3. 检查是否通过 Homebrew 安装（即使不在 PATH 中）
	if m.hasHomebrew() {
		// 检查是否通过 Homebrew 安装
		cmd := exec.Command("brew", "list", "--formula", "ollama")
		if err := cmd.Run(); err == nil {
			// 已通过 Homebrew 安装
			cmd = exec.Command("brew", "info", "--formula", "ollama")
			output, err := cmd.Output()
			if err == nil {
				// 尝试从输出中提取版本信息
				lines := strings.Split(string(output), "\n")
				for _, line := range lines {
					if strings.Contains(line, "ollama:") {
						parts := strings.Fields(line)
						if len(parts) > 1 {
							return true, parts[1], nil
						}
					}
				}
			}
			return true, "installed", nil
		}

		// 检查是否通过 Homebrew Cask 安装
		cmd = exec.Command("brew", "list", "--cask", "ollama")
		if err := cmd.Run(); err == nil {
			return true, "installed", nil
		}
	}

	return false, "", nil // 未安装
}

// checkService 检查 Ollama 服务是否运行
func (m *OllamaManager) checkService() bool {
	// 检查端口 11434 是否监听
	conn, err := net.Dial("tcp", "localhost:11434")
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Install 安装 Ollama
func (m *OllamaManager) Install() error {
	// 1. 检查是否已安装
	if installed, version, _ := m.Check(); installed {
		fmt.Printf(lang.T("msg.ollama.already_installed")+"\n", version)

		// 检查服务是否运行
		if !m.checkService() {
			fmt.Println(lang.T("msg.ollama.service_not_running"))
			fmt.Println(lang.T("msg.ollama.start_service_hint"))
		} else {
			fmt.Println(lang.T("msg.ollama.service_running"))
		}
		return nil // 已安装，返回 nil 而不是错误
	}

	// 2. 检查 Homebrew 是否可用
	if !m.hasHomebrew() {
		// 如果没有 Homebrew，使用官方安装脚本
		return m.installViaScript()
	}

	// 3. 询问用户使用哪种方式安装
	fmt.Println("🦙 检测到已安装 Homebrew，可以选择以下安装方式：")
	fmt.Println("  1. Homebrew (推荐，快速安装)")
	fmt.Println("  2. 官方安装脚本 (直接安装)")
	fmt.Print("请选择 [1/2] (默认: 1): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "2" {
		return m.installViaScript()
	}

	// 4. 通过 Homebrew 安装
	fmt.Println(lang.T("msg.ollama.mac.installing_via_homebrew"))
	fmt.Println("💡 注意: 这将通过 Homebrew 安装 Ollama")

	// 使用 brew install ollama (formula，不是 cask)
	cmd := exec.Command("brew", "install", "ollama")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.ollama.install_failed"), err)
	}

	// 5. 等待安装完成并验证（增加重试机制）
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		time.Sleep(2 * time.Second)
		if installed, version, _ := m.Check(); installed {
			fmt.Printf(lang.T("msg.ollama.install_success")+"\n", version)

			// 检查服务是否运行
			if m.checkService() {
				fmt.Println(lang.T("msg.ollama.service_running"))
			} else {
				fmt.Println(lang.T("msg.ollama.service_not_running"))
				fmt.Println(lang.T("msg.ollama.start_service_hint"))
			}

			// 检查是否在 PATH 中，如果不在则提示用户
			if _, err := exec.LookPath("ollama"); err != nil {
				fmt.Println(lang.T("msg.ollama.path_not_updated"))
				fmt.Println(lang.T("msg.ollama.restart_terminal_hint"))
			}

			return nil
		}
	}

	// 安装可能成功但 PATH 未更新，提示用户
	fmt.Println(lang.T("msg.ollama.install_maybe_success"))
	fmt.Println(lang.T("msg.ollama.restart_terminal_hint"))
	return fmt.Errorf("%s", lang.T("msg.ollama.install_verify_failed"))
}

// installViaScript 通过官方安装脚本安装
func (m *OllamaManager) installViaScript() error {
	fmt.Println(lang.T("msg.ollama.mac.installing_via_script"))
	fmt.Println("💡 注意: 这将使用官方安装脚本安装 Ollama")

	// 下载并运行官方安装脚本
	cmd := exec.Command("sh", "-c", "curl -fsSL https://ollama.com/install.sh | sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.ollama.install_failed"), err)
	}

	// 等待安装完成并验证
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		time.Sleep(2 * time.Second)
		if installed, version, _ := m.Check(); installed {
			fmt.Printf(lang.T("msg.ollama.install_success")+"\n", version)

			// 检查服务是否运行
			if m.checkService() {
				fmt.Println(lang.T("msg.ollama.service_running"))
			} else {
				fmt.Println(lang.T("msg.ollama.service_not_running"))
				fmt.Println(lang.T("msg.ollama.start_service_hint"))
			}

			return nil
		}
	}

	return fmt.Errorf("%s", lang.T("msg.ollama.install_verify_failed"))
}

// Uninstall 卸载 Ollama
func (m *OllamaManager) Uninstall() error {
	// 1. 检查是否已安装
	if installed, _, _ := m.Check(); !installed {
		return fmt.Errorf("%s", lang.T("msg.ollama.not_installed"))
	}

	// 2. 检测安装方式
	installMethod := m.detectInstallMethod()

	// 3. 询问确认
	fmt.Println(lang.T("msg.ollama.mac.uninstall_prompt"))
	confirmed, err := m.askForConfirmation(lang.T("msg.ollama.mac.uninstall_confirm"))
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println(lang.T("msg.ollama.mac.uninstall_cancelled"))
		return nil
	}

	// 4. 根据安装方式卸载
	fmt.Println(lang.T("msg.ollama.uninstalling"))

	switch installMethod {
	case "homebrew":
		return m.uninstallViaHomebrew()
	case "script":
		return m.uninstallViaScript()
	default:
		return m.uninstallGeneric()
	}
}

// detectInstallMethod 检测安装方式
func (m *OllamaManager) detectInstallMethod() string {
	// 检查是否通过 Homebrew 安装
	if m.hasHomebrew() {
		cmd := exec.Command("brew", "list", "--formula", "ollama")
		if err := cmd.Run(); err == nil {
			return "homebrew"
		}
		cmd = exec.Command("brew", "list", "--cask", "ollama")
		if err := cmd.Run(); err == nil {
			return "homebrew"
		}
	}

	// 检查是否通过脚本安装（检查常见路径）
	if _, err := exec.LookPath("ollama"); err == nil {
		// 在 PATH 中，可能是脚本安装
		return "script"
	}

	// 检查常见安装路径
	commonPaths := []string{
		"/usr/local/bin/ollama",
		"/opt/homebrew/bin/ollama",
		filepath.Join(os.Getenv("HOME"), ".local/bin/ollama"),
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return "script"
		}
	}

	return "unknown"
}

// uninstallViaHomebrew 通过 Homebrew 卸载
func (m *OllamaManager) uninstallViaHomebrew() error {
	// 先尝试 formula
	cmd := exec.Command("brew", "uninstall", "ollama")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err == nil {
		fmt.Println(lang.T("msg.ollama.uninstall_success"))
		return nil
	}

	// 如果 formula 失败，尝试 cask
	cmd = exec.Command("brew", "uninstall", "--cask", "ollama")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.ollama.uninstall_failed"), err)
	}

	fmt.Println(lang.T("msg.ollama.uninstall_success"))
	return nil
}

// uninstallViaScript 卸载通过脚本安装的 Ollama
func (m *OllamaManager) uninstallViaScript() error {
	// 官方脚本安装的 Ollama 通常需要手动卸载
	// 检查是否有卸载脚本
	uninstallScript := filepath.Join(os.Getenv("HOME"), ".ollama", "uninstall.sh")
	if _, err := os.Stat(uninstallScript); err == nil {
		cmd := exec.Command("sh", uninstallScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %w", lang.T("msg.ollama.uninstall_failed"), err)
		}
		fmt.Println(lang.T("msg.ollama.uninstall_success"))
		return nil
	}

	// 如果没有卸载脚本，手动清理
	fmt.Println(lang.T("msg.ollama.mac.manual_uninstall_hint"))
	fmt.Println("  请手动执行以下步骤：")
	fmt.Println("  1. 删除 ~/.ollama 目录")
	fmt.Println("  2. 从 PATH 中移除 ollama 命令（通常在 ~/.zshrc 或 ~/.bash_profile）")
	fmt.Println("  3. 删除 /usr/local/bin/ollama 或 ~/.local/bin/ollama（如果存在）")

	return fmt.Errorf("%s", lang.T("msg.ollama.uninstall_verify_failed"))
}

// uninstallGeneric 通用卸载方法
func (m *OllamaManager) uninstallGeneric() error {
	fmt.Println(lang.T("msg.ollama.mac.unknown_install_method"))
	fmt.Println("请手动卸载 Ollama：")
	fmt.Println("  - Homebrew: brew uninstall ollama")
	fmt.Println("  - 脚本安装: 删除 ~/.ollama 目录和相关 PATH 配置")

	return fmt.Errorf("%s", lang.T("msg.ollama.uninstall_failed"))
}
