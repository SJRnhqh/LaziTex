// core/tools/lock.go
// 编译锁和任务管理，防止并发编译冲突，支持任务抢占

package tools

import (
	"context"
	"fmt"
	"sync"
	"time"

	lang "github.com/SJRnhqh/lazitex/lang"
)

// CompileManager 编译管理器（单例模式）
// 负责管理编译任务的并发控制和任务抢占
type CompileManager struct {
	mu          sync.Mutex
	currentTask *CompileTask
	timeout     time.Duration // 编译超时时间
}

// CompileTask 编译任务
type CompileTask struct {
	ctx      context.Context
	cancel   context.CancelFunc
	filePath string
}

// globalManager 全局编译管理器实例
var globalManager = &CompileManager{
	timeout: 5 * time.Minute, // 默认超时时间：5分钟
}

// GetCompileManager 获取全局编译管理器实例
func GetCompileManager() *CompileManager {
	return globalManager
}

// SetTimeout 设置编译超时时间
func (m *CompileManager) SetTimeout(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.timeout = duration
}

// GetTimeout 获取编译超时时间
func (m *CompileManager) GetTimeout() time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.timeout
}


// StartBuild 开始编译任务（带锁和任务抢占）
// 如果已有编译任务在运行，会取消旧任务并启动新任务
// 返回 context 和 cancel 函数，用于控制编译过程
func (m *CompileManager) StartBuild(filePath string) (context.Context, context.CancelFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果已有任务在运行，取消它
	if m.currentTask != nil {
		if m.currentTask.cancel != nil {
			fmt.Println(lang.T("msg.cancelling_previous_task"))
			m.currentTask.cancel()
		}
	}

	// 创建带超时的 context
	// 先创建一个可取消的 context
	ctx, cancel := context.WithCancel(context.Background())

	// 如果有设置超时，包装成带超时的 context
	if m.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, m.timeout)
	}

	// 保存新任务
	m.currentTask = &CompileTask{
		ctx:      ctx,
		cancel:   cancel,
		filePath: filePath,
	}

	return ctx, cancel
}

// FinishBuild 完成编译任务（释放锁）
func (m *CompileManager) FinishBuild() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 清除当前任务
	m.currentTask = nil
}

// IsBuilding 检查是否正在编译
func (m *CompileManager) IsBuilding() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.currentTask != nil
}

// GetCurrentTask 获取当前任务（用于调试）
func (m *CompileManager) GetCurrentTask() *CompileTask {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.currentTask
}
