# LaziTex 🧸

跨平台即时 LaTex 编译工具 - 基于 Go 构建，集成本地 AI 助力写作与优化

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://github.com/SJRnhqh/lazitex)
[![Status](https://img.shields.io/badge/Status-Active%20Development-brightgreen)](https://github.com/SJRnhqh/lazitex)

中文文档 | [English](README.md)

---

## ✨ 特性

- 🔍 **智能环境检测** - 自动检测 LaTex 安装并提供详细诊断报告
- 📦 **自动安装/更新** - 一键安装或更新 LaTex 环境（macOS 已支持）
- 🗑️ **自动卸载** - 一键卸载 LaTex 环境，智能检测安装方式（macOS 已支持）
- 🌍 **跨平台支持** - 无缝支持 Windows、Linux 和 macOS
- 🎯 **零配置启动** - 开箱即用，完美兼容 TeX Live、MiKTeX 和 MacTeX
- 🚀 **快速轻量** - 基于 Go 构建，编译速度极快
- 💬 **交互式 REPL** - 支持检查、安装、卸载、语言切换等交互命令
- 🔄 **自适应多轮编译** - 自动检测并处理交叉引用、目录、参考文献、索引、术语表等需要多次编译的场景
- 📦 **自动包管理** - 从编译错误中自动检测缺失的包并通过 tlmgr/mpm 安装
- 🧠 **AI 增强** (即将推出) - 本地 AI 助手，帮助撰写和优化 LaTex 文档
- 🎨 **多种模式** - 支持 TUI（终端界面）和 REPL 模式，适应不同工作流

---

## 📦 安装

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
|------|------|------|
| `-c, --check` | 检查 LaTex 环境 | `lazitex -c` |
| `-i, --install` | 安装/更新 LaTex | `lazitex -i` |
| `-u, --uninstall` | 卸载 LaTex | `lazitex -u` |
| `-b, --build` | 构建 LaTex 文档 (支持 `-o` 输出, `-s` 编译后展示) | `lazitex -b main.tex [-o out/] [-s]` |
| `-p, --preview` | 实时预览 PDF 文档 (监听保存动作并自动刷新) | `lazitex -p main.tex` |
| `-r, --repl` | 启动 REPL 模式 | `lazitex -r` |
| `-l, --lang` | 设置语言 | `lazitex -l zh` |
| `-h, --help` | 显示帮助 | `lazitex -h` |
| `-v, --version` | 显示版本 | `lazitex -v` |

---

## 📖 使用指南

### 语言设置

LaziTex 支持**中文**和**英文**两种语言，可通过以下两种方式设置：

方式一：运行时更改

```bash
lazitex -l zh --check    # 本次命令使用中文
lazitex --lang en --check  # 本次命令使用英文
```

使用 `-l` 或 `--lang` 参数时，语言偏好会自动保存。

方式二：配置文件

| 平台 | 配置文件位置 |
|------|------------|
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

### 构建与预览

**功能特点：**

- 一键构建 - 自动调用 XeLaTex 编译器
- 自动包检测与安装 - 从编译错误中检测缺失包并自动安装
- 自适应多轮编译 - 自动处理交叉引用、目录、参考文献、索引等场景
- 智能预览 - 构建成功后自动打开 PDF（macOS 优先使用 Skim，Windows 优先使用 SumatraPDF）
- 自定义输出 - 支持指定输出目录或完整路径

**构建示例：**

```bash
$ lazitex -b report.tex -s
🚀 正在构建 LaTex 文档: report.tex
📁 工作目录: /Users/user/projects/paper
... (编译器输出) ...
✨ 构建成功！
(自动打开展示窗口)
```

### REPL 模式

```bash
lazitex> help                      # 显示可用命令
lazitex> check                     # 检查 LaTex 环境
lazitex> install                   # 安装或更新 LaTex 环境
lazitex> uninstall                 # 卸载 LaTex 环境
lazitex> build main.tex -o out/ -s # 构建 LaTex 文档，指定输出目录并带展示
lazitex> lang zh                   # 切换到中文
lazitex> lang en                   # 切换到英文
lazitex> lang                      # 查看当前语言
lazitex> version                   # 显示版本号
lazitex> exit                      # 退出 REPL
lazitex> cd /tmp                   # 切换目录（不带参数回到用户主目录）
lazitex> ls                        # 列出当前内容（支持参数）
lazitex> pwd                       # 显示当前路径
lazitex> clear                     # 清屏
lazitex> cat file.log              # 查看文件内容（支持 .tex, .log, .aux）
```

**现代化的 REPL 体验：**

- 🕒 **命令历史** - 支持上下方向键翻找历史执行过的命令
- ⌨️ **光标移动** - 支持左右方向键、Home/End、Ctrl+A/E 等标准光标操作
- 🎯 **智能补全** - 强大的 Tab 自动补全功能：
  - **命令补全**：自动补全内置命令
  - **上下文感知**：`cd` 只补全文件夹；`build` 只补全 `.tex`；`cat` 补全 `.tex/.log/.aux`
- 🌐 **全方位国际化** - 所有命令帮助、错误提示、用法说明均支持中英文切换

REPL 中的语言切换即时生效，并且会保存偏好设置供下次使用！

### 可检测的 LaTex 工具

LaziTex 可以检测 **23 个 LaTex 工具**，涵盖 7 大类别：

| 类别 | 工具 | 数量 |
|------|------|------|
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

---

## 🗂️ 项目结构

LaziTex 采用清晰的模块化架构设计，易于扩展和维护。

```txt
lazitex/
├── 🗂️  go.mod                      # Go 模块定义
├── 🔒 go.sum                      # Go 依赖锁定
├── 🛠️  build.sh                   # Linux/macOS 一键编译脚本
├── 🛠️  build.bat                  # Windows 一键编译脚本
│
├── 🧠 core/                       # 核心逻辑 - LaziTex 的大脑
│   ├── 🌍 env.go                  # 环境检测、安装器接口定义及安装/卸载逻辑
│   ├── 🔨 build.go                # 跨平台编译工作流（策略模式，自适应多轮编译）
│   ├── 👀 watcher.go              # 文件监听（实时预览）
│   └── 🚨 errors/                 # 编译错误处理模块
│       ├── 📦 package.go           # 自动包检测与安装
│       └── 🔄 passes.go            # 自适应多轮编译检测
│
├── 🌐 lang/                       # 国际化模块
│   └── i18n.go                    # 多语言支持（中英文）
│
├── ⚙️  config/                     # 配置管理
│   └── config.go                  # 用户偏好与设置
│
├── 🎯 target/                     # 平台特定实现
│   ├── 🪟 win/                    # Windows 平台
│   │   ├── 🔍 checker.go          # 检测 Windows 上的 TeX Live 和 MiKTeX
│   │   ├── 📦 installer.go        # Windows 安装器（开发中）
│   │   └── 🛠️  builder.go         # Windows 编译逻辑（即将推出）
│   ├── 🐧 linux/                  # Linux 平台
│   │   ├── 🔍 checker.go          # 检测各发行版的 TeX 安装
│   │   ├── 📦 installer.go        # Linux 安装器（开发中）
│   │   └── 🛠️  builder.go         # Linux 编译逻辑（即将推出）
│   └── 🍏 mac/                    # macOS 平台
│       ├── 🔍 checker.go          # 检测 MacTeX、Homebrew、MacPorts
│       ├── 📦 installer.go        # macOS 安装器（✅ 已实现）
│       └── 🛠️  builder.go         # macOS 编译与预览逻辑（✅ 已实现）
│
├── 🖥️  cmd/                       # 命令行界面
│   └── 🚀 lazitex-cli/
│       ├── 📄 main.go             # CLI 入口：解析命令并路由到不同模式
│       ├── 📋 tasks/              # 任务执行层（统一管理命令执行逻辑）
│       │   ├── 🌍 env.go          # 环境操作（检查、安装、卸载）
│       │   ├── 🔨 build.go        # 构建功能
│       │   ├── 👀 preview.go      # 实时预览功能
│       │   └── 📖 help.go         # 帮助信息
│       └── 🎨 ui/                 # 用户界面层（交互模式实现）
│           ├── 🖼️  tui.go         # 基于 Bubble Tea 的终端 UI 模式
│           └── 💬 repl.go         # 交互式 REPL 模式（命令解析、补全、历史）
│
├── 📖 README.md                   # 英文文档
├── 📖 README_zh.md                # 中文文档（当前文件）
└── 📄 LICENSE                     # Apache 许可证 2.0
```

### 架构亮点

- **编译策略模式** - 采用策略模式（`compileStrategy` 接口）处理不同编译场景，默认策略自动处理包错误和多轮编译，未来可轻松扩展 AI 驱动策略，无需修改核心编译循环

- **编译错误处理模块** - `core/errors/` 模块化处理编译问题：自动包检测与安装、自适应多轮编译检测（覆盖目录、交叉引用、参考文献、索引、术语表、PDF 书签等 6 种场景）

- **工具优先级系统** - 工具按优先级分类（⭐ 核心、🔹 重要、🔸 可选），帮助用户快速识别关键工具

---

## 🛠️ 平台特定功能

### Windows

- 检测 **TeX Live** 和 **MiKTeX** 安装
- 搜索常见安装路径（C:\texlive、C:\Program Files\MiKTeX 等）
- 检查 AppData 中的用户特定安装
- **自动预览支持** - 构建成功后自动调用系统默认 PDF 阅读器预览
- **智能预览增强** - 优先检测并调用 **SumatraPDF**（支持静默刷新与 -reuse-instance），若未安装则回退至系统默认关联程序
- 提供详细的 Windows 安装指南

### Linux

- 自动识别 Linux 发行版（Ubuntu、Fedora、Arch 等）
- 提供针对不同发行版的安装命令（apt、dnf、pacman）
- 区分包管理器安装和手动安装
- 支持 x86_64 和 aarch64 架构

### macOS

- 检测 **MacTeX**、**Homebrew** 和 **MacPorts** 安装
- 显示 macOS 版本信息
- **自动安装/更新** - 通过 Homebrew 一键安装 BasicTeX
- **自动卸载** - 智能检测安装方式并一键卸载
- 根据检测到的包管理器提供优化的安装指南
- 支持 Intel 和 Apple Silicon 芯片

---

## 🌟 开发路线图

### 已完成 ✅

- [x] **跨平台环境检测** - Windows、Linux、macOS 环境检测与 23 个工具检查
- [x] **macOS 环境管理** - 自动安装/更新/卸载，支持 Homebrew、MacPorts、MacTeX
- [x] **REPL 交互模式** - 历史记录、Tab 补全、Shell 快捷命令、全语言国际化
- [x] **编译与预览系统** - 一键编译、智能预览（macOS Skim/Windows SumatraPDF）、实时预览监听
- [x] **智能编译优化** - 自动包检测与安装、自适应多轮编译（交叉引用、目录、参考文献等）

### 进行中 🚧

- [ ] **Web 预览模式** - 本地 HTTP 服务器 + 浏览器预览，提供 Overleaf 风格的 Web 预览体验
- [ ] **实时预览优化** - 任务抢占、编译锁、错误反馈增强
- [ ] **Linux 环境管理** - 自动安装/更新/卸载

### 计划中 📋

- [ ] **AI 集成（CLI 内部）** - 编译错误诊断、代码生成与优化、智能补全
- [ ] **多文档项目支持** - 自动处理文件依赖关系和编译顺序
- [ ] **自定义编译配置** - 项目级别配置文件（编译器选择、输出目录、编译参数等）

### 后续计划 🔮

- [ ] **项目初始化** - `lazitex init` 极简启动，快速拉取优质 LaTex 模板
- [ ] **GUI 图形界面** - 现代化图形界面支持
- [ ] **Windows 环境管理** - 自动安装/更新/卸载（优先级较低）
- [ ] **AI 原生创作流** - Prompt2PDF 体验，多智能体协同工作流
- [ ] **云端垂直生态** - 服务器端部署，在线交互式生成服务

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
- 灵感来源于 LaTex 社区对优美排版的执着追求

---

## 💡 常见问题

### 1. 为什么检测不到我的 LaTex 安装？

- **Windows**：确保 LaTex 安装路径已添加到 PATH 环境变量
- **Linux**：尝试运行 `which pdflatex` 确认安装位置
- **macOS**：如果使用 Homebrew，确保已运行 `brew link` 命令

### 2. 支持哪些 LaTex 发行版？

目前支持：

- TeX Live（所有平台）
- MiKTeX（Windows）
- MacTeX（macOS）

### 3. 如何安装缺失的工具？

方法一：使用 LaziTex 自动安装（推荐）

```bash
lazitex -i    # 自动安装或更新 LaTex 环境（目前支持 macOS）
```

方法二：手动安装

运行 `lazitex --check` 后，如果有工具未安装，系统会自动显示针对你平台的安装指南。

### 4. 如何卸载 LaTex 环境？

使用 LaziTex 一键卸载（推荐，目前支持 macOS）：

```bash
lazitex -u    # 自动卸载 LaTex 环境（目前支持 macOS）
```

**功能特点：**

- 自动检测安装方式（Homebrew、MacPorts、MacTeX 官方安装等）
- 交互式确认，防止误操作
- 未安装时会提示可以使用 `-i` 安装

**手动卸载：**

- Homebrew: `brew uninstall --cask basictex` 或 `brew uninstall --cask mactex`
- MacPorts: `sudo port uninstall texlive`
- MacTeX 官方: 手动删除 `/Library/TeX/` 目录

### 5. 如何查看当前使用的语言？

```bash
# 查看配置文件
# Windows
type %APPDATA%\lazitex\config.json
# Linux/macOS
cat ~/.config/lazitex/config.json
```

### 6. 如何重置为默认语言（英文）？

```bash
lazitex -l en --help
# 或使用长选项
lazitex --lang en --help
# 或删除配置文件
```

### 7. 可以添加其他语言吗？

当前只支持中文和英文。如果你需要其他语言支持，欢迎提交 Issue 或 Pull Request！

### 8. 未来会支持哪些功能？

请查看上方的 [开发路线图](#-开发路线图) 部分。

---

**用 ❤️ 制作 by [SJRnhqh](https://github.com/SJRnhqh)**
