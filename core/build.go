// core/build.go
// 核心业务：编译LaTex编译器相关业务

package core

import (
	// 外部包
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	// 内部包
	errors "github.com/SJRnhqh/lazitex/core/errors"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// BuildOptions 编译选项，方便后续扩展
type BuildOptions struct {
	InputPath  string // 输入的 .tex 文件路径
	OutputPath string // 输出的 PDF 文件路径
	Show       bool   // 编译完是否展示PDF
}

// compileContext 编译上下文，封装编译过程中的所有信息
type compileContext struct {
	workDir   string // 工作目录（源文件所在目录）
	outDir    string // 输出目录
	jobName   string // 作业名称（不含扩展名）
	fileName  string // 源文件名
	passCount int    // 当前编译次数
	firstPass bool   // 是否是第一次编译
}

// compileStrategy 编译策略接口，定义不同场景的处理方式
// 此接口设计支持未来扩展（如AI策略），当前只实现默认策略
type compileStrategy interface {
	// handleError 处理编译错误
	// 返回 shouldContinue: 是否应该继续编译, newErr: 新的错误（如果需要）
	handleError(ctx *compileContext, logStr string, err error) (shouldContinue bool, newErr error)

	// shouldRunIntermediateTools 判断是否应该运行中间工具（bibtex、biber等）
	shouldRunIntermediateTools(ctx *compileContext, passInfo errors.PassInfo) bool

	// shouldContinue 判断是否应该继续编译
	shouldContinue(ctx *compileContext, logStr string, passInfo errors.PassInfo) bool
}

// Build 执行 LaTex 编译
func Build(opts BuildOptions) (string, error) {
	// 1-3. 准备编译上下文
	ctx, err := prepareContext(opts)
	if err != nil {
		return "", err
	}

	// 使用默认策略（处理包缺失和多轮编译）
	strategy := &defaultStrategy{}

	// 执行编译流程
	return buildWithStrategy(ctx, strategy)
}

// prepareContext 准备编译上下文
func prepareContext(opts BuildOptions) (*compileContext, error) {
	// 1. 获取输入文件的绝对路径
	absPath, err := filepath.Abs(opts.InputPath)
	if err != nil {
		return nil, fmt.Errorf(lang.T("msg.err_abs_path"), opts.InputPath)
	}

	// 检查是否是 .tex 文件
	if strings.ToLower(filepath.Ext(absPath)) != ".tex" {
		return nil, fmt.Errorf(lang.T("msg.err_invalid_ext"), opts.InputPath)
	}

	// 2. 检查文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf(lang.T("msg.err_file_not_found"), absPath)
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
				return nil, fmt.Errorf(lang.T("msg.err_mkdir"), outDir)
			}
		} else {
			// 如果是完整 PDF 路径
			outDir = filepath.Dir(outPath)
			// 从自定义路径中提取 JobName
			jobName = strings.TrimSuffix(filepath.Base(outPath), ".pdf")
			if err := os.MkdirAll(outDir, 0755); err != nil {
				return nil, fmt.Errorf(lang.T("msg.err_mkdir"), outDir)
			}
		}
	}

	return &compileContext{
		workDir:   workDir,
		outDir:    outDir,
		jobName:   jobName,
		fileName:  fileName,
		passCount: 0,
		firstPass: true,
	}, nil
}

// buildWithStrategy 使用策略执行编译流程
func buildWithStrategy(ctx *compileContext, strategy compileStrategy) (string, error) {
	for {
		// 更新编译次数
		ctx.passCount++

		// 执行编译
		logStr, err := compileOnce(ctx, ctx.firstPass)

		// 处理编译错误
		if err != nil {
			shouldContinue, newErr := strategy.handleError(ctx, logStr, err)
			if !shouldContinue {
				// 特殊错误：包已安装，需要重新开始
				if newErr != nil && newErr.Error() == "package_installed_retry" {
					fmt.Println(lang.T("msg.package_retry_build"))
					// 重置上下文，重新开始编译循环
					ctx.passCount = 0
					ctx.firstPass = true
					continue
				}
				return "", newErr
			}
		}

		// 检测多次编译需求
		needsMultiple, passInfo, _ := errors.HandleMultiplePasses(ctx.outDir, ctx.jobName, logStr)

		// 运行中间工具（如果需要）
		if strategy.shouldRunIntermediateTools(ctx, passInfo) {
			runIntermediateTools(ctx.workDir, ctx.jobName, passInfo)
		}

		// 判断是否继续编译
		if !needsMultiple || !strategy.shouldContinue(ctx, logStr, passInfo) {
			fmt.Println(lang.T("msg.build_success"))
			return filepath.Join(ctx.outDir, ctx.jobName+".pdf"), nil
		}

		// 准备下一次编译
		ctx.firstPass = false
		fmt.Printf(lang.T("msg.compile_pass_n")+"\n", ctx.passCount+1)
	}
}

