<script setup>
import { computed, onMounted } from 'vue'
import { useAppModeStore } from './stores/appMode'
import LaziViewLayout from './layouts/LaziViewLayout.vue'
import LaziWorkspaceLayout from './layouts/LaziWorkspaceLayout.vue'

// 获取 store 实例
const appModeStore = useAppModeStore()

// 根据模式选择对应的布局组件
const currentLayout = computed(() => {
  return appModeStore.isViewMode ? LaziViewLayout : LaziWorkspaceLayout
})

// 初始化模式（从后端 API 或 URL 参数获取）
onMounted(async () => {
  await appModeStore.initMode()
})
</script>

<template>
    <!-- 根据模式渲染对应的布局组件 -->
    <component :is="currentLayout" />
</template>
