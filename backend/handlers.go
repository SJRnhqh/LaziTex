// backend/handlers.go

package backend

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/SJRnhqh/lazitex/web"
)

// HandleIndex 处理首页请求
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(web.IndexHTML)
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
