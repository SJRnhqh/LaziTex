// target/mac/installer.go

package mac

import (
	// 外部包
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	// 内部包
	"github.com/SJRnhqh/lazitex/lang"
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
			fmt.Println(lang.T("msg.mac.invalid_input"))
		}
	}
}

// Install 安装或更新 LaTeX 环境
// 如果已安装则更新，如果未安装则安装
func (i *Installer) Install() error {
	// 1. 检查是否安装了 Homebrew
	if _, err := exec.LookPath("brew"); err != nil {
		return errors.New(lang.T("msg.mac.brew_not_installed"))
	}

	// 2. 检查是否已经安装了 LaTeX
	if i.isLaTeXInstalled() {
		// 已安装，询问是否更新
		fmt.Println(lang.T("msg.mac.latex_already_installed"))

		// 列出要更新的包（附带包名，便于后续只更新必要包）
		packagesOutput, pkgNames, err := i.listUpdatablePackages()
		if err != nil {
			// 如果无法列出包，仍然询问是否更新
			fmt.Println(lang.T("msg.mac.cannot_list_packages"))
		} else {
			if len(pkgNames) == 0 {
				fmt.Println(lang.T("msg.mac.all_up_to_date"))
				return nil
			}
			fmt.Println(lang.T("msg.mac.updatable_packages"))
			fmt.Println(packagesOutput)
		}

		// 询问是否更新
		confirmed, err := i.askForConfirmation(lang.T("msg.mac.confirm_update"))
		if err != nil {
			return err
		}

		if !confirmed {
			fmt.Println(lang.T("msg.mac.update_cancelled"))
			return nil
		}

		// 执行更新
		if err := i.updateTlmgr(); err != nil {
			return fmt.Errorf("%s: %w", lang.T("msg.mac.update_failed"), err)
		}

		fmt.Println(lang.T("msg.mac.update_success"))
		return nil
	}

	// 3. 未安装，询问是否安装
	fmt.Println(lang.T("msg.mac.install_prompt"))
	fmt.Println(lang.T("msg.mac.install_info"))

	// 询问是否安装
	confirmed, err := i.askForConfirmation(lang.T("msg.mac.confirm_install"))
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println(lang.T("msg.mac.install_cancelled"))
		return nil
	}

	// 4. 执行安装
	fmt.Println(lang.T("msg.mac.installing_basictex"))
	fmt.Println(lang.T("msg.mac.install_warning"))

	// 使用 Homebrew 安装 BasicTeX
	cmd := exec.Command("brew", "install", "--cask", "basictex")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.mac.install_failed"), err)
	}

	// 5. 安装后自动更新 tlmgr
	fmt.Println()
	fmt.Println(lang.T("msg.mac.updating_tlmgr"))
	if err := i.updateTlmgr(); err != nil {
		// 更新失败不影响安装成功，只显示警告
		fmt.Printf(lang.T("msg.mac.update_warning")+": %v\n", err)
	}

	// 6. 显示成功信息和后续步骤
	fmt.Println()
	fmt.Println(lang.T("msg.mac.install_success"))
	fmt.Println(lang.T("msg.mac.next_steps"))

	return nil
}

