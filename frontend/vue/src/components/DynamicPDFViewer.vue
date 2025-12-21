<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
// 导入 CSS Modules
import styles from '../styles/dynamic-pdf-viewer.module.css'

// 在顶层导入 PDF.js（避免动态导入问题）
import * as pdfjsLib from 'pdfjs-dist'

// 导入 PDF 工具函数
import { parsePDF, calculateScale, cleanupPDF, formatPDFError } from '../utils/pdf'

// 导入 SSE API
import { createSSEConnection, closeSSEConnection } from '../api/sse'

// 导入 PDF API
import { fetchPDF, testBackendConnection } from '../api/pdf'

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

// 加载 PDF
// isInitialLoad: 是否是首次加载（首次加载显示 loading，刷新时不显示）
async function loadPDF(showLoading = true) {
  try {
    // 只在首次加载时显示 loading
    if (showLoading) {
      loading.value = true
    }
    error.value = ''
    
    // 使用 API 下载 PDF
    console.log('⬇️ 正在下载 PDF 文件...')
    const pdfBlob = await fetchPDF()
    console.log('✅ PDF 文件下载成功，大小:', (pdfBlob.size / 1024).toFixed(2), 'KB')
    
    // 将 Blob 转换为 ArrayBuffer
    console.log('📥 开始通过 PDF.js 解析 PDF...')
    const pdfArrayBuffer = await pdfBlob.arrayBuffer()
    
    // 使用工具函数解析 PDF
    const pdf = await parsePDF(pdfArrayBuffer)
    
    // 清理旧的 PDF 对象（如果存在）
    if (currentPdf) {
      cleanupPDF(currentPdf)
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
    
    // 使用工具函数格式化错误信息
    error.value = formatPDFError(err)
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
  const containerWidth = container ? container.clientWidth : 800
  const baseViewport = page.getViewport({ scale: 1.0 })
  
  // 使用工具函数计算自适应缩放比例
  const scale = calculateScale(containerWidth, baseViewport.width)
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
  eventSource = createSSEConnection({
    onMessage: async (event) => {
      if (event.data === 'refresh') {
        console.log('🔄 PDF 已更新，正在刷新...')
        await loadPDF(false)
      }
    },
    onError: (err) => {
      console.error('SSE 连接错误:', err)
    },
    onOpen: () => {
      console.log('SSE 连接已建立')
    }
  })
}

// 检查后端连接（使用 API）
async function checkBackendConnection() {
  const { success, error: errorMsg } = await testBackendConnection()
  if (!success) {
    error.value = errorMsg
    return false
  }
  return true
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
  const connected = await checkBackendConnection()
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
  closeSSEConnection(eventSource)
  // 使用工具函数清理 PDF 对象
  cleanupPDF(currentPdf)
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

<template>
  <div :class="styles.viewer">
    <!-- 只在首次加载时显示 loading -->
    <div v-if="loading && isInitialLoad" :class="styles.loading">正在加载 PDF...</div>
    <div v-if="error" :class="styles.error">{{ error }}</div>
    <!-- Canvas 始终显示，这样刷新时不会出现黑屏 -->
    <div :class="styles.canvasContainer" v-show="!error">
      <canvas ref="canvas" :class="styles.canvas"></canvas>
    </div>
  </div>
</template>
