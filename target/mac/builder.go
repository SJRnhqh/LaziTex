// target/mac/builder.go

package mac

import (
	"os"
	"os/exec"
)

// isSkimInstalled 检查 macOS 上是否安装了 Skim.app
func isSkimInstalled() bool {
	// 方式 1：检查标准应用路径 (最快)
	if _, err := os.Stat("/Applications/Skim.app"); err == nil {
		return true
	}
	// 方式 2：使用 mdfind 查找 Bundle ID (更准，能找到安装在非标准目录的 Skim)
	cmd := exec.Command("mdfind", "kMDItemCFBundleIdentifier == 'net.sourceforge.skim-app.skim'")
	output, err := cmd.Output()
	return err == nil && len(output) > 0
}

// PreviewPDF 在 macOS 上打开 PDF 预览
// 目前使用系统默认程序（open 命令），后续可扩展为优先调用 Skim
func PreviewPDF(pdfPath string) error {
	if isSkimInstalled() {
		return exec.Command("open", "-a", "Skim.app", pdfPath).Run()
	}

	// 如果 Skim 未安装，则使用系统默认程序打开
	return exec.Command("open", pdfPath).Run()
}
