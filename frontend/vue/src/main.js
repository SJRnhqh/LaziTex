// main.js
// main.js 是 Vue 应用的入口文件，负责创建应用实例并挂载到 DOM 上
import { createApp } from 'vue'// 导入Vue
import { createPinia } from 'pinia'// 导入Pinia
import App from './App.vue'// 导入根组件

const app = createApp(App)// 创建应用

// 使用 Pinia（状态管理）
app.use(createPinia())// 使用Pinia

app.mount('#lazihub')// 挂载应用到ID为lazihub的挂载点
