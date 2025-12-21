<!-- layouts/LaziWorkspaceLayout.vue -->
<!-- LaziWorkspaceLayout.vue 是 LaziTeX 的工作区布局组件，负责管理工作区的布局和状态 -->

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import DynamicPDFViewer from '../components/DynamicPDFViewer.vue'
import ModeSwitcher from '../components/ModeSwitcher.vue'

// 左侧代码区域宽度百分比（默认50%）
const codeWidthPercent = ref(50)

// 拖拽状态
const isResizing = ref(false)
const startX = ref(0)
const startCodeWidth = ref(0)

// 从 localStorage 加载保存的宽度
onMounted(() => {
    const savedWidth = localStorage.getItem('lazihub-code-width-percent')
    if (savedWidth) {
        const width = parseFloat(savedWidth)
        // 确保加载的值在有效范围内
        const minWidth = 10
        const maxWidth = 90
        codeWidthPercent.value = Math.max(minWidth, Math.min(maxWidth, width))
    }
})

// 保存宽度到 localStorage
function saveWidth() {
    localStorage.setItem('lazihub-code-width-percent', codeWidthPercent.value.toString())
}

// 开始调整
function startResize(e) {
    isResizing.value = true
    startX.value = e.clientX
    startCodeWidth.value = codeWidthPercent.value
    document.addEventListener('mousemove', handleResize)
    document.addEventListener('mouseup', stopResize)
    e.preventDefault()
}

// 调整宽度
function handleResize(e) {
    if (!isResizing.value) return
    const workspace = document.querySelector('.workspace')
    if (!workspace) return
    
    const workspaceWidth = workspace.clientWidth
    const diff = e.clientX - startX.value
    const diffPercent = (diff / workspaceWidth) * 100
    const newWidth = startCodeWidth.value + diffPercent
    
    // 范围限制：10% - 90%（确保两侧都有最小可见区域）
    const minWidth = 10
    const maxWidth = 90
    codeWidthPercent.value = Math.max(minWidth, Math.min(maxWidth, newWidth))
}

// 停止调整
function stopResize() {
    isResizing.value = false
    saveWidth()
    document.removeEventListener('mousemove', handleResize)
    document.removeEventListener('mouseup', stopResize)
}

// 清理事件监听器
onUnmounted(() => {
    document.removeEventListener('mousemove', handleResize)
    document.removeEventListener('mouseup', stopResize)
})
</script>

<template>
    <div class="laziworkspace-layout">
        <!-- 极简标题栏 -->
        <header>
            <h1>LaziHub</h1>
            <ModeSwitcher />
        </header>
        
        <!-- 主工作区：左右两列布局 -->
        <div class="workspace">
        <!-- 左侧：代码编辑器区域 -->
        <div class="code" :style="{ width: codeWidthPercent + '%' }">
            <!-- 代码编辑器将在这里 -->
        </div>
        
        <!-- 可调整的分隔条 -->
        <div 
            class="resize-handle"
            :class="{ 'resizing': isResizing }"
            @mousedown="startResize"
        ></div>
        
        <!-- 右侧：PDF 预览区域 -->
        <div class="pdf-preview" :style="{ width: (100 - codeWidthPercent) + '%' }">
            <DynamicPDFViewer />
        </div>
        </div>
    </div>
</template>

<style scoped>
/* Nord 配色方案 */
/* Polar Night (深色背景) */
/* nord0: #2E3440 - 最深背景 */
/* nord1: #3B4252 - 面板背景 */
/* nord2: #434C5E - 稍浅面板 */
/* nord3: #4C566A - 边框、分隔线 */
/* Snow Storm (浅色文字) */
/* nord4: #D8DEE9 - 主要文字 */
/* nord5: #E5E9F0 - 次要文字 */
/* nord6: #ECEFF4 - 最浅文字 */

.laziworkspace-layout {
    width: 100vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: #2E3440; /* nord0 - 主背景 */
    color: #D8DEE9; /* nord4 - 主要文字 */
}

header {
    padding: 20px;
    background: #3B4252; /* nord1 - header 背景 */
    border-bottom: 1px solid #4C566A; /* nord3 - 边框 */
    display: flex;
    justify-content: space-between;
    align-items: center;
}

header h1 {
    margin: 0;
    font-size: 24px;
    color: #ECEFF4; /* nord6 - 标题文字 */
}

/* 主工作区：左右两列布局 */
.workspace {
    flex: 1;
    display: flex;
    overflow: hidden;
}

/* 左侧：代码编辑器区域 */
.code {
    background: #2E3440; /* nord0 - 编辑器背景 */
    overflow: hidden;
    display: flex;
    flex-direction: column;
    flex-shrink: 0; /* 不允许缩小，由宽度控制 */
}

/* 可调整的分隔条 */
.resize-handle {
    width: 4px;
    background: #4C566A; /* nord3 - 分隔线颜色，与背景区分 */
    cursor: col-resize;
    flex-shrink: 0;
    position: relative;
    transition: background 0.2s;
    user-select: none;
    z-index: 10; /* 确保在其他元素之上 */
}

/* 增加可点击区域（视觉上仍然是 4px，但点击区域更大） */
.resize-handle::before {
    content: '';
    position: absolute;
    left: -2px;
    right: -2px;
    top: 0;
    bottom: 0;
    cursor: col-resize;
}

.resize-handle:hover {
    background: #5E81AC; /* nord10 - 悬停时变亮 */
}

.resize-handle.resizing {
    background: #81A1C1; /* nord9 - 拖拽时更亮 */
}

/* 右侧：PDF 预览区域 */
.pdf-preview {
    background: #3B4252; /* nord1 - 面板背景 */
    overflow: hidden;
    display: flex;
    flex-direction: column;
    flex-shrink: 0; /* 不允许缩小，由宽度控制 */
}
</style>

