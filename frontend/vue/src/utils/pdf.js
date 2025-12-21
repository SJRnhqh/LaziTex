// utils/pdf.js
// pdf.js 是 LaziTeX 的 PDF 渲染工具，负责管理 PDF 的渲染和更新

import * as pdfjsLib from 'pdfjs-dist'

/**
 * 解析 PDF 文件
 * @param {ArrayBuffer} arrayBuffer - PDF 文件的 ArrayBuffer
 * @param {number} timeout - 超时时间（毫秒），默认 60 秒
 * @returns {Promise<Object>} PDF 对象
 */
export async function parsePDF(arrayBuffer, timeout = 60000) {
    const loadingTask = pdfjsLib.getDocument({
        data: arrayBuffer,
        verbosity: 1,
    })
    
    const parseTimeoutPromise = new Promise((_, reject) => {
        setTimeout(() => {
            loadingTask.destroy()
            reject(new Error(`PDF 解析超时（${timeout / 1000}秒）`))
        }, timeout)
    })
    
    return await Promise.race([
        loadingTask.promise,
        parseTimeoutPromise
    ])
}

/**
 * 计算 PDF 缩放比例
 * @param {number} containerWidth - 容器宽度
 * @param {number} pdfWidth - PDF 页面宽度
 * @param {number} maxScale - 最大缩放比例，默认 1.5
 * @param {number} padding - 容器内边距，默认 40
 * @returns {number} 缩放比例
 */
export function calculateScale(containerWidth, pdfWidth, maxScale = 1.5, padding = 40) {
    const availableWidth = containerWidth - padding
    return Math.min(maxScale, (availableWidth / pdfWidth) * 0.95)
}

/**
 * 清理 PDF 对象
 * @param {Object} pdf - PDF 对象
 */
export function cleanupPDF(pdf) {
    if (pdf) {
        try {
            console.log('🧹 清理 PDF 对象...')
            pdf.destroy()
        } catch (error) {
            console.warn('❌ 清理 PDF 时出错:', error)
        }
    }
}

/**
 * 格式化 PDF 错误信息
 * @param {Error} err - 错误对象
 * @returns {string} 格式化的错误信息
 */
export function formatPDFError(err) {
    return `PDF 加载失败: ${err.message || err.toString()}`
}

/**
 * 渲染单个 PDF 页面到 canvas
 * @param {Object} pdf - PDF 对象
 * @param {number} pageNum - 页码（从 1 开始）
 * @param {HTMLCanvasElement} canvas - Canvas 元素
 * @param {number} containerWidth - 容器宽度
 * @returns {Promise<{width: number, height: number}>} 渲染后的页面尺寸
 */
export async function renderPageToCanvas(pdf, pageNum, canvas, containerWidth) {
    const page = await pdf.getPage(pageNum)
    const baseViewport = page.getViewport({ scale: 1.0 })
    
    // 计算缩放比例
    const scale = calculateScale(containerWidth, baseViewport.width)
    const viewport = page.getViewport({ scale: scale })
    
    // 设置 canvas 尺寸
    canvas.height = viewport.height
    canvas.width = viewport.width
    
    // 获取 canvas 上下文
    const context = canvas.getContext('2d')
    
    // 填充白色背景
    context.fillStyle = '#ffffff'
    context.fillRect(0, 0, canvas.width, canvas.height)
    
    // 渲染 PDF 页面
    const renderContext = {
        canvasContext: context,
        viewport: viewport
    }
    
    await page.render(renderContext).promise
    
    return {
        width: viewport.width,
        height: viewport.height
    }
}

/**
 * 计算所有页面的高度（用于虚拟滚动占位）
 * @param {Object} pdf - PDF 对象
 * @param {number} containerWidth - 容器宽度
 * @returns {Promise<Array<number>>} 每页的高度数组
 */
export async function calculatePageHeights(pdf, containerWidth) {
    const heights = []
    for (let i = 1; i <= pdf.numPages; i++) {
        const page = await pdf.getPage(i)
        const baseViewport = page.getViewport({ scale: 1.0 })
        const scale = calculateScale(containerWidth, baseViewport.width)
        const viewport = page.getViewport({ scale: scale })
        heights.push(viewport.height)
    }
    return heights
}