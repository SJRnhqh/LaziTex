// frontend/embed.go
// 嵌入静态前端资源

package frontend

import (
	_ "embed"
)

//go:embed static/index.html
var IndexHTML []byte
