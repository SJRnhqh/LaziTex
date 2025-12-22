<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
// 导入 CSS Modules
import styles from '../styles/dynamic-pdf-viewer.module.css'

// 在顶层导入 PDF.js（避免动态导入问题）
import * as pdfjsLib from 'pdfjs-dist'

// 导入 PDF 工具函数
import { 
  parsePDF, 
  cleanupPDF, 
  formatPDFError,
  renderPageToCanvas,
  calculatePageHeights
} from '../utils/pdf'

// 导入 PDF 显示模式工具函数
import {
  initVirtualScrollMode,
  renderVirtualScrollPage,
  initPaginationMode,
  renderPaginationPage,
  loadDisplayModeConfig,
  switchDisplayMode as switchModeUtil
} from '../utils/pdfDisplay'

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
const displayMode = ref('virtual') // 默认虚拟滚动，会从后端配置加载

// 分页模式相关
const canvas = ref(null) // 分页模式需要
const currentPageNum = ref(1) // 当前页码（分页模式）
const totalPages = ref(0) // 总页数（分页模式需要）
const pageInput = ref('') // 页码输入框的值

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
    
    // 保存总页数（两种模式都需要）
    totalPages.value = pdf.numPages

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

  // 重置状态
  pageContainers.value = []
  renderedPages.value.clear()

  // 页面可见时的回调函数
  const onPageVisible = (pageNum) => {
    if (!renderedPages.value.has(pageNum)) {
      renderPageVirtual(pageNum)
    }
  }

  // 使用工具函数初始化虚拟滚动
  const { pageHeights: heights, containers, observer } = await initVirtualScrollMode(
    pdf,
    container.value,
    onPageVisible,
    {
      pageContainer: styles.pageContainer,
      placeholder: styles.pagePlaceholder,
      canvas: styles.canvas
    }
  )

  // 保存状态
  pageHeights.value = heights
  pageContainers.value = containers
  pageObserver = observer

  // 初始渲染可见页面（第一页）
  if (pageContainers.value.length > 0) {
    await renderPageVirtual(1)
  }
  
  console.log('✨ 虚拟滚动初始化完成')
}

// 渲染单个页面（虚拟滚动版）
async function renderPageVirtual(pageNum) {
  if (!currentPdf || renderedPages.value.has(pageNum)) {
    return
  }

  // 立即标记为正在渲染，防止重复渲染（在异步操作之前）
  renderedPages.value.add(pageNum)
  
  try {
    const pageDiv = pageContainers.value[pageNum - 1]
    if (!pageDiv) {
      renderedPages.value.delete(pageNum) // 如果失败，移除标记
      return
    }
    
    // 获取容器宽度
    const containerWidth = container.value.clientWidth - 40
    
    // 使用工具函数渲染页面
    await renderVirtualScrollPage(
      currentPdf,
      pageNum,
      pageDiv,
      containerWidth,
      styles.pagePlaceholder,
      styles.canvas
    )
    
    console.log(`✅ 页面 ${pageNum} 渲染完成`)
  } catch (err) {
    console.error(`❌ 渲染页面 ${pageNum} 失败:`, err)
    renderedPages.value.delete(pageNum) // 如果失败，移除标记，允许重试
  }
}

// ========== 分页导航模式 ==========
async function initPagination(pdf) {
  console.log('🎨 初始化分页导航模式...')
  
  await nextTick()
  if (!container.value) {
    throw new Error('容器元素尚未挂载')
  }

  // 保存总页数
  totalPages.value = pdf.numPages
  currentPageNum.value = 1
  pageInput.value = '1'

  // 使用工具函数初始化分页导航
  const { canvas: canvasEl } = await initPaginationMode(
    pdf,
    container.value,
    {
      pageContainer: styles.pageContainer,
      canvas: styles.canvas
    }
  )

  // 保存 canvas 引用（用于后续操作）
  canvas.value = canvasEl

  console.log('✨ 分页导航模式初始化完成')
}

