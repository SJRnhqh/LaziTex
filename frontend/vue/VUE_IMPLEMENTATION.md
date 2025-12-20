# Vue 前端实现文档

本文档详细说明 LaziTex Vue 前端的实现思路和与后端的交互方式。

## 📋 目录

- [架构概述](#架构概述)
- [技术栈](#技术栈)
- [核心组件](#核心组件)
- [与后端交互](#与后端交互)
- [关键实现细节](#关键实现细节)

---

## 架构概述

### 整体架构

```txt
┌─────────────────────────────────────────────────────────┐
│                    Vue 前端 (Vite)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   PDFViewer  │  │   App.vue    │  │   main.js    │ │
│  │  Component   │  │              │  │              │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────┘
                         ↕ HTTP/SSE
┌─────────────────────────────────────────────────────────┐
│              Go 后端 (HTTP Server)                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │ /pdf     │  │ /events  │  │ /        │             │
│  │ Route    │  │ Route    │  │ Route    │             │
│  └──────────┘  └──────────┘  └──────────┘             │
└─────────────────────────────────────────────────────────┘
                         ↕ File System
                    LaTeX Source Files
```

### 开发模式 vs 生产模式

**开发模式：**

- Vue 前端运行在 `localhost:5173` (Vite Dev Server)
- Go 后端运行在 `localhost:8080`
- Vite 代理配置将 `/pdf`、`/events` 等请求转发到后端

**生产模式：**

- Vue 前端构建后通过 Go `embed` 嵌入到二进制文件
- 单一二进制文件，前端资源通过 Go HTTP 服务器提供

---

## 技术栈

### 核心依赖

- **Vue 3** - 渐进式 JavaScript 框架
- **Vite** - 下一代前端构建工具
- **Pinia** - Vue 官方状态管理库
- **PDF.js** - PDF 渲染库（Mozilla）

### 开发工具

- **Node.js** - JavaScript 运行时
- **npm** - 包管理器

---

## 核心组件

### PDFViewer.vue

PDF 预览的核心组件，负责：

1. PDF 文件的加载和渲染
2. 实时更新监听（通过 SSE）
3. 窗口大小自适应
4. 平滑刷新体验

#### 关键功能实现

##### 1. PDF.js Worker 配置

```javascript
// 使用 Vite 的 ?url 导入语法获取 worker 文件 URL
import workerUrl from 'pdfjs-dist/build/pdf.worker.min.mjs?url'
pdfjsLib.GlobalWorkerOptions.workerSrc = workerUrl
```

- 使用 `?url` 后缀让 Vite 正确处理 worker 文件
- 避免动态导入导致的模块加载问题

##### 2. PDF 加载流程

```javascript
async function loadPDF(showLoading = true) {
  // 1. 通过 fetch 下载 PDF 为 Blob
  const response = await fetch('/pdf?t=' + Date.now())
  const pdfBlob = await response.blob()
  
  // 2. 转换为 ArrayBuffer
  const pdfArrayBuffer = await pdfBlob.arrayBuffer()
  
  // 3. 使用 PDF.js 解析
  const loadingTask = pdfjsLib.getDocument({ data: pdfArrayBuffer })
  const pdf = await loadingTask.promise
  
  // 4. 渲染到 Canvas
  await renderPage(pdf, 1)
}
```

##### 3. 自适应缩放

```javascript
// 根据容器宽度计算缩放比例
const containerWidth = container.clientWidth - 40
const scale = Math.min(1.5, (containerWidth / baseViewport.width) * 0.95)
const viewport = page.getViewport({ scale: scale })
```

##### 4. 平滑刷新

- 首次加载显示 loading 提示
- 刷新时保持旧内容可见，后台渲染新内容
- 使用 `v-show` 而非 `v-if` 保持 DOM 元素存在

##### 5. 窗口大小监听

```javascript
// 使用 ResizeObserver 监听容器大小变化
resizeObserver = new ResizeObserver(() => {
  handleResize() // 重新渲染 PDF
})
```

---

## 与后端交互

### HTTP 路由

后端提供以下 HTTP 路由：

#### 1. `/pdf` - PDF 文件服务

**请求：**

```
GET /pdf?t=1234567890
```

**响应：**

- Content-Type: `application/pdf`
- 返回 PDF 文件的二进制内容

**实现位置：**

- 后端：`backend/handlers.go` → `Server.HandlePDF()`
- 前端：`PDFViewer.vue` → `loadPDF()` 函数

#### 2. `/events` - Server-Sent Events (SSE)

**请求：**

```
GET /events
```

**响应：**

- Content-Type: `text/event-stream`
- 保持长连接，推送事件

**事件格式：**

```txt
data: refresh\n\n
```

**实现位置：**

- 后端：`backend/sse.go` → `Server.HandleSSE()`
- 前端：`PDFViewer.vue` → `setupSSE()` 函数

**工作流程：**

1. **前端建立连接：**

```javascript
eventSource = new EventSource('/events')
eventSource.onmessage = async (event) => {
  if (event.data === 'refresh') {
    await loadPDF(false) // 刷新 PDF，不显示 loading
  }
}
```

2. **后端监听文件变化：**

```go
// cmd/lazitex-cli/tasks/preview.go
core.WatchAndAction(absPath, func() {
    // 重新编译 LaTeX
    newPdfPath, err := core.Build(opts)
    if newPdfPath != "" {
        server.SetPDFPath(newPdfPath)
        server.BroadcastSSE("refresh") // 广播刷新消息
    }
})
```

3. **后端广播消息：**

```go
// backend/sse.go
func (s *Server) BroadcastSSE(message string) {
    // 遍历所有连接的客户端，发送消息
    for clientChan := range s.sseClients {
        clientChan <- message
    }
}
```

### 代理配置

在开发模式下，Vite 配置代理将前端请求转发到后端：

```javascript
// vite.config.js
server: {
  proxy: {
    '/pdf': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
    '/events': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      ws: true, // WebSocket 支持
    },
  },
}
```

---

## 关键实现细节

### 1. PDF.js Worker 配置问题解决

**问题：** PDF.js 5.x 版本在使用动态导入时，Vite 会拦截 worker 请求并添加 `?import` 参数，导致 worker 加载失败。

**解决方案：** 使用 Vite 的 `?url` 导入语法：

```javascript
import workerUrl from 'pdfjs-dist/build/pdf.worker.min.mjs?url'
pdfjsLib.GlobalWorkerOptions.workerSrc = workerUrl
```

### 2. Canvas 渲染优化

**避免黑屏：**

- 使用 `v-show` 而非 `v-if`，保持 canvas 元素始终存在
- 刷新时不显示 loading，保持旧内容可见

**资源清理：**

```javascript
// 清理旧的 PDF 对象
if (currentPdf) {
  currentPdf.destroy()
}

// 组件卸载时清理
onUnmounted(() => {
  if (currentPdf) currentPdf.destroy()
  if (eventSource) eventSource.close()
  if (resizeObserver) resizeObserver.disconnect()
})
```

### 3. 自适应和滚动

**实现：**

- Canvas 容器设置 `overflow: auto` 支持滚动
- 根据容器宽度动态计算缩放比例
- 使用 `ResizeObserver` 监听容器大小变化，自动重新渲染

**CSS 关键样式：**

```css
.canvas-container {
  flex: 1;
  overflow: auto;
  scroll-behavior: smooth;
}

canvas {
  max-width: 100%;
  height: auto;
}
```

---

## 文件结构

```txt
frontend/vue/
├── public/
│   └── pdf.worker.min.mjs    # PDF.js worker 文件（复制自 node_modules）
├── src/
│   ├── components/
│   │   └── PDFViewer.vue     # PDF 预览组件
│   ├── App.vue               # 根组件
│   ├── main.js               # 入口文件
│   └── style.css             # 全局样式
├── index.html                # HTML 模板
├── vite.config.js            # Vite 配置
└── package.json              # 依赖配置
```

---

## 开发流程

### 1. 启动开发环境

**终端 1 - Vue 前端：**

```bash
cd frontend/vue
npm install
npm run dev
# 运行在 http://localhost:5173
```

**终端 2 - Go 后端：**

```bash
go run ./cmd/lazitex-cli -p test/simple.tex
# 运行在 http://localhost:8080
```

### 2. 构建生产版本

```bash
cd frontend/vue
npm run build
# 构建输出到 frontend/dist/
```

### 3. 集成到 Go 二进制

（待实现）使用 Go `embed` 将构建后的前端文件嵌入到二进制中。

---

## 未来计划

### 待实现功能

1. **文件管理模块**
   - 文件列表/树形结构
   - 文件创建/删除/重命名
   - 文件切换

2. **代码编辑器（Monaco Editor）**
   - LaTeX 语法高亮
   - 代码补全
   - 错误提示
   - 多文件编辑

3. **终端集成（xterm.js）**
   - 内置终端
   - 执行 LaTeX 命令
   - 查看编译日志

4. **AI 聊天框**
   - 集成 Ollama
   - AI 辅助写作
   - 错误诊断建议

5. **UI 布局优化**
   - 左右分栏（编辑器 + PDF 预览）
   - 可调整面板大小
   - 响应式布局

---

## 参考资料

- [Vue 3 文档](https://vuejs.org/)
- [Vite 文档](https://vitejs.dev/)
- [PDF.js 文档](https://mozilla.github.io/pdf.js/)
- [Server-Sent Events (SSE)](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)

