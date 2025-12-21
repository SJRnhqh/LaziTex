// backend/handlers.go

package backend

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	frontend "github.com/SJRnhqh/lazitex/frontend"
)

// HandleIndex 处理首页请求
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(frontend.IndexHTML)
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

// HandleConfig 处理配置请求，返回应用模式
func (s *Server) HandleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 返回 JSON 格式的配置信息
	response := map[string]string{
		"mode": s.mode,
	}

	// 将 map 转换为 JSON
	json.NewEncoder(w).Encode(response)
}
