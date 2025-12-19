# LaziTex 🧸

Instant LaTex compilation across platforms — powered by Go with local AI to help you write and refine.

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://github.com/SJRnhqh/lazitex)
[![Status](https://img.shields.io/badge/Status-Active%20Development-brightgreen)](https://github.com/SJRnhqh/lazitex)

[中文文档](README_zh.md) | English

---

## ✨ Features

- 🔍 **Smart Environment Detection** - Automatically detects your LaTex installation and provides detailed diagnostics
- 📦 **Auto Install/Update** - One-click installation or update of LaTex environments (macOS supported)
- 🗑️ **Auto Uninstall** - One-click uninstallation of LaTex environments with smart installation method detection (macOS supported)
- 🌍 **Cross-Platform** - Seamless support for Windows, Linux, and macOS
- 🎯 **Zero Configuration** - Works out of the box with TeX Live, MiKTeX, and MacTeX
- 🚀 **Fast & Lightweight** - Built with Go for blazing-fast compilation
- 💬 **Interactive REPL** - Interactive commands for check, install, uninstall, language switching, and more
- 🔄 **Adaptive Multi-Pass Compilation** - Automatically detects and handles multiple compilation passes for cross-references, table of contents, bibliographies, indexes, and glossaries
- 📦 **Auto Package Management** - Automatically detects missing packages from compilation errors and installs them via tlmgr/mpm
- 🌐 **Web Preview Mode** - Local HTTP server + browser preview, Overleaf-style web preview experience
- 🧠 **AI-Powered** (Coming Soon) - Local AI assistance for writing and refining LaTex documents
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
| `-c, --check` | Check LaTex environment | `lazitex -c` |
| `-i, --install` | Install/update LaTex | `lazitex -i` |
| `-u, --uninstall` | Uninstall LaTex | `lazitex -u` |
| `-b, --build` | Build LaTex document (supports `-o` output, `-s` show) | `lazitex -b main.tex [-o out/] [-s]` |
| `-p, --preview` | Live preview PDF (watches for saves and refreshes) | `lazitex -p main.tex` |
| `-r, --repl` | Start REPL mode | `lazitex -r` |
| `-l, --lang` | Set language | `lazitex -l zh` |
| `-h, --help` | Show help | `lazitex -h` |
| `-v, --version` | Show version | `lazitex -v` |

---

## 📖 Usage Guide

### Language Settings

LaziTex supports **English** and **Chinese** languages. You can set the language in two ways:

Method 1: Runtime Change

```bash
lazitex -l zh --check    # Use Chinese for this command
lazitex --lang en --check  # Use English for this command
```

Language preference is automatically saved when using `-l` or `--lang` parameters.

Method 2: Configuration File

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

### Build & Preview

**Features:**

- One-click build - Automatically invokes XeLaTex compiler
- Auto package detection & installation - Detects missing packages from compilation errors and installs them automatically
- Adaptive multi-pass compilation - Automatically handles cross-references, TOC, bibliographies, indexes, etc.
- Smart preview - Automatically opens PDF after successful build (macOS prioritizes Skim, Windows prioritizes SumatraPDF)
- Custom output - Supports specifying output directory or full path

**Build Example:**

```bash
$ lazitex -b report.tex -s
🚀 Building LaTex document: report.tex
📁 Working directory: /Users/user/projects/paper
... (compiler output) ...
✨ Build successful!
(Automatically opening show window)
```

### REPL Mode

```bash
lazitex> help                      # Show available commands
lazitex> check                     # Check LaTex environment
lazitex> install                   # Install or update LaTex environment
lazitex> uninstall                 # Uninstall LaTex environment
lazitex> build main.tex -o out/ -s # Build LaTex document with show and output path
lazitex> lang zh                   # Switch to Chinese
lazitex> lang en                   # Switch to English
lazitex> lang                      # Show current language
lazitex> version                   # Show version
lazitex> exit                      # Exit REPL
lazitex> cd /tmp                   # Change directory (default to home)
lazitex> ls                        # List contents (supports path arg)
lazitex> pwd                       # Print working directory
lazitex> clear                     # Clear the screen
lazitex> cat file.log              # View file content (supports .tex, .log, .aux)
```

**Modern REPL Experience:**

- 🕒 **Command History** - Use Up/Down arrows to navigate through previous commands
- ⌨️ **Cursor Control** - Full support for Left/Right arrows, Home/End, Ctrl+A/E, etc.
- 🎯 **Smart Completion** - Context-aware Tab completion:
  - **Command Completion**: Auto-completes built-in commands
  - **Context-Aware Pathing**: `cd` suggests directories; `build` suggests `.tex` files; `cat` suggests `.tex/.log/.aux`
- 🌐 **Full i18n Support** - All command help, error messages, and usage tips support instant language switching

Language switching in REPL is instant and persists for future sessions!

### Detected LaTex Tools

LaziTex can detect **23 LaTex tools** across 7 categories:

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
├── 🗂️  go.mod                      # Go module definition
├── 🔒 go.sum                      # Go dependencies (lock file)
├── 🛠️  build.sh                   # Linux/macOS build script
├── 🛠️  build.bat                  # Windows build script
│
├── 🧠 core/                       # Core logic - the brain of LaziTex
│   ├── 🌍 env.go                  # Environment detection, installer interface & install/uninstall logic
│   ├── 🔨 build.go                # Build workflow with Strategy Pattern (cross-platform, adaptive multi-pass)
│   ├── 👀 watcher.go              # File watching for live preview
│   └── 🚨 errors/                 # Compilation error handling module
│       ├── 📦 package.go           # Auto package detection & installation
│       └── 🔄 passes.go            # Adaptive multi-pass compilation detection
│
├── 🌐 lang/                       # Internationalization module
│   └── i18n.go                    # Multi-language support (English/Chinese)
│
├── ⚙️  config/                     # Configuration management
│   └── config.go                  # User preferences & settings
│
├── 🎯 target/                     # Platform-specific implementations
│   ├── 🪟 win/                    # Windows-specific detection
│   │   ├── 🔍 checker.go          # Detects TeX Live & MiKTeX on Windows
│   │   ├── 📦 installer.go        # Windows installer (In Development)
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
├── 🖥️  cmd/                       # Command-line interface
│   └── 🚀 lazitex-cli/
│       ├── 📄 main.go             # CLI entry point: parses commands and routes to modes
│       ├── 📋 tasks/              # Task execution layer (unified command execution logic)
│       │   ├── 🌍 env.go          # Environment operations (check, install, uninstall)
│       │   ├── 🔨 build.go        # Build functionality
│       │   ├── 👀 preview.go      # Live preview functionality
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
- Provides optimized installation guides based on detected package managers
- Supports both Intel and Apple Silicon chips

---

## 🌟 Roadmap

### Completed ✅

- [x] **Cross-Platform Environment Detection** - Windows, Linux, macOS environment detection and 23 tools checking
- [x] **macOS Environment Management** - Auto install/update/uninstall, supports Homebrew, MacPorts, MacTeX
- [x] **REPL Interactive Mode** - Command history, Tab completion, Shell shortcuts, full i18n support
- [x] **Build & Preview System** - One-click compilation, smart preview (macOS Skim/Windows SumatraPDF), live preview watching
- [x] **Smart Compilation Optimization** - Auto package detection & installation, adaptive multi-pass compilation (cross-refs, TOC, bibliographies, etc.)
- [x] **Web Preview Mode** - Local HTTP server + browser preview, Overleaf-style web preview experience

### In Progress 🚧

- [ ] **Live Preview Optimization** - Task preemption, compilation locks, enhanced error feedback
- [ ] **Linux Environment Management** - Auto install/update/uninstall

### Planned 📋

- [ ] **AI Integration (CLI Internal)** - Compilation error diagnosis, code generation & optimization, smart completion
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
- Inspired by the LaTex community's dedication to beautiful typesetting

---

## 💡 FAQ

### 1. Why can't LaziTex detect my LaTex installation?

- **Windows**: Make sure the LaTex installation path is added to your PATH environment variable
- **Linux**: Try running `which pdflatex` to confirm the installation location
- **macOS**: If using Homebrew, make sure you've run `brew link` command

### 2. Which LaTex distributions are supported?

Currently supported:

- TeX Live (all platforms)
- MiKTeX (Windows)
- MacTeX (macOS)

### 3. How to install missing tools?

Method 1: Use LaziTex auto-install (Recommended)

```bash
lazitex -i    # Auto-install or update LaTex environment (macOS supported)
```

Method 2: Manual installation

After running `lazitex --check`, if tools are missing, the system will automatically display installation guides for your platform.

### 4. How to uninstall LaTex environment?

Use LaziTex one-click uninstall (Recommended, currently macOS):

```bash
lazitex -u    # Auto-uninstall LaTex environment (macOS supported)
```

**Features:**

- Automatically detects installation method (Homebrew, MacPorts, MacTeX official, etc.)
- Interactive confirmation to prevent accidental removal
- Suggests using `-i` to install when LaTex is not installed

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
