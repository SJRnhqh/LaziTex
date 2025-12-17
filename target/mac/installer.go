// target/mac/installer.go

package mac

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/SJRnhqh/lazitex/core"
)

// Installer macOS 平台安装器
type Installer struct{}

// NewInstaller 创建 macOS 安装器
func NewInstaller() *Installer {
	return &Installer{}
}

// askForConfirmation 询问用户确认 (y/n)，其他键重新询问
func (i *Installer) askForConfirmation(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.TrimSpace(strings.ToLower(response))

		switch response {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			// 其他键，重新询问
			fmt.Println(core.T("msg.mac.invalid_input"))
		}
	}
}

// Install 安装或更新 LaTeX 环境
// 如果已安装则更新，如果未安装则安装
func (i *Installer) Install() error {
	// 1. 检查是否安装了 Homebrew
	if _, err := exec.LookPath("brew"); err != nil {
		return errors.New(core.T("msg.mac.brew_not_installed"))
	}

	// 2. 检查是否已经安装了 LaTeX
	if i.isLaTeXInstalled() {
		// 已安装，询问是否更新
		fmt.Println(core.T("msg.mac.latex_already_installed"))

		// 列出要更新的包
		packages, err := i.listUpdatablePackages()
		if err != nil {
			// 如果无法列出包，仍然询问是否更新
			fmt.Println(core.T("msg.mac.cannot_list_packages"))
		} else {
			fmt.Println(core.T("msg.mac.updatable_packages"))
			fmt.Println(packages)
		}

		// 询问是否更新
		confirmed, err := i.askForConfirmation(core.T("msg.mac.confirm_update"))
		if err != nil {
			return err
		}

		if !confirmed {
			fmt.Println(core.T("msg.mac.update_cancelled"))
			return nil
		}

		// 执行更新
		if err := i.updateTlmgr(); err != nil {
			return fmt.Errorf("%s: %w", core.T("msg.mac.update_failed"), err)
		}

		fmt.Println(core.T("msg.mac.update_success"))
		return nil
	}

	// 3. 未安装，询问是否安装
	fmt.Println(core.T("msg.mac.install_prompt"))
	fmt.Println(core.T("msg.mac.install_info"))

	// 询问是否安装
	confirmed, err := i.askForConfirmation(core.T("msg.mac.confirm_install"))
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println(core.T("msg.mac.install_cancelled"))
		return nil
	}

	// 4. 执行安装
	fmt.Println(core.T("msg.mac.installing_basictex"))
	fmt.Println(core.T("msg.mac.install_warning"))

	// 使用 Homebrew 安装 BasicTeX
	cmd := exec.Command("brew", "install", "--cask", "basictex")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", core.T("msg.mac.install_failed"), err)
	}

	// 5. 安装后自动更新 tlmgr
	fmt.Println()
	fmt.Println(core.T("msg.mac.updating_tlmgr"))
	if err := i.updateTlmgr(); err != nil {
		// 更新失败不影响安装成功，只显示警告
		fmt.Printf(core.T("msg.mac.update_warning")+": %v\n", err)
	}

	// 6. 显示成功信息和后续步骤
	fmt.Println()
	fmt.Println(core.T("msg.mac.install_success"))
	fmt.Println(core.T("msg.mac.next_steps"))

	return nil
}

// listUpdatablePackages 列出可更新的包
func (i *Installer) listUpdatablePackages() (string, error) {
	// 查找 tlmgr 路径
	tlmgrPath, err := exec.LookPath("tlmgr")
	if err != nil {
		tlmgrPath = "/Library/TeX/texbin/tlmgr"
		if _, err := os.Stat(tlmgrPath); err != nil {
			return "", fmt.Errorf("tlmgr not found")
		}
	}

	// 获取可更新的包列表
	cmd := exec.Command(tlmgrPath, "update", "--list")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// isLaTeXInstalled 检查是否已安装 LaTeX
func (i *Installer) isLaTeXInstalled() bool {
	// 检查 pdflatex 是否存在
	if _, err := exec.LookPath("pdflatex"); err == nil {
		return true
	}
	// 检查 tlmgr 是否存在
	if _, err := exec.LookPath("tlmgr"); err == nil {
		return true
	}
	// 检查标准安装路径
	if _, err := os.Stat("/Library/TeX/texbin"); err == nil {
		return true
	}
	return false
}

// updateTlmgr 更新 tlmgr 和所有包
func (i *Installer) updateTlmgr() error {
	// 查找 tlmgr 路径
	tlmgrPath, err := exec.LookPath("tlmgr")
	if err != nil {
		// 如果找不到 tlmgr，尝试标准路径
		tlmgrPath = "/Library/TeX/texbin/tlmgr"
		if _, err := os.Stat(tlmgrPath); err != nil {
			return fmt.Errorf("tlmgr not found")
		}
	}

	// 更新 tlmgr 自身
	fmt.Println(core.T("msg.mac.updating_tlmgr_self"))
	cmd := exec.Command("sudo", tlmgrPath, "update", "--self")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update tlmgr: %w", err)
	}

	// 更新所有包（可选，可能需要较长时间）
	fmt.Println(core.T("msg.mac.updating_packages"))
	cmd = exec.Command("sudo", tlmgrPath, "update", "--all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 这里不检查错误，因为更新所有包可能需要很长时间，用户可能中断
	_ = cmd.Run()

	return nil
}

// detectInstallationMethod 检测 LaTeX 的安装方式
func (i *Installer) detectInstallationMethod() string {
	// 检查是否通过 Homebrew 安装
	if _, err := exec.LookPath("brew"); err == nil {
		// 检查是否是 Homebrew cask 安装的 basictex
		cmd := exec.Command("brew", "list", "--cask", "basictex")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			return "homebrew_basictex"
		}
		// 检查是否是 Homebrew cask 安装的 mactex
		cmd = exec.Command("brew", "list", "--cask", "mactex")
		output, err = cmd.Output()
		if err == nil && len(output) > 0 {
			return "homebrew_mactex"
		}
	}

	// 检查是否是 MacTeX 官方安装（标准路径）
	if _, err := os.Stat("/Library/TeX/texbin"); err == nil {
		// 检查是否有 GUI 应用（MacTeX 完整版）
		if _, err := os.Stat("/Applications/TeX/TeXShop.app"); err == nil {
			return "mactex_full"
		}
		return "mactex_basic"
	}

	// 检查是否是 MacPorts 安装
	if _, err := exec.LookPath("port"); err == nil {
		cmd := exec.Command("port", "installed", "texlive")
		if err := cmd.Run(); err == nil {
			return "macports"
		}
	}

	// 默认返回未知
	return "unknown"
}

