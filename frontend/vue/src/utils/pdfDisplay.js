// utils/pdfDisplay.js
// PDF 显示模式管理工具（虚拟滚动、分页导航、配置管理）

import { renderPageToCanvas, calculatePageHeights } from './pdf'
import { getConfig, setConfig } from '../api/config'

/**
 * 创建虚拟滚动的页面容器
 * @param {HTMLElement} container - 容器元素
 * @param {number} totalPages - 总页数
 * @param {Array<number>} pageHeights - 每页高度数组
 * @param {string} pageContainerClass - 页面容器 CSS 类名
 * @param {string} placeholderClass - 占位符 CSS 类名
 * @returns {Array<HTMLElement>} 页面容器数组
 */
export function createVirtualScrollContainers(container, totalPages, pageHeights, pageContainerClass, placeholderClass) {
    container.innerHTML = ''
    const containers = []

    for (let i = 1; i <= totalPages; i++) {
        const pageDiv = document.createElement('div')
        pageDiv.className = pageContainerClass
        pageDiv.dataset.pageNum = i

        // 设置占位高度
        const height = pageHeights[i - 1] || 800
        pageDiv.style.height = `${height}px`
        pageDiv.style.minHeight = `${height}px`

        // 创建占位符
        const placeholder = document.createElement('div')
        placeholder.className = placeholderClass
        placeholder.style.width = '100%'
        placeholder.style.height = '100%'
        placeholder.textContent = `页面 ${i}`
        pageDiv.appendChild(placeholder)

        container.appendChild(pageDiv)
        containers.push(pageDiv)
    }

    return containers
}

/**
 * 设置 Intersection Observer（虚拟滚动用）
 * @param {HTMLElement} container - 滚动容器
 * @param {Array<HTMLElement>} pageContainers - 页面容器数组
 * @param {Function} onPageVisible - 页面可见时的回调函数 (pageNum) => void
 * @param {number} rootMargin - 预加载边距，默认 200px
 * @returns {IntersectionObserver} 观察器实例
 */
export function setupVirtualScrollObserver(container, pageContainers, onPageVisible, rootMargin = 200) {
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
        const pageNum = parseInt(entry.target.dataset.pageNum)
        if (entry.isIntersecting) {
            onPageVisible(pageNum)
        }
        })
    }, {
        root: container,
        rootMargin: `${rootMargin}px`,
        threshold: 0
    })

    // 观察所有页面容器
    pageContainers.forEach(pageDiv => {
        observer.observe(pageDiv)
    })

    return observer
}

/**
 * 渲染虚拟滚动中的单个页面
 * @param {Object} pdf - PDF 对象
 * @param {number} pageNum - 页码
 * @param {HTMLElement} pageDiv - 页面容器元素
 * @param {number} containerWidth - 容器宽度
 * @param {string} placeholderClass - 占位符 CSS 类名
 * @param {string} canvasClass - Canvas CSS 类名
 * @returns {Promise<HTMLCanvasElement>} 渲染后的 canvas 元素
 */
export async function renderVirtualScrollPage(pdf, pageNum, pageDiv, containerWidth, placeholderClass, canvasClass) {
    // 创建 canvas
    const canvas = document.createElement('canvas')
    canvas.className = canvasClass

    // 渲染页面
    await renderPageToCanvas(pdf, pageNum, canvas, containerWidth)

    // 移除所有占位符
    const placeholders = pageDiv.querySelectorAll(`.${placeholderClass}`)
    placeholders.forEach(placeholder => {
        pageDiv.removeChild(placeholder)
    })

    // 移除可能存在的旧 canvas
    const oldCanvas = pageDiv.querySelector(`canvas.${canvasClass}`)
    if (oldCanvas) {
        pageDiv.removeChild(oldCanvas)
    }

    // 添加新的 canvas
    pageDiv.appendChild(canvas)

    return canvas
}

/**
 * 初始化虚拟滚动模式
 * @param {Object} pdf - PDF 对象
 * @param {HTMLElement} container - 容器元素
 * @param {Function} onPageVisible - 页面可见回调
 * @param {Object} styleClasses - 样式类名对象 {pageContainer, placeholder, canvas}
 * @returns {Promise<{pageHeights: Array<number>, containers: Array<HTMLElement>, observer: IntersectionObserver}>}
 */
