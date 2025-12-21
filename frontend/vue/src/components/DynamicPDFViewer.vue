<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
// 导入 CSS Modules
import styles from '../styles/dynamic-pdf-viewer.module.css'

// 在顶层导入 PDF.js（避免动态导入问题）
import * as pdfjsLib from 'pdfjs-dist'

// 导入 PDF 工具函数
import { 
  parsePDF, 
  calculateScale, 
  cleanupPDF, 
  formatPDFError,
  renderPageToCanvas,
  calculatePageHeights
} from '../utils/pdf'

// 导入 SSE API
import { createSSEConnection, closeSSEConnection } from '../api/sse'

// 导入 PDF API
import { fetchPDF, testBackendConnection } from '../api/pdf'

// 使用 Vite 的 ?url 导入语法获取 worker 文件 URL
import workerUrl from 'pdfjs-dist/build/pdf.worker.min.mjs?url'

// 配置 worker - 使用 Vite 的 URL 导入
pdfjsLib.GlobalWorkerOptions.workerSrc = workerUrl
console.log('✅ 配置 PDF.js worker (Vite URL):', workerUrl)

// ========== 响应式数据 ==========
const loading = ref(false)
const error = ref('')
const isInitialLoad = ref(true)

// 显示模式：'virtual' (虚拟滚动) 或 'pagination' (分页导航)
// 预留接口，未来可以从配置或用户设置中读取
const displayMode = ref('virtual') // 默认虚拟滚动

// 分页模式相关（预留）
const canvas = ref(null) // 分页模式需要
const currentPageNum = ref(1) // 当前页码（分页模式）

// PDF 相关
let eventSource = null
let currentPdf = null
let resizeObserver = null
let pageObserver = null // Intersection Observer

// 页面状态管理（虚拟滚动）
const pageHeights = ref([]) // 每页的高度
const renderedPages = ref(new Set()) // 已渲染的页面
const pageContainers = ref([]) // 页面容器 DOM 引用

// DOM 引用
const container = ref(null) // 滚动容器

// ========== 清理资源 ==========
function cleanup() {
  // 清理 Intersection Observer
  if (pageObserver) {
    pageObserver.disconnect()
    pageObserver = null
  }
  
  // 清理 PDF
  cleanupPDF(currentPdf)
  currentPdf = null
  
  // 重置状态
  pageHeights.value = []
  renderedPages.value.clear()
  pageContainers.value = []
  
  // 清空容器
  if (container.value) {
    container.value.innerHTML = ''
  }
}

