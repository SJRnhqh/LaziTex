// backend/sse.go

package backend

import (
	"fmt"
	"net/http"
)

// HandleSSE 处理 SSE 连接
// 当浏览器访问 /events 时，这个函数会被调用
func (s *Server) HandleSSE(w http.ResponseWriter, r *http.Request) {
	// 1. 设置 SSE 响应头（告诉浏览器这是流式响应）
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 2. 为这个客户端创建一个通道（用于接收消息）
	clientChan := make(chan string, 1)

	// 3. 将这个客户端添加到服务器管理的客户端列表中
	s.sseMu.Lock()                  // 加锁保护
	s.sseClients[clientChan] = true // 注册客户端
	s.sseMu.Unlock()                // 解锁

	// 4. 清理函数：当客户端断开时，从列表中移除
	defer func() {
		s.sseMu.Lock()
		delete(s.sseClients, clientChan) // 移除客户端
		close(clientChan)                // 关闭通道
		s.sseMu.Unlock()
	}()

	// 5. 获取请求的上下文（用于检测客户端是否断开）
	ctx := r.Context()

	// 6. 发送初始连接消息（告诉客户端连接成功）
	fmt.Fprintf(w, "data: connected\n\n")
	if f, ok := w.(http.Flusher); ok {
		f.Flush() // 立即发送，不等函数结束
	}

	// 7. 持续监听，等待消息或客户端断开
	for {
		select {
		case <-ctx.Done():
			// 客户端断开连接，退出循环
			return
		case msg := <-clientChan:
			// 收到要发送给客户端的消息
			// SSE 格式：data: 消息内容\n\n
			fmt.Fprintf(w, "data: %s\n\n", msg)
			if f, ok := w.(http.Flusher); ok {
				f.Flush() // 立即发送给客户端
			}
		}
	}
}

// BroadcastSSE 广播消息给所有 SSE 客户端
// 当 PDF 更新时，调用这个函数通知所有浏览器刷新
func (s *Server) BroadcastSSE(message string) {
	s.sseMu.Lock()         // 加锁保护
	defer s.sseMu.Unlock() // 函数结束时解锁

	// 遍历所有连接的客户端，发送消息
	for clientChan := range s.sseClients {
		select {
		case clientChan <- message:
			// 成功发送消息到客户端通道
		default:
			// 客户端通道已满，跳过（避免阻塞）
			// 这种情况很少发生，因为通道容量是 1
		}
	}
}
