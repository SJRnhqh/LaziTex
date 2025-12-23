// target/win/ollama.go

package win

import (
	// 外部包
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	// 内部包
	lang "github.com/SJRnhqh/lazitex/lang"
)

// OllamaManager Windows Ollama 管理器
type OllamaManager struct{}

// NewOllamaManager 创建 Windows Ollama 管理器
func NewOllamaManager() *OllamaManager {
	return &OllamaManager{}
}

// hasWinget 检查 winget 是否可用
func (m *OllamaManager) hasWinget() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	cmd := exec.Command("winget", "--version")
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

// installWinget 安装 winget
func (m *OllamaManager) installWinget() error {
	fmt.Println(lang.T("msg.ollama.installing_winget"))

	// 方法 1: 尝试通过 Microsoft Store 安装 App Installer
	// 使用 start ms-windows-store://pdp/?ProductId=9NBLGGH4NNS1
	storeURL := "ms-windows-store://pdp/?ProductId=9NBLGGH4NNS1"

	// 尝试打开 Microsoft Store
	cmd := exec.Command("cmd", "/c", "start", storeURL)
	if err := cmd.Run(); err == nil {
		fmt.Println(lang.T("msg.ollama.winget_store_opened"))
		fmt.Println(lang.T("msg.ollama.winget_store_instruction"))

		// 等待用户安装
		fmt.Print(lang.T("msg.ollama.winget_wait_prompt"))
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')

		// 验证安装
		time.Sleep(2 * time.Second)
		if m.hasWinget() {
			fmt.Println(lang.T("msg.ollama.winget_install_success"))
			return nil
		}
	}

	// 方法 2: 尝试下载安装包（如果 Store 方法失败）
	return m.installWingetViaDownload()
}

// installWingetViaDownload 通过下载安装包安装 winget
func (m *OllamaManager) installWingetViaDownload() error {
	fmt.Println(lang.T("msg.ollama.downloading_winget"))

	// 创建临时目录
	tempDir := os.TempDir()
	installerPath := filepath.Join(tempDir, "AppInstaller.msixbundle")

	// 下载 URL（GitHub Releases）
	downloadURL := "https://aka.ms/getwinget"

	// 使用 PowerShell 下载
	downloadCmd := fmt.Sprintf(
		"Invoke-WebRequest -Uri '%s' -OutFile '%s' -UseBasicParsing",
		downloadURL, installerPath,
	)

	cmd := exec.Command("powershell", "-Command", downloadCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.ollama.winget_download_failed"), err)
	}

	// 安装
	fmt.Println(lang.T("msg.ollama.installing_winget_package"))
	installCmd := fmt.Sprintf("Add-AppxPackage -Path '%s'", installerPath)
	cmd = exec.Command("powershell", "-Command", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// 清理下载的文件
		os.Remove(installerPath)
		return fmt.Errorf("%s: %w", lang.T("msg.ollama.winget_install_package_failed"), err)
	}

	// 清理下载的文件
	os.Remove(installerPath)

	// 验证安装
	time.Sleep(2 * time.Second)
	if m.hasWinget() {
		fmt.Println(lang.T("msg.ollama.winget_install_success"))
		return nil
	}

	return fmt.Errorf("%s", lang.T("msg.ollama.winget_install_verify_failed"))
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
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Ollama", "ollama.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Ollama", "ollama.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Ollama", "ollama.exe"),
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

	// 2. 检查 winget 是否可用
	if !m.hasWinget() {
		fmt.Println(lang.T("msg.ollama.winget_not_available"))
		fmt.Println(lang.T("msg.ollama.winget_install_prompt"))

		confirmed, err := m.askForConfirmation(lang.T("msg.ollama.winget_install_confirm"))
		if err != nil {
			return err
		}

		if !confirmed {
			fmt.Println(lang.T("msg.ollama.winget_install_cancelled"))
			return fmt.Errorf("%s", lang.T("msg.ollama.winget_required"))
		}

		// 尝试安装 winget
		if err := m.installWinget(); err != nil {
			return fmt.Errorf("%s: %w", lang.T("msg.ollama.winget_install_error"), err)
		}

		// 安装后再次检查
		if !m.hasWinget() {
			return fmt.Errorf("%s", lang.T("msg.ollama.winget_install_verify_failed"))
		}
	}

	// 3. 通过 winget 安装（注意：winget 安装的是完整版，包含 GUI）
	// Ollama 官方在 Windows 上只提供包含 GUI 的安装包
	// 但 CLI 工具（ollama.exe）在安装后是可用的，GUI 只是额外的界面
	fmt.Println(lang.T("msg.ollama.installing"))
	fmt.Println(lang.T("msg.ollama.install_note_gui"))

	cmd := exec.Command(
		"winget", "install", "Ollama.Ollama",
		"--silent",
		"--accept-package-agreements",
		"--accept-source-agreements",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.ollama.install_failed"), err)
	}

	// 4. 等待安装完成并验证（增加重试机制）
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

	// 2. 检查 winget 是否可用
	if !m.hasWinget() {
		return fmt.Errorf("%s", lang.T("msg.ollama.winget_uninstall_not_available"))
	}

	// 3. 通过 winget 卸载
	fmt.Println(lang.T("msg.ollama.uninstalling"))

	cmd := exec.Command("winget", "uninstall", "Ollama.Ollama", "--silent")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.ollama.uninstall_failed"), err)
	}

	// 4. 验证卸载（增加重试机制）
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
