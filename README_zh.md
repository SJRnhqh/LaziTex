# LaziTex 🧸

跨平台即时 LaTeX 编译工具 - 基于 Go 构建，集成本地 AI 助力写作与优化

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://github.com/SJRnhqh/lazitex)
[![Status](https://img.shields.io/badge/Status-Active%20Development-brightgreen)](https://github.com/SJRnhqh/lazitex)

中文文档 | [English](README.md)

---

## ✨ 特性

- 🔍 **智能环境检测** - 自动检测 LaTeX 安装
- 📦 **自动安装/更新** - 一键安装或更新 LaTeX 环境
- 🗑️ **自动卸载** - 一键卸载 LaTeX 环境
- 🌍 **跨平台支持** - 支持 Windows、Linux 和 macOS
- 🎯 **零配置启动** - 开箱即用
- 🚀 **快速轻量** - 基于 Go 构建
- 🔄 **自适应多轮编译** - 自动检测并处理多轮编译场景
- 📦 **自动包管理** - 自动检测并安装缺失的包
- 🌐 **Web 预览模式** - Web 实时预览
- 🦙 **Ollama 管理** - 一键管理 Ollama
- 🤖 **LLM 管理** - 智能 LLM 提供商管理和配置
- ⚙️ **配置管理** - 跨平台配置文件管理
- 🧠 **AI 架构** - 模块化 AI 智能体和工具系统 (即将推出)
- 🎨 **多种模式** - 支持 TUI 和 REPL 模式

---

## 📦 安装

### 构建依赖

在从源码构建之前，请确保已安装以下工具：

- **Node.js** (包含 npm)
  - macOS: `brew install node`
  - Linux: 使用发行版的包管理器
  - Windows: 从 [nodejs.org](https://nodejs.org/) 下载安装
  
- **Go** (版本 1.25+)
  - macOS: `brew install go`
  - Linux: 使用发行版的包管理器
  - Windows: 从 [go.dev](https://go.dev/dl/) 下载安装
  
- **Go 代理配置** (中国大陆用户推荐)

  ```bash
  go env -w GOPROXY=https://goproxy.cn,direct
  ```

### 运行时依赖

- **LaTeX 发行版** (如 MacTeX、TeX Live、MiKTeX)
  - macOS: `brew install --cask mactex` (或使用 Homebrew: `brew install basictex`)
  - Linux: 使用发行版的包管理器
  - Windows: 下载安装 [TeX Live](https://www.tug.org/texlive/) 或 [MiKTeX](https://miktex.org/)

> **注意**：构建脚本会自动检查构建依赖。如果缺少任何依赖，您将收到有用的安装说明。

### 从源码编译

#### 方式一：使用一键编译脚本（推荐）

项目提供了多平台一键编译脚本，会自动在 `bin/` 目录下生成各平台的二进制文件：

- **Linux/macOS**:

  ```bash
  chmod +x build.sh
  ./build.sh
  ```
  
- **Windows**:

  ```batch
  build.bat
  ```

#### 方式二：手动编译

```bash
git clone https://github.com/SJRnhqh/lazitex.git
cd lazitex
go build -o lazitex ./cmd/lazitex-cli
```

### 下载预编译版本

从 [Releases](https://github.com/SJRnhqh/lazitex/releases) 下载预编译二进制文件（即将推出）

---

## 🚀 快速开始

| 命令 | 说明 | 示例 |
| ------ | ------ | ------ |
| `-h, --help` | 显示此帮助 | `lazitex -h` |
| `-v, --version` | 显示版本号 | `lazitex -v` |
| `-r, --repl` | 启动 REPL 模式 | `lazitex -r` |
| `-t, --tui` | 启动终端 UI | `lazitex -t` |
| `-x, --latex` | LaTeX 环境管理 (使用 `-c/-i/-u` 标志) | `lazitex -x -c` |
| `-o, --ollama` | Ollama 管理 (使用 `-c/-i/-u` 标志) | `lazitex -o -c` |
| `-b, --build` | 构建 LaTeX 文档 (支持 `-o` 输出, `-s` 编译后展示, `-q` 静默模式, `-t` 清理辅助文件) | `lazitex -b main.tex [-o out/] [-s] [-q] [-t]` |
| `-p, --preview` | 实时预览 PDF 文档 (支持 `-q` 静默模式, `-t` 清理辅助文件, `:端口号` 自定义端口) | `lazitex -p main.tex [-q] [-t] [:端口号]` |
| `-m, --llm` | LLM 管理 (add/link/list/test/remove) | `lazitex -m list` 或 `lazitex -m add llama3.2` |
| `config` | 打开配置文件 | `lazitex config` |
| `-l, --lang` | 设置语言 (zh/en) | `lazitex -l zh` |

---

## 📖 使用指南

### 语言设置

LaziTex 支持**中文**和**英文**两种语言，可通过以下两种方式设置：

#### 方式一：运行时更改

```bash
lazitex -l zh --check    # 本次命令使用中文
lazitex --lang en --check  # 本次命令使用英文
```

使用 `-l` 或 `--lang` 参数时，语言偏好会自动保存。

#### 方式二：配置文件

| 平台 | 配置文件位置 |
| ------ | ------------ |
| Windows | `%APPDATA%\lazitex\config.json` |
| Linux | `~/.config/lazitex/config.json` |
| macOS | `~/Library/Application Support/lazitex/config.json` |

配置文件示例：

```json
{
  "language": "zh"
}
```

直接编辑配置文件或使用 `--lang` 参数自动更新。

### LaTeX 环境管理

跨平台 LaTeX 环境管理，支持自动检测、安装和卸载。

```bash
# CLI 模式
$ lazitex -x -c        # 检查
$ lazitex -x -i        # 安装
$ lazitex -x -u        # 卸载

# REPL 模式
lazitex> latex -c      # 检查
lazitex> latex -i      # 安装
lazitex> latex -u      # 卸载
```

### Ollama 管理

跨平台 Ollama 管理，支持检查、安装和卸载。

```bash
# CLI 模式
$ lazitex -o -c        # 检查
$ lazitex -o -i        # 安装
$ lazitex -o -u        # 卸载

# REPL 模式
lazitex> ollama -c     # 检查
lazitex> ollama -i     # 安装
lazitex> ollama -u     # 卸载
```

### 构建与预览

一键构建 LaTeX 文档，支持自动包检测与安装、自适应多轮编译、智能预览和 Web 实时预览。支持静默模式（`-q`）和清理模式（`-t`）。

```bash
# CLI 模式
$ lazitex -b report.tex -s           # 构建并预览
$ lazitex -b report.tex -q           # 静默构建
$ lazitex -b report.tex -t            # 构建并清理辅助文件
$ lazitex -p report.tex              # Web 预览模式（LaziView）
$ lazitex -p report.tex :3000        # Web 预览模式（自定义端口）
$ lazitex -p report.tex -d           # Web 预览模式（开发模式）

# REPL 模式
$ lazitex -r
lazitex> build main.tex -o out/ -s    # 构建文档
lazitex> preview main.tex -q          # 实时预览
lazitex> help                         # 显示帮助
```

**Web 预览模式：**

- **LaziView 模式** (`-p`)：仅 PDF 预览界面，适合多屏幕场景或使用外部编辑器
- **LaziWorkspace 模式** (`-w`)：完整工作区（即将推出）

### 可检测的 LaTeX 工具

LaziTex 可以检测 **23 个 LaTeX 工具**，涵盖 7 大类别：

| 类别 | 工具 | 数量 |
| ------ | ------ | ------ |
| 🔨 **编译器** | pdflatex, xelatex, lualatex, latex, pdftex, tex, etex | 7 |
| 📚 **文献管理** | bibtex, biber | 2 |
| 📇 **索引工具** | makeindex, xindy, texindy | 3 |
| 🔄 **格式转换** | dvipdfmx, dvips, ps2pdf, dvisvgm | 4 |
| ⚙️ **自动化工具** | latexmk | 1 |
| 📦 **包管理器** | tlmgr, mpm | 2 |
| 🔧 **实用工具** | kpsewhich, texdoc, texhash, updmap | 4 |

**工具优先级说明：**

- ⭐ **核心工具（Priority 0）**：必须安装，是基本功能的基础
  - 示例：pdflatex, xelatex, lualatex, bibtex, latexmk
  
- 🔹 **重要工具（Priority 1）**：推荐安装，提供常用功能
  - 示例：biber, makeindex, dvipdfmx, tlmgr
  
- 🔸 **可选工具（Priority 2）**：进阶使用，特定场景需要
  - 示例：xindy, dvips, texdoc, updmap

### 平台支持

LaziTex 支持 Windows、Linux 和 macOS 三大平台：

- **Windows**：检测 TeX Live 和 MiKTeX，支持自动预览（优先 SumatraPDF），通过 winget 管理 Ollama
- **Linux**：自动识别发行版，提供针对性的安装指南，通过多种包管理器或官方脚本管理 Ollama
- **macOS**：支持 MacTeX、Homebrew、MacPorts，支持自动安装/更新/卸载（Intel 和 Apple Silicon），通过 Homebrew 或官方脚本管理 Ollama

---

## 🗂️ 项目结构

LaziTex 采用清晰的模块化架构设计，易于扩展和维护。

```txt
lazitex/
├── 🐹 go.mod                      # Go 模块定义
├── 🔒 go.sum                      # Go 依赖锁定
├── 🐚 build.sh                    # Linux/macOS 一键编译脚本
├── 🪟 build.bat                   # Windows 一键编译脚本
│
├── 🧠 core/                       # 核心逻辑 - LaziTex 的大脑
│   ├── env.go                     # 环境检测、安装器接口定义及安装/卸载逻辑
│   ├── build.go                    # 跨平台编译工作流（策略模式，自适应多轮编译）
│   ├── watcher.go                  # 文件监听（实时预览）
│   ├── ollama.go                   # Ollama 管理器接口定义和包装函数
│   ├── 🚨 errors/                 # 编译错误处理模块
│   │   ├── formatter.go            # 错误信息提取和格式化
│   │   ├── package.go              # 自动包检测与安装
│   │   └── passes.go               # 自适应多轮编译检测
│   ├── ⚡ performance/             # 编译性能优化模块
│   │   ├── concurrent.go           # 并发优化（中间工具并发执行）
│   │   └── lock.go                 # 编译锁与任务管理（任务抢占、超时控制）
│   └── 🤖 ai/                     # AI 架构模块（模块化设计）
│       ├── agent/                 # AI 智能体管理
│       │   └── registry.go        # 智能体注册和管理
│       ├── llm/                   # LLM 提供商抽象
│       │   └── registry.go        # LLM 提供商注册和管理
│       ├── providers/             # LLM 提供商实现
│       │   └── ollama.go          # Ollama 提供商实现
│       └── tools/                 # AI 工具系统
│           └── registry.go        # 工具注册和管理
│
├── 📚 lang/                       # 国际化模块
│   └── i18n.go                    # 多语言支持（中英文）
│
├── ⚙️  config/                     # 配置管理
│   ├── common.go                  # 应用配置管理
│   └── llm.go                     # LLM 提供商配置管理
│
├── 🗄️  backend/                    # Web 服务器后端
│   ├── server.go                  # HTTP 服务器核心（支持模式配置）
│   ├── handlers.go                # HTTP 请求处理（首页、PDF、配置 API）
│   ├── routes.go                  # 路由注册（包含 /api/config）
│   └── sse.go                     # Server-Sent Events 实时推送
│
├── 🌐 frontend/                   # 前端资源
│   ├── embed.go                   # 静态资源嵌入（旧版）
│   ├── static/                    # 旧版静态文件
│   │   └── index.html             # 旧版预览页面（双缓冲优化）
│   └── 🟢 vue/                    # Vue 3 前端（新版）
│       ├── src/                   # Vue 源码文件
│       │   ├── api/               # API 层 - 统一管理后端 API 调用
│       │   │   ├── pdf.js         # PDF 相关 API（下载、连接测试）
│       │   │   ├── sse.js         # SSE 相关 API（实时推送）
│       │   │   └── config.js      # 配置相关 API（应用模式）
│       │   ├── utils/             # 工具函数 - 纯函数工具库
│       │   │   ├── pdf.js         # PDF 工具函数（解析、缩放、清理、错误格式化）
│       │   │   └── pdfDisplay.js  # PDF 显示管理（虚拟滚动、分页导航、配置管理）
│       │   ├── components/        # Vue 组件
│       │   │   ├── DynamicPDFViewer.vue  # PDF 预览组件
│       │   │   └── ModeSwitcher.vue      # 模式切换组件
│       │   ├── layouts/           # 布局组件
│       │   │   ├── LaziViewLayout.vue      # 仅 PDF 预览布局
│       │   │   └── LaziWorkspaceLayout.vue # 完整工作区布局
│       │   ├── stores/            # Pinia 状态管理
│       │   │   └── appMode.js     # 应用模式状态管理
│       │   ├── styles/            # 样式文件
│       │   │   ├── theme.css      # 主题配色（Nord 配色方案）
│       │   │   └── dynamic-pdf-viewer.module.css  # PDF 预览组件样式（CSS Modules）
│       │   ├── App.vue            # 根组件（模式分发器）
│       │   └── main.js            # 入口文件
│       ├── public/                # 公共资源
│       ├── index.html             # HTML 模板
│       ├── ⚡ vite.config.js       # Vite 配置
│       └── 📦 package.json        # 依赖配置
│
├── 🎯 target/                     # 平台特定实现
│   ├── 🪟 win/                    # Windows 平台
│   │   ├── checker.go             # 检测 Windows 上的 TeX Live 和 MiKTeX
│   │   ├── ollama.go              # Windows Ollama 管理器（✅ 已实现）
│   │   └── builder.go             # Windows 编译逻辑（即将推出）
│   ├── 🐧 linux/                  # Linux 平台
│   │   ├── checker.go             # 检测各发行版的 TeX 安装
│   │   ├── installer.go           # Linux 安装器（开发中）
│   │   ├── ollama.go              # Linux Ollama 管理器（✅ 已实现）
│   │   └── builder.go             # Linux 编译逻辑（即将推出）
│   └── 🍏 mac/                    # macOS 平台
│       ├── checker.go              # 检测 MacTeX、Homebrew、MacPorts
│       ├── installer.go            # macOS 安装器（✅ 已实现）
│       └── builder.go             # macOS 编译与预览逻辑（✅ 已实现）
│
├── 💻 cmd/                       # 命令行界面
│   └── 🚀 lazitex-cli/
│       ├── main.go                # CLI 入口：解析命令并路由到不同模式
│       ├── 📋 tasks/              # 任务执行层（统一管理命令执行逻辑）
│       │   ├── build.go           # 构建功能
│       │   ├── preview.go         # 实时预览功能
│       │   ├── ollama.go          # Ollama 管理功能
│       │   ├── latex.go           # LaTeX 环境管理功能
│       │   ├── llm.go             # LLM 管理功能
│       │   ├── config.go          # 配置文件管理功能
│       │   └── help.go            # 帮助信息
│       ├── 🧠 internal/           # 内部命令处理（统一动作处理器）
│       │   ├── common.go          # 通用工具和动作处理器
│       │   ├── latex.go           # LaTeX 命令处理逻辑
│       │   └── ollama.go          # Ollama 命令处理逻辑
│       └── 🎨 ui/                 # 用户界面层（交互模式实现）
│           ├── tui.go             # 基于 Bubble Tea 的终端 UI 模式
│           └── repl.go            # 交互式 REPL 模式（命令解析、补全、历史）
│
├── 📖 README.md                   # 英文文档
├── 📖 README_zh.md                # 中文文档（当前文件）
└── 📄 LICENSE                     # Apache 许可证 2.0
```

### 架构亮点

- **自适应多轮编译检测** - 自动检测并处理 6 种需要多轮编译的场景（目录、交叉引用、参考文献、索引、术语表、PDF 书签），无需手动指定编译次数

- **自动包检测与安装** - 从编译错误中智能提取缺失的包名，自动通过 tlmgr/mpm 安装，实现真正的零配置编译体验

- **工具优先级系统** - 23 个 LaTeX 工具按优先级分类（⭐ 核心、🔹 重要、🔸 可选），帮助用户快速识别关键工具，优化安装建议

- **模块化 AI 架构** - AI 智能体、LLM 提供商和工具的清晰分离，采用可扩展的注册系统实现无缝集成

- **智能 LLM 管理** - 智能提供商注册、配置持久化和当前 LLM 跟踪，支持跨平台配置

---

## 🌟 开发路线图

### 已完成 ✅

- [x] **跨平台环境检测与管理** - Windows、Linux、macOS 环境检测，macOS 自动安装/更新/卸载
- [x] **编译系统** - 一键编译、智能预览、自动包检测与安装、自适应多轮编译、性能优化
- [x] **Web 预览前端** - Vue 3 前端、SSE 实时刷新、PDF 虚拟滚动、双模式架构
- [x] **Ollama 管理** - 跨平台检查、安装、卸载
- [x] **AI 架构** - LLM 和 Agent 模块化架构、智能提供商管理
- [x] **交互模式** - REPL 交互模式（历史记录、Tab 补全、国际化）

### 进行中 🚧

- [ ] **实时预览优化** - 错误反馈增强
- [ ] **Linux 环境管理** - 自动安装/更新/卸载
- [ ] **PDF 预览增强** - 分页导航、缩放控制

### 计划中 📋

- [ ] **AI 集成** - CLI 内部错误诊断、代码生成、智能补全
- [ ] **Web 工作区** - 文件管理、代码编辑器、终端集成、AI 聊天
- [ ] **项目功能** - 项目初始化、多文档支持、自定义配置
- [ ] **平台扩展** - Windows 环境管理、TUI 完善、GUI 支持

---

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的更改 (`git commit -m '添加某个很棒的特性'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

### 开发指南

- 遵循 Go 代码规范
- 为新功能添加测试
- 更新相关文档
- 确保跨平台兼容性

---

## 📄 许可证

Apache License 2.0 - 包含明确的专利授权，保护开发者与用户。

详见 [LICENSE](LICENSE) 文件。

---

## 🙏 致谢

- 使用 [Bubble Tea](https://github.com/charmbracelet/bubbletea) 构建精美的终端 UI
- 灵感来源于 LaTeX 社区对优美排版的执着追求

---

## 💡 常见问题

### 为什么检测不到我的 LaTeX 安装？

- **Windows**：确保 LaTeX 安装路径已添加到 PATH 环境变量
- **Linux**：尝试运行 `which pdflatex` 确认安装位置
- **macOS**：如果使用 Homebrew，确保已运行 `brew link` 命令

### 可以添加其他语言吗？

当前只支持中文和英文。如果你需要其他语言支持，欢迎提交 Issue 或 Pull Request！

---

**用 ❤️ 制作 by [SJRnhqh](https://github.com/SJRnhqh)**
