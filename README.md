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

### 📋 Command Reference

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

### Check Your LaTex Environment

Before compiling, verify your LaTex installation:

```bash
lazitex --check    # Full command
lazitex -c         # Short command
```

This will detect:

- ✅ All installed LaTex compilers (pdflatex, xelatex, lualatex, etc.)
- ✅ Bibliography tools (bibtex, biber)
- ✅ Conversion utilities (dvipdfmx, ps2pdf, etc.)
- ✅ Package managers (tlmgr, mpm)
- ✅ Your LaTex distribution (TeX Live, MiKTeX, MacTeX)

**Sample Output:**

``` txt
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
        LaTex Environment Check Results
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🖥️  Operating System: Windows (windows)
📦 LaTex Distribution: TeX Live 2025

✅ Installed: 22/23 tools (core tools: 6/6)

🔨 Compilers (7/7)
────────────────────────────────────────────────────────────
  ⭐ ✓ PDFLaTex        [Installed]
      Path: C:\texlive\2025\bin\windows\pdflatex.exe
      Version: pdfTeX 3.141592653-2.6-1.40.28 (TeX Live 2025)
  ⭐ ✓ XeLaTex         [Installed]
      Path: C:\texlive\2025\bin\windows\xelatex.exe
      Version: XeTeX 3.141592653-2.6-0.999997 (TeX Live 2025)
  ...
```

---

## 🌍 Language Settings

LaziTex is a **bilingual tool** supporting both **English** and **Chinese (中文)**.

### Default Language

- **English** is the default language
- All output (help messages, environment check results, errors) will be in English by default

### Quick Start

**Use a specific language once:**

```bash
lazitex --lang en --check    # Use English for this command
lazitex -l zh --check        # Use Chinese (中文) for this command (short form)
```

**Switch default language permanently:**

```bash
# When you use --lang or -l, your choice is automatically saved
lazitex -l zh --help         # Switch to Chinese and save as default
lazitex --help               # Future commands now use Chinese

# Switch back anytime (use either form)
lazitex --lang en --help     # Switch to English and save as default
lazitex --check              # Now uses English
```

### Advanced Usage

**Temporary override with environment variable:**

```bash
# Windows PowerShell
$env:LAZITEX_LANG="zh"; lazitex --check

# Linux/macOS/Git Bash
export LAZITEX_LANG=zh
lazitex --check
```

**Supported language codes:**

- `en`, `english` → English
- `zh`, `chinese` → 中文 (Chinese)

**Command formats:**

- Long form: `--lang <code>`
- Short form: `-l <code>`

Both forms work identically and save your preference.

### Language Priority (highest to lowest)

1. **Command-line flag** `-l <lang>` or `--lang <lang>` (saves preference)
2. **Environment variable** `LAZITEX_LANG=<lang>`
3. **Saved user preference** (config file)
4. **System default** (English)

### Configuration File

Your language preference is stored in:

| Platform | Config Location |
|----------|----------------|
| **Windows** | `%APPDATA%\lazitex\config.json` |
| **Linux** | `~/.config/lazitex/config.json` |
| **macOS** | `~/Library/Application Support/lazitex/config.json` |

Example config file:

```json
{
  "language": "zh"
}
```

### Examples

```bash
# Check environment in English (default)
lazitex --check

# Check environment in Chinese once (short form)
lazitex -l zh --check

# Switch to Chinese permanently
lazitex --lang zh --help

# All future commands use Chinese
lazitex --check
lazitex --version

# Switch back to English (short form)
lazitex -l en --help

# Mix short and long forms freely
lazitex -l zh --check       # Short form
lazitex --lang en --check   # Long form

# Switch language in REPL mode
lazitex --repl
lazitex> lang zh            # Switch to Chinese interactively
lazitex> help               # Help is now in Chinese
```

**Output comparison:**

English:

```txt
✅ Installed: 22/23 tools (core tools: 6/6)
🔨 Compilers (7/7)
  ⭐ ✓ PDFLaTex        [Installed]
  ⭐ ✓ XeLaTex         [Installed]
```

中文 (Chinese):

```txt
✅ 已安装: 22/23 工具 (核心工具: 6/6)
🔨 编译器 (7/7)
  ⭐ ✓ PDFLaTex        [已安装]
  ⭐ ✓ XeLaTex         [已安装]
```

### Install/Update LaTex Environment

LaziTex supports automatic installation and updates for LaTex environments (currently macOS):

```bash
lazitex --install    # Install or update LaTex environment
lazitex -i           # Short form
```

**Features:**

- 🔍 **Smart Detection** - Automatically detects if LaTex is installed
- 📦 **Auto Install** - Installs BasicTeX automatically when not installed (macOS)
- 🔄 **Smart Update** - Checks and updates packages when already installed
- 💬 **Interactive Confirmation** - Asks for confirmation before install/update
- 📋 **Package List** - Shows list of updatable packages during update

**macOS Installation Example:**

