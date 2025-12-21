// cmd/lazitex-cli/tasks/preview.go
// 实时预览业务管理

package tasks

import (
	//外部包
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	//内部包
	"github.com/SJRnhqh/lazitex/backend"
	core "github.com/SJRnhqh/lazitex/core"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// openBrowser 打开浏览器（跨平台）
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}

// StartLivePreview 启动实时预览模式（Web模式）
// filePath: 要预览的 .tex 文件路径
// quiet: 是否安静模式（不输出编译日志）
// tidy: 是否清理辅助文件
func StartLivePreview(filePath string, port int, quiet bool, tidy bool) {
	// 1. 获取绝对路径，确保监听准确
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Printf(lang.T("msg.err_abs_path")+"\n", filePath)
		return
	}

	// 2. 启动后立即执行一次"初次构建"
	opts := core.BuildOptions{
		InputPath:  absPath,
		OutputPath: "",
		Show:       false, // Web模式下不打开本地PDF
		Quiet:      quiet, // 使用传入的 quiet 参数
		Tidy:       tidy,  // 使用传入的 tidy 参数
	}

	pdfPath, err := core.Build(opts)
	if err != nil {
		fmt.Printf(lang.T("msg.build_failed")+": %v\n", err)
		return
	}

	if pdfPath == "" {
		fmt.Println(lang.T("msg.build_failed"))
		return
	}

	// 3. 创建并启动Web服务器
	server := backend.NewServer(port, pdfPath, lang.T("msg.pdf_not_found"), "view")

	// 在goroutine中启动服务器
	go func() {
		fmt.Printf(lang.T("msg.server_starting")+"\n", port)
		if err := server.Start(); err != nil {
			fmt.Printf(lang.T("msg.server_error")+"\n", err)
		}
	}()

	// 4. 打开浏览器（Vue 前端地址）
	// 注意：后端在 :8080，Vue 前端在 :5173（开发模式）或 :8080（生产模式）
	// 开发模式：打开 Vue 前端
	// 生产模式：前端会嵌入到后端，打开后端地址
	vueDevURL := "http://localhost:5173"
	backendURL := fmt.Sprintf("http://localhost:%d", port)

	// 尝试打开 Vue 前端（开发模式）
	// 如果 Vue 前端未运行，用户可以手动打开
	fmt.Printf("💡 提示：Vue 前端地址: %s\n", vueDevURL)
	fmt.Printf("💡 后端 API 地址: %s\n", backendURL)
	fmt.Printf("🔗 正在打开浏览器: %s\n", vueDevURL)
	time.Sleep(500 * time.Millisecond) // 等待服务器启动
	if err := openBrowser(vueDevURL); err != nil {
		fmt.Printf(lang.T("msg.browser_error")+"\n", err)
		fmt.Printf("请手动打开: %s\n", vueDevURL)
	}

	// 5. 打印监听提示
	fmt.Printf(lang.T("msg.watching_file")+"\n", filepath.Base(absPath))

	// 6. 监听文件变化并重新编译
	err = core.WatchAndAction(absPath, func() {
		currentTime := time.Now().Format(time.TimeOnly)
		fmt.Printf("\n🔄 [%s] %s\n", currentTime, lang.T("msg.building_doc"))

		// 重新编译（使用相同的 quiet 和 tidy 选项）
		opts := core.BuildOptions{
			InputPath:  absPath,
			OutputPath: "",
			Show:       false,
			Quiet:      quiet, // 使用相同的 quiet 选项
			Tidy:       tidy,  // 使用相同的 tidy 选项
		}

		newPdfPath, err := core.Build(opts)
		if err != nil {
			fmt.Printf(lang.T("msg.build_failed")+": %v\n", err)
			return
		}

		// 更新服务器中的PDF路径
		if newPdfPath != "" {
			server.SetPDFPath(newPdfPath)
			// 广播消息给所有SSE客户端
			server.BroadcastSSE("refresh")
		}
	})

	// 7. 错误处理
	if err != nil {
		fmt.Printf(lang.T("msg.watcher_error")+": %v\n", err)
	}

	// 8. 保持运行（等待 Ctrl+C）
	// 注意：这里需要处理信号来优雅关闭服务器
	// 暂时先保持简单，后续可以添加信号处理
	select {} // 阻塞等待
}
