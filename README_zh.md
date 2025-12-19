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

### 📋 命令速查

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

### 检查 LaTex 环境

在编译之前，先验证你的 LaTex 安装：

```bash
lazitex --check    # 完整命令
lazitex -c         # 简短命令
```

这将检测：

- ✅ 所有已安装的 LaTex 编译器（pdflatex、xelatex、lualatex 等）
- ✅ 文献管理工具（bibtex、biber）
- ✅ 格式转换工具（dvipdfmx、ps2pdf 等）
- ✅ 包管理器（tlmgr、mpm）
- ✅ 你的 LaTex 发行版（TeX Live、MiKTeX、MacTeX）

**示例输出：**

```txt
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
           LaTex 编译环境检测结果 (Environment Check)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🖥️  操作系统: Windows (windows)
📦 LaTex 发行版: TeX Live 2025

✅ 已安装: 22/23 工具 (核心工具: 6/6)

🔨 编译器 (7/7)
────────────────────────────────────────────────────────────
  ⭐ ✓ PDFLaTex        [已安装]
      路径: C:\texlive\2025\bin\windows\pdflatex.exe
      版本: pdfTeX 3.141592653-2.6-1.40.28 (TeX Live 2025)
  ⭐ ✓ XeLaTex         [已安装]
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
  ⭐ ✓ PDFLaTex        [Installed]
      Path: C:\texlive\2025\bin\windows\pdflatex.exe
  ⭐ ✓ XeLaTex         [Installed]
🎉 Excellent! All core compilers are installed
```

**中文:**

```txt
✅ 已安装: 22/23 工具 (核心工具: 6/6)
🔨 编译器 (7/7)
  ⭐ ✓ PDFLaTex        [已安装]
      路径: C:\texlive\2025\bin\windows\pdflatex.exe
  ⭐ ✓ XeLaTex         [已安装]
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

### 安装/更新 LaTex 环境

LaziTex 支持自动安装和更新 LaTex 环境（目前支持 macOS）：

```bash
lazitex --install    # 安装或更新 LaTex 环境
lazitex -i           # 简短命令
```

**功能特点：**

- 🔍 **智能检测** - 自动检测是否已安装 LaTex
- 📦 **自动安装** - 未安装时自动安装 BasicTeX（macOS）
- 🔄 **智能更新** - 已安装时自动检查并更新包
- 💬 **交互确认** - 安装/更新前会询问确认
- 📋 **包列表** - 更新时会显示可更新的包列表

**macOS 安装示例：**

```bash
$ lazitex -i
未检测到 LaTex 环境
将使用 Homebrew 安装 BasicTeX (轻量版，约 1.5 GB)
安装命令: brew install --cask basictex
是否要安装? [y/n]: y
正在通过 Homebrew 安装 BasicTeX...
✅ BasicTeX 安装成功！
```

**更新示例：**

```bash
$ lazitex -i
检测到 LaTex 已安装
以下包可以更新:

  ...
