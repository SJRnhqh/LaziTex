import { createApp } from 'vue'// 导入Vue
import { createPinia } from 'pinia'// 导入Pinia
import App from './App.vue'// 导入根组件

const app = createApp(App)// 创建应用

// 使用 Pinia（状态管理）
app.use(createPinia())// 使用Pinia

app.mount('#lazihub')// 挂载应用到ID为lazihub的挂载点
