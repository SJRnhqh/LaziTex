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