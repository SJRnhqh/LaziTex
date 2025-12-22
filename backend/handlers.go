// backend/handlers.go

package backend

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	cfg "github.com/SJRnhqh/lazitex/config"
	frontend "github.com/SJRnhqh/lazitex/frontend"
)

// HandleIndex 处理首页请求（生产模式：返回嵌入的 Vue 前端）
func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	html, err := frontend.GetVueIndexHTML()
	if err != nil {
		// 如果获取 Vue 前端失败，回退到旧的静态页面
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(frontend.IndexHTML)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(html)
}

// HandleIndexDev 处理首页请求（开发模式：返回提示信息）
func (s *Server) HandleIndexDev(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<title>LaziTex - Development Mode</title>
	<meta charset="utf-8">
	<style>
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
			margin: 0;
			background: #2E3440;
			color: #ECEFF4;
		}
		.container {
			text-align: center;
			max-width: 600px;
			padding: 20px;
		}
		h1 {
			margin-bottom: 20px;
		}
		p {
			margin: 10px 0;
			line-height: 1.6;
		}
		a {
			color: #5E81AC;
			text-decoration: none;
			font-weight: bold;
		}
		a:hover {
			text-decoration: underline;
		}
		.code {
			background: #3B4252;
			padding: 8px 12px;
			border-radius: 4px;
			font-family: 'Courier New', monospace;
			display: inline-block;
			margin: 5px 0;
		}
		.tip {
			margin-top: 30px;
			padding: 15px;
			background: #434C5E;
			border-radius: 4px;
			font-size: 14px;
		}
		code {
			background: #3B4252;
			padding: 2px 6px;
			border-radius: 3px;
			font-family: 'Courier New', monospace;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>🔧 Development Mode</h1>
		<p>Please access the Vue dev server at:</p>
		<p><a href="http://localhost:5173">http://localhost:5173</a></p>
		
		<div class="tip">
			<p><strong>💡 Tip:</strong> If the Vue dev server is not running, start it with:</p>
			<p class="code">cd frontend/vue && npm run dev</p>
			<p style="margin-top: 15px; font-size: 13px;">Or if you want to use the embedded frontend instead, restart without the <code>-d</code> flag.</p>
		</div>
	</div>
</body>
</html>`))
}

// HandlePDF 处理PDF文件请求
// 注意：这是 Server 的方法，通过 s.pdfPath 获取PDF路径
func (s *Server) HandlePDF(w http.ResponseWriter, r *http.Request) {
	pdfPath := s.pdfPath

	// 检查文件是否存在
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		http.Error(w, s.errorMsg, http.StatusNotFound)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(pdfPath))

	// 读取并返回PDF文件
	http.ServeFile(w, r, pdfPath)
}

// HandleConfig 处理配置请求
// GET: 返回用户配置和应用模式
// POST: 更新用户配置
func (s *Server) HandleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// GET 请求：返回配置
		userConfig, err := cfg.LoadConfig()
		if err != nil {
			// 如果加载失败，使用默认值
			userConfig = &cfg.Config{
				Language:    "en",
				DisplayMode: "virtual",
			}
		}

		response := map[string]string{
			"mode":        s.mode,                 // 服务器模式（LaziView/LaziWorkspace）
			"language":    userConfig.Language,    // 用户语言偏好
			"displayMode": userConfig.DisplayMode, // PDF 显示模式
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	case http.MethodPost:
		// POST 请求：更新配置
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// 使用 UpdateConfig 部分更新配置
		if err := cfg.UpdateConfig(updates); err != nil {
			http.Error(w, "Failed to save config", http.StatusInternalServerError)
			return
		}

		// 返回更新后的配置
		userConfig, _ := cfg.LoadConfig()
		response := map[string]string{
			"mode":        s.mode,
			"language":    userConfig.Language,
			"displayMode": userConfig.DisplayMode,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
