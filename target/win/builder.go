// target/win/builder.go

package win

import (
	"os"
	"os/exec"
	"path/filepath"
)

// findSumatraPDF 尝试在常见路径和 PATH 中查找 SumatraPDF
func findSumatraPDF() string {
	// 1. 首先尝试在系统 PATH 中查找
	if path, err := exec.LookPath("SumatraPDF.exe"); err == nil {
		return path
	}

	// 2. 检查 Windows 常见的安装目录
	// 考虑 64 位和 32 位程序目录，以及用户本地目录
	programFiles := os.Getenv("ProgramFiles")
	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	localAppData := os.Getenv("LocalAppData")

	commonPaths := []string{
		filepath.Join(programFiles, "SumatraPDF", "SumatraPDF.exe"),
		filepath.Join(programFilesX86, "SumatraPDF", "SumatraPDF.exe"),
		filepath.Join(localAppData, "SumatraPDF", "SumatraPDF.exe"),
	}

	for _, p := range commonPaths {
		if p == "" {
			continue
		}
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

// PreviewPDF 在 Windows 上打开 PDF 预览
func PreviewPDF(pdfPath string) error {
	sumatraPath := findSumatraPDF()

	if sumatraPath != "" {
		// 使用 SumatraPDF 打开
		// -reuse-instance: 如果已经打开了，就在同一个窗口刷新，不新开窗口
		// 使用 Start() 而不是 Run()，这样预览器打开后不会阻塞命令行
		cmd := exec.Command(sumatraPath, "-reuse-instance", pdfPath)
		return cmd.Start()
	}

	// 如果没找到 SumatraPDF，回退到系统默认关联程序
	return exec.Command("cmd", "/c", "start", "", pdfPath).Run()
}
