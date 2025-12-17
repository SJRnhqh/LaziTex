# LaziTex 🧸

跨平台即时 LaTeX 编译工具 - 基于 Go 构建，集成本地 AI 助力写作与优化

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://github.com/SJRnhqh/lazitex)
[![Status](https://img.shields.io/badge/Status-Active%20Development-brightgreen)](https://github.com/SJRnhqh/lazitex)

中文文档 | [English](README.md)

---

## ✨ 特性

- 🔍 **智能环境检测** - 自动检测 LaTeX 安装并提供详细诊断报告
- 📦 **自动安装/更新** - 一键安装或更新 LaTeX 环境（macOS 已支持）
- 🗑️ **自动卸载** - 一键卸载 LaTeX 环境，智能检测安装方式（macOS 已支持）
- 🌍 **跨平台支持** - 无缝支持 Windows、Linux 和 macOS
- 🎯 **零配置启动** - 开箱即用，完美兼容 TeX Live、MiKTeX 和 MacTeX
- 🚀 **快速轻量** - 基于 Go 构建，编译速度极快
- 💬 **交互式 REPL** - 支持检查、安装、卸载、语言切换等交互命令
- 🧠 **AI 增强** (即将推出) - 本地 AI 助手，帮助撰写和优化 LaTeX 文档
- 🎨 **多种模式** - 支持 TUI（终端界面）和 REPL 模式，适应不同工作流

---

## 📦 安装

### 从源码编译

```bash
git clone https://github.com/SJRnhqh/lazitex.git
cd lazitex
go build -o lazitex ./cmd/lazitex-cli
```

### 下载预编译版本

从 [Releases](https://github.com/SJRnhqh/lazitex/releases) 下载预编译二进制文件（即将推出）

---

## 🚀 快速开始

### 📋 命令速查

| 命令 | 说明 | 示例 |
|------|------|------|
| `-c, --check` | 检查 LaTeX 环境 | `lazitex -c` |
| `-i, --install` | 安装/更新 LaTeX | `lazitex -i` |
| `-u, --uninstall` | 卸载 LaTeX | `lazitex -u` |
| `-r, --repl` | 启动 REPL 模式 | `lazitex -r` |
| `-l, --lang` | 设置语言 | `lazitex -l zh` |
| `-h, --help` | 显示帮助 | `lazitex -h` |
| `-v, --version` | 显示版本 | `lazitex -v` |

### 检查 LaTeX 环境

在编译之前，先验证你的 LaTeX 安装：

```bash
lazitex --check    # 完整命令
lazitex -c         # 简短命令
```

这将检测：

- ✅ 所有已安装的 LaTeX 编译器（pdflatex、xelatex、lualatex 等）
- ✅ 文献管理工具（bibtex、biber）
- ✅ 格式转换工具（dvipdfmx、ps2pdf 等）
- ✅ 包管理器（tlmgr、mpm）
- ✅ 你的 LaTeX 发行版（TeX Live、MiKTeX、MacTeX）

**示例输出：**

```txt
╔═══════════════════════════════════════════════════════════╗
║           LaTeX 编译环境检测结果                          ║
╚═══════════════════════════════════════════════════════════╝

🖥️  操作系统: Windows (windows)
📦 LaTeX 发行版: TeX Live 2025

✅ 已安装: 22/23 工具 (核心工具: 6/6)

🔨 编译器 (7/7)
────────────────────────────────────────────────────────────
  ⭐ ✓ PDFLaTeX        [已安装]
      路径: C:\texlive\2025\bin\windows\pdflatex.exe
      版本: pdfTeX 3.141592653-2.6-1.40.28 (TeX Live 2025)
  ⭐ ✓ XeLaTeX         [已安装]
      路径: C:\texlive\2025\bin\windows\xelatex.exe
      版本: XeTeX 3.141592653-2.6-0.999997 (TeX Live 2025)
  ...
```

---

## 🌍 语言设置

LaziTex 是一个**双语工具**，完整支持**中文**和**英文 (English)**。

### 默认语言

- **英文 (English)** 是默认语言
- 所有输出（帮助信息、环境检测结果、错误提示）默认使用英文

### 快速开始

**一次性使用特定语言：**

```bash
lazitex --lang zh --check    # 本次命令使用中文
lazitex -l en --check        # 本次命令使用英文（短选项形式）
```

**永久切换默认语言：**

```bash
# 使用 --lang 或 -l 时，你的选择会自动保存
lazitex -l zh --help         # 切换到中文并保存为默认语言
lazitex --help               # 以后的命令都使用中文

# 随时切换回英文（两种形式都可以）
lazitex --lang en --help     # 切换到英文并保存为默认语言
lazitex --check              # 现在使用英文
```

### 高级用法

**使用环境变量临时覆盖：**

```bash
# Windows PowerShell
$env:LAZITEX_LANG="zh"; lazitex --check

# Linux/macOS/Git Bash
export LAZITEX_LANG=zh
lazitex --check
```

**支持的语言代码：**

- `zh`, `chinese` → 中文
- `en`, `english` → English (英文)

**命令格式：**

- 长选项: `--lang <语言代码>`
- 短选项: `-l <语言代码>`

两种形式完全等效，都会保存你的偏好设置。

### 语言优先级（从高到低）

1. **命令行参数** `-l <语言>` 或 `--lang <语言>` （会保存偏好）
2. **环境变量** `LAZITEX_LANG=<语言>`
3. **已保存的用户偏好**（配置文件）
4. **系统默认**（英文）

### 配置文件

你的语言偏好保存在：

| 平台 | 配置文件位置 |
|------|------------|
| **Windows** | `%APPDATA%\lazitex\config.json` |
| **Linux** | `~/.config/lazitex/config.json` |
| **macOS** | `~/Library/Application Support/lazitex/config.json` |

配置文件示例：

```json
{
  "language": "zh"
}
```

你可以：

- 直接编辑这个文件来更改默认语言
- 使用 `--lang` 参数自动更新
- 删除配置文件恢复默认（英文）

### 使用示例

```bash
# 使用默认语言（英文）检查环境
lazitex --check

# 使用中文检查一次（短选项形式）
lazitex -l zh --check

# 永久切换到中文
lazitex --lang zh --help

# 以后所有命令都使用中文
lazitex --check
lazitex --version

# 切换回英文（短选项形式）
lazitex -l en --help

# 灵活使用长短选项
lazitex -l zh --check       # 短选项
lazitex --lang en --check   # 长选项

# 在 REPL 模式中切换语言
lazitex --repl
lazitex> lang zh            # 交互式切换到中文
lazitex> help               # 帮助信息现在是中文
```

### 输出对比

**英文 (English):**

```txt
✅ Installed: 22/23 tools (core tools: 6/6)
🔨 Compilers (7/7)
  ⭐ ✓ PDFLaTeX        [Installed]
      Path: C:\texlive\2025\bin\windows\pdflatex.exe
  ⭐ ✓ XeLaTeX         [Installed]
🎉 Excellent! All core compilers are installed
```

**中文:**

```txt
✅ 已安装: 22/23 工具 (核心工具: 6/6)
🔨 编译器 (7/7)
  ⭐ ✓ PDFLaTeX        [已安装]
      路径: C:\texlive\2025\bin\windows\pdflatex.exe
  ⭐ ✓ XeLaTeX         [已安装]
🎉 太棒了！所有核心编译器都已安装
```

### 常见问题

**Q: 如何查看当前使用的语言？**

```bash
# 查看配置文件
# Windows
type %APPDATA%\lazitex\config.json
# Linux/macOS
cat ~/.config/lazitex/config.json
```

**Q: 如何重置为默认语言（英文）？**

```bash
lazitex -l en --help
# 或使用长选项
lazitex --lang en --help
# 或删除配置文件
```

**Q: 可以添加其他语言吗？**

当前只支持中文和英文。如果你需要其他语言支持，欢迎提交 Issue 或 Pull Request！

### 安装/更新 LaTeX 环境

LaziTex 支持自动安装和更新 LaTeX 环境（目前支持 macOS）：

```bash
lazitex --install    # 安装或更新 LaTeX 环境
lazitex -i           # 简短命令
```

**功能特点：**

- 🔍 **智能检测** - 自动检测是否已安装 LaTeX
- 📦 **自动安装** - 未安装时自动安装 BasicTeX（macOS）
- 🔄 **智能更新** - 已安装时自动检查并更新包
- 💬 **交互确认** - 安装/更新前会询问确认
- 📋 **包列表** - 更新时会显示可更新的包列表

**macOS 安装示例：**

```bash
$ lazitex -i
未检测到 LaTeX 环境
将使用 Homebrew 安装 BasicTeX (轻量版，约 1.5 GB)
安装命令: brew install --cask basictex
是否要安装? [y/n]: y
正在通过 Homebrew 安装 BasicTeX...
✅ BasicTeX 安装成功！
```

**更新示例：**

```bash
$ lazitex -i
检测到 LaTeX 已安装
以下包可以更新:

  ...
是否要更新这些包? [y/n]: y
✅ 更新成功！
```

### 卸载 LaTeX 环境

LaziTex 支持一键卸载 LaTeX 环境（目前支持 macOS）：

```bash
lazitex --uninstall  # 卸载 LaTeX 环境
lazitex -u           # 简短命令
```

**功能特点：**

- 🔍 **智能检测** - 自动检测是否已安装 LaTeX
- 🗑️ **自动卸载** - 自动检测安装方式并执行相应卸载
- 💬 **交互确认** - 卸载前会询问确认，防止误操作
- 📋 **友好提示** - 未安装时会提示可以使用 `-i` 一键安装

**支持的卸载方式：**

- ✅ Homebrew BasicTeX - 自动卸载
- ✅ Homebrew MacTeX - 自动卸载
- ✅ MacPorts - 自动卸载
- ℹ️ MacTeX 官方安装 - 提供手动卸载指南

**macOS 卸载示例：**

```bash
$ lazitex -u
检测到 LaTeX 已安装
将卸载 LaTeX 环境（包括所有已安装的包）
是否要卸载? [y/n]: y
正在卸载 LaTeX 环境...
✅ LaTeX 环境卸载成功！
```

**未安装时的提示：**

```bash
$ lazitex -u
未检测到 LaTeX 环境
提示: 可以使用 'lazitex -i' 一键安装 LaTeX 环境
```

### 其他命令

```bash
lazitex --help     # 显示帮助信息
lazitex --version  # 显示版本号
lazitex --check    # 检查 LaTeX 环境
lazitex --repl     # REPL 交互模式
lazitex --tui      # 终端界面模式（即将推出）
```

**REPL 模式** - 交互式命令：

```bash
lazitex> help              # 显示可用命令
lazitex> check             # 检查 LaTeX 环境
lazitex> install           # 安装或更新 LaTeX 环境
lazitex> uninstall         # 卸载 LaTeX 环境
lazitex> lang zh           # 切换到中文
lazitex> lang en           # 切换到英文
lazitex> lang              # 显示当前语言
lazitex> version           # 显示版本号
lazitex> exit              # 退出 REPL
```

REPL 中的语言切换即时生效，并且会保存偏好设置供下次使用！

---

## 🗂️ 项目结构

LaziTex 采用清晰的模块化架构设计，易于扩展和维护。

```txt
lazitex/
├── 📦 go.mod                      # Go 模块定义
├── 📦 go.sum                      # Go 依赖锁定
│
├── 🧠 core/                       # 核心逻辑 - LaziTex 的大脑
│   ├── 🔍 checker.go              # 环境检测引擎
│   ├── 📦 installer.go            # 安装器接口定义
│   ├── 🌐 i18n.go                 # 国际化支持（中英文）
│   └── 🔨 build.go                # 跨平台编译工作流
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
│       └── 🛠️  builder.go         # macOS 编译逻辑（即将推出）
│
├── 🖥️  cmd/                       # 命令行界面
│   └── 🚀 lazitex-cli/
│       ├── 📄 main.go             # CLI 入口：解析命令并路由到不同模式
│       └── 🎨 modes/              # 不同的交互模式
│           ├── 🖼️  tui.go         # 基于 Bubble Tea 的终端 UI 模式
│           └── 💬 repl.go         # 交互式 REPL 模式（支持检查、安装、卸载、语言切换）
│
├── 📖 README.md                   # 英文文档
├── 📖 README_zh.md                # 中文文档（当前文件）
└── 📄 LICENSE                     # MIT 许可证
```

### 架构亮点

- **依赖注入**：平台特定的检查器在 `main.go` 中创建并注入到核心逻辑，避免循环依赖
- **接口驱动**：`EnvironmentChecker` 接口使平台扩展变得简单
- **优先级系统**：工具按优先级分类（⭐ 核心、🔹 重要、🔸 可选）

---

## 🔍 可检测的 LaTeX 工具

LaziTex 可以检测 **23 个 LaTeX 工具**，涵盖 7 大类别：

| 类别 | 工具 | 数量 |
|------|------|------|
| 🔨 **编译器** | pdflatex, xelatex, lualatex, latex, pdftex, tex, etex | 7 |
| 📚 **文献管理** | bibtex, biber | 2 |
| 📇 **索引工具** | makeindex, xindy, texindy | 3 |
| 🔄 **格式转换** | dvipdfmx, dvips, ps2pdf, dvisvgm | 4 |
| ⚙️ **自动化工具** | latexmk | 1 |
| 📦 **包管理器** | tlmgr, mpm | 2 |
| 🔧 **实用工具** | kpsewhich, texdoc, texhash, updmap | 4 |

### 工具优先级说明

- ⭐ **核心工具（Priority 0）**：必须安装，是基本功能的基础
  - 示例：pdflatex, xelatex, lualatex, bibtex, latexmk
  
- 🔹 **重要工具（Priority 1）**：推荐安装，提供常用功能
  - 示例：biber, makeindex, dvipdfmx, tlmgr
  
- 🔸 **可选工具（Priority 2）**：进阶使用，特定场景需要
  - 示例：xindy, dvips, texdoc, updmap

---

## 🛠️ 平台特定功能

### Windows

- 检测 **TeX Live** 和 **MiKTeX** 安装
- 搜索常见安装路径（C:\texlive、C:\Program Files\MiKTeX 等）
- 检查 AppData 中的用户特定安装
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

- [x] Windows、Linux、macOS 环境检测
- [x] 全面的 LaTeX 工具检查（23 个工具）
- [x] macOS LaTeX 环境自动安装/更新
- [x] macOS LaTeX 环境自动卸载
- [x] REPL 交互模式（支持检查、安装、卸载、语言切换）

### 进行中 🚧

- [ ] Linux LaTeX 环境自动安装/更新
- [ ] Linux LaTeX 环境自动卸载

### 计划中 📋

- [ ] LaTeX 编译引擎（编译 .tex 文件为 PDF，支持多轮编译、交叉引用、错误处理）
- [ ] 带实时预览的 TUI 模式
- [ ] 本地 AI 集成，提供 LaTeX 写作辅助
- [ ] 多文档项目支持
- [ ] 自定义编译配置文件

### 后续计划 🔮

- [ ] Windows LaTeX 环境自动安装/更新（Windows 安装方式复杂，优先级较低）
- [ ] Windows LaTeX 环境自动卸载

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

MIT 许可证 - 自由使用、修改和分享。

详见 [LICENSE](LICENSE) 文件。

---

## 🙏 致谢

- 使用 [Bubble Tea](https://github.com/charmbracelet/bubbletea) 构建精美的终端 UI
- 灵感来源于 LaTeX 社区对优美排版的执着追求

---

## 💡 常见问题

### 1. 为什么检测不到我的 LaTeX 安装？

- **Windows**：确保 LaTeX 安装路径已添加到 PATH 环境变量
- **Linux**：尝试运行 `which pdflatex` 确认安装位置
- **macOS**：如果使用 Homebrew，确保已运行 `brew link` 命令

### 2. 支持哪些 LaTeX 发行版？

目前支持：

- TeX Live（所有平台）
- MiKTeX（Windows）
- MacTeX（macOS）

### 3. 如何安装缺失的工具？

方法一：使用 LaziTex 自动安装（推荐）

```bash
lazitex -i    # 自动安装或更新 LaTeX 环境（目前支持 macOS）
```

**方法二：手动安装**
运行 `lazitex --check` 后，如果有工具未安装，系统会自动显示针对你平台的安装指南。

### 4. 如何卸载 LaTeX 环境？

使用 LaziTex 一键卸载（推荐，目前支持 macOS）：

```bash
lazitex -u    # 自动卸载 LaTeX 环境（目前支持 macOS）
```

**功能特点：**

- 自动检测安装方式（Homebrew、MacPorts、MacTeX 官方安装等）
- 交互式确认，防止误操作
- 未安装时会提示可以使用 `-i` 安装

**手动卸载：**

- Homebrew: `brew uninstall --cask basictex` 或 `brew uninstall --cask mactex`
- MacPorts: `sudo port uninstall texlive`
- MacTeX 官方: 手动删除 `/Library/TeX/` 目录

### 5. 未来会支持哪些功能？

请查看上方的 [开发路线图](#-开发路线图) 部分。

---

**用 ❤️ 制作 by [SJRnhqh](https://github.com/SJRnhqh)**
