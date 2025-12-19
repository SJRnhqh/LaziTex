// backend/server.go

package backend

import (
	"context"
	"fmt"
	"net/http"
)

// Server HTTP服务器
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	pdfPath    string
	errorMsg   string // PDF文件不存在的错误消息（支持国际化）
}

// NewServer 创建新的HTTP服务器
func NewServer(port int, pdfPath string, errorMsg string) *Server {
	mux := http.NewServeMux()
	// 如果没有提供错误消息，使用默认英文
	if errorMsg == "" {
		errorMsg = "PDF file not found"
	}
	server := &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		mux:      mux,
		pdfPath:  pdfPath,
		errorMsg: errorMsg,
	}

	// 注册路由
	server.RegisterRoutes()

	return server
}

// GetMux 获取路由管理器
func (s *Server) GetMux() *http.ServeMux {
	return s.mux
}

// SetPDFPath 设置PDF路径（用于更新）
func (s *Server) SetPDFPath(pdfPath string) {
	s.pdfPath = pdfPath
}

// Start 启动服务器
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Stop 停止服务器
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
