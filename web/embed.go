// web/embed.go
// 嵌入静态前端资源

package web

import (
	_ "embed"
)

//go:embed static/index.html
var IndexHTML []byte
