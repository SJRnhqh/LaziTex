// api/pdf.js
// PDF API，负责管理 PDF 文件的下载和连接测试

/**
 * 下载 PDF 文件
 * @param {number} timestamp - 时间戳（用于缓存破坏），默认当前时间
 * @param {number} timeout - 超时时间（毫秒），默认 30 秒
 * @returns {Promise<Blob>} PDF 文件 Blob
 */
export async function fetchPDF(timestamp = Date.now(), timeout = 30000) {
    const pdfUrl = `/pdf?t=${timestamp}`
    const fetchController = new AbortController()
    const fetchTimeout = setTimeout(() => fetchController.abort(), timeout)
    
    try {
        const response = await fetch(pdfUrl, {
        signal: fetchController.signal,
        })
        clearTimeout(fetchTimeout)
        
        if (!response.ok) {
        throw new Error(`PDF 文件请求失败: ${response.status} ${response.statusText}`)
        }
        
        return await response.blob()
    } catch (fetchErr) {
        clearTimeout(fetchTimeout)
        if (fetchErr.name === 'AbortError') {
        throw new Error(`PDF 下载超时（${timeout / 1000}秒），文件可能过大或网络连接缓慢`)
        }
        throw fetchErr
    }
}

/**
 * 测试后端连接
 * @returns {Promise<{success: boolean, error?: string}>}
 */
export async function testBackendConnection() {
    try {
        const response = await fetch('/pdf', { method: 'HEAD' })
        
        if (response.status === 404) {
        return { 
            success: false, 
            error: 'PDF 文件不存在，请确保 LaTeX 文件已成功编译' 
        }
        }
        
        if (!response.ok) {
            return { 
                success: false, 
                error: `后端返回错误: ${response.status} ${response.statusText}` 
            }
        }
        
        return { success: true }
    } catch (err) {
        return { 
            success: false, 
            error: '无法连接到后端服务器（:8080），请确保 Go 后端正在运行' 
        }
    }
}
