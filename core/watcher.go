// core/watcher.go

package core

import (
	// 外部包
	"path/filepath"
	"time"

	fsnotify "github.com/fsnotify/fsnotify"
)

// WatchAndAction 监听文件变动并执行指定的动作
// filePath: 要监听的 .tex 文件路径
// action: 文件变动后要执行的函数（比如编译预览函数）
func WatchAndAction(filePath string, action func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	absPath, _ := filepath.Abs(filePath)
	// 监听文件所在的目录，这在 Windows/macOS 上比直接监听单个文件更稳定
	err = watcher.Add(filepath.Dir(absPath))
	if err != nil {
		return err
	}

	// 防抖计时器，避免保存时瞬间触发多次编译
	var timer *time.Timer
	const debounceDuration = 200 * time.Millisecond

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			// 只在目标文件发生写入（Write）事件时触发
			if event.Op&fsnotify.Write == fsnotify.Write && event.Name == absPath {
				if timer != nil {
					timer.Stop()
				}
				// 延迟执行动作，实现防抖
				timer = time.AfterFunc(debounceDuration, action)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			return err
		}
	}
}
