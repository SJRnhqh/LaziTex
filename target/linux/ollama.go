// target/linux/ollama.go
// Linux Ollama 管理器实现

package linux

import (
	// 外部包
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	// 内部包
	lang "github.com/SJRnhqh/lazitex/lang"
)

// OllamaManager Linux Ollama 管理器
type OllamaManager struct{}

// NewOllamaManager 创建 Linux Ollama 管理器
func NewOllamaManager() *OllamaManager {
	return &OllamaManager{}
}

// detectPackageManager 检测可用的包管理器
func (m *OllamaManager) detectPackageManager() string {
	packageManagers := []struct {
		name     string
		checkCmd []string
	}{
		{"apt", []string{"apt", "--version"}},
		{"yum", []string{"yum", "--version"}},
		{"dnf", []string{"dnf", "--version"}},
		{"pacman", []string{"pacman", "--version"}},
		{"zypper", []string{"zypper", "--version"}},
		{"snap", []string{"snap", "--version"}},
	}

	for _, pm := range packageManagers {
		cmd := exec.Command(pm.checkCmd[0], pm.checkCmd[1:]...)
		if err := cmd.Run(); err == nil {
			return pm.name
		}
	}

	return ""
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
	homeDir := os.Getenv("HOME")
	commonPaths := []string{
		"/usr/local/bin/ollama",
		"/usr/bin/ollama",
		"/opt/ollama/ollama",
		filepath.Join(homeDir, ".local/bin/ollama"),
		filepath.Join(homeDir, ".ollama/bin/ollama"),
		"/snap/bin/ollama",
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

	// 3. 检查是否通过 snap 安装
	if m.hasSnap() {
		cmd := exec.Command("snap", "list", "ollama")
		if err := cmd.Run(); err == nil {
			// 尝试获取版本
			cmd = exec.Command("snap", "info", "ollama")
			output, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(output), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "installed:") {
						parts := strings.Fields(line)
						if len(parts) > 1 {
							return true, parts[1], nil
						}
					}
				}
			}
			return true, "installed", nil
		}
	}

	return false, "", nil // 未安装
}

// hasSnap 检查 snap 是否可用
func (m *OllamaManager) hasSnap() bool {
	cmd := exec.Command("snap", "--version")
	return cmd.Run() == nil
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

	// 2. 检测包管理器
	pm := m.detectPackageManager()

	// 3. 询问用户使用哪种方式安装
	fmt.Println(lang.T("msg.ollama.linux.install_method_prompt"))
	if pm != "" {
		fmt.Printf(lang.T("msg.ollama.linux.install_method_option1")+"\n", pm)
	}
	fmt.Println(lang.T("msg.ollama.linux.install_method_option2"))
	fmt.Print(lang.T("msg.ollama.linux.install_method_choice"))

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// 4. 根据用户选择安装
	if choice == "1" && pm != "" {
		return m.installViaPackageManager(pm)
	}

	// 默认使用官方脚本安装
	return m.installViaScript()
}

// installViaPackageManager 通过包管理器安装
func (m *OllamaManager) installViaPackageManager(pm string) error {
	fmt.Printf(lang.T("msg.ollama.linux.installing_via_pm")+"\n", pm)

	var cmd *exec.Cmd
	switch pm {
	case "apt":
		// 更新包列表
		updateCmd := exec.Command("sudo", "apt", "update")
		updateCmd.Stdout = os.Stdout
		updateCmd.Stderr = os.Stderr
		if err := updateCmd.Run(); err != nil {
			return fmt.Errorf("%s: %w", lang.T("msg.ollama.linux.pm_update_failed"), err)
		}
		cmd = exec.Command("sudo", "apt", "install", "-y", "ollama")
	case "yum":
		cmd = exec.Command("sudo", "yum", "install", "-y", "ollama")
	case "dnf":
		cmd = exec.Command("sudo", "dnf", "install", "-y", "ollama")
	case "pacman":
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "ollama")
	case "zypper":
		cmd = exec.Command("sudo", "zypper", "install", "-y", "ollama")
	case "snap":
		cmd = exec.Command("sudo", "snap", "install", "ollama")
	default:
		return fmt.Errorf("%s: %s", lang.T("msg.ollama.linux.unsupported_pm"), pm)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.ollama.install_failed"), err)
	}

	// 验证安装
	return m.verifyInstall()
}

// installViaScript 通过官方安装脚本安装
func (m *OllamaManager) installViaScript() error {
	fmt.Println(lang.T("msg.ollama.linux.installing_via_script"))
	fmt.Println(lang.T("msg.ollama.linux.install_note_script"))

	// 下载并运行官方安装脚本
	cmd := exec.Command("sh", "-c", "curl -fsSL https://ollama.com/install.sh | sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.ollama.install_failed"), err)
	}

	// 验证安装
	return m.verifyInstall()
}