是否要更新这些包? [y/n]: y
✅ 更新成功！
```

### 卸载 LaTex 环境

LaziTex 支持一键卸载 LaTex 环境（目前支持 macOS）：

```bash
lazitex --uninstall  # 卸载 LaTex 环境
lazitex -u           # 简短命令
```

**功能特点：**

- 🔍 **智能检测** - 自动检测是否已安装 LaTex
- 🗑️ **自动卸载** - 自动检测安装方式并执行相应卸载
- 💬 **交互确认** - 卸载前会询问确认，防止误操作
- 📋 **友好提示** - 未安装时会提示可以使用 `-i` 一键安装
- 🧹 **残留清理** - 卸载成功后尝试删除常见残留目录

**支持的卸载方式：**

- ✅ Homebrew BasicTeX - 自动卸载
- ✅ Homebrew MacTeX - 自动卸载
- ✅ MacPorts（texlive*，含 basic/latex/full）- 自动卸载
- ✅ MacTeX 官方安装 - 尝试运行官方卸载脚本，缺失时给出需手动删除的目录提示

**macOS 卸载示例：**

```bash
$ lazitex -u
检测到 LaTex 已安装
将卸载 LaTex 环境（包括所有已安装的包）
是否要卸载? [y/n]: y
正在卸载 LaTex 环境...
✅ LaTex 环境卸载成功！
```

**未安装时的提示：**

```bash
$ lazitex -u
未检测到 LaTex 环境
提示: 可以使用 'lazitex -i' 一键安装 LaTex 环境
```

### 构建与预览 LaTex 文档

LaziTex 支持一键将 LaTex 文档编译为 PDF，并提供智能预览功能：

```bash
lazitex -b main.tex          # 仅构建 LaTex 文档
lazitex -b main.tex -s       # 构建后自动打开展示
lazitex -b main.tex -o out/  # 指定输出目录
lazitex -b main.tex -o res.pdf # 指定输出完整路径
```

**功能特点：**

- 🚀 **一键构建** - 自动调用编译器（XeLaTex）并配置最优参数
- 🔍 **自动包检测与安装** - 自动从编译错误中检测缺失的 LaTex 包（`.sty` 和 `.cls` 文件），并通过 `tlmgr` 或 `mpm` 提供安装
  - **智能包名解析** - 自动查找正确的包名（例如：`xeCJK.sty` → `xecjk` 包）
  - **批量安装** - 一次性安装所有缺失的包，减少密码输入次数
  - **智能权限处理** - 先尝试不使用 sudo 安装，仅在需要时提示输入密码
  - **自动重试** - 包安装成功后自动重新编译
- 👁️ **智能展示 (`-s`)** - 构建成功后自动打开 PDF。
  - **macOS**: 优先检测并使用 **Skim.app**（支持静默刷新），若未安装则自动回退至系统默认浏览器或预览程序。
  - **Windows**: 优先检测并使用 **SumatraPDF**（支持静默刷新与实例复用），若未安装则自动回退至系统默认关联程序。
- 📁 **自定义输出 (`-o`)** - 支持指定输出目录（若路径不存在将自动递归创建）或指定完整的输出文件名
- 📁 **智能默认** - 若未指定 `-o`，则自动将生成的 PDF 和日志文件存放在与 `.tex` 源文件相同的目录下
- 📍 **路径支持** - 支持当前目录下的文件名，也支持绝对路径或相对路径
- 🔍 **类型安全** - 自动校验文件后缀，确保源文件存在
- 🔄 **逻辑复用** - 预览逻辑在 CLI 和 REPL 模式下完全一致，并为未来的实时监听模式打下了基础

**构建示例：**

```bash
$ lazitex -b report.tex -s
🚀 正在构建 LaTex 文档: report.tex
📁 工作目录: /Users/user/projects/paper
... (编译器输出) ...
✨ 构建成功！
(自动打开展示窗口)
```

**自动包安装示例：**

```bash
$ lazitex -b document.tex
🚀 正在构建 LaTex 文档: document.tex
... (编译失败) ...
🔍 检测到缺失的包: xeCJK
💡 是否自动安装？[Y/n]: y
📦 正在安装 xeCJK...
需要管理员权限，请输入密码...
✅ 安装成功！
🔄 重新编译中...
✨ 构建成功！
```

**自动包安装示例：**

```bash
$ lazitex -b document.tex
🚀 正在构建 LaTex 文档: document.tex
... (编译失败) ...
🔍 检测到缺失的包: xeCJK
💡 是否自动安装？[Y/n]: y
📦 正在安装 xeCJK...
需要管理员权限，请输入密码...
✅ 安装成功！
🔄 重新编译中...
✨ 构建成功！
```

### 其他命令

```bash
lazitex --help     # 显示帮助信息
lazitex --version  # 显示版本号
lazitex --check    # 检查 LaTex 环境
lazitex --repl     # REPL 交互模式
lazitex --tui      # 终端界面模式（即将推出）
```

**REPL 模式** - 交互式命令：

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

---

## 🗂️ 项目结构

LaziTex 采用清晰的模块化架构设计，易于扩展和维护。

```txt
lazitex/
├── 📦 go.mod                      # Go 模块定义
├── 📦 go.sum                      # Go 依赖锁定
├── 🛠️  build.sh                   # Linux/macOS 一键编译脚本
├── 🛠️  build.bat                  # Windows 一键编译脚本
├── 📁 bin/                        # 编译产物目录
│
├── 🧠 core/                       # 核心逻辑 - LaziTex 的大脑
│   ├── 🌍 env.go                  # 环境检测、安装器接口定义及安装/卸载逻辑
│   ├── 🔨 build.go                # 跨平台编译工作流
│   ├── 📦 package.go              # 自动包检测与安装
│   └── 👀 watcher.go              # 文件监听（实时预览）
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
│       ├── 🎯 tasks/              # 任务执行层（统一管理命令执行逻辑）
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

- **分层架构**：清晰的依赖层次 - `config`/`lang`（支持层）→ `core`（业务逻辑层）→ `tasks`（命令执行层）→ `ui`（交互层）

