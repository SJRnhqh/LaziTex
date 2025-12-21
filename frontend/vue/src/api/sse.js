// api/sse.js
// sse.js 是 LaziTeX 的 SSE API，负责管理 SSE 的连接和消息

/**
 * 创建 SSE 连接
 * @param {Object} options - 配置选项
 * @param {Function} options.onMessage - 消息回调函数，接收 (event) => void
 * @param {Function} [options.onError] - 错误回调函数，接收 (err) => void
 * @param {Function} [options.onOpen] - 连接打开回调函数，接收 () => void
 * @param {string} [options.url='/events'] - SSE 端点 URL
 * @returns {EventSource} EventSource 实例
 */
export function createSSEConnection({ onMessage, onError, onOpen, url = '/events' }) {
    const eventSource = new EventSource(url)
    
    if (onMessage) {
        eventSource.onmessage = onMessage
    }
    
    if (onError) {
        eventSource.onerror = onError
    }
    
    if (onOpen) {
        eventSource.onopen = onOpen
    }
    
    return eventSource
}

/**
 * 关闭 SSE 连接
 * @param {EventSource} eventSource - EventSource 实例
 */
export function closeSSEConnection(eventSource) {
    if (eventSource) {
        eventSource.close()
    }
}