```bash
$ lazitex -i
LaTex environment not detected
Will install BasicTeX via Homebrew (lightweight version, ~1.5 GB)
Install command: brew install --cask basictex
Do you want to install? [y/n]: y
Installing BasicTeX via Homebrew...
✅ BasicTeX installed successfully!
```

**Update Example:**

```bash
$ lazitex -i
LaTex is already installed
The following packages can be updated:
  ...
Do you want to update these packages? [y/n]: y
✅ Update successful!
```

### Uninstall LaTex Environment

LaziTex supports one-click uninstallation of LaTex environments (currently macOS):

```bash
lazitex --uninstall  # Uninstall LaTex environment
lazitex -u           # Short form
```

**Features:**

- 🔍 **Smart Detection** - Automatically detects if LaTex is installed
- 🗑️ **Auto Uninstall** - Automatically detects installation method and performs appropriate uninstallation
- 💬 **Interactive Confirmation** - Asks for confirmation before uninstalling to prevent accidental removal
- 📋 **Friendly Hints** - Suggests using `-i` to install when LaTex is not installed
- 🧹 **Residual Cleanup** - Attempts to remove common leftover paths after successful uninstall

**Supported Uninstallation Methods:**

- ✅ Homebrew BasicTeX - Automatic uninstallation
- ✅ Homebrew MacTeX - Automatic uninstallation
- ✅ MacPorts (texlive*, including basic/latex/full) - Automatic uninstallation
- ✅ MacTeX Official Installation - Attempts the official uninstall script; if missing, provides manual paths to remove

**macOS Uninstall Example:**

```bash
$ lazitex -u
LaTex is already installed
Will uninstall LaTex environment (including all installed packages)
Do you want to uninstall? [y/n]: y
Uninstalling LaTex environment...
✅ LaTex environment uninstalled successfully!
```

**When Not Installed:**

```bash
$ lazitex -u
LaTex environment not detected
Tip: You can use 'lazitex -i' to install LaTex environment with one click
```

### Build & Preview LaTex Document

LaziTex supports one-click compilation of LaTex documents into PDF with smart previewing:

```bash
lazitex -b main.tex             # Build LaTex document only
lazitex -b main.tex -s             # Build LaTex document and show
lazitex -b main.tex -o out/     # Specify output directory
lazitex -b main.tex -o res.pdf  # Specify custom output filename
```

**Features:**

- 🚀 **One-Click Build** - Automatically runs the compiler (XeLaTex) with optimal settings
- 🔍 **Auto Package Detection & Installation** - Automatically detects missing LaTex packages (`.sty` and `.cls` files) from compilation errors and offers to install them via `tlmgr` or `mpm`
  - **Smart Package Name Resolution** - Automatically finds the correct package name (e.g., `xeCJK.sty` → `xecjk` package)
  - **Batch Installation** - Installs all missing packages at once to minimize password prompts
  - **Intelligent Permission Handling** - Tries installation without sudo first, then prompts for password only when needed
  - **Auto Retry** - Automatically recompiles after successful package installation
- 👁️ **Smart Show (`-s`)** - Opens the PDF automatically upon successful build.
  - **macOS**: Prioritizes **Skim.app** (supporting silent refresh) if installed; otherwise, falls back to the system default browser or viewer.
  - **Windows**: Prioritizes **SumatraPDF** (supporting silent refresh and instance reuse) if installed; otherwise, falls back to the system default viewer.
- 📁 **Custom Output (`-o`)** - Supports specifying an output directory (automatically and recursively created if missing) or a complete output filename
- 📁 **Smart Default** - If `-o` is not specified, it automatically places PDF and log files in the same directory as the source `.tex` file
- 📍 **Path Support** - Supports both filenames in current directory and absolute/relative paths
- 🔍 **Type Safety** - Automatically validates file extensions and ensures source existence
- 🔄 **Consistent Logic** - Preview behavior is identical across CLI and REPL modes, laying the groundwork for future live-watch features

**Build Example:**

```bash
$ lazitex -b report.tex -s
🚀 Building LaTex document: report.tex
📁 Working directory: /Users/user/projects/paper
... (compiler output) ...
✨ Build successful!
(Automatically opening show window)
```

**Auto Package Installation Example:**

```bash
$ lazitex -b document.tex
🚀 Building LaTex document: document.tex
... (compilation fails) ...
🔍 Missing package detected: xeCJK
💡 Auto-install? [Y/n]: y
📦 Installing xeCJK...
Administrator privileges required, please enter password...
✅ Installation successful!
🔄 Retrying build...
✨ Build successful!
```

### Other Commands

```bash
lazitex --help     # Show help
lazitex --version  # Show version
lazitex --check    # Check LaTex environment
lazitex --repl     # REPL mode (interactive)
lazitex --tui      # Terminal UI mode (Coming Soon)
```

**REPL Mode** - Interactive commands:

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

---

## 🗂️ Project Structure

LaziTex is built with a clean, modular architecture — making it easy to extend and maintain.

