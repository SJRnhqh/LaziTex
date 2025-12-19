// core/errors/passes.go
// 编译问题处理：多次编译需求检测与执行

package errors

import (
	// 外部包
	"os"
	"path/filepath"
	"strings"
)

// PassInfo 编译次数信息
type PassInfo struct {
	NeedsMultiplePasses bool // 是否需要多次编译
	MaxPasses           int  // 最大编译次数
	NeedsBibtex         bool // 是否需要运行 bibtex
	NeedsBiber          bool // 是否需要运行 biber
	NeedsMakeindex      bool // 是否需要运行 makeindex
	NeedsMakeglossaries bool // 是否需要运行 makeglossaries
}

// ensureMaxPasses 确保 MaxPasses 至少为指定值
func (info *PassInfo) ensureMaxPasses(minPasses int) {
	if info.MaxPasses < minPasses {
		info.MaxPasses = minPasses
	}
}

// markNeedsMultiplePasses 标记需要多次编译并设置最小编译次数
func (info *PassInfo) markNeedsMultiplePasses(minPasses int) {
	info.NeedsMultiplePasses = true
	info.ensureMaxPasses(minPasses)
}

// HandleMultiplePasses 检测并处理多次编译需求
// 返回是否需要多次编译，以及编译信息
func HandleMultiplePasses(workDir, jobName string, logOutput string) (bool, PassInfo, error) {
	// 检测编译需求
	passInfo := detectCompilePasses(workDir, jobName, logOutput)

	if !passInfo.NeedsMultiplePasses {
		return false, passInfo, nil
	}

	// 如果需要多次编译，返回信息供调用者处理
	return true, passInfo, nil
}

// detectCompilePasses 检测需要多少次编译（私有方法）
func detectCompilePasses(workDir, jobName string, logOutput string) PassInfo {
	info := PassInfo{MaxPasses: 1}

	// 预先读取 aux 文件（避免重复读取）
	auxPath := buildAuxFilePath(workDir, jobName, ".aux")
	auxContent, _ := readAuxFile(auxPath)

	// 阶段1：基础检测（目录、交叉引用、PDF书签）
	info = detectBasicPasses(workDir, jobName, info, auxContent)

	// 阶段2：参考文献检测
	info = detectBibliographyPasses(workDir, jobName, info, auxContent)

	// 阶段3：高级功能检测（索引、术语表）
	info = detectAdvancedPasses(workDir, jobName, info)

	// 备用：从日志中检测提示
	if logOutput != "" {
		info = detectFromLog(logOutput, info)
	}

	return info
}

// detectBasicPasses 检测基础多次编译需求（目录、交叉引用、PDF书签）
func detectBasicPasses(workDir, jobName string, info PassInfo, auxContent string) PassInfo {
	// 1. 检查目录 (.toc)
	if fileExistsAndNotEmpty(buildAuxFilePath(workDir, jobName, ".toc")) {
		info.markNeedsMultiplePasses(2)
	}

	// 2. 检查交叉引用 (.aux 中的标记)
	if auxContent != "" {
		if strings.Contains(auxContent, "\\newlabel") ||
			strings.Contains(auxContent, "\\@writefile{toc}") {
			info.markNeedsMultiplePasses(2)
		}
	}

	// 3. 检查 PDF 书签 (.out)
	if fileExistsAndNotEmpty(buildAuxFilePath(workDir, jobName, ".out")) {
		info.markNeedsMultiplePasses(2)
	}

	return info
}

// detectBibliographyPasses 检测参考文献编译需求
func detectBibliographyPasses(workDir, jobName string, info PassInfo, auxContent string) PassInfo {
	// 检查是否有引用
	if auxContent == "" || !strings.Contains(auxContent, "\\citation{") {
		return info
	}

	info.NeedsMultiplePasses = true

	// 检查是否有 .bib 文件
	bibPath := buildAuxFilePath(workDir, jobName, ".bib")
	if !fileExistsAndNotEmpty(bibPath) {
		// 使用内嵌的 thebibliography，只需要 2 次编译
		info.ensureMaxPasses(2)
		return info
	}

	// 有 .bib 文件，需要运行 bibtex/biber
	bblPath := buildAuxFilePath(workDir, jobName, ".bbl")
	if fileExistsAndNotEmpty(bblPath) {
		// .bbl 已存在，检查是 bibtex 还是 biber
		if bblContent, err := readAuxFile(bblPath); err == nil {
			if strings.Contains(bblContent, "\\begin{thebibliography}") {
				info.NeedsBibtex = true
			} else {
				info.NeedsBiber = true
			}
		}
	} else {
		// .bbl 不存在，默认尝试 bibtex
		info.NeedsBibtex = true
	}

	// 有 .bib 文件需要 4 次编译
	info.ensureMaxPasses(4)
	return info
}

// detectAdvancedPasses 检测高级功能编译需求（索引、术语表）
func detectAdvancedPasses(workDir, jobName string, info PassInfo) PassInfo {
	// 检查索引 (.idx)
	if fileExistsAndNotEmpty(buildAuxFilePath(workDir, jobName, ".idx")) {
		info.NeedsMultiplePasses = true
		info.NeedsMakeindex = true
		info.ensureMaxPasses(3)
	}

	// 检查术语表 (.glo)
	if fileExistsAndNotEmpty(buildAuxFilePath(workDir, jobName, ".glo")) {
		info.NeedsMultiplePasses = true
		info.NeedsMakeglossaries = true
		info.ensureMaxPasses(3)
	}

	return info
}

// detectFromLog 从日志中检测提示（备用方法）
func detectFromLog(logOutput string, info PassInfo) PassInfo {
	if strings.Contains(logOutput, "Rerun to get") ||
		strings.Contains(logOutput, "Rerun LaTeX") ||
		strings.Contains(logOutput, "undefined reference") {
		info.markNeedsMultiplePasses(2)
	}
	return info
}

// ==================== 辅助函数 ====================

// fileExistsAndNotEmpty 检查文件是否存在且非空
func fileExistsAndNotEmpty(filePath string) bool {
	stat, err := os.Stat(filePath)
	return err == nil && stat.Size() > 0
}

// buildAuxFilePath 构建辅助文件路径
func buildAuxFilePath(workDir, jobName, ext string) string {
	return filepath.Join(workDir, jobName+ext)
}

// readAuxFile 读取辅助文件内容
func readAuxFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
