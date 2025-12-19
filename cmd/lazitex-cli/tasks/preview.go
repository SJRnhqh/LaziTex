// cmd/lazitex-cli/tasks/preview.go
// 实时预览业务管理

package tasks

import (
	//外部包
	"fmt"
	"path/filepath"
	"time"

	//内部包
	core "github.com/SJRnhqh/lazitex/core"
	lang "github.com/SJRnhqh/lazitex/lang"
)

// StartLivePreview 启动实时预览模式
// filePath: 要预览的 .tex 文件路径
func StartLivePreview(filePath string) {
	// 1. 获取绝对路径，确保监听准确
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Printf(lang.T("msg.err_abs_path")+"\n", filePath)
		return
	}

	// 2. 启动后立即执行一次“初次构建并展示”
	// 这里调用我们已有的 BuildLaTex 函数，设置展示标志为 true
	BuildLaTex(absPath, "", true)

	// 3. 打印监听提示（文案已在 i18n 中定义）
	fmt.Printf(lang.T("msg.watching_file")+"\n", filepath.Base(absPath))

	// 4. 调用 core 层的监听引擎
	// 当文件变动时，它会回调执行我们定义的闭包函数
	err = core.WatchAndAction(absPath, func() {
		// 这里是文件变动后的动作
		currentTime := time.Now().Format(time.TimeOnly)
		fmt.Printf("\n🔄 [%s] %s\n", currentTime, lang.T("msg.building_doc"))

		// 重新执行编译和展示逻辑
		BuildLaTex(absPath, "", true)
	})

	// 5. 错误处理（如果监听器意外崩溃）
	if err != nil {
		fmt.Printf(lang.T("msg.watcher_error")+": %v\n", err)
	}
}