// listUpdatablePackages 列出可更新的包
func (i *Installer) listUpdatablePackages() (string, []string, error) {
	tlmgrPath, err := i.findTlmgrPath()
	if err != nil {
		return "", nil, err
	}

	// 获取可更新的包列表
	cmd := exec.Command(tlmgrPath, "update", "--list")
	output, err := cmd.Output()
	if err != nil {
		return "", nil, err
	}

	lines := strings.Split(string(output), "\n")
	var pkgNames []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// tlmgr --list 格式示例: "[ update: collection-basictex ]"
		if strings.HasPrefix(line, "[") && strings.Contains(line, "update:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				// fields[1] 形如 "update:"，fields[2] 是包名或带右括号，做一下清理
				name := strings.Trim(fields[2], "[]")
				name = strings.TrimSuffix(name, "]")
				name = strings.TrimSpace(name)
				if name != "" {
					pkgNames = append(pkgNames, name)
				}
			}
		}
	}

	return string(output), pkgNames, nil
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
	tlmgrPath, err := i.findTlmgrPath()
	if err != nil {
		return err
	}

	// 可选：允许通过环境变量指定镜像，加速国内更新
	// 示例: LAZITEX_TLMGR_REPO="https://mirrors.tuna.tsinghua.edu.cn/CTAN/systems/texlive/tlnet"
	if repo := os.Getenv("LAZITEX_TLMGR_REPO"); repo != "" {
		fmt.Printf("Using TeX Live mirror: %s\n", repo)
		cmdRepo := exec.Command("sudo", tlmgrPath, "option", "repository", repo)
		cmdRepo.Stdout = os.Stdout
		cmdRepo.Stderr = os.Stderr
		_ = cmdRepo.Run()
	}

	// 更新 tlmgr 自身
	fmt.Println(lang.T("msg.mac.updating_tlmgr_self"))
	cmd := exec.Command("sudo", tlmgrPath, "update", "--self")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update tlmgr: %w", err)
	}

	// 获取需要更新的包列表，避免对全部包执行 --all
	packagesOutput, pkgNames, err := i.listUpdatablePackages()
	if err != nil {
		fmt.Println(lang.T("msg.mac.cannot_list_packages"))
		// 如果无法列出，退回全量更新，但会更慢
		fmt.Println(lang.T("msg.mac.updating_packages"))
		cmdAll := exec.Command("sudo", tlmgrPath, "update", "--all", "--no-doc", "--no-src")
		cmdAll.Stdout = os.Stdout
		cmdAll.Stderr = os.Stderr
		_ = cmdAll.Run()
		return nil
	}

	if len(pkgNames) == 0 {
		fmt.Println(lang.T("msg.mac.all_up_to_date"))
		return nil
	}

	fmt.Println(lang.T("msg.mac.updating_packages"))
	fmt.Println(packagesOutput)

	args := append([]string{"update", "--no-doc", "--no-src"}, pkgNames...)
	cmd = exec.Command("sudo", append([]string{tlmgrPath}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	if len(pkgNames) > 0 {
		fmt.Printf(lang.T("msg.mac.updated_packages")+"\n", strings.Join(pkgNames, ", "))
	}

	return nil
}

// findTlmgrPath 查找 tlmgr 的路径
func (i *Installer) findTlmgrPath() (string, error) {
	if path, err := exec.LookPath("tlmgr"); err == nil {
		return path, nil
	}
	// 如果找不到 tlmgr，尝试标准路径
	tlmgrPath := "/Library/TeX/texbin/tlmgr"
	if _, err := os.Stat(tlmgrPath); err == nil {
		return tlmgrPath, nil
	}
	return "", fmt.Errorf("tlmgr not found")
}

// detectInstallationMethod 检测 LaTeX 的安装方式
func (i *Installer) detectInstallationMethod() string {
	// 检查是否通过 Homebrew 安装
	if _, err := exec.LookPath("brew"); err == nil {
		if out, _ := exec.Command("brew", "list", "--cask", "basictex").Output(); len(out) > 0 {
			return "homebrew_basictex"
		}
		if out, _ := exec.Command("brew", "list", "--cask", "mactex").Output(); len(out) > 0 {
			return "homebrew_mactex"
		}
	}

	// 检查是否是 MacTeX 官方安装（标准路径）
	if _, err := os.Stat("/Library/TeX/texbin"); err == nil {
		if _, err := os.Stat("/Applications/TeX/TeXShop.app"); err == nil {
			return "mactex_full"
		}
		return "mactex_basic"
	}

	// 检查是否是 MacPorts 安装（覆盖常见端口）
	if _, err := exec.LookPath("port"); err == nil {
		ports := []string{"texlive-full", "texlive-latex", "texlive-basic", "texlive"}
		for _, p := range ports {
			if err := exec.Command("port", "installed", p).Run(); err == nil {
				return "macports"
			}
		}
	}

	// 默认返回未知
	return "unknown"
}

// cleanupResidual 清理常见残留目录（尽力而为，忽略错误）
func (i *Installer) cleanupResidual() {
	home := os.Getenv("HOME")
	paths := []struct {
		path    string
		useSudo bool
	}{
		{"/usr/local/texlive", true},
		{"/Library/TeX", true},
		{home + "/Library/TeX", false},
		{home + "/Library/texlive", false},
	}

	for _, p := range paths {
		if _, err := os.Stat(p.path); err == nil {
			fmt.Printf(lang.T("msg.mac.removing_residual")+"\n", p.path)
			var cmd *exec.Cmd
			if p.useSudo {
				cmd = exec.Command("sudo", "rm", "-rf", p.path)
			} else {
				cmd = exec.Command("rm", "-rf", p.path)
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Run()
		}
	}
}

// Uninstall 卸载 LaTeX 环境
func (i *Installer) Uninstall() error {
	// 1. 检查是否已安装 LaTeX
	if !i.isLaTeXInstalled() {
		// 未安装，提示用户
		fmt.Println(lang.T("msg.mac.latex_not_installed"))
		fmt.Println()
		fmt.Println(lang.T("msg.mac.uninstall_install_hint"))
		return nil
	}

	// 2. 已安装，询问是否卸载
	fmt.Println(lang.T("msg.mac.uninstall_prompt"))
	fmt.Println(lang.T("msg.mac.uninstall_info"))
	fmt.Println()

	// 询问是否卸载
	confirmed, err := i.askForConfirmation(lang.T("msg.mac.confirm_uninstall"))
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println(lang.T("msg.mac.uninstall_cancelled"))
		return nil
	}

	// 3. 检测安装方式并执行卸载
	installMethod := i.detectInstallationMethod()
	fmt.Println(lang.T("msg.mac.uninstalling"))

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
		return fmt.Errorf("%s: %w", lang.T("msg.mac.uninstall_failed"), err)
	}

	i.cleanupResidual()
	fmt.Println(lang.T("msg.mac.uninstall_success"))
	return nil
}