// 渲染 PDF 页面（分页模式使用）
async function renderPage(pdf, pageNum) {
  if (!pdf || pageNum < 1 || pageNum > pdf.numPages) {
    return
  }

  await nextTick()
  if (!canvas.value) {
    throw new Error('Canvas 元素尚未挂载，无法渲染 PDF')
  }
  
  const containerWidth = canvas.value.parentElement 
    ? canvas.value.parentElement.clientWidth - 40 
    : 800
  
  // 使用工具函数渲染页面
  await renderPaginationPage(pdf, pageNum, canvas.value, containerWidth)
  currentPageNum.value = pageNum
  pageInput.value = pageNum.toString()
  
  console.log(`✅ PDF 页面 ${pageNum} 渲染完成`)
}

// 分页导航：上一页
function goToPreviousPage() {
  if (currentPdf && currentPageNum.value > 1) {
    renderPage(currentPdf, currentPageNum.value - 1).catch(err => {
      console.error('跳转到上一页失败:', err)
    })
  }
}

// 分页导航：下一页
function goToNextPage() {
  if (currentPdf && currentPageNum.value < totalPages.value) {
    renderPage(currentPdf, currentPageNum.value + 1).catch(err => {
      console.error('跳转到下一页失败:', err)
    })
  }
}

// 分页导航：跳转到指定页
function goToPage(pageNum) {
  if (!currentPdf) return
  
  const targetPage = parseInt(pageNum)
  if (targetPage >= 1 && targetPage <= totalPages.value) {
    renderPage(currentPdf, targetPage).catch(err => {
      console.error('跳转到指定页失败:', err)
    })
  } else {
    // 如果输入无效，恢复当前页码
    pageInput.value = currentPageNum.value.toString()
  }
}

// 处理页码输入框回车
function handlePageInputEnter(event) {
  if (event.key === 'Enter') {
    goToPage(pageInput.value)
  }
}

// ========== 模式切换 ==========
async function switchDisplayMode(newMode) {
  await switchModeUtil(
    newMode,
    displayMode.value,
    (mode) => { displayMode.value = mode },
    async (pdf, mode) => {
      // 重新初始化函数
      if (mode === 'virtual') {
        await initVirtualScroll(pdf)
      } else {
        await initPagination(pdf)
      }
    },
    currentPdf
  )
}

// 加载配置
async function loadConfig() {
  await loadDisplayModeConfig(
    (mode) => { displayMode.value = mode },
    'virtual'
  )
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
    console.log('✅ 后端连接正常')
    
    // 加载配置（在加载 PDF 之前）
    await loadConfig()
    
    // 开始加载 PDF
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

    <!-- 模式切换控件 -->
    <div v-if="!error && totalPages > 0" :class="styles.modeControls">
      <button
        :class="[styles.modeButton, displayMode === 'virtual' ? styles.modeButtonActive : '']"
        @click="switchDisplayMode('virtual')"
        title="虚拟滚动模式"
      >
        📄 虚拟滚动
      </button>
      <button
        :class="[styles.modeButton, displayMode === 'pagination' ? styles.modeButtonActive : '']"
        @click="switchDisplayMode('pagination')"
        title="分页导航模式"
      >
        📑 分页导航
      </button>
    </div>

    <!-- 分页导航控件（仅在分页模式下显示） -->
    <div v-if="!error && displayMode === 'pagination' && totalPages > 0" :class="styles.paginationControls">
      <button 
        :class="styles.paginationButton"
        :disabled="currentPageNum <= 1"
        @click="goToPreviousPage"
      >
        ← 上一页
      </button>
      
      <div :class="styles.pageInfo">
        <input
          :class="styles.pageInput"
          type="number"
          :min="1"
          :max="totalPages"
          v-model="pageInput"
          @keyup="handlePageInputEnter"
          @blur="goToPage(pageInput)"
        />
        <span :class="styles.pageSeparator">/</span>
        <span :class="styles.totalPages">{{ totalPages }}</span>
      </div>
      
      <button 
        :class="styles.paginationButton"
        :disabled="currentPageNum >= totalPages"
        @click="goToNextPage"
      >
        下一页 →
      </button>
    </div>
  </div>
</template>
