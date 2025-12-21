# LaziTex 🧸

Instant LaTeX compilation across platforms — powered by Go with local AI to help you write and refine.

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://github.com/SJRnhqh/lazitex)
[![Status](https://img.shields.io/badge/Status-Active%20Development-brightgreen)](https://github.com/SJRnhqh/lazitex)

[中文文档](README_zh.md) | English

---

## ✨ Features

- 🔍 **Smart Environment Detection** - Automatically detects your LaTeX installation and provides detailed diagnostics
- 📦 **Auto Install/Update** - One-click installation or update of LaTeX environments (macOS supported)
- 🗑️ **Auto Uninstall** - One-click uninstallation of LaTeX environments with smart installation method detection (macOS supported)
- 🌍 **Cross-Platform** - Seamless support for Windows, Linux, and macOS
- 🎯 **Zero Configuration** - Works out of the box with TeX Live, MiKTeX, and MacTeX
- 🚀 **Fast & Lightweight** - Built with Go for blazing-fast compilation
- 💬 **Interactive REPL** - Interactive commands for check, install, uninstall, language switching, and more
- 🔄 **Adaptive Multi-Pass Compilation** - Automatically detects and handles multiple compilation passes for cross-references, table of contents, bibliographies, indexes, and glossaries
- 📦 **Auto Package Management** - Automatically detects missing packages from compilation errors and installs them via tlmgr/mpm
- 🌐 **Web Preview Mode** - Local HTTP server + browser preview, SSE real-time auto-refresh, double-buffering optimization, Overleaf-style web preview experience
- 🦙 **Ollama Management** - One-click check, install, and uninstall Ollama (Windows & macOS supported), providing foundation for AI features
- 🧠 **AI-Powered** (Coming Soon) - Local AI assistance for writing and refining LaTeX documents
- 🎨 **Multiple Modes** - TUI (Terminal UI) and REPL modes for different workflows

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

**Features:**

- One-click build - Automatically invokes XeLaTeX compiler
- Auto package detection & installation - Detects missing packages from compilation errors and installs them automatically
- Adaptive multi-pass compilation - Automatically handles cross-references, TOC, bibliographies, indexes, etc.
- Smart preview - Automatically opens PDF after successful build (macOS prioritizes Skim, Windows prioritizes SumatraPDF)
- Web preview mode - Local HTTP server + browser preview, SSE real-time auto-refresh, zero-latency update experience
- Custom output - Supports specifying output directory or full path
- Quiet mode - Use `-q` flag to suppress compiler verbose output, showing only key information and errors
- Tidy mode - Use `-t` flag to automatically clean auxiliary files (.log, .aux, .toc, .out, etc.) after successful compilation, keeping only the PDF

**Build Example:**

```bash
$ lazitex -b report.tex -s
🚀 Building LaTeX document: report.tex
📁 Working directory: /Users/user/projects/paper
... (compiler output) ...
✨ Build successful!
(Automatically opening show window)
```

**Quiet Mode Example:**

```bash
$ lazitex -b report.tex -q
🚀 Building LaTeX document: report.tex
📁 Working directory: /Users/user/projects/paper
✨ Build successful!
# Compiler verbose output is suppressed, only key information is shown
```

**Tidy Mode Example:**

```bash
$ lazitex -b report.tex -t
🚀 Building LaTeX document: report.tex
📁 Working directory: /Users/user/projects/paper
✨ Build successful!
🧹 Auxiliary files cleaned
# Automatically removes .log, .aux, .toc, .out, etc., keeping only the PDF
```

**Note:** The `-t/--tidy` flag only cleans auxiliary files when compilation succeeds. If compilation fails, all files are preserved for debugging.

**Web Preview Mode:**

LaziTex frontend (LaziHub) provides two modes:

- **LaziView Mode** (`-p`): PDF-only preview interface, perfect for multi-screen setups or when you prefer external editors like VS Code or Neovim
- **LaziWorkspace Mode** (`-w`): Full workspace with code editor, AI chat, file browser, and PDF preview (coming soon)

```bash
# Preview mode (LaziView - PDF only)
$ lazitex -p report.tex
🚀 Building LaTeX document: report.tex
✨ Build successful!
💡 Vue Frontend: http://localhost:5173
💡 Backend API: http://localhost:8080
🔗 Opening browser: http://localhost:5173
🌐 Server starting on port 8080
👀 Watching file: report.tex (press Ctrl+C to stop)

# Preview mode with quiet compilation
$ lazitex -p report.tex -q

# Preview mode with tidy mode (clean auxiliary files after compilation)
$ lazitex -p report.tex -t

# Preview mode with custom port (default is 8080)
$ lazitex -p report.tex :3000
$ lazitex -p report.tex report.tex:3000  # Alternative format
$ lazitex -p report.tex -q -t :3000      # Port can be anywhere in arguments

# Preview mode with development mode (-d flag)
$ lazitex -p report.tex -d

# After modifying the file, the browser will automatically refresh to show the latest PDF
```

**Frontend Architecture:**

- **Dual-Mode System**: LaziHub frontend supports both LaziView (PDF-only) and LaziWorkspace (full workspace) modes
- **Mode Switching**: Users can switch between modes at runtime via the UI
- **Backend API**: The `/api/config` endpoint returns the initial mode based on CLI parameters (`-p` for view mode)
- **State Management**: Uses Pinia for reactive state management
- **Layout Components**: Modular layout architecture (LaziViewLayout, LaziWorkspaceLayout)

### REPL Mode

Interactive command-line interface with command history, Tab completion, and full i18n support.

```bash
$ lazitex -r
lazitex> build main.tex -o out/ -s    # Build document
lazitex> preview main.tex -q          # Live preview
lazitex> ollama -c                    # Check Ollama
lazitex> help                         # Show help
```

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

---

## 🗂️ Project Structure

LaziTex is built with a clean, modular architecture — making it easy to extend and maintain.

```txt
lazitex/
├── 🗂️  go.mod                     # Go module definition
├── 🔒  go.sum                     # Go dependencies
├── 🛠️  build.sh                   # Linux/macOS build script
├── 🛠️  build.bat                  # Windows build script
│
├── 🧠 core/                       # Core logic - the brain of LaziTex
│   ├── 🌍 env.go                  # Environment detection, installer interface & install/uninstall logic
│   ├── 🔨 build.go                # Build workflow with Strategy Pattern (cross-platform, adaptive multi-pass)
│   ├── 👀 watcher.go              # File watching for live preview
│   ├── 🦙 ollama.go               # Ollama manager interface definition and wrapper functions
│   └── 🚨 errors/                 # Compilation error handling module
│       ├── 📦 package.go           # Auto package detection & installation
│       └── 🔄 passes.go            # Adaptive multi-pass compilation detection
│
├── 📚 lang/                       # Internationalization module
│   └── i18n.go                    # Multi-language support (English/Chinese)
│
├── ⚙️  config/                     # Configuration management
│   └── config.go                  # User preferences & settings
│
├── 🗄️  backend/                    # Web server backend
│   ├── server.go                  # HTTP server core (with mode configuration)
│   ├── handlers.go                # HTTP request handlers (index, PDF, config API)
│   ├── routes.go                  # Route registration (includes /api/config)
│   └── sse.go                     # Server-Sent Events real-time push
│
├── 🌐 frontend/                   # Frontend resources
│   ├── embed.go                   # Static resource embedding (legacy)
│   ├── static/                    # Legacy static files
│   │   └── index.html             # Legacy preview page (double-buffering optimized)
│   └── vue/                       # Vue 3 frontend (modern)
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
│       ├── vite.config.js         # Vite configuration
│       └── package.json           # Dependencies
│
├── 🎯 target/                     # Platform-specific implementations
│   ├── 🪟 win/                    # Windows-specific detection
│   │   ├── 🔍 checker.go          # Detects TeX Live & MiKTeX on Windows
│   │   ├── 📦 installer.go        # Windows installer (In Development)
│   │   ├── 🦙 ollama.go           # Windows Ollama manager (✅ Implemented)
│   │   └── 🛠️  builder.go         # Windows compilation logic (Coming Soon)
│   ├── 🐧 linux/                  # Linux-specific detection
│   │   ├── 🔍 checker.go          # Detects distro-specific TeX installations
│   │   ├── 📦 installer.go        # Linux installer (In Development)
│   │   └── 🛠️  builder.go         # Linux compilation logic (Coming Soon)
│   └── 🍏 mac/                    # macOS-specific detection
│       ├── 🔍 checker.go          # Detects MacTeX, Homebrew, MacPorts
│       ├── 📦 installer.go        # macOS installer (✅ Implemented)
│       └── 🛠️  builder.go         # macOS build & preview logic (✅ Implemented)
│
├── 💻 cmd/                       # Command-line interface
│   └── 🚀 lazitex-cli/
│       ├── 📄 main.go             # CLI entry point: parses commands and routes to modes
│       ├── 📋 tasks/              # Task execution layer (unified command execution logic)
│       │   ├── 🌍 env.go          # Environment operations (check, install, uninstall)
│       │   ├── 🔨 build.go        # Build functionality
│       │   ├── 👀 preview.go      # Live preview functionality
│       │   ├── 🦙 ollama.go       # Ollama management functionality
│       │   └── 📖 help.go         # Help information
│       └── 🎨 ui/                 # User interface layer (interaction mode implementation)
│           ├── 🖼️  tui.go         # Terminal UI mode with Bubble Tea
│           └── 💬 repl.go         # Interactive REPL mode (command parsing, completion, history)
│
├── 📖 README.md                   # English documentation (You're reading it!)
├── 📖 README_zh.md                # Chinese documentation
└── 📄 LICENSE                     # Apache License 2.0
```

### Architecture Highlights

- **Compilation Strategy Pattern** - Uses Strategy Pattern (`compileStrategy` interface) to handle different compilation scenarios. Default strategy automatically handles package errors and multi-pass compilation. Future AI-powered strategies can be easily added without modifying the core compilation loop

- **Compilation Error Handling Module** - `core/errors/` modularly handles compilation issues: automatic package detection & installation, adaptive multi-pass compilation detection (covers 6 scenarios: TOC, cross-refs, bibliographies, indexes, glossaries, PDF bookmarks)

- **Web Preview Architecture** - Lightweight HTTP server based on Go standard library, using Server-Sent Events (SSE) for real-time push, frontend double-buffering eliminates refresh flicker, providing smooth preview experience. LaziHub frontend supports dual-mode architecture: LaziView (PDF-only) and LaziWorkspace (full workspace)
- **Frontend State Management** - Uses Pinia for reactive state management, enabling seamless mode switching and component communication
- **Frontend Code Organization** - Clean directory structure: `api/` for unified backend API calls, `utils/` for pure function utilities, `styles/` for centralized styling, enabling code reuse and maintainability

- **Tool Priority System** - Tools categorized by priority (⭐ Core, 🔹 Important, 🔸 Optional), helping users quickly identify critical tools

---

## 🛠️ Platform-Specific Features

### Windows

- Detects **TeX Live** and **MiKTeX** installations
- Searches common install paths (C:\texlive, C:\Program Files\MiKTeX, etc.)
- Checks user-specific installations in AppData
- **Auto Preview Support** - Automatically invokes the system default PDF viewer after a successful build
- **Smart Preview Enhancement** - Prioritizes detection and invocation of **SumatraPDF** (supporting silent refresh with -reuse-instance), falling back to the system default if not found

### Linux

- Identifies Linux distribution automatically
- Provides distro-specific installation instructions (apt, dnf, pacman)
- Detects package manager vs manual installations
- Supports both x86_64 and aarch64 architectures

### macOS

- Detects **MacTeX**, **Homebrew**, and **MacPorts** installations
- Shows macOS version information
- **Auto Install/Update** - One-click BasicTeX installation via Homebrew
- **Auto Uninstall** - Smart detection of installation method with one-click uninstallation
- **Ollama Management** - One-click check, install, and uninstall Ollama via Homebrew or official script
- Provides optimized installation guides based on detected package managers
- Supports both Intel and Apple Silicon chips

---

## 🌟 Roadmap

### Completed ✅

- [x] **Cross-Platform Environment Detection** - Windows, Linux, macOS environment detection and 23 tools checking
- [x] **macOS Environment Management** - Auto install/update/uninstall, supports Homebrew, MacPorts, MacTeX
- [x] **REPL Interactive Mode** - Command history, Tab completion, Shell shortcuts, full i18n support
- [x] **Build System** - One-click compilation, smart preview (macOS Skim/Windows SumatraPDF), auto package detection & installation, adaptive multi-pass compilation (cross-refs, TOC, bibliographies, etc.)
- [x] **Web Preview & Frontend** - LaziHub frontend (Vue 3 + PDF.js), web preview mode with SSE real-time refresh, dual-mode architecture (LaziView/LaziWorkspace), preview mode supports `-q/-t/-d/:port` flags
- [x] **Frontend Code Refactoring** - Code organization optimization: `api/` directory for unified backend API calls, `utils/` directory for pure function utilities, `styles/` directory for centralized styling, enabling code reuse and maintainability
- [x] **Ollama Management (Windows & macOS)** - One-click check, install, and uninstall Ollama, smart detection of installation status and service running status

### In Progress 🚧

- [ ] **Live Preview Optimization** - Task preemption, compilation locks, enhanced error feedback
- [ ] **Linux Environment Management** - Auto install/update/uninstall

### Planned 📋

- [ ] **AI Agent Environment Management** - Manage AI Agents like conda environments, supporting multi-agent creation, switching, and configuration. Use `ai` command to manage agents (create, list, switch LLM), and `agent` command to invoke agents (ask, chat, generate). Each agent can be configured with independent LLM providers, tool sets, and system prompts, enabling quick switching between different scenarios (compiler expert, writing assistant, error diagnosis, etc.)
- [ ] **AI Integration (CLI Internal)** - Compilation error diagnosis, code generation & optimization, smart completion. CLI provides quick access via `-a` flag using the currently configured agent, while REPL offers full agent management and invocation capabilities
- [ ] **Quiet Mode Error Formatting** - In quiet mode, extract and format compilation errors with LaziTex-style error messages (with icons) instead of raw compiler output
- [ ] **Multi-document Project Support** - Automatically handle file dependencies and compilation order
- [ ] **Custom Compilation Profiles** - Project-level configuration files (compiler selection, output directories, compilation parameters, etc.)

### Future 🔮

- [ ] **Project Initialization** - `lazitex init` minimalist startup, fast template fetching
- [ ] **Modern GUI Support** - Modern graphical interface
- [ ] **Windows Environment Management** - Auto install/update/uninstall (lower priority)
- [ ] **AI-Native Authoring Flow** - Prompt2PDF experience, multi-agent collaborative workflows
- [ ] **Cloud Vertical Ecosystem** - Server-side deployment, online interactive generation services

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

### 1. Why can't LaziTex detect my LaTeX installation?

- **Windows**: Make sure the LaTeX installation path is added to your PATH environment variable
- **Linux**: Try running `which pdflatex` to confirm the installation location
- **macOS**: If using Homebrew, make sure you've run `brew link` command

### 2. Which LaTeX distributions are supported?

Currently supported:

- TeX Live (all platforms)
- MiKTeX (Windows)
- MacTeX (macOS)

### 3. How to install missing tools?

Method 1: Use LaziTex auto-install (Recommended)

```bash
lazitex -i    # Auto-install or update LaTeX environment (macOS supported)
```

Method 2: Manual installation

After running `lazitex --check`, if tools are missing, the system will automatically display installation guides for your platform.

### 4. How to uninstall LaTeX environment?

Use LaziTex one-click uninstall (Recommended, currently macOS):

```bash
lazitex -u    # Auto-uninstall LaTeX environment (macOS supported)
```

**Features:**

- Automatically detects installation method (Homebrew, MacPorts, MacTeX official, etc.)
- Interactive confirmation to prevent accidental removal
- Suggests using `-i` to install when LaTeX is not installed

**Manual uninstallation:**

- Homebrew: `brew uninstall --cask basictex` or `brew uninstall --cask mactex`
- MacPorts: `sudo port uninstall texlive`
- MacTeX Official: Manually delete `/Library/TeX/` directory

### 5. How to check current language?

```bash
# View config file
# Windows
type %APPDATA%\lazitex\config.json
# Linux/macOS
cat ~/.config/lazitex/config.json
```

### 6. How to reset to default language (English)?

```bash
lazitex -l en --help
# Or use long form
lazitex --lang en --help
# Or delete the config file
```

### 7. Can I add other languages?

Currently only English and Chinese are supported. If you need other language support, please submit an Issue or Pull Request!

### 8. What features will be supported in the future?

Please check the [Roadmap](#-roadmap) section above.

---

**Made with ❤️ by [SJRnhqh](https://github.com/SJRnhqh)**
