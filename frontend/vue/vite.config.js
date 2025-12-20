import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  
  // 路径别名（方便导入）
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  
  // 优化配置（确保 PDF.js worker 正确加载）
  optimizeDeps: {
    include: ['pdfjs-dist'],
  },
  
  // 开发服务器配置
  server: {
    port: 5173,
    open: true,
    proxy: {
      // 代理 API 请求到 Go 后端
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // 代理 PDF 请求
      '/pdf': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // 代理 SSE 事件流
      '/events': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        ws: true, // WebSocket 支持
      },
    },
  },
  
  // 构建配置
  build: {
    outDir: '../dist',  // 输出到 frontend/dist/
    emptyOutDir: true,
    assetsDir: 'assets',
  },
})
