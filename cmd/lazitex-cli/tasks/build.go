// cmd/lazitex-cli/tasks/build.go
// 构建PDF相关业务管理

package tasks

import (
	"fmt"
	"runtime"

	core "github.com/SJRnhqh/lazitex/core"
	lang "github.com/SJRnhqh/lazitex/lang"
	mac "github.com/SJRnhqh/lazitex/target/mac"
	win "github.com/SJRnhqh/lazitex/target/win"
)

// showPDF 根据平台分发展示任务
func showPDF(pdfPath string) {
	var err error
	switch runtime.GOOS {
	case "darwin":
		err = mac.ShowPDF(pdfPath) // macOS 展示 PDF
	case "windows":
		err = win.ShowPDF(pdfPath) // Windows 展示 PDF
	case "linux":
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
	default:
		// 如果不支持，就打印个提示，不强求
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	if err != nil {
		fmt.Printf(lang.T("msg.show_failed")+"\n", err)
	}
}

// 构建 LaTeX 文档
func BuildLaTeX(filePath string, outputPath string, show bool, quiet bool) {
	opts := core.BuildOptions{
		InputPath:  filePath,
		OutputPath: outputPath,
		Show:       show,
		Quiet:      quiet,
	}

	pdfPath, err := core.Build(opts)
	if err != nil {
		fmt.Printf(lang.T("msg.build_failed")+": %v\n", err)
		return
	}

	if show && pdfPath != "" {
		showPDF(pdfPath)
	}
}
