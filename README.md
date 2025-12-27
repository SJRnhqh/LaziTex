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
- 🤖 **LLM Management** - Complete LLM provider management (register, link, switch, test, remove, list)
- ⚙️ **Configuration Management** - Cross-platform configuration file management
- 🎨 **Multiple Modes** - Supports TUI and REPL modes

---

## 📦 Installation

### Build Dependencies

Before building from source, ensure you have the following tools installed:

- **Node.js** (includes npm)
  - macOS: `brew install node`
  - Linux: Use your distribution's package manager
  - Windows: Download from [nodejs.org](https://nodejs.org/)
  
- **Go** (version 1.25+)
  - macOS: `brew install go`
  - Linux: Use your distribution's package manager
  - Windows: Download from [go.dev](https://go.dev/dl/)
  
- **Go Proxy Configuration** (Recommended for users in China)

  ```bash
  go env -w GOPROXY=https://goproxy.cn,direct
  ```

### Runtime Dependencies

- **LaTeX Distribution** (e.g., MacTeX, TeX Live, MiKTeX)
  - macOS: `brew install --cask mactex` (or use Homebrew: `brew install basictex`)
  - Linux: Use your distribution's package manager
  - Windows: Download and install [TeX Live](https://www.tug.org/texlive/) or [MiKTeX](https://miktex.org/)

> **Note**: The build scripts will automatically check for build dependencies. If any are missing, you'll receive helpful installation instructions.

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
| --------- | ------------- | --------- |
| `-h, --help` | Show this help | `lazitex -h` |
| `-v, --version` | Show version | `lazitex -v` |
| `-r, --repl` | Start REPL mode | `lazitex -r` |
| `-t, --tui` | Start terminal UI | `lazitex -t` |
| `-x, --latex` | LaTeX environment management (use `-c/-i/-u` flags) | `lazitex -x -c` |
| `-o, --ollama` | Ollama management (use `-c/-i/-u` flags) | `lazitex -o -c` |
| `-b, --build` | Build LaTeX document (supports `-o` output, `-s` show, `-q` quiet, `-t` tidy) | `lazitex -b main.tex [-o out/] [-s] [-q] [-t]` |
| `-p, --preview` | Live preview PDF (supports `-q` quiet mode, `-t` tidy mode, `:port` custom port) | `lazitex -p main.tex [-q] [-t] [:port]` |
| `config` | Open configuration file | `lazitex config` |
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
| ---------- | ---------------- |
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

### LaTeX Environment Management

Cross-platform LaTeX environment management with automatic detection, installation, and uninstallation.

```bash
# CLI Mode
$ lazitex -x -c        # Check
$ lazitex -x -i        # Install
$ lazitex -x -u        # Uninstall

# REPL Mode
lazitex> latex -c      # Check
lazitex> latex -i      # Install
lazitex> latex -u      # Uninstall
```

### Ollama Management

Cross-platform Ollama management with check, install, and uninstall support.

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

### LLM Management

Complete LLM provider management with registration, activation, switching, testing, and removal. Supports interactive selection and multi-activation.

```bash
# REPL Mode
lazitex> llm link              # Register new LLM provider
lazitex> llm link <name>       # Activate LLM provider
lazitex> llm unlink            # Deactivate current LLM
lazitex> llm unlink <name>     # Deactivate specific LLM
lazitex> llm switch            # Switch foreground LLM (interactive)
lazitex> llm switch <name>     # Switch to specific LLM
lazitex> llm test              # Test current LLM
lazitex> llm test <name>       # Test specific LLM
lazitex> llm remove            # Remove LLM provider (interactive)
lazitex> llm list              # List all registered LLMs
lazitex> llm list <name>       # List specific LLM
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
| ---------- | ------- | ------- |
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

- **Windows**: Detects TeX Live and MiKTeX, supports auto preview (prioritizes SumatraPDF), Ollama management via winget
- **Linux**: Auto-identifies distributions, provides distro-specific installation guides, Ollama management via multiple package managers or official script
- **macOS**: Supports MacTeX, Homebrew, MacPorts, supports auto install/update/uninstall (Intel and Apple Silicon), Ollama management via Homebrew or official script

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
│   ├── latex.go                   # LaTeX manager interface definition, environment detection & wrapper functions
│   ├── build.go                   # Build workflow with Strategy Pattern (cross-platform, adaptive multi-pass)
│   ├── watcher.go                 # File watching for live preview
│   ├── ollama.go                  # Ollama manager interface definition and wrapper functions
│   ├── 🚨 errors/                 # Compilation error handling module
│   │   ├── formatter.go           # Error extraction and formatting
│   │   ├── package.go             # Auto package detection & installation
│   │   └── passes.go              # Adaptive multi-pass compilation detection
│   ├── ⚡ performance/             # Compilation performance optimization module
│   │   ├── concurrent.go          # Concurrent optimization (parallel execution of intermediate tools)
│   │   └── lock.go                # Compilation lock & task management (task preemption, timeout control)
│   └── 🤖 ai/                      # AI & LLM management module
│       └── llm/                    # LLM provider management
│           ├── registry.go         # LLM provider registration
│           ├── link.go             # LLM activation & linking
│           ├── unlink.go           # LLM deactivation
│           ├── switch.go           # Foreground LLM switching
│           ├── test.go             # LLM connection testing
│           ├── remove.go           # LLM provider removal
│           ├── list.go             # LLM provider listing (table format)
│           └── select.go           # Interactive selection utilities
│
├── 📚 lang/                       # Internationalization module
│   └── i18n.go                    # Multi-language support (English/Chinese)
│
├── ⚙️  config/                     # Configuration management
│   ├── common.go                   # Application configuration management
│   └── llm.go                      # LLM provider configuration management
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
│       │   │   ├── pdf.js         # PDF utility functions (parse, scale, cleanup, error formatting)
│       │   │   └── pdfDisplay.js  # PDF display management (virtual scroll, pagination, config)
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
│   ├── 🪟 win/                    # Windows platform
│   │   ├── checker.go             # Detects TeX Live & MiKTeX on Windows
│   │   ├── installer.go           # Windows LaTeX installer (In Development)
│   │   ├── latex.go               # Windows LaTeX manager (unified interface)
│   │   ├── ollama.go              # Windows Ollama manager (✅ Implemented)
│   │   └── builder.go             # Windows compilation logic (Coming Soon)
│   ├── 🐧 linux/                  # Linux platform
│   │   ├── checker.go             # Detects distro-specific TeX installations
│   │   ├── installer.go           # Linux LaTeX installer (In Development)
│   │   ├── latex.go               # Linux LaTeX manager (unified interface)
│   │   ├── ollama.go              # Linux Ollama manager (✅ Implemented)
│   │   └── builder.go             # Linux compilation logic (Coming Soon)
│   └── 🍏 mac/                    # macOS platform
│       ├── checker.go              # Detects MacTeX, Homebrew, MacPorts
│       ├── installer.go           # macOS LaTeX installer (✅ Implemented)
│       ├── latex.go               # macOS LaTeX manager (unified interface)
│       └── builder.go             # macOS build & preview logic (✅ Implemented)
│
├── 💻 cmd/                       # Command-line interface
│   └── 🚀 lazitex-cli/
│       ├── main.go                # CLI entry point: parses commands and routes to modes
│       ├── 📋 tasks/              # Task execution layer (unified command execution logic)
│       │   ├── build.go           # Build functionality
│       │   ├── preview.go         # Live preview functionality
│       │   ├── ollama.go          # Ollama management functionality
│       │   ├── latex.go           # LaTeX environment management
│       │   ├── llm.go             # LLM management functionality
│       │   ├── config.go          # Configuration file management
│       │   └── help.go            # Help information
│       ├── 🧠 internal/           # Internal command processing (unified action handlers)
│       │   ├── common.go          # Common utilities and action handlers
│       │   ├── latex.go           # LaTeX command processing logic
│       │   └── ollama.go          # Ollama command processing logic
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

- **Unified Management Architecture** - Consistent architecture pattern for LaTeX and Ollama management: unified manager interface (LaTeXManager/OllamaManager) with platform-specific implementations, following the same flow from command parsing to task execution to core logic to platform implementation

---

## 🌟 Roadmap

### Completed ✅

- [x] Cross-platform environment detection & management (macOS auto install/update/uninstall)
- [x] Build system with adaptive multi-pass compilation & auto package management
- [x] Web preview frontend (Vue 3, SSE real-time refresh, PDF virtual scrolling)
- [x] Ollama management (cross-platform check, install, uninstall)
- [x] LLM management (register, link, switch, test, remove, list with full i18n support)
- [x] Interactive REPL mode (command history, Tab completion, i18n)

### In Progress 🚧

- [ ] AI application features (error diagnosis, code generation, smart completion)
- [ ] Agent system framework
- [ ] Linux environment auto management

### Planned 📋

- [ ] Web workspace (file management, editor, terminal, AI chat)
- [ ] Project features (initialization, multi-document support)
- [ ] Windows environment auto management

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
