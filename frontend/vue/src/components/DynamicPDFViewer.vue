<template>
  <div class="pdf-viewer">
    <!-- 只在首次加载时显示 loading -->
    <div v-if="loading && isInitialLoad" class="loading">正在加载 PDF...</div>
    <div v-if="error" class="error">{{ error }}</div>
    <!-- Canvas 始终显示，这样刷新时不会出现黑屏 -->
    <div class="canvas-container" v-show="!error">
      <canvas ref="canvas"></canvas>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
// 在顶层导入 PDF.js（避免动态导入问题）
import * as pdfjsLib from 'pdfjs-dist'
// 使用 Vite 的 ?url 导入语法获取 worker 文件 URL
import workerUrl from 'pdfjs-dist/build/pdf.worker.min.mjs?url'

// 配置 worker - 使用 Vite 的 URL 导入
// 这样 Vite 会正确处理 worker 文件
pdfjsLib.GlobalWorkerOptions.workerSrc = workerUrl
console.log('✅ 配置 PDF.js worker (Vite URL):', workerUrl)

const canvas = ref(null)
const loading = ref(false)
const error = ref('')
const isInitialLoad = ref(true) // 标记是否是首次加载
let eventSource = null
let currentPdf = null // 保存当前 PDF 对象，用于清理
let resizeObserver = null // 窗口大小变化监听器

// 加载 PDF.js（现在直接使用导入的模块）
async function ensurePDFJS() {
  return pdfjsLib
}

// 加载 PDF
// isInitialLoad: 是否是首次加载（首次加载显示 loading，刷新时不显示）
async function loadPDF(showLoading = true) {
  try {
    // 只在首次加载时显示 loading
    if (showLoading) {
      loading.value = true
    }
    error.value = ''
    
    console.log('📦 开始加载 PDF.js...')
    
    // 确保 PDF.js 已加载
    const pdfjs = await ensurePDFJS()
    console.log('✅ PDF.js 已加载，版本:', pdfjs.version)
    
    const pdfUrl = '/pdf?t=' + Date.now()
    console.log('📄 尝试加载 PDF:', pdfUrl)
    
    // 方法：先通过 fetch 下载 PDF 为 Blob，然后传给 PDF.js
    // 这样可以更好地处理错误和超时
    console.log('⬇️ 正在下载 PDF 文件...')
    const fetchController = new AbortController()
    const fetchTimeout = setTimeout(() => fetchController.abort(), 30000) // 30秒超时
    
    let pdfBlob
    try {
      const response = await fetch(pdfUrl, {
        signal: fetchController.signal,
      })
      clearTimeout(fetchTimeout)
      
      console.log('🔍 PDF 文件响应状态:', response.status, response.statusText)
      console.log('📊 Content-Type:', response.headers.get('Content-Type'))
      console.log('📏 Content-Length:', response.headers.get('Content-Length'))
      
      if (!response.ok) {
        throw new Error(`PDF 文件请求失败: ${response.status} ${response.statusText}`)
      }
      
      pdfBlob = await response.blob()
      console.log('✅ PDF 文件下载成功，大小:', (pdfBlob.size / 1024).toFixed(2), 'KB')
    } catch (fetchErr) {
      clearTimeout(fetchTimeout)
      console.error('❌ PDF 文件下载失败:', fetchErr)
      if (fetchErr.name === 'AbortError') {
        throw new Error('PDF 下载超时（30秒），文件可能过大或网络连接缓慢')
      }
      throw fetchErr
    }
    
    // 将 Blob 转换为 ArrayBuffer 或 URL
    console.log('📥 开始通过 PDF.js 解析 PDF...')
    const pdfArrayBuffer = await pdfBlob.arrayBuffer()
    
    // 使用 ArrayBuffer 加载 PDF（更可靠）
    const loadingTask = pdfjs.getDocument({
      data: pdfArrayBuffer,
      verbosity: 1,
    })
    
    // 添加超时（60秒，因为已经下载了）
    const parseTimeoutPromise = new Promise((_, reject) => {
      setTimeout(() => {
        loadingTask.destroy()
        reject(new Error('PDF 解析超时（60秒）'))
      }, 60000)
    })
    
    const pdf = await Promise.race([
      loadingTask.promise,
      parseTimeoutPromise
    ])
    
    // 清理旧的 PDF 对象（如果存在）
    if (currentPdf) {
      console.log('🧹 清理旧的 PDF 对象...')
      try {
        currentPdf.destroy()
      } catch (e) {
        console.warn('清理旧 PDF 时出错:', e)
      }
      currentPdf = null
    }
    
    // 保存新的 PDF 对象
    currentPdf = pdf
    
    console.log('✅ PDF 解析成功，页数:', pdf.numPages)
    console.log('🎨 开始渲染第一页...')
    
    // 渲染第一页
    await renderPage(pdf, 1)
    
    console.log('✨ PDF 渲染完成')
    loading.value = false
    isInitialLoad.value = false // 标记首次加载完成
  } catch (err) {
    console.error('❌ PDF 加载错误:', err)
    console.error('📋 错误详情:', err.message)
    if (err.stack) {
      console.error('📚 错误堆栈:', err.stack)
    }
    
    // 提供更友好的错误信息
    let errorMsg = 'PDF 加载失败: '
    if (err.message) {
      errorMsg += err.message
    } else {
      errorMsg += err.toString()
    }
    
    error.value = errorMsg
    loading.value = false
  }
}