```txt
lazitex/
├── 📦 go.mod                      # Go module definition
├── 📦 go.sum                      # Go dependencies
├── 🛠️  build.sh                   # Linux/macOS build script
├── 🛠️  build.bat                  # Windows build script
├── 📁 bin/                        # Build output directory
│
├── 🧠 core/                       # Core logic - the brain of LaziTex
│   ├── 🌍 env.go                  # Environment detection, installer interface & install/uninstall logic
│   ├── 🔨 build.go                # Build workflow (cross-platform)
│   ├── 📦 package.go              # Auto package detection & installation
│   └── 👀 watcher.go              # File watching for live preview
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
│       ├── 🎯 tasks/              # Task execution layer (unified command execution logic)
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

- **Layered Architecture**: Clean dependency hierarchy - `config`/`lang` (support) → `core` (business logic) → `tasks` (command execution) → `ui` (interaction)

- **Separation of Concerns**:
  - `core/` - Core business logic (environment detection, build, package management), completely independent of CLI, reusable by other applications
  - `tasks/` - Command execution layer, unified management of all command execution logic (environment, build, preview, help), shared by `main.go` and `ui/repl.go`
  - `ui/` - User interface layer, responsible for interaction mode implementation (REPL command parsing, completion, history, TUI state management)
  - `lang/` - Internationalization module, manages multi-language support and language preferences (contains `Language` type definition)
  - `config/` - Configuration management, handles user preference settings

- **Code Reusability**: The refactored architecture eliminates code duplication. Functions in `tasks/` package (like `CheckEnvironment()`, `BuildLaTex()`, `InstallLaTexEnvironment()`) are shared by both CLI mode and REPL mode

- **Modular Design**: Independent packages for i18n and configuration with clear responsibilities, easy to maintain and extend

- **Dependency Injection**: Platform-specific checkers and installers are created in `tasks/` and injected into core logic, avoiding circular dependencies

- **Interface-Driven**: `EnvironmentChecker` and `EnvironmentInstaller` interfaces allow easy platform extensibility. Adding a new platform only requires implementing the interfaces

- **Priority System**: Tools are categorized by priority (⭐ Core, 🔹 Important, 🔸 Optional), helping users understand tool importance

- **Type Safety**: Shared data structures (like `Language`) are defined in `lang/` package, simplifying architecture and avoiding unnecessary abstraction layers

---

## 🔍 Detected LaTex Tools

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

- [x] Environment detection for Windows, Linux, macOS
- [x] Comprehensive LaTex tool checking (23 tools)
- [x] macOS LaTex environment auto-install/update
- [x] macOS LaTex environment auto-uninstall
- [x] Modern REPL interactive mode (History, Tab completion, Shell shortcuts, Full i18n, Live Preview command)
- [x] LaTex Build Engine (Basic: One-click .tex to PDF, smart working directory)
- [x] Smart Preview System (macOS Skim/Web & Windows SumatraPDF/Default auto-adapter)
- [x] Live Preview Mode - Auto-trigger millisecond-level compilation and PDF refresh on save
- [x] Auto Package Detection & Installation - Automatically detects missing packages from compilation errors and installs them via tlmgr/mpm

### In Progress 🚧

- [ ] Live Preview experience optimization (task preemption, compilation locks, enhanced error feedback, etc.)
- [ ] Auto-Package Completion (Detect missing packages and install automatically)
- [ ] Multi-pass Compilation (Handle cross-references and bibliographies)
- [ ] Linux LaTex environment auto-install/update
- [ ] Linux LaTex environment auto-uninstall

### Planned 📋

- [ ] Local AI integration for LaTex assistance
- [ ] Multi-document project support
- [ ] Custom compilation profiles

### Future 🔮

- [ ] `lazitex init` Project Initialization - Minimalist startup experience similar to `uv init`, supporting fast fetching of high-quality LaTex templates from Gitee/GitHub
- [ ] TUI mode with live preview
- [ ] Modern GUI support
- [ ] Windows LaTex environment auto-install/update (Windows installation is complex, lower priority)
- [ ] Windows LaTex environment auto-uninstall
- [ ] **AI-Native Authoring Flow**: Bridging human creativity and AI intelligence with high-efficiency interaction and instantaneous feedback, realizing a "DocuGen-style Prompt2PDF" experience. Supporting end-to-end real-time PDF generation from natural language prompts, powered by Multi-Agent workflows for literature retrieval, outline planning, and content refinement
- [ ] **LaziTex Server & Cloud Vertical Ecosystem**: Exploring server-side deployment solutions to provide "Prompt2PDF" online interactive generation services for vertical sectors such as finance, medical, and research. Supporting multi-modal recognition (e.g., hand-written formulas/charts to LaTex), bridging unstructured intents to professional PDF documents, and building a distributed document platform that is "Free at Local, Intelligent in Cloud".

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

**Method 2: Manual installation**
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

### 5. What features will be supported in the future?

Please check the [Roadmap](#-roadmap) section above.

---

**Made with ❤️ by [SJRnhqh](https://github.com/SJRnhqh)**
