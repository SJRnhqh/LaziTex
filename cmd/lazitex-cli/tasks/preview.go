// cmd/lazitex-cli/tasks/preview.go
// 实时预览业务管理

package tasks

import (
	//外部包
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	//内部包
	backend "github.com/SJRnhqh/lazitex/backend"
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

// checkVueDevServer 检查 Vue dev server 是否可用
func checkVueDevServer(url string, timeout time.Duration) bool {
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	// 如果返回 200-399 状态码，认为服务可用
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

// isPortAvailable 检查端口是否可用（未被占用）
func isPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false // 端口被占用
	}
	listener.Close()
	return true // 端口可用
}

// StartLivePreview 启动实时预览模式（Web模式）
// filePath: 要预览的 .tex 文件路径
// quiet: 是否安静模式（不输出编译日志）
// tidy: 是否清理辅助文件
func StartLivePreview(filePath string, port int, quiet bool, tidy bool, devMode bool) {
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

	// 3. 如果是开发模式，检测 Vue dev server 是否可用，并检查端口兼容性
	actualDevMode := devMode
	actualPort := port
	vueDevURL := "http://localhost:5173"
	const defaultBackendPort = 8080
	if devMode {
		fmt.Printf(lang.T("msg.checking_vue_dev_server")+"\n", vueDevURL)
		if checkVueDevServer(vueDevURL, 2*time.Second) {
			// Vue dev server 可用，检查端口兼容性
			if port != defaultBackendPort {
				// 使用了自定义端口，需要检查 8080 是否可用
				if isPortAvailable(defaultBackendPort) {
					// 8080 可用，强制使用 8080 以保持开发模式
					fmt.Printf(lang.T("msg.dev_mode_custom_port_warning")+"\n", port)
					actualPort = defaultBackendPort
					fmt.Println(lang.T("msg.vue_dev_server_available"))
					actualDevMode = true
				} else {
					// 8080 被占用，回退到生产模式并使用用户指定的端口
					fmt.Printf(lang.T("msg.port_8080_occupied")+"\n", port)
					actualDevMode = false
				}
			} else {
				// 已经使用默认端口 8080
				fmt.Println(lang.T("msg.vue_dev_server_available"))
				actualDevMode = true
			}
		} else {
			fmt.Println(lang.T("msg.vue_dev_server_unavailable"))
			actualDevMode = false
		}
	}

	// 4. 创建并启动Web服务器（使用实际检测到的模式和端口）
	server := backend.NewServer(actualPort, pdfPath, lang.T("msg.pdf_not_found"), "view", actualDevMode)

	// 在goroutine中启动服务器
	go func() {
		fmt.Printf(lang.T("msg.server_starting")+"\n", actualPort)
		if err := server.Start(); err != nil {
			fmt.Printf(lang.T("msg.server_error")+"\n", err)
		}
	}()

	// 5. 根据实际模式打开不同的 URL
	backendURL := fmt.Sprintf("http://localhost:%d", actualPort)

	if actualDevMode {
		// 开发模式：打开 Vue 开发服务器地址
		fmt.Printf(lang.T("msg.dev_mode_vue_url")+"\n", vueDevURL)
		fmt.Printf(lang.T("msg.dev_mode_backend_url")+"\n", backendURL)
		fmt.Printf(lang.T("msg.opening_browser")+"\n", vueDevURL)
		time.Sleep(500 * time.Millisecond) // 等待服务器启动
		if err := openBrowser(vueDevURL); err != nil {
			fmt.Printf(lang.T("msg.browser_error")+"\n", err)
			fmt.Printf(lang.T("msg.manual_open")+"\n", vueDevURL)
		}
	} else {
		// 生产模式：打开后端地址（前端已嵌入）
		fmt.Printf(lang.T("msg.prod_mode_embedded")+"\n", backendURL)
		fmt.Printf(lang.T("msg.opening_browser")+"\n", backendURL)
		time.Sleep(500 * time.Millisecond) // 等待服务器启动
		if err := openBrowser(backendURL); err != nil {
			fmt.Printf(lang.T("msg.browser_error")+"\n", err)
			fmt.Printf(lang.T("msg.manual_open")+"\n", backendURL)
		}
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