// 渲染 PDF 页面
async function renderPage(pdf, pageNum) {
  // 确保 canvas 已经挂载（使用 nextTick 等待 DOM 更新）
  await nextTick()
  if (!canvas.value) {
    throw new Error('Canvas 元素尚未挂载，无法渲染 PDF')
  }
  
  const page = await pdf.getPage(pageNum)
  
  // 获取容器宽度，用于自适应缩放
  const container = canvas.value.parentElement
  const containerWidth = container ? container.clientWidth - 40 : 800 // 减去padding
  const baseViewport = page.getViewport({ scale: 1.0 })
  
  // 计算自适应缩放比例（确保PDF能完整显示在容器内）
  const scale = Math.min(1.5, (containerWidth / baseViewport.width) * 0.95)
  const viewport = page.getViewport({ scale: scale })
  
  // 设置 canvas 尺寸（这会自动清除 canvas）
  canvas.value.height = viewport.height
  canvas.value.width = viewport.width
  
  // 获取 canvas 上下文
  const context = canvas.value.getContext('2d')
  
  // 确保 canvas 可见（如果被隐藏，设置背景色）
  context.fillStyle = '#ffffff'
  context.fillRect(0, 0, canvas.value.width, canvas.value.height)
  
  // 渲染 PDF 页面到 canvas
  const renderContext = {
    canvasContext: context,
    viewport: viewport
  }
  await page.render(renderContext).promise
  
  console.log('✅ PDF 页面渲染完成，尺寸:', canvas.value.width, 'x', canvas.value.height)
}

// 监听 SSE 事件（PDF 更新）
function setupSSE() {
  eventSource = new EventSource('/events')
  
  eventSource.onmessage = async (event) => {
    if (event.data === 'refresh') {
      console.log('🔄 PDF 已更新，正在刷新...')
      // 刷新时不显示 loading（保持旧内容可见，直到新内容渲染完成）
      // 这样用户体验更平滑，不会出现黑屏
      await loadPDF(false) // 传入 false 表示不是首次加载
    }
  }
  
  eventSource.onerror = (err) => {
    console.error('SSE 连接错误:', err)
    // 可以在这里添加重连逻辑
  }
  
  eventSource.onopen = () => {
    console.log('SSE 连接已建立')
  }
}

// 测试后端连接
async function testBackendConnection() {
  try {
    console.log('🔍 测试后端连接...')
    console.log('📍 PDF 请求 URL: /pdf')
    
    const response = await fetch('/pdf', { method: 'HEAD' })
    console.log('✅ 后端响应状态:', response.status, response.statusText)
    console.log('📋 响应头:', Object.fromEntries(response.headers.entries()))
    
    if (response.status === 404) {
      error.value = 'PDF 文件不存在，请确保 LaTeX 文件已成功编译'
      return false
    }
    
    if (!response.ok) {
      error.value = `后端返回错误: ${response.status} ${response.statusText}`
      return false
    }
    
    return true
  } catch (err) {
    console.error('❌ 后端连接测试失败:', err)
    console.error('💡 请确保：')
    console.error('   1. Go 后端正在运行（go run ./cmd/lazitex-cli -p test/001.tex）')
    console.error('   2. 后端运行在 http://localhost:8080')
    console.error('   3. Vite 代理配置正确（vite.config.js）')
    error.value = '无法连接到后端服务器（:8080），请确保 Go 后端正在运行'
    return false
  }
}

// 处理窗口大小变化，重新渲染PDF以适应新尺寸
function handleResize() {
  if (currentPdf && canvas.value) {
    console.log('🔄 窗口大小变化，重新渲染 PDF...')
    renderPage(currentPdf, 1).catch(err => {
      console.error('重新渲染失败:', err)
    })
  }
}

onMounted(async () => {
  console.log('🚀 PDFViewer 组件已挂载')
  
  // 监听窗口大小变化
  if (typeof window !== 'undefined') {
    window.addEventListener('resize', handleResize)
    // 使用 ResizeObserver 监听容器大小变化（更精确）
    await nextTick()
    if (canvas.value?.parentElement) {
      resizeObserver = new ResizeObserver(() => {
        handleResize()
      })
      resizeObserver.observe(canvas.value.parentElement)
    }
  }
  
  // 先测试后端连接
  const connected = await testBackendConnection()
  if (connected) {
    console.log('✅ 后端连接正常，开始加载 PDF')
    setupSSE()
    loadPDF()
  } else {
    console.log('❌ 后端连接失败，停止加载')
    loading.value = false
  }
})

onUnmounted(() => {
  // 清理资源
  if (eventSource) {
    eventSource.close()
  }
  if (currentPdf) {
    try {
      currentPdf.destroy()
    } catch (e) {
      console.warn('清理 PDF 时出错:', e)
    }
  }
  // 移除窗口大小监听
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', handleResize)
  }
  // 断开 ResizeObserver
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
})
</script>

<style scoped>
/* Nord 配色方案 */
.pdf-viewer {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #2E3440; /* nord0 - 背景 */
  overflow: hidden;
}

.loading,
.error {
  color: #D8DEE9; /* nord4 - 文字 */
  font-size: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.error {
  color: #BF616A; /* nord11 - 错误颜色（红色） */
}

.canvas-container {
  flex: 1;
  overflow: auto;
  padding: 20px;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  background: #2E3440; /* nord0 - 背景 */
  /* 平滑滚动 */
  scroll-behavior: smooth;
}

canvas {
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  background: #ECEFF4; /* nord6 - canvas 背景（浅色以显示 PDF） */
  display: block;
  /* 确保canvas不会超出容器 */
  max-width: 100%;
  height: auto;
}
</style>