// verifyInstall 验证安装是否成功
func (m *OllamaManager) verifyInstall() error {
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

// Uninstall 卸载 Ollama
func (m *OllamaManager) Uninstall() error {
	// 1. 检查是否已安装
	if installed, _, _ := m.Check(); !installed {
		return fmt.Errorf("%s", lang.T("msg.ollama.not_installed"))
	}

	// 2. 检测安装方式
	installMethod := m.detectInstallMethod()

	// 3. 询问确认
	fmt.Println(lang.T("msg.ollama.linux.uninstall_prompt"))
	confirmed, err := m.askForConfirmation(lang.T("msg.ollama.linux.uninstall_confirm"))
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println(lang.T("msg.ollama.linux.uninstall_cancelled"))
		return nil
	}

	// 4. 根据安装方式卸载
	fmt.Println(lang.T("msg.ollama.uninstalling"))

	switch installMethod {
	case "snap":
		return m.uninstallViaSnap()
	case "package":
		pm := m.detectPackageManager()
		if pm != "" {
			return m.uninstallViaPackageManager(pm)
		}
		fallthrough
	default:
		return m.uninstallGeneric()
	}
}

// detectInstallMethod 检测安装方式
func (m *OllamaManager) detectInstallMethod() string {
	// 检查是否通过 snap 安装
	if m.hasSnap() {
		cmd := exec.Command("snap", "list", "ollama")
		if err := cmd.Run(); err == nil {
			return "snap"
		}
	}

	// 检查是否通过包管理器安装
	// 尝试检查常见包管理器
	packageManagers := []string{"apt", "yum", "dnf", "pacman", "zypper"}
	for _, pm := range packageManagers {
		cmd := exec.Command("which", pm)
		if err := cmd.Run(); err == nil {
			// 检查 ollama 是否通过此包管理器安装
			var checkCmd *exec.Cmd
			switch pm {
			case "apt":
				checkCmd = exec.Command("dpkg", "-l", "ollama")
			case "yum", "dnf":
				checkCmd = exec.Command("rpm", "-q", "ollama")
			case "pacman":
				checkCmd = exec.Command("pacman", "-Q", "ollama")
			case "zypper":
				checkCmd = exec.Command("zypper", "se", "-i", "ollama")
			}

			if checkCmd != nil && checkCmd.Run() == nil {
				return "package"
			}
		}
	}

	return "script"
}

// uninstallViaSnap 通过 snap 卸载
func (m *OllamaManager) uninstallViaSnap() error {
	cmd := exec.Command("sudo", "snap", "remove", "ollama")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.ollama.uninstall_failed"), err)
	}

	return m.verifyUninstall()
}

// uninstallViaPackageManager 通过包管理器卸载
func (m *OllamaManager) uninstallViaPackageManager(pm string) error {
	var cmd *exec.Cmd
	switch pm {
	case "apt":
		cmd = exec.Command("sudo", "apt", "remove", "-y", "ollama")
	case "yum":
		cmd = exec.Command("sudo", "yum", "remove", "-y", "ollama")
	case "dnf":
		cmd = exec.Command("sudo", "dnf", "remove", "-y", "ollama")
	case "pacman":
		cmd = exec.Command("sudo", "pacman", "-R", "--noconfirm", "ollama")
	case "zypper":
		cmd = exec.Command("sudo", "zypper", "remove", "-y", "ollama")
	default:
		return fmt.Errorf("%s: %s", lang.T("msg.ollama.linux.unsupported_pm"), pm)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.ollama.uninstall_failed"), err)
	}

	return m.verifyUninstall()
}

// uninstallGeneric 通用卸载方法（脚本安装）
func (m *OllamaManager) uninstallGeneric() error {
	// 官方脚本安装的 Ollama 通常需要手动卸载
	// 检查是否有卸载脚本
	homeDir := os.Getenv("HOME")
	uninstallScript := filepath.Join(homeDir, ".ollama", "uninstall.sh")
	if _, err := os.Stat(uninstallScript); err == nil {
		cmd := exec.Command("sh", uninstallScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %w", lang.T("msg.ollama.uninstall_failed"), err)
		}
		return m.verifyUninstall()
	}

	// 如果没有卸载脚本，提供手动卸载提示
	fmt.Println(lang.T("msg.ollama.linux.manual_uninstall_hint"))
	fmt.Println(lang.T("msg.ollama.linux.manual_uninstall_steps"))

	return fmt.Errorf("%s", lang.T("msg.ollama.uninstall_verify_failed"))
}

// verifyUninstall 验证卸载是否成功
func (m *OllamaManager) verifyUninstall() error {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		time.Sleep(2 * time.Second)
		if installed, _, _ := m.Check(); !installed {
			fmt.Println(lang.T("msg.ollama.uninstall_success"))
			return nil
		}
	}

	return fmt.Errorf("%s", lang.T("msg.ollama.uninstall_verify_failed"))
}

// Status 检查 Ollama 服务状态
// 返回: (是否运行, 错误)
func (m *OllamaManager) Status() (bool, error) {
	// 1. 检查是否已安装
	installed, version, err := m.Check()
	if err != nil {
		return false, err
	}

	if !installed {
		fmt.Println(lang.T("msg.ollama.not_installed"))
		return false, nil
	}

	// 2. 检查服务状态
	running := m.checkService()

	// 3. 输出状态信息
	if version != "" && version != "installed" {
		fmt.Printf(lang.T("msg.ollama.status_installed")+"\n", version)
	} else {
		fmt.Println(lang.T("msg.ollama.status_installed_no_version"))
	}

	if running {
		fmt.Println(lang.T("msg.ollama.service_running"))
	} else {
		fmt.Println(lang.T("msg.ollama.service_not_running"))
		fmt.Println(lang.T("msg.ollama.start_service_hint"))
	}

	return running, nil
}