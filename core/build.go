// core/build.go

package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BuildOptions 编译选项，方便后续扩展
type BuildOptions struct {
	InputPath string // 输入的 .tex 文件路径
}

// Build 执行 LaTeX 编译
func Build(opts BuildOptions) error {
	// 1. 获取输入文件的绝对路径
	absPath, err := filepath.Abs(opts.InputPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// 新增：检查是否是 .tex 文件
	if strings.ToLower(filepath.Ext(absPath)) != ".tex" {
		return fmt.Errorf(T("msg.err_invalid_ext"), opts.InputPath)
	}

	// 2. 检查文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf(T("msg.err_file_not_found"), absPath)
	}

	// 3. 提取工作目录和文件名
	workDir := filepath.Dir(absPath)
	fileName := filepath.Base(absPath)

	// 4. 构造命令 (暂时硬编码使用 xelatex，它是中文最稳的选择)
	// -interaction=nonstopmode 确保遇到错误不卡住
	cmd := exec.Command("xelatex", "-interaction=nonstopmode", fileName)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout // 将编译进度直接打印到控制台
	cmd.Stderr = os.Stderr

	fmt.Printf(T("msg.building_doc")+"\n", fileName)
	fmt.Printf(T("msg.working_dir")+"\n", workDir)

	// 5. 执行编译
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println(T("msg.build_success"))
	return nil
}