- **职责分离**：
  - `core/` - 核心业务逻辑（环境检测、构建、包管理），完全独立于 CLI，可被其他应用复用
  - `tasks/` - 命令执行层，统一管理所有命令的执行逻辑（环境、构建、预览、帮助），被 `main.go` 和 `ui/repl.go` 共享使用
  - `ui/` - 用户界面层，负责交互模式实现（REPL 的命令解析、补全、历史，TUI 的状态管理）
  - `lang/` - 国际化模块，管理多语言支持和语言偏好（包含 `Language` 类型定义）
  - `config/` - 配置管理，处理用户偏好设置

- **代码复用**：重构后的架构消除了代码重复，`tasks/` 包中的函数（如 `CheckEnvironment()`, `BuildLaTex()`, `InstallLaTexEnvironment()`）被命令行模式和 REPL 模式共享使用

- **模块化设计**：独立的包管理国际化、配置，职责清晰，易于维护和扩展

- **依赖注入**：平台特定的检查器和安装器在 `tasks/` 中创建并注入到核心逻辑，避免循环依赖

- **接口驱动**：`EnvironmentChecker` 和 `EnvironmentInstaller` 接口使平台扩展变得简单，新增平台只需实现接口

- **优先级系统**：工具按优先级分类（⭐ 核心、🔹 重要、🔸 可选），便于用户了解工具重要性

- **类型安全**：共享数据结构（如 `Language`）定义在 `lang/` 包中，简化架构，避免不必要的抽象层

---

## 🔍 可检测的 LaTex 工具

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

- [x] Windows、Linux、macOS 环境检测
- [x] 全面的 LaTex 工具检查（23 个工具）
- [x] macOS LaTex 环境自动安装/更新
- [x] macOS LaTex 环境自动卸载
- [x] 现代化的 REPL 交互模式（支持历史记录、Tab 补全、Shell 快捷命令、全语言国际化、实时预览命令）
- [x] LaTex 编译引擎（基础功能：支持 .tex 到 PDF 的一键转换，自动处理工作目录）
- [x] 智能预览系统（macOS Skim/Web 与 Windows SumatraPDF/系统默认自动适配）
- [x] 实时预览监听 (Live Preview Mode) - 支持监听文件保存并自动触发毫秒级编译与 PDF 刷新
- [x] 自动包检测与安装 - 从编译错误中自动检测缺失的包并通过 tlmgr/mpm 安装

### 进行中 🚧

- [ ] 实时预览体验优化（包括任务抢占、编译锁、错误反馈增强等细节）
- [ ] 自动宏包补全（检测缺失宏包并提示自动安装）
- [ ] 多轮编译支持（自动处理交叉引用和参考文献）
- [ ] Linux LaTex 环境自动安装/更新
- [ ] Linux LaTex 环境自动卸载

### 计划中 📋

- [ ] 本地 AI 集成，提供 LaTex 写作辅助
- [ ] 多文档项目支持
- [ ] 自定义编译配置文件

### 后续计划 🔮

- [ ] `lazitex init` 项目初始化 - 类似 `uv init` 的极简启动体验，支持从 Gitee/GitHub 快速拉取国内外优质 LaTex 模板
- [ ] 带实时预览的 TUI 模式
- [ ] 现代化的 GUI 图形界面支持
- [ ] Windows LaTex 环境自动安装/更新（Windows 安装方式复杂，优先级较低）
- [ ] Windows LaTex 环境自动卸载
- [ ] **AI 原生创作流**：构建人类与 AI 高效交互与瞬时反馈的桥梁，实现“DocuGen 式的 Prompt2PDF”体验。支持从自然语言意图到 PDF 文档的端到端实时生成与交互式编辑，并引入多智能体 (Multi-Agent) 协同驱动的文献检索、大纲规划与内容润色工作流
- [ ] **LaziTex Server & 云端垂直生态**：探索服务器端部署方案，为金融、医疗、科研等垂直领域提供“Prompt2PDF”在线交互式生成服务。支持多模态识别（如手写公式/图表转 LaTex），将非结构化意图一键桥接至专业级 PDF 文档，打造“本地免费、云端智能”的分布式文档平台。

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

**方法二：手动安装**
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

### 5. 未来会支持哪些功能？

请查看上方的 [开发路线图](#-开发路线图) 部分。

---

**用 ❤️ 制作 by [SJRnhqh](https://github.com/SJRnhqh)**