// ========== PDF 加载 ==========
async function loadPDF(showLoading = true) {
  try {
    // 只在首次加载时显示 loading
    if (showLoading) {
      loading.value = true
    }
    error.value = ''
    
    // 清理旧的 PDF 和观察器
    cleanup()
    
    // 下载 PDF
    console.log('⬇️ 正在下载 PDF 文件...')
    const pdfBlob = await fetchPDF()
    console.log('✅ PDF 文件下载成功，大小:', (pdfBlob.size / 1024).toFixed(2), 'KB')
    
    // 解析 PDF
    console.log('📥 开始通过 PDF.js 解析 PDF...')
    const pdfArrayBuffer = await pdfBlob.arrayBuffer()
    const pdf = await parsePDF(pdfArrayBuffer)
    
    // 保存新的 PDF 对象
    currentPdf = pdf
    
    console.log('✅ PDF 解析成功，页数:', pdf.numPages)
    
    // 根据显示模式初始化
    if (displayMode.value === 'virtual') {
      await initVirtualScroll(pdf)
    } else {
      await initPagination(pdf)
    }
    
    loading.value = false
    isInitialLoad.value = false
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

// ========== 虚拟滚动模式 ==========
async function initVirtualScroll(pdf) {
  console.log('🎨 初始化虚拟滚动模式...')
  
  await nextTick()
  if (!container.value) {
    throw new Error('容器元素尚未挂载')
  }
  
  const containerWidth = container.value.clientWidth - 40 // 减去 padding
  
  // 计算所有页面的高度
  console.log('📏 计算页面高度...')
  pageHeights.value = await calculatePageHeights(pdf, containerWidth)
  
  // 创建页面容器
  await createPageContainers(pdf.numPages)
  
  // 设置 Intersection Observer
  setupPageObserver()
  
  // 初始渲染可见页面
  await renderVisiblePages()
  
  console.log('✨ 虚拟滚动初始化完成')
}

// 创建页面容器（占位符）
async function createPageContainers(totalPages) {
  await nextTick()
  if (!container.value) return
  
  // 清空容器
  container.value.innerHTML = ''
  pageContainers.value = []
  renderedPages.value.clear()
  
  // 为每一页创建容器
  for (let i = 1; i <= totalPages; i++) {
    const pageDiv = document.createElement('div')
    pageDiv.className = styles.pageContainer
    pageDiv.dataset.pageNum = i
    
    // 设置占位高度
    const height = pageHeights.value[i - 1] || 800
    pageDiv.style.height = `${height}px`
    pageDiv.style.minHeight = `${height}px`
    
    // 创建占位符
    const placeholder = document.createElement('div')
    placeholder.className = styles.pagePlaceholder
    placeholder.style.width = '100%'
    placeholder.style.height = '100%'
    placeholder.textContent = `页面 ${i}`
    pageDiv.appendChild(placeholder)
    
    container.value.appendChild(pageDiv)
    pageContainers.value.push(pageDiv)
  }
}

// 设置 Intersection Observer
function setupPageObserver() {
  // 清理旧的观察器
  if (pageObserver) {
    pageObserver.disconnect()
  }
  
  // 创建新的观察器
  // rootMargin: '200px' 表示提前 200px 预加载
  pageObserver = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      const pageNum = parseInt(entry.target.dataset.pageNum)
      
      if (entry.isIntersecting) {
        // 页面进入视口，渲染它
        if (!renderedPages.value.has(pageNum)) {
          renderPageVirtual(pageNum)
        }
      }
      // 页面离开视口时，保留已渲染的页面（不清理，提升体验）
    })
  }, {
    root: container.value,
    rootMargin: '200px', // 提前 200px 预加载
    threshold: 0 // 一进入就触发
  })
  
  // 观察所有页面容器
  pageContainers.value.forEach(pageDiv => {
    pageObserver.observe(pageDiv)
  })
}

// 渲染单个页面（虚拟滚动版）
async function renderPageVirtual(pageNum) {
  if (!currentPdf || renderedPages.value.has(pageNum)) {
    return
  }
  
  try {
    const pageDiv = pageContainers.value[pageNum - 1]
    if (!pageDiv) return
    
    // 创建 canvas
    const canvas = document.createElement('canvas')
    canvas.className = styles.canvas
    
    // 获取容器宽度
    const containerWidth = container.value.clientWidth - 40
    
    // 渲染页面
    await renderPageToCanvas(currentPdf, pageNum, canvas, containerWidth)
    
    // 替换占位符
    const placeholder = pageDiv.querySelector(`.${styles.pagePlaceholder}`)
    if (placeholder) {
      pageDiv.removeChild(placeholder)
    }
    pageDiv.appendChild(canvas)
    
    // 标记为已渲染
    renderedPages.value.add(pageNum)
    
    console.log(`✅ 页面 ${pageNum} 渲染完成`)
  } catch (err) {
    console.error(`❌ 渲染页面 ${pageNum} 失败:`, err)
  }
}

// 初始渲染可见页面
async function renderVisiblePages() {
  // 渲染第一页（通常可见）
  if (pageContainers.value.length > 0) {
    await renderPageVirtual(1)
  }
}

