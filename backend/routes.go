// backend/routes.go

package backend

// RegisterRoutes 注册所有路由
func (s *Server) RegisterRoutes() {
	// 首页路由
	s.mux.HandleFunc("/", HandleIndex)

	// PDF文件路由（使用 Server 的方法）
	s.mux.HandleFunc("/pdf", s.HandlePDF)

	// SSE 事件流路由（新增）
	s.mux.HandleFunc("/events", s.HandleSSE)
}
