// core/build.go

package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BuildOptions 编译选项，方便后续扩展
type BuildOptions struct {
	InputPath  string // 输入的 .tex 文件路径
	OutputPath string // 输出的 PDF 文件路径
	Preview    bool   // 是否预览
}

// Build 执行 LaTeX 编译
func Build(opts BuildOptions) (string, error) {
	// 1. 获取输入文件的绝对路径
	absPath, err := filepath.Abs(opts.InputPath)
	if err != nil {
		return "", fmt.Errorf(T("msg.err_abs_path"), opts.InputPath)
	}

	// 检查是否是 .tex 文件
	if strings.ToLower(filepath.Ext(absPath)) != ".tex" {
		return "", fmt.Errorf(T("msg.err_invalid_ext"), opts.InputPath)
	}

	// 2. 检查文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf(T("msg.err_file_not_found"), absPath)
	}

	// 3. 提取工作目录和文件名
	workDir := filepath.Dir(absPath)
	fileName := filepath.Base(absPath)

	outDir := workDir // 默认输出目录就是源文件所在目录
	// 默认的 JobName (不带后缀的文件名)
	jobName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	if opts.OutputPath != "" {
		outPath, _ := filepath.Abs(opts.OutputPath)
		// 简单的逻辑：如果路径不包含 .pdf，视为目录
		if !strings.HasSuffix(strings.ToLower(outPath), ".pdf") {
			outDir = outPath
			// 确保目录存在
			if err := os.MkdirAll(outDir, 0755); err != nil {
				return "", fmt.Errorf(T("msg.err_mkdir"), outDir)
			}
		} else {
			// 如果是完整 PDF 路径
			outDir = filepath.Dir(outPath)
			// 从自定义路径中提取 JobName
			jobName = strings.TrimSuffix(filepath.Base(outPath), ".pdf")
			if err := os.MkdirAll(outDir, 0755); err != nil {
				return "", fmt.Errorf(T("msg.err_mkdir"), outDir)
			}
		}
	}

	// 4. 构造命令
	// -jobname 指定输出文件的主名称
	// -output-directory 指定输出目录
	args := []string{
		"-interaction=nonstopmode",
		"-jobname=" + jobName,
		"-output-directory=" + outDir,
		fileName,
	}
	cmd := exec.Command("xelatex", args...)
	cmd.Dir = workDir

	fmt.Printf(T("msg.building_doc")+"\n", fileName)
	fmt.Printf(T("msg.working_dir")+"\n", workDir)

	// 5. 捕获编译输出（同时显示给用户和保存用于错误分析）
	var logOutput bytes.Buffer
	multiWriter := io.MultiWriter(os.Stdout, &logOutput)
	cmd.Stdout = multiWriter
	cmd.Stderr = multiWriter

	// 执行编译
	err = cmd.Run()
	if err != nil {
		// 编译失败，尝试检测缺失的包
		logStr := logOutput.String()
		handled, handleErr := handleMissingPackages(logStr)

		if handleErr != nil {
			// 处理包管理时出错，返回原始编译错误
			return "", err
		}

		if handled {
			// 包已安装，重新编译
			fmt.Println(T("msg.package_retry_build"))
			return Build(opts) // 递归调用，但只重试一次（因为 handleMissingPackages 已经处理了）
		}

		// 没有缺失包或用户取消安装，返回原始编译错误
		return "", err
	}

	fmt.Println(T("msg.build_success"))

	// 返回生成的 PDF 完整路径
	return filepath.Join(outDir, jobName+".pdf"), nil
}
