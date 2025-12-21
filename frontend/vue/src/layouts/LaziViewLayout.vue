<!-- layouts/LaziViewLayout.vue -->
<!-- LaziViewLayout.vue 是 LaziTeX 的视图布局组件，负责管理视图的布局和状态 -->

<script setup>
import DynamicPDFViewer from '../components/DynamicPDFViewer.vue'
import ModeSwitcher from '../components/ModeSwitcher.vue'
</script>

<template>
    <div class="laziview-layout">
        <!-- 极简标题栏 -->
        <header>
            <h1>LaziHub</h1>
            <ModeSwitcher />
        </header>
        <!-- 预览模式：全屏 PDF 预览 -->
        <DynamicPDFViewer />
    </div>
</template>

<style scoped>
.laziview-layout {
    width: 100vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: #2E3440; /* nord0 - 主背景 */
    overflow: hidden;
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
</style>