// Uninstall 卸载 LaTeX 环境
func (i *Installer) Uninstall() error {
	// 1. 检查是否已安装 LaTeX
	if !i.isLaTeXInstalled() {
		// 未安装，提示用户
		fmt.Println(core.T("msg.mac.latex_not_installed"))
		fmt.Println()
		fmt.Println(core.T("msg.mac.uninstall_install_hint"))
		return nil
	}

	// 2. 已安装，询问是否卸载
	fmt.Println(core.T("msg.mac.uninstall_prompt"))
	fmt.Println(core.T("msg.mac.uninstall_info"))
	fmt.Println()

	// 询问是否卸载
	confirmed, err := i.askForConfirmation(core.T("msg.mac.confirm_uninstall"))
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println(core.T("msg.mac.uninstall_cancelled"))
		return nil
	}

	// 3. 检测安装方式并执行卸载
	installMethod := i.detectInstallationMethod()
	fmt.Println(core.T("msg.mac.uninstalling"))

	switch installMethod {
	case "homebrew_basictex":
		return i.uninstallHomebrewBasicTeX()
	case "homebrew_mactex":
		return i.uninstallHomebrewMacTeX()
	case "mactex_full", "mactex_basic":
		return i.uninstallMacTeX()
	case "macports":
		return i.uninstallMacPorts()
	default:
		// 未知安装方式，尝试通用卸载
		return i.uninstallGeneric()
	}
}

// uninstallHomebrewBasicTeX 卸载 Homebrew 安装的 BasicTeX
func (i *Installer) uninstallHomebrewBasicTeX() error {
	cmd := exec.Command("brew", "uninstall", "--cask", "basictex")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", core.T("msg.mac.uninstall_failed"), err)
	}

	fmt.Println(core.T("msg.mac.uninstall_success"))
	return nil
}

// uninstallHomebrewMacTeX 卸载 Homebrew 安装的 MacTeX
func (i *Installer) uninstallHomebrewMacTeX() error {
	cmd := exec.Command("brew", "uninstall", "--cask", "mactex")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", core.T("msg.mac.uninstall_failed"), err)
	}

	fmt.Println(core.T("msg.mac.uninstall_success"))
	return nil
}

// uninstallMacTeX 卸载 MacTeX 官方安装
func (i *Installer) uninstallMacTeX() error {
	// MacTeX 官方安装需要手动删除
	// 提示用户使用官方卸载工具或手动删除
	fmt.Println("检测到 MacTeX 官方安装")
	fmt.Println("请使用以下方法卸载：")
	fmt.Println("  1. 打开 /Library/TeX/Distributions/.DefaultTeX/Contents/Library/texlive/bin/universal-darwin/")
	fmt.Println("  2. 运行卸载脚本，或手动删除以下目录：")
	fmt.Println("     - /Library/TeX/")
	fmt.Println("     - /usr/local/texlive/")
	fmt.Println("     - ~/Library/TeX/")
	fmt.Println("     - /Applications/TeX/ (如果存在)")
	fmt.Println()
	fmt.Println("或者使用 Homebrew 重新安装后，可以使用 'lazitex -u' 一键卸载")

	return fmt.Errorf("MacTeX 官方安装需要手动卸载")
}

// uninstallMacPorts 卸载 MacPorts 安装的 TeX Live
func (i *Installer) uninstallMacPorts() error {
	cmd := exec.Command("sudo", "port", "uninstall", "texlive")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", core.T("msg.mac.uninstall_failed"), err)
	}

	fmt.Println(core.T("msg.mac.uninstall_success"))
	return nil
}

// uninstallGeneric 通用卸载方法（当无法确定安装方式时）
func (i *Installer) uninstallGeneric() error {
	fmt.Println("警告: 无法确定 LaTeX 的安装方式")
	fmt.Println("请手动卸载，或使用以下命令：")
	fmt.Println("  - Homebrew: brew uninstall --cask basictex 或 brew uninstall --cask mactex")
	fmt.Println("  - MacPorts: sudo port uninstall texlive")
	fmt.Println("  - MacTeX: 手动删除 /Library/TeX/ 目录")

	return fmt.Errorf("无法自动卸载，请手动卸载")
}
