// core/errors/formatter.go
// 错误信息提取和格式化

package errors

import (
	"regexp"
	"strings"

	lang "github.com/SJRnhqh/lazitex/lang"
)

// extractErrors 从编译日志中提取错误信息
func ExtractErrors(logOutput string) []string {
	var errors []string

	// 匹配 LaTeX 错误格式：! LaTeX Error: ...
	errorPattern := regexp.MustCompile(`! LaTeX Error: (.+?)(?:\n|$)`)
	matches := errorPattern.FindAllStringSubmatch(logOutput, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			errorMsg := strings.TrimSpace(match[1])
			if errorMsg != "" {
				errors = append(errors, errorMsg)
			}
		}
	}

	return errors
}

// FormatError 格式化错误信息为 LaziTex 风格
func FormatError(errorMsg string) string {
	// 提取关键错误信息
	if strings.Contains(errorMsg, "File") && strings.Contains(errorMsg, "not found") {
		// 文件未找到错误
		return lang.T("msg.build_failed") + ": " + errorMsg
	}

	// 通用错误格式
	return "❌ " + lang.T("msg.build_failed") + ": " + errorMsg
}
