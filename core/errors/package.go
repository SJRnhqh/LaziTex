// core/errors/package.go
// 核心业务内部，编译LaTeX文档时检测并处理缺失的包

package errors

import (
	// 外部包
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	// 内部包
	lang "github.com/SJRnhqh/lazitex/lang"
)

// extractMissingPackages 从编译日志中提取缺失的包名
func extractMissingPackages(logOutput string) []string {
	// 正则表达式：匹配 ! LaTeX Error: File 'xxx.sty' not found 或 File `xxx.cls` not found
	// 注意：LaTeX 错误信息可能使用单引号或反引号，末尾可能有句号
	// 使用不区分大小写匹配以兼容不同格式
	pattern := "(?i)! LaTeX Error: File ['`](.+?)\\.(sty|cls)['`] not found\\.?"
	re := regexp.MustCompile(pattern)

	// 查找所有匹配
	matches := re.FindAllStringSubmatch(logOutput, -1)

	// 使用 map 去重
	packages := make(map[string]bool)
	for _, match := range matches {
		if len(match) >= 2 {
			fileName := match[1] // 提取文件名（不含扩展名）
			packages[fileName] = true
		}
	}

	// 转换为切片
	result := make([]string, 0, len(packages))
	for pkg := range packages {
		result = append(result, pkg)
	}

	return result
}

// findPackageManager 查找可用的包管理器
func findPackageManager() (string, string, error) {
	// 优先查找 tlmgr（TeX Live，所有平台通用）
	if path, err := exec.LookPath("tlmgr"); err == nil {
		return path, "tlmgr", nil
	}

	// 回退到 mpm（MiKTeX，主要是 Windows）
	if path, err := exec.LookPath("mpm"); err == nil {
		return path, "mpm", nil
	}

	return "", "", fmt.Errorf("%s", lang.T("msg.package_no_manager"))
}

// findPackageName 使用 tlmgr search 查找文件对应的正确包名
func findPackageName(managerPath string, fileName string) (string, error) {
	// 使用 tlmgr search --global --file 查找包含该文件的包
	cmd := exec.Command(managerPath, "search", "--global", "--file", fileName+".sty")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// 解析输出，格式可能是：
	// package_name:path/to/file.sty
	// 或者多行输出，每行一个匹配
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 跳过明显的非包名行（如 "tlmgr:" 开头的行）
		if strings.HasPrefix(line, "tlmgr:") {
			continue
		}

		// 提取包名（冒号前的部分）
		if idx := strings.Index(line, ":"); idx > 0 {
			pkgName := strings.TrimSpace(line[:idx])
			// 确保包名不是 "tlmgr" 或其他明显错误的值
			if pkgName != "" && pkgName != "tlmgr" && !strings.Contains(pkgName, " ") {
				return pkgName, nil
			}
		}
	}

	// 如果找不到，返回空（将使用原始文件名）
	return "", nil
}

// installPackages 安装指定的包
func installPackages(packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	managerPath, managerType, err := findPackageManager()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	var args []string

	switch managerType {
	case "tlmgr":
		// 先尝试直接使用包名安装
		args = append([]string{"install"}, packages...)
	case "mpm":
		args = []string{}
		for _, pkg := range packages {
			args = append(args, "--install="+pkg)
		}
	default:
		return fmt.Errorf(lang.T("msg.package_unsupported_manager"), managerType)
	}

	// 先尝试不使用 sudo
	cmd = exec.Command(managerPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	// 如果失败，且是 tlmgr
	if err != nil && managerType == "tlmgr" {
		// 先尝试查找正确的包名
		resolvedPackages := make([]string, 0, len(packages))
		for _, pkg := range packages {
			// 尝试查找包含该文件的包
			if resolvedPkg, findErr := findPackageName(managerPath, pkg); findErr == nil && resolvedPkg != "" {
				resolvedPackages = append(resolvedPackages, resolvedPkg)
			} else {
				// 如果查找失败，尝试小写版本
				resolvedPackages = append(resolvedPackages, strings.ToLower(pkg))
			}
		}

		// 如果找到了不同的包名，重新尝试安装（不用 sudo）
		if len(resolvedPackages) > 0 {
			packagesChanged := false
			for i := range packages {
				if !strings.EqualFold(packages[i], resolvedPackages[i]) {
					packagesChanged = true
					break
				}
			}

			if packagesChanged {
				resolvedArgs := append([]string{"install"}, resolvedPackages...)
				cmd = exec.Command(managerPath, resolvedArgs...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
			}
		}

		// 如果还是失败，尝试使用 sudo
		if err != nil {
			fmt.Println(lang.T("msg.package_need_sudo"))
			// 确定要安装的包列表：优先使用 resolvedPackages，否则使用原始 packages
			packagesToInstall := resolvedPackages
			if len(packagesToInstall) == 0 {
				packagesToInstall = packages
			}

			sudoArgs := append([]string{managerPath, "install"}, packagesToInstall...)
			sudoCmd := exec.Command("sudo", sudoArgs...)
			sudoCmd.Stdout = os.Stdout
			sudoCmd.Stderr = os.Stderr
			sudoCmd.Stdin = os.Stdin // 允许输入密码
			err = sudoCmd.Run()
		}
	}

	return err
}

// handleMissingPackages 检测并处理缺失的包
func HandleMissingPackages(logOutput string) (bool, error) {
	// 提取缺失的包
	packages := extractMissingPackages(logOutput)
	if len(packages) == 0 {
		return false, nil // 没有缺失包
	}

	// 显示检测到的缺失包
	fmt.Printf(lang.T("msg.package_missing")+"\n", strings.Join(packages, ", "))

	// 询问用户是否安装
	fmt.Print(lang.T("msg.package_install_prompt"))

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.TrimSpace(strings.ToLower(response))

	// 如果用户确认（Y 或直接回车）
	if response == "" || response == "y" || response == "yes" {
		// 一次性安装所有包（减少密码输入次数）
		fmt.Printf(lang.T("msg.package_installing")+"\n", strings.Join(packages, ", "))
		if err := installPackages(packages); err != nil {
			fmt.Printf(lang.T("msg.package_install_failed")+"\n", err)
			return false, err
		}
		fmt.Println(lang.T("msg.package_install_success"))
		return true, nil
	}

	// 用户取消安装
	return false, nil
}
