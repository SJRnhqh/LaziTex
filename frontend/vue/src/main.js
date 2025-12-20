import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './style.css'

const app = createApp(App)

// 使用 Pinia（状态管理）
app.use(createPinia())

app.mount('#app')
