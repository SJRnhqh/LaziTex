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

// openPDF 根据平台分发预览任务
func openPDF(pdfPath string) {
	var err error
	switch runtime.GOOS {
	case "darwin":
		err = mac.PreviewPDF(pdfPath) // macOS 预览 PDF
	case "windows":
		err = win.PreviewPDF(pdfPath) // Windows 预览 PDF
	case "linux":
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
	default:
		// 如果不支持，就打印个提示，不强求
		fmt.Printf(lang.T("msg.unsupported_os")+"\n", runtime.GOOS)
		return
	}

	if err != nil {
		fmt.Printf(lang.T("msg.preview_failed")+"\n", err)
	}
}

// 构建 LaTex 文档
func BuildLaTex(filePath string, outputPath string, show bool) {
	opts := core.BuildOptions{
		InputPath:  filePath,
		OutputPath: outputPath,
		Show:       show,
	}

	pdfPath, err := core.Build(opts)
	if err != nil {
		fmt.Printf(lang.T("msg.build_failed")+": %v\n", err)
		return
	}

	if show && pdfPath != "" {
		openPDF(pdfPath)
	}
}
