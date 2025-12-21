// stores/appMode.js
// 应用模式状态管理（预览模式或工作区模式）

import {defineStore} from 'pinia'
import {ref, computed} from 'vue'
import { getConfig } from '../api/config'

// ========== 常量定义 ==========
// 模式常量
export const APP_MODE = {
    VIEW: 'view',
    WORKSPACE: 'workspace'
}

// ========== Store 定义 ==========
export const useAppModeStore = defineStore('appMode', () => {
    // 当前模式：初始值由后端决定（通过 initMode 初始化）
    // 默认值设为 VIEW，如果初始化失败则使用此默认值
    const currentMode = ref(APP_MODE.VIEW)

    // 计算属性：是否预览模式
    const isViewMode = computed(() => currentMode.value === APP_MODE.VIEW)

    // 计算属性：是否工作区模式
    const isWorkspaceMode = computed(() => currentMode.value === APP_MODE.WORKSPACE)

    // 设置当前模式（前端切换模式使用）
    function setMode(mode) {
        const validModes = [APP_MODE.VIEW, APP_MODE.WORKSPACE]
        if (validModes.includes(mode)) {
            currentMode.value = mode
        } else {
            console.warn(`无效的模式值: ${mode}，必须是 ${validModes.join(' 或 ')}`)
        }
    }

    // 初始化模式（从 URL 参数或后端 API 获取）
    async function initMode() {
        // 方式1：从后端 API 获取模式（主要方式）
        // 后端lazitex -p 参数会返回 {mode: "view"}
        // 后端lazitex -w 参数会返回 {mode: "workspace"}
        try {
            const config = await getConfig()
            if (config.mode && [APP_MODE.VIEW, APP_MODE.WORKSPACE].includes(config.mode)) {
                currentMode.value = config.mode
                console.log('✅ 从后端获取模式:', config.mode)
                return // ✅ 立即退出函数，不执行后面的代码
            }
        } catch (error) {
            // API 不存在或失败，使用备用方案（URL 参数）
            console.log('⚠️ 无法从后端获取模式配置，尝试使用 URL 参数')
        }
        
        // 方式2：从 URL 参数获取（备用方案，开发时方便测试）
        const urlParams = new URLSearchParams(window.location.search)
        const modeFromUrl = urlParams.get('mode')

        if (modeFromUrl === APP_MODE.WORKSPACE || modeFromUrl === 'hub') {
            currentMode.value = APP_MODE.WORKSPACE
            console.log('✅ 从 URL 参数获取模式: workspace')
        } else if (modeFromUrl === APP_MODE.VIEW) {
            currentMode.value = APP_MODE.VIEW
            console.log('✅ 从 URL 参数获取模式: view')
        } else {
            // 如果都没有，使用默认值 VIEW
            console.warn('❌ URL 参数未指定模式，使用默认模式: view')
        }
    }

    // ========== 返回 Store 的公开接口 ==========
    // return 返回的是"外部可以访问的内容"
    // 只有 return 的内容，其他组件才能使用
    return {
        // 状态（响应式变量）
        currentMode,
        // 计算属性（只读的计算值）
        isViewMode,
        isWorkspaceMode,
        // 方法（函数）
        setMode,    // 前端切换模式（运行时切换）
        initMode    // 初始化模式（由后端决定）
    }
})