// frontend/embed.go
// 嵌入静态前端资源

package frontend

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/index.html
var IndexHTML []byte

// 嵌入 Vue 前端构建产物（整个 dist 目录及其子目录）
//go:embed dist
var vueDistFS embed.FS

// GetVueIndexHTML 获取 Vue 前端的 index.html 内容
func GetVueIndexHTML() ([]byte, error) {
	return vueDistFS.ReadFile("dist/index.html")
}

// ServeVueAssets 返回一个 HTTP Handler，用于服务 Vue 前端的静态资源（JS/CSS 等）
// 这会处理 /assets/ 路径的请求
func ServeVueAssets() http.Handler {
	// 获取 dist 子目录作为文件系统的根
	distFS, err := fs.Sub(vueDistFS, "dist")
	if err != nil {
		// 如果出错，返回 404 处理器
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
	}
	
	// 再次获取 assets 子目录，这样文件服务器的根就是 assets/
	assetsFS, err := fs.Sub(distFS, "assets")
	if err != nil {
		// 如果出错，返回 404 处理器
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
	}
	
	// 返回文件服务器，根目录是 assets/，请求 /assets/xxx 会被 StripPrefix 后变成 /xxx
	// 这样就能正确找到 assets/xxx 文件了
	return http.FileServer(http.FS(assetsFS))
}