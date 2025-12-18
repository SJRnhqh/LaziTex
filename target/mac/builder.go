// target/mac/builder.go

package mac

import "os/exec"

// PreviewPDF 在 macOS 上打开 PDF 预览
// 目前使用系统默认程序（open 命令），后续可扩展为优先调用 Skim
func PreviewPDF(pdfPath string) error {
	// macOS 的 open 命令非常强大：
	// 1. 如果 PDF 没打开，它会用默认程序（通常是预览.app 或浏览器）打开
	// 2. 如果 PDF 已经打开，它会将对应窗口提到最前
	return exec.Command("open", pdfPath).Run()
}