// defaultStrategy 默认策略实现，处理包缺失和多轮编译
type defaultStrategy struct{}

// handleError 处理编译错误
func (s *defaultStrategy) handleError(ctx *compileContext, logStr string, err error) (bool, error) {
	// 第一次编译失败：处理包缺失
	if ctx.firstPass {
		handled, handleErr := errors.HandleMissingPackages(logStr)
		if handleErr != nil {
			return false, err
		}
		if handled {
			// 包已安装，返回特殊错误触发重新编译
			return false, fmt.Errorf("package_installed_retry")
		}
	}

	// 检测是否需要多次编译（即使编译失败也可能需要）
	needsMultiple, _, _ := errors.HandleMultiplePasses(ctx.outDir, ctx.jobName, logStr)
	if needsMultiple {
		// 需要多次编译，允许继续
		return true, nil
	}

	// 后续编译失败：允许继续（可能是引用问题，多次编译能解决）
	if !ctx.firstPass {
		return true, nil
	}

	// 第一次编译失败且不需要多次编译，返回错误
	return false, err
}

// shouldRunIntermediateTools 判断是否应该运行中间工具
func (s *defaultStrategy) shouldRunIntermediateTools(ctx *compileContext, passInfo errors.PassInfo) bool {
	// 第一次编译后，如果需要多次编译，运行中间工具
	return ctx.passCount == 1 && (passInfo.NeedsBibtex || passInfo.NeedsBiber ||
		passInfo.NeedsMakeindex || passInfo.NeedsMakeglossaries)
}

// shouldContinue 判断是否应该继续编译
func (s *defaultStrategy) shouldContinue(ctx *compileContext, logStr string, passInfo errors.PassInfo) bool {
	// 检测是否还需要继续编译
	needsMore, _, _ := errors.HandleMultiplePasses(ctx.outDir, ctx.jobName, logStr)

	// 检查是否达到最大编译次数
	if ctx.passCount >= passInfo.MaxPasses {
		return false
	}

	return needsMore
}

// compileOnce 执行单次编译，返回日志和错误
func compileOnce(ctx *compileContext, isFirstPass bool) (string, error) {
	args := []string{
		"-interaction=nonstopmode",
		"-jobname=" + ctx.jobName,
		"-output-directory=" + ctx.outDir,
		ctx.fileName,
	}
	cmd := exec.Command("xelatex", args...)
	cmd.Dir = ctx.workDir

	// 只在第一次编译时打印详细信息
	if isFirstPass {
		fmt.Printf(lang.T("msg.building_doc")+"\n", ctx.fileName)
		fmt.Printf(lang.T("msg.working_dir")+"\n", ctx.workDir)
	}

	var logOutput bytes.Buffer
	multiWriter := io.MultiWriter(os.Stdout, &logOutput)
	cmd.Stdout = multiWriter
	cmd.Stderr = multiWriter

	err := cmd.Run()
	return logOutput.String(), err
}

// runIntermediateTools 运行中间工具（bibtex、biber、makeindex等）
func runIntermediateTools(workDir, jobName string, passInfo errors.PassInfo) {
	if passInfo.NeedsBibtex {
		fmt.Println(lang.T("msg.running_bibtex"))
		runTool("bibtex", workDir, jobName)
	} else if passInfo.NeedsBiber {
		fmt.Println(lang.T("msg.running_biber"))
		runTool("biber", workDir, jobName)
	}

	if passInfo.NeedsMakeindex {
		fmt.Println(lang.T("msg.running_makeindex"))
		runTool("makeindex", workDir, jobName)
	}

	if passInfo.NeedsMakeglossaries {
		fmt.Println(lang.T("msg.running_makeglossaries"))
		runTool("makeglossaries", workDir, jobName)
	}
}

// runTool 运行单个工具
func runTool(toolName, workDir, jobName string) error {
	cmd := exec.Command(toolName, jobName)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf(lang.T("msg.tool_failed")+"\n", toolName, err)
		return err
	}
	return nil
}
