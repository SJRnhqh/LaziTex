// backend/routes.go

package backend

import (
	"net/http"

	frontend "github.com/SJRnhqh/lazitex/frontend"
)

// RegisterRoutes 注册所有路由
func (s *Server) RegisterRoutes() {
	// 根据是否是开发模式决定如何处理前端
	if s.devMode {
		// 开发模式：首页返回提示信息
		s.mux.HandleFunc("/", s.HandleIndexDev)
	} else {
		// 生产模式：首页返回嵌入的 Vue 前端
		s.mux.HandleFunc("/", s.HandleIndex)
		// 服务 Vue 前端的静态资源（JS/CSS 等）
		// 注意：Vue 构建产物中的资源路径是 /assets/xxx，所以这里注册 /assets/ 路由
		s.mux.Handle("/assets/", http.StripPrefix("/assets", frontend.ServeVueAssets()))
	}
	// 开发模式下不注册根路径，用户应该访问 Vue dev server (:5173)

	// PDF文件路由（使用 Server 的方法）
	s.mux.HandleFunc("/pdf", s.HandlePDF)

	// SSE 事件流路由（新增）
	s.mux.HandleFunc("/events", s.HandleSSE)

	// 配置API路由（返回应用模式）
	s.mux.HandleFunc("/api/config", s.HandleConfig)
}
