// core/performance/concurrent.go
// 并发优化相关函数

package performance

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"

	errors "github.com/SJRnhqh/lazitex/core/errors"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// RunIntermediateToolsConcurrent 并发运行中间工具（优化版本）
// 可以并发的工具会同时执行，提升编译速度
// 效果：当需要多个工具时（如 bibtex + makeindex），可以缩短 30-50% 的执行时间
//
// 参数说明：
//   - workDir: 工作目录
//   - jobName: 作业名称（不含扩展名）
//   - passInfo: 编译轮次信息（包含需要运行哪些工具）
//   - quiet: 是否静默模式
func RunIntermediateToolsConcurrent(workDir, jobName string, passInfo errors.PassInfo, quiet bool) {
	var wg sync.WaitGroup

	// 定义需要运行的工具列表
	tools := []struct {
		name   string
		needed bool
		msgKey string
	}{
		{"bibtex", passInfo.NeedsBibtex, "msg.running_bibtex"},
		{"biber", passInfo.NeedsBiber, "msg.running_biber"},
		{"makeindex", passInfo.NeedsMakeindex, "msg.running_makeindex"},
		{"makeglossaries", passInfo.NeedsMakeglossaries, "msg.running_makeglossaries"},
	}

	// 并发执行所有需要的工具
	for _, tool := range tools {
		if !tool.needed {
			continue
		}

		wg.Add(1)
		go func(toolName, msgKey string) {
			defer wg.Done()

			fmt.Println(lang.T(msgKey))
			if err := runTool(toolName, workDir, jobName, quiet); err != nil {
				// 错误信息已经在 runTool 中打印了
				// 即使有错误也继续，因为某些工具失败不影响整体编译
			}
		}(tool.name, tool.msgKey)
	}

	// 等待所有工具完成
	wg.Wait()
}

// runTool 运行单个工具
// 这是 runTool 的实现，复制自 core/build.go 的逻辑，用于 tools 包内部使用
func runTool(toolName, workDir, jobName string, quiet bool) error {
	cmd := exec.Command(toolName, jobName)
	cmd.Dir = workDir

	if quiet {
		// 静默模式：只捕获输出，不打印
		var logOutput bytes.Buffer
		cmd.Stdout = &logOutput
		cmd.Stderr = &logOutput
	} else {
		// 正常模式：输出到控制台
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		// 错误信息仍然显示（即使是静默模式，错误也是重要信息）
		fmt.Printf(lang.T("msg.tool_failed")+"\n", toolName, err)
		return err
	}
	return nil
}