export async function initVirtualScrollMode(pdf, container, onPageVisible, styleClasses) {
    const containerWidth = container.clientWidth - 40

    // 计算所有页面的高度
    const pageHeights = await calculatePageHeights(pdf, containerWidth)

    // 创建页面容器
    const containers = createVirtualScrollContainers(
        container,
        pdf.numPages,
        pageHeights,
        styleClasses.pageContainer,
        styleClasses.placeholder
    )

    // 设置观察器
    const observer = setupVirtualScrollObserver(container, containers, onPageVisible)

    return {
        pageHeights,
        containers,
        observer
    }
}

/**
 * 初始化分页导航模式
 * @param {Object} pdf - PDF 对象
 * @param {HTMLElement} container - 容器元素
 * @param {Object} styleClasses - 样式类名对象 {pageContainer, canvas}
 * @returns {Promise<{canvas: HTMLCanvasElement, pageDiv: HTMLElement}>}
 */
export async function initPaginationMode(pdf, container, styleClasses) {
    container.innerHTML = ''
    
    // 创建单页容器
    const pageDiv = document.createElement('div')
    pageDiv.className = styleClasses.pageContainer
    pageDiv.dataset.pageNum = 1

    // 创建 canvas
    const canvasEl = document.createElement('canvas')
    canvasEl.className = styleClasses.canvas

    const containerWidth = container.clientWidth - 40
    await renderPageToCanvas(pdf, 1, canvasEl, containerWidth)

    pageDiv.appendChild(canvasEl)
    container.appendChild(pageDiv)

    return {
        canvas: canvasEl,
        pageDiv: pageDiv
    }
}

// ========== 配置管理 ==========

/**
 * 加载显示模式配置
 * @param {Function} setDisplayMode - 设置显示模式的函数 (mode: string) => void
 * @param {string} defaultMode - 默认模式，默认 'virtual'
 * @returns {Promise<string>} 加载的显示模式
 */
export async function loadDisplayModeConfig(setDisplayMode, defaultMode = 'virtual') {
    try {
        const config = await getConfig()
        if (config.displayMode && (config.displayMode === 'virtual' || config.displayMode === 'pagination')) {
            setDisplayMode(config.displayMode)
            console.log('✅ 从后端加载显示模式:', config.displayMode)
            return config.displayMode
        }
    } catch (err) {
        console.warn('⚠️ 加载配置失败（使用默认模式）:', err)
    }
    
    // 使用默认值
    setDisplayMode(defaultMode)
    return defaultMode
}

/**
 * 保存显示模式配置
 * @param {string} displayMode - 显示模式
 * @returns {Promise<void>}
 */
export async function saveDisplayModeConfig(displayMode) {
    try {
        await setConfig({ displayMode })
        console.log('✅ 显示模式已保存到配置')
    } catch (err) {
        console.warn('⚠️ 保存显示模式配置失败（后端可能尚未支持）:', err)
        // 静默失败，不影响用户体验
    }
}

/**
 * 切换显示模式
 * @param {string} newMode - 新模式
 * @param {string} currentMode - 当前模式
 * @param {Function} setMode - 设置模式的函数
 * @param {Function} reinitialize - 重新初始化的函数 (pdf, mode) => Promise<void>
 * @param {Object} pdf - PDF 对象（如果已加载）
 * @returns {Promise<void>}
 */
export async function switchDisplayMode(newMode, currentMode, setMode, reinitialize, pdf = null) {
    if (newMode === currentMode) {
        return // 已经是该模式，无需切换
    }

    // 如果 PDF 已加载，需要重新初始化
    if (pdf) {
        setMode(newMode)
        await reinitialize(pdf, newMode)
    } else {
        // PDF 未加载，直接切换模式
        setMode(newMode)
    }

    // 保存配置到后端
    await saveDisplayModeConfig(newMode)
}

/**
 * 渲染分页导航模式中的页面
 * @param {Object} pdf - PDF 对象
 * @param {number} pageNum - 页码
 * @param {HTMLCanvasElement} canvas - Canvas 元素
 * @param {number} containerWidth - 容器宽度
 * @returns {Promise<void>}
 */
export async function renderPaginationPage(pdf, pageNum, canvas, containerWidth) {
    if (!pdf || pageNum < 1 || pageNum > pdf.numPages) {
        throw new Error(`无效的页码: ${pageNum}`)
    }

    if (!canvas) {
        throw new Error('Canvas 元素未提供')
    }

    await renderPageToCanvas(pdf, pageNum, canvas, containerWidth)
}