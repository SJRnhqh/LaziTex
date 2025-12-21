// backend/server.go

package backend

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// Server HTTP服务器
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	pdfPath    string
	errorMsg   string // PDF文件不存在的错误消息（支持国际化）

	// SSE 客户端管理（新增）
	sseClients map[chan string]bool // SSE 客户端通道映射
	sseMu      sync.Mutex           // 保护 sseClients 的互斥锁

	// 应用模式（新增）
	mode string
}

// NewServer 创建新的HTTP服务器
func NewServer(port int, pdfPath string, errorMsg string, mode string) *Server {
	mux := http.NewServeMux()
	// 如果没有提供错误消息，使用默认英文
	if errorMsg == "" {
		errorMsg = "PDF file not found"
	}
	// 如果没有提供模式，使用默认值 "view"
	if mode == "" {
		mode = "workspace"
	}
	server := &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		mux:        mux,
		pdfPath:    pdfPath,
		errorMsg:   errorMsg,
		sseClients: make(map[chan string]bool),
		sseMu:      sync.Mutex{},
		mode:       mode,
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
