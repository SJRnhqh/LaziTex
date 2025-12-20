<script setup>
import PDFViewer from './components/PDFViewer.vue'
</script>

<template>
  <div id="app">
    <header>
      <h1>LaziTex 编辑器</h1>
      <p>本地 LaTeX 编译工具</p>
    </header>
    <main>
      <PDFViewer />
    </main>
  </div>
</template>

<style scoped>
#app {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
  color: #fff;
}

header {
  padding: 20px;
  background: #2d2d2d;
  border-bottom: 1px solid #444;
}

header h1 {
  margin: 0;
  font-size: 24px;
}

header p {
  margin: 5px 0 0 0;
  color: #aaa;
  font-size: 14px;
}

main {
  flex: 1;
  overflow: hidden;
}
</style>