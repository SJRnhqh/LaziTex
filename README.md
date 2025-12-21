# LaziTex 🧸

Instant LaTeX compilation across platforms — powered by Go with local AI to help you write and refine.

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://github.com/SJRnhqh/lazitex)
[![Status](https://img.shields.io/badge/Status-Active%20Development-brightgreen)](https://github.com/SJRnhqh/lazitex)

[中文文档](README_zh.md) | English

---

## ✨ Features

- 🔍 **Smart Environment Detection** - Automatically detects LaTeX installation
- 📦 **Auto Install/Update** - One-click installation or update of LaTeX environments
- 🗑️ **Auto Uninstall** - One-click uninstallation of LaTeX environments
- 🌍 **Cross-Platform** - Supports Windows, Linux, and macOS
- 🎯 **Zero Configuration** - Works out of the box
- 🚀 **Fast & Lightweight** - Built with Go
- 🔄 **Adaptive Multi-Pass Compilation** - Automatically detects and handles multi-pass compilation scenarios
- 📦 **Auto Package Management** - Automatically detects and installs missing packages
- 🌐 **Web Preview Mode** - Web real-time preview
- 🦙 **Ollama Management** - One-click Ollama management
- 🧠 **AI-Powered** (Coming Soon) - Local AI assistant
- 🎨 **Multiple Modes** - Supports TUI and REPL modes

---

## 📦 Installation

### From Source

#### Method 1: Using Build Scripts (Recommended)

The project provides multi-platform build scripts that automatically generate binaries for all platforms in the `bin/` directory:

- **Linux/macOS**:

  ```bash
  chmod +x build.sh
  ./build.sh
  ```

- **Windows**:

  ```batch
  build.bat
  ```

#### Method 2: Manual Build

```bash
git clone https://github.com/SJRnhqh/lazitex.git
cd lazitex
go build -o lazitex ./cmd/lazitex-cli
```

### Binary Releases

Download pre-built binaries from [Releases](https://github.com/SJRnhqh/lazitex/releases) (Coming Soon)

---

## 🚀 Quick Start

| Command | Description | Example |
|---------|-------------|---------|
| `-h, --help` | Show this help | `lazitex -h` |
| `-v, --version` | Show version | `lazitex -v` |
| `-c, --check` | Check LaTeX environment | `lazitex -c` |
| `-r, --repl` | Start REPL mode | `lazitex -r` |
| `-t, --tui` | Start terminal UI | `lazitex -t` |
| `-i, --install` | Install or update LaTeX environment | `lazitex -i` |
| `-u, --uninstall` | Uninstall LaTeX environment | `lazitex -u` |
| `-o, --ollama` | Ollama management (use `-c/-i/-u` flags) | `lazitex -o -c` |
| `-b, --build` | Build LaTeX document (supports `-o` output, `-s` show, `-q` quiet, `-t` tidy) | `lazitex -b main.tex [-o out/] [-s] [-q] [-t]` |
| `-p, --preview` | Live preview PDF (supports `-q` quiet mode, `-t` tidy mode, `:port` custom port) | `lazitex -p main.tex [-q] [-t] [:port]` |
| `-l, --lang` | Set language (zh/en) | `lazitex -l zh` |

---

## 📖 Usage Guide

### Language Settings

LaziTex supports **English** and **Chinese** languages. You can set the language in two ways:

#### Method 1: Runtime Change

```bash
lazitex -l zh --check    # Use Chinese for this command
lazitex --lang en --check  # Use English for this command
```

Language preference is automatically saved when using `-l` or `--lang` parameters.

#### Method 2: Configuration File

| Platform | Config Location |
|----------|----------------|
| Windows | `%APPDATA%\lazitex\config.json` |
| Linux | `~/.config/lazitex/config.json` |
| macOS | `~/Library/Application Support/lazitex/config.json` |

Example config file:

```json
{
  "language": "zh"
}
```

Edit the config file directly or use `--lang` parameter to update automatically.

### Ollama Management

One-click check, install, and uninstall Ollama (Windows & macOS). Windows uses winget, macOS supports Homebrew or official script.

```bash
# CLI Mode
$ lazitex -o -c        # Check
$ lazitex -o -i        # Install
$ lazitex -o -u        # Uninstall

# REPL Mode
lazitex> ollama -c     # Check
lazitex> ollama -i     # Install
lazitex> ollama -u     # Uninstall
```

### Build & Preview

One-click LaTeX document compilation with auto package detection & installation, adaptive multi-pass compilation, smart preview, and web real-time preview. Supports quiet mode (`-q`) and tidy mode (`-t`).

```bash
# CLI mode
$ lazitex -b report.tex -s           # Build and preview
$ lazitex -b report.tex -q           # Quiet build
$ lazitex -b report.tex -t           # Build and clean auxiliary files
$ lazitex -p report.tex              # Web preview mode (LaziView)
$ lazitex -p report.tex :3000       # Web preview mode (custom port)
$ lazitex -p report.tex -d          # Web preview mode (development mode)

# REPL mode
$ lazitex -r
lazitex> build main.tex -o out/ -s    # Build document
lazitex> preview main.tex -q          # Live preview
lazitex> help                         # Show help
```

**Web Preview Modes:**

- **LaziView Mode** (`-p`): PDF-only preview interface, perfect for multi-screen setups or external editors
- **LaziWorkspace Mode** (`-w`): Full workspace (coming soon)

### Detected LaTeX Tools

LaziTex can detect **23 LaTeX tools** across 7 categories:

| Category | Tools | Count |
|----------|-------|-------|
| 🔨 **Compilers** | pdflatex, xelatex, lualatex, latex, pdftex, tex, etex | 7 |
| 📚 **Bibliography** | bibtex, biber | 2 |
| 📇 **Indexing** | makeindex, xindy, texindy | 3 |
| 🔄 **Converters** | dvipdfmx, dvips, ps2pdf, dvisvgm | 4 |
| ⚙️ **Automation** | latexmk | 1 |
| 📦 **Package Managers** | tlmgr, mpm | 2 |
| 🔧 **Utilities** | kpsewhich, texdoc, texhash, updmap | 4 |

**Tool Priority:**

- ⭐ **Core Tools (Priority 0)**: Must install, foundation for basic functionality
  - Examples: pdflatex, xelatex, lualatex, bibtex, latexmk
  
- 🔹 **Important Tools (Priority 1)**: Recommended, provides common features
  - Examples: biber, makeindex, dvipdfmx, tlmgr
  
- 🔸 **Optional Tools (Priority 2)**: Advanced usage, specific scenarios
  - Examples: xindy, dvips, texdoc, updmap

### Platform Support

LaziTex supports Windows, Linux, and macOS:

- **Windows**: Detects TeX Live and MiKTeX, supports auto preview (prioritizes SumatraPDF)
- **Linux**: Auto-identifies distributions, provides distro-specific installation guides
- **macOS**: Supports MacTeX, Homebrew, MacPorts, supports auto install/update/uninstall (Intel and Apple Silicon)

---

## 🗂️ Project Structure

LaziTex is built with a clean, modular architecture — making it easy to extend and maintain.

```txt
lazitex/
├── 🐹 go.mod                     # Go module definition
├── 🔒 go.sum                     # Go dependencies
├── 🐚 build.sh                   # Linux/macOS build script
├── 🪟 build.bat                  # Windows build script
│
├── 🧠 core/                       # Core logic - the brain of LaziTex
│   ├── env.go                     # Environment detection, installer interface & install/uninstall logic
│   ├── build.go                   # Build workflow with Strategy Pattern (cross-platform, adaptive multi-pass)
│   ├── watcher.go                 # File watching for live preview
│   ├── ollama.go                  # Ollama manager interface definition and wrapper functions
│   ├── 🚨 errors/                 # Compilation error handling module
│   │   ├── package.go             # Auto package detection & installation
│   │   └── passes.go              # Adaptive multi-pass compilation detection
│   └── ⚡ performance/             # Compilation performance optimization module
│       ├── concurrent.go          # Concurrent optimization (parallel execution of intermediate tools)
│       └── lock.go                # Compilation lock & task management (task preemption, timeout control)
│
├── 📚 lang/                       # Internationalization module
│   └── i18n.go                    # Multi-language support (English/Chinese)
│
├── ⚙️  config/                     # Configuration management
│   └── config.go                  # User preferences & settings
│
├── 🗄️  backend/                    # Web server backend
│   ├── server.go                  # HTTP server core (with mode configuration)
│   ├── handlers.go               # HTTP request handlers (index, PDF, config API)
│   ├── routes.go                 # Route registration (includes /api/config)
│   └── sse.go                     # Server-Sent Events real-time push
│
├── 🌐 frontend/                   # Frontend resources
│   ├── embed.go                   # Static resource embedding (legacy)
│   ├── static/                    # Legacy static files
│   │   └── index.html             # Legacy preview page (double-buffering optimized)
│   └── 🟢 vue/                    # Vue 3 frontend (modern)
│       ├── src/                   # Vue source files
│       │   ├── api/               # API layer - unified backend API calls
│       │   │   ├── pdf.js         # PDF-related APIs (download, connection test)
│       │   │   ├── sse.js         # SSE-related APIs (real-time push)
│       │   │   └── config.js      # Config-related APIs (application mode)
│       │   ├── utils/             # Utility functions - pure function library
│       │   │   └── pdf.js         # PDF utility functions (parse, scale, cleanup, error formatting)
│       │   ├── components/        # Vue components
│       │   │   ├── DynamicPDFViewer.vue  # PDF preview component
│       │   │   └── ModeSwitcher.vue      # Mode switching component
│       │   ├── layouts/           # Layout components
│       │   │   ├── LaziViewLayout.vue      # PDF-only preview layout
│       │   │   └── LaziWorkspaceLayout.vue # Full workspace layout
│       │   ├── stores/            # Pinia state management
│       │   │   └── appMode.js     # Application mode store
│       │   ├── styles/            # Style files
│       │   │   ├── theme.css      # Theme colors (Nord color scheme)
│       │   │   └── dynamic-pdf-viewer.module.css  # PDF preview component styles (CSS Modules)
│       │   ├── App.vue            # Root component (mode dispatcher)
│       │   └── main.js            # Entry point
│       ├── public/                # Public assets
│       ├── index.html             # HTML template
│       ├── ⚡ vite.config.js       # Vite configuration
│       └── 📦 package.json        # Dependencies
│
├── 🎯 target/                     # Platform-specific implementations
│   ├── 🪟 win/                    # Windows-specific detection
│   │   ├── checker.go             # Detects TeX Live & MiKTeX on Windows
│   │   ├── installer.go           # Windows installer (In Development)
│   │   ├── ollama.go              # Windows Ollama manager (✅ Implemented)
│   │   └── builder.go             # Windows compilation logic (Coming Soon)
│   ├── 🐧 linux/                  # Linux-specific detection
│   │   ├── checker.go             # Detects distro-specific TeX installations
│   │   ├── installer.go           # Linux installer (In Development)
│   │   └── builder.go             # Linux compilation logic (Coming Soon)
│   └── 🍏 mac/                    # macOS-specific detection
│       ├── checker.go              # Detects MacTeX, Homebrew, MacPorts
│       ├── installer.go           # macOS installer (✅ Implemented)
│       └── builder.go             # macOS build & preview logic (✅ Implemented)
│
├── 💻 cmd/                       # Command-line interface
│   └── 🚀 lazitex-cli/
│       ├── main.go                # CLI entry point: parses commands and routes to modes
│       ├── 📋 tasks/              # Task execution layer (unified command execution logic)
│       │   ├── env.go             # Environment operations (check, install, uninstall)
│       │   ├── build.go           # Build functionality
│       │   ├── preview.go         # Live preview functionality
│       │   ├── ollama.go          # Ollama management functionality
│       │   └── help.go            # Help information
│       └── 🎨 ui/                 # User interface layer (interaction mode implementation)
│           ├── tui.go             # Terminal UI mode with Bubble Tea
│           └── repl.go            # Interactive REPL mode (command parsing, completion, history)
│
├── 📖 README.md                   # English documentation (You're reading it!)
├── 📖 README_zh.md                # Chinese documentation
└── 📄 LICENSE                     # Apache License 2.0
```

### Architecture Highlights

- **Adaptive Multi-Pass Compilation Detection** - Automatically detects and handles 6 scenarios requiring multiple compilation passes (TOC, cross-refs, bibliographies, indexes, glossaries, PDF bookmarks), no manual pass count specification needed

- **Automatic Package Detection & Installation** - Intelligently extracts missing package names from compilation errors and automatically installs them via tlmgr/mpm, delivering true zero-configuration compilation experience

- **Tool Priority System** - 23 LaTeX tools categorized by priority (⭐ Core, 🔹 Important, 🔸 Optional), helping users quickly identify critical tools and optimize installation recommendations

---

## 🌟 Roadmap

### Completed ✅

- [x] **Cross-Platform Environment Detection** - Windows, Linux, macOS environment detection
- [x] **macOS Environment Management** - Auto install/update/uninstall
- [x] **REPL Interactive Mode** - Command history, Tab completion, full i18n support
- [x] **Build System** - One-click compilation, smart preview, auto package detection & installation, adaptive multi-pass compilation
- [x] **Web Preview & Frontend** - Vue 3 frontend, SSE real-time refresh, dual-mode architecture
- [x] **PDF Virtual Scrolling** - Multi-page PDF virtual scrolling preview, improved performance for large documents
- [x] **Ollama Management** - One-click check, install, and uninstall Ollama
- [x] **Vue Frontend Production Build** - Single binary deployment via build script

### In Progress 🚧

- [ ] **Live Preview Optimization** - Enhanced error feedback
- [x] **Compilation Performance Optimization** - Concurrent execution, task preemption, timeout control
- [ ] **Linux Environment Management** - Auto install/update/uninstall
- [ ] **PDF Preview Enhancement** - Pagination mode, page navigation, zoom controls
- [ ] **AI Architecture Implementation** - Modular architecture for LLM and Agent management

### Planned 📋

- [ ] **AI Agent Environment Management** - Multi-agent creation, switching, and configuration, similar to conda environment management
- [ ] **AI Integration (CLI Internal)** - Compilation error diagnosis, code generation & optimization, smart completion
- [ ] **Web Workspace Enhancement** - File management, code editor, terminal integration, AI chat interface
- [ ] **Project Feature Enhancement** - Project initialization, multi-document support, custom compilation profiles, collaborative editing
- [ ] **TUI Enhancement** - Terminal UI feature enhancement and interaction optimization
- [ ] **Platform Support Extension** - Windows environment management, modern GUI support
- [ ] **AI-Native Authoring Flow** - Prompt2PDF, multi-agent collaboration, cloud ecosystem

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

Apache License 2.0 - Includes explicit patent grants to protect developers and users.

See [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for beautiful terminal UIs
- Inspired by the LaTeX community's dedication to beautiful typesetting

---

## 💡 FAQ

### Why can't LaziTex detect my LaTeX installation?

- **Windows**: Make sure the LaTeX installation path is added to your PATH environment variable
- **Linux**: Try running `which pdflatex` to confirm the installation location
- **macOS**: If using Homebrew, make sure you've run `brew link` command

### Can I add other languages?

Currently only English and Chinese are supported. If you need other language support, please submit an Issue or Pull Request!

---

**Made with ❤️ by [SJRnhqh](https://github.com/SJRnhqh)**
