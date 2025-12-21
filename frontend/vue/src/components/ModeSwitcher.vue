<script setup>
import { computed } from 'vue'
import { useAppModeStore, APP_MODE } from '../stores/appMode'

const appModeStore = useAppModeStore()

// 根据当前模式显示对应的切换按钮
const switchLabel = computed(() => {
  return appModeStore.isViewMode ? '切换到工作区' : '切换到预览'
})

const switchTargetMode = computed(() => {
  return appModeStore.isViewMode ? APP_MODE.WORKSPACE : APP_MODE.VIEW
})

// 切换模式
function switchMode() {
  appModeStore.setMode(switchTargetMode.value)
}
</script>

<template>
  <button class="mode-switch-btn" @click="switchMode">
    {{ switchLabel }}
  </button>
</template>

<style scoped>
.mode-switch-btn {
  padding: 8px 16px;
  background: #5E81AC; /* nord10 - 按钮背景 */
  color: #ECEFF4; /* nord6 - 按钮文字 */
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.2s;
  font-family: inherit;
}

.mode-switch-btn:hover {
  background: #81A1C1; /* nord9 - 悬停时变亮 */
}

.mode-switch-btn:active {
  background: #5E81AC; /* nord10 - 点击时恢复 */
}
</style>