// ========== 分页导航模式（预留） ==========
async function initPagination(pdf) {
  console.log('🎨 初始化分页导航模式...')
  
  await nextTick()
  if (!container.value) {
    throw new Error('容器元素尚未挂载')
  }
  
  // 创建单页容器
  container.value.innerHTML = ''
  const pageDiv = document.createElement('div')
  pageDiv.className = styles.pageContainer
  pageDiv.dataset.pageNum = 1
  
  // 创建 canvas
  const canvasEl = document.createElement('canvas')
  canvasEl.className = styles.canvas
  
  const containerWidth = container.value.clientWidth - 40
  await renderPageToCanvas(pdf, 1, canvasEl, containerWidth)
  
  pageDiv.appendChild(canvasEl)
  container.value.appendChild(pageDiv)
  
  // 保存 canvas 引用（用于后续操作）
  canvas.value = canvasEl
  
  currentPageNum.value = 1
  
  console.log('✨ 分页导航模式初始化完成（当前仅显示第一页）')
}

// 渲染 PDF 页面（分页模式使用）
async function renderPage(pdf, pageNum) {
  await nextTick()
  if (!canvas.value) {
    throw new Error('Canvas 元素尚未挂载，无法渲染 PDF')
  }
  
  const containerWidth = canvas.value.parentElement 
    ? canvas.value.parentElement.clientWidth - 40 
    : 800
  
  await renderPageToCanvas(pdf, pageNum, canvas.value, containerWidth)
  
  console.log('✅ PDF 页面渲染完成，尺寸:', canvas.value.width, 'x', canvas.value.height)
}

// ========== 窗口大小变化处理 ==========
function handleResize() {
  if (currentPdf) {
    if (displayMode.value === 'virtual') {
      // 虚拟滚动模式：重新计算页面高度
      console.log('🔄 窗口大小变化，重新计算页面高度...')
      const containerWidth = container.value.clientWidth - 40
      calculatePageHeights(currentPdf, containerWidth).then(heights => {
        pageHeights.value = heights
        // 更新每个容器的高度
        pageContainers.value.forEach((pageDiv, index) => {
          const height = heights[index] || 800
          pageDiv.style.height = `${height}px`
          pageDiv.style.minHeight = `${height}px`
        })
        // 重新渲染已渲染的页面
        renderedPages.value.forEach(pageNum => {
          renderPageVirtual(pageNum)
        })
      })
    } else {
      // 分页模式：重新渲染当前页
      if (canvas.value) {
        console.log('🔄 窗口大小变化，重新渲染当前页...')
        renderPage(currentPdf, currentPageNum.value).catch(err => {
          console.error('重新渲染失败:', err)
        })
      }
    }
  }
}

// ========== SSE 事件监听 ==========
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

// ========== 后端连接检查 ==========
async function checkBackendConnection() {
  const { success, error: errorMsg } = await testBackendConnection()
  if (!success) {
    error.value = errorMsg
    return false
  }
  return true
}

// ========== 生命周期 ==========
onMounted(async () => {
  console.log('🚀 PDFViewer 组件已挂载')
  
  // 监听窗口大小变化
  if (typeof window !== 'undefined') {
    window.addEventListener('resize', handleResize)
    // 使用 ResizeObserver 监听容器大小变化（更精确）
    await nextTick()
    if (container.value) {
      resizeObserver = new ResizeObserver(() => {
        handleResize()
      })
      resizeObserver.observe(container.value)
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
  cleanup()
  
  // 移除监听器
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', handleResize)
  }
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
})
</script>

<template>
  <div :class="styles.viewer">
    <!-- 加载提示 -->
    <div v-if="loading && isInitialLoad" :class="styles.loading">
      正在加载 PDF...
    </div>
    
    <!-- 错误提示 -->
    <div v-if="error" :class="styles.error">{{ error }}</div>
    
    <!-- PDF 容器 -->
    <div 
      v-show="!error" 
      ref="container"
      :class="[
        styles.canvasContainer,
        displayMode === 'virtual' ? styles.multiPage : styles.singlePage
      ]"
    >
      <!-- 页面会通过 JavaScript 动态创建 -->
    </div>
  </div>
</template>