// uninstallHomebrewMacTeX 卸载 Homebrew 安装的 MacTeX
func (i *Installer) uninstallHomebrewMacTeX() error {
	cmd := exec.Command("brew", "uninstall", "--cask", "mactex")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.mac.uninstall_failed"), err)
	}

	i.cleanupResidual()
	fmt.Println(lang.T("msg.mac.uninstall_success"))
	return nil
}

// uninstallMacTeX 卸载 MacTeX 官方安装
func (i *Installer) uninstallMacTeX() error {
	scriptCandidates := []string{
		"/Library/TeX/Distributions/.DefaultTeX/Contents/Library/texlive/bin/universal-darwin/uninstall-texlive.sh",
		"/Library/TeX/Distributions/.DefaultTeX/Contents/Library/texlive/bin/universal-darwin/uninstall-texlive",
	}

	for _, p := range scriptCandidates {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			fmt.Println(lang.T("msg.mac.uninstall_mactex_running"))
			cmd := exec.Command("sudo", p)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("%s: %w", lang.T("msg.mac.uninstall_failed"), err)
			}
			i.cleanupResidual()
			fmt.Println(lang.T("msg.mac.uninstall_success"))
			return nil
		}
	}

	// 找不到脚本则提示手动卸载
	fmt.Println(lang.T("msg.mac.uninstall_mactex_manual"))
	fmt.Println("  - /Library/TeX/")
	fmt.Println("  - /usr/local/texlive/")
	fmt.Println("  - ~/Library/TeX/")
	fmt.Println("  - /Applications/TeX/ (如果存在)")
	return fmt.Errorf("MacTeX official uninstall script not found")
}

// uninstallMacPorts 卸载 MacPorts 安装的 TeX Live
func (i *Installer) uninstallMacPorts() error {
	ports := []string{"texlive-full", "texlive-latex", "texlive-basic", "texlive"}
	var lastErr error

	for _, p := range ports {
		cmd := exec.Command("sudo", "port", "uninstall", p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			lastErr = err
			continue
		}
		i.cleanupResidual()
		fmt.Println(lang.T("msg.mac.uninstall_success"))
		return nil
	}

	if lastErr != nil {
		return fmt.Errorf("%s: %w", lang.T("msg.mac.uninstall_failed"), lastErr)
	}
	return fmt.Errorf("%s: no texlive-* port found", lang.T("msg.mac.uninstall_failed"))
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
