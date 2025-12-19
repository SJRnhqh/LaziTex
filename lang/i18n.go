// lang/i18n.go
// 中英文支持管理模块

package lang

import (
	// 外部包
	"os"
	"strings"

	// 内部包
	"github.com/SJRnhqh/lazitex/config"
)

// Language 语言代码
type Language string

const (
	LangZH Language = "zh" // 中文
	LangEN Language = "en" // English
)

// 全局当前语言
var currentLang Language = LangZH

// SaveLanguagePreference 保存语言偏好
func SaveLanguagePreference(lang Language) error {
	jsonConfig := &config.Config{
		Language: string(lang),
	}
	return config.SaveConfig(jsonConfig)
}

// LoadLanguagePreference 加载语言偏好
func LoadLanguagePreference() Language {
	jsonConfig, err := config.LoadConfig()
	if err != nil {
		return LangEN // 默认英文
	}

	if jsonConfig.Language == "zh" {
		return LangZH
	}
	return LangEN
}

// SetLanguage 设置当前语言
func SetLanguage(lang Language) {
	currentLang = lang
}

// GetLanguage 获取当前语言
func GetLanguage() Language {
	return currentLang
}

// GetCurrentLanguageName 获取当前语言的友好名称
func GetCurrentLanguageName() string {
	switch GetLanguage() {
	case LangZH:
		return "中文 (Chinese)"
	case LangEN:
		return "English"
	default:
		return "English"
	}
}

// DetectSystemLanguage 检测系统语言
func DetectSystemLanguage() Language {
	// 优先级：
	// 1. 环境变量 LAZITEX_LANG
	// 2. 用户配置文件
	// 3. 系统语言环境变量
	// 4. 默认英文

	// 1. 检查环境变量 LAZITEX_LANG
	if lang := os.Getenv("LAZITEX_LANG"); lang != "" {
		if strings.HasPrefix(strings.ToLower(lang), "zh") {
			return LangZH
		}
		return LangEN
	}

	// 2. 检查用户配置文件
	userLang := LoadLanguagePreference()
	if userLang != "" {
		return userLang
	}

	// 3. 检查系统环境变量
	for _, env := range []string{"LANG", "LANGUAGE", "LC_ALL", "LC_MESSAGES"} {
		if val := os.Getenv(env); val != "" {
			if strings.HasPrefix(strings.ToLower(val), "zh") {
				return LangZH
			}
			if strings.HasPrefix(strings.ToLower(val), "en") {
				return LangEN
			}
		}
	}

	// 4. 默认英文
	return LangEN
}

// I18n 国际化消息
type I18n struct {
	messages map[Language]map[string]string
}

// 全局 i18n 实例
var i18n = NewI18n()

// NewI18n 创建国际化实例
func NewI18n() *I18n {
	return &I18n{
		messages: make(map[Language]map[string]string),
	}
}

// T 翻译函数 (Translate)
func T(key string) string {
	return i18n.Get(key, currentLang)
}

// Get 获取翻译
func (i *I18n) Get(key string, lang Language) string {
	if langMap, ok := i.messages[lang]; ok {
		if msg, ok := langMap[key]; ok {
			return msg
		}
	}
	// 回退到中文
	if lang != LangZH {
		if langMap, ok := i.messages[LangZH]; ok {
			if msg, ok := langMap[key]; ok {
				return msg
			}
		}
	}
	// 如果都没有，返回 key
	return key
}

// Register 注册翻译
func (i *I18n) Register(lang Language, messages map[string]string) {
	i.messages[lang] = messages
}

// 初始化所有翻译
func init() {
	// 注册中文翻译
	i18n.Register(LangZH, map[string]string{
		// 标题和边框
		"title.env_check":  "LaTex 编译环境检测结果",
		"label.os":         "🖥️  操作系统",
		"label.distro":     "📦 LaTex 发行版",
		"label.installed":  "✅ 已安装",
		"label.tools":      "工具",
		"label.core_tools": "核心工具",
		"status.installed": "[已安装]",
		"status.missing":   "[未安装]",
		"label.path":       "路径",
		"label.version":    "版本",

		// 工具分类
		"category.compiler":        "🔨 编译器",
		"category.bibliography":    "📚 文献管理",
		"category.index":           "📇 索引工具",
		"category.converter":       "🔄 格式转换",
		"category.automation":      "⚙️ 自动化工具",
		"category.package_manager": "📦 包管理器",
		"category.utility":         "🔧 实用工具",

		// 提示消息
		"msg.unknown_command":    "未知命令: %s",
		"msg.no_compilers":       "⚠️  警告: 未检测到任何核心 LaTex 编译器",
		"msg.partial_install":    "💡 提示: 部分核心工具未安装，但基本功能可用",
		"msg.all_installed":      "🎉 太棒了！所有核心编译器都已安装",
		"msg.install_guide":      "📖 安装指南:",
		"msg.unknown_distro":     "未知",
		"msg.tip":                "💡 提示",
		"msg.warning":            "⚠️  警告",
		"msg.unsupported_os":     "错误: 不支持的操作系统 '%s'",
		"msg.build_usage":        "用法: lazitex -b <文件名.tex> [-o 输出路径] [-s]",
		"msg.checking_env":       "正在检查 LaTex 环境...",
		"msg.install_failed":     "安装失败",
		"msg.install_success":    "✅ LaTex 环境安装成功",
		"msg.uninstall_failed":   "卸载失败",
		"msg.uninstall_success":  "✅ LaTex 环境卸载成功",
		"msg.building_doc":       "🚀 正在构建 LaTex 文档: %s",
		"msg.working_dir":        "📁 工作目录: %s",
		"msg.build_success":      "✨ 构建成功！",
		"msg.build_failed":       "❌ 构建失败",
		"msg.err_abs_path":       "错误: 无法获取文件 '%s' 的绝对路径",
		"msg.err_mkdir":          "错误: 无法创建输出目录 '%s'",
		"msg.err_invalid_ext":    "错误: 无效的文件类型 '%s' (仅支持 .tex 文件)",
		"msg.err_file_not_found": "错误: 找不到文件 '%s'",
		"msg.preview_usage":      "用法: lazitex -p <文件名.tex>",
		"msg.preview_failed":     "⚠️  警告: 预览失败: %v",
		"msg.watching_file":      "👀 正在实时监听文件: %s (按 Ctrl+C 退出监听)",
		"msg.watcher_error":      "错误: 监听器意外崩溃: %v",

		// 包管理相关
		"msg.package_missing":             "🔍 检测到缺失的包: %s",
		"msg.package_install_prompt":      "💡 是否自动安装？[Y/n]: ",
		"msg.package_installing":          "📦 正在安装 %s...",
		"msg.package_install_success":     "✅ 安装成功！",
		"msg.package_install_failed":      "❌ 安装失败: %v",
		"msg.package_retry_build":         "🔄 重新编译中...",
		"msg.package_no_manager":          "⚠️  未找到包管理器 (tlmgr/mpm)，无法自动安装",
		"msg.package_need_sudo":           "需要管理员权限，请输入密码...",
		"msg.package_unsupported_manager": "不支持的包管理器: %s",

		// 多次编译相关
		"msg.compile_pass_n":         "第 %d 次编译（解决交叉引用）...",
		"msg.running_bibtex":         "运行 bibtex...",
		"msg.running_biber":          "运行 biber...",
		"msg.running_makeindex":      "运行 makeindex...",
		"msg.running_makeglossaries": "运行 makeglossaries...",
		"msg.tool_failed":            "运行 %s 失败: %v（继续编译）",

		// macOS 安装相关消息
		"msg.mac.brew_not_installed":       "错误: 未安装 Homebrew。请先安装 Homebrew: https://brew.sh",
		"msg.mac.latex_already_installed":  "检测到 LaTex 已安装",
		"msg.mac.checking_updates":         "正在检查更新...",
		"msg.mac.update_failed":            "更新失败",
		"msg.mac.update_success":           "✅ 更新成功！",
		"msg.mac.installing_basictex":      "正在通过 Homebrew 安装 BasicTeX...",
		"msg.mac.install_warning":          "这可能需要一些时间（约 1.5 GB 下载）...",
		"msg.mac.install_failed":           "安装失败",
		"msg.mac.updating_tlmgr":           "正在更新 TeX Live 包管理器...",
		"msg.mac.update_warning":           "警告: 更新过程中出现错误",
		"msg.mac.install_success":          "✅ BasicTeX 安装成功！",
		"msg.mac.next_steps":               "提示: 你可能需要重启终端或运行: eval $(/usr/libexec/path_helper)",
		"msg.mac.updating_tlmgr_self":      "正在更新 tlmgr 自身...",
		"msg.mac.updating_packages":        "正在更新包（这可能需要较长时间）...",
		"msg.mac.all_up_to_date":           "所有包都是最新的。",
		"msg.mac.updated_packages":         "已更新的包: %s",
		"msg.mac.invalid_input":            "无效输入，请输入 y (是) 或 n (否)",
		"msg.mac.cannot_list_packages":     "无法列出可更新的包列表",
		"msg.mac.updatable_packages":       "以下包可以更新:",
		"msg.mac.confirm_update":           "是否要更新这些包? [y/n]: ",
		"msg.mac.update_cancelled":         "更新已取消",
		"msg.mac.install_prompt":           "未检测到 LaTex 环境",
		"msg.mac.install_info":             "将使用 Homebrew 安装 BasicTeX (轻量版，约 1.5 GB)\n安装命令: brew install --cask basictex",
		"msg.mac.confirm_install":          "是否要安装? [y/n]: ",
		"msg.mac.install_cancelled":        "安装已取消",
		"msg.mac.latex_not_installed":      "未检测到 LaTex 环境",
		"msg.mac.uninstall_prompt":         "检测到 LaTex 已安装",
		"msg.mac.uninstall_info":           "将卸载 LaTex 环境（包括所有已安装的包）",
		"msg.mac.confirm_uninstall":        "是否要卸载? [y/n]: ",
		"msg.mac.uninstall_cancelled":      "卸载已取消",
		"msg.mac.uninstalling":             "正在卸载 LaTex 环境...",
		"msg.mac.uninstall_success":        "✅ LaTex 环境卸载成功！",
		"msg.mac.uninstall_failed":         "卸载失败",
		"msg.mac.uninstall_install_hint":   "提示: 可以使用 'lazitex -i' 一键安装 LaTex 环境",
		"msg.mac.uninstall_mactex_running": "检测到 MacTeX 官方安装，尝试运行卸载脚本...",
		"msg.mac.uninstall_mactex_manual":  "检测到 MacTeX 官方安装，但未找到卸载脚本，请手动删除：",
		"msg.mac.removing_residual":        "正在移除残留路径: %s",

		// 命令行帮助
		"help.usage":           "使用方法:",
		"help.commands":        "命令:",
		"help.show_help":       "显示此帮助",
		"help.show_version":    "显示版本号",
		"help.check_env":       "检查 LaTex 环境",
		"help.start_repl":      "启动 REPL 模式",
		"help.start_tui":       "启动终端 UI",
		"help.install_latex":   "安装或更新 LaTex 环境",
		"help.uninstall_latex": "卸载 LaTex 环境",
		"help.build_latex":     "一键构建 LaTex 文档为 PDF (使用 -o 指定输出, -s 编译后展示)",
		"help.live_preview":    "一键实时预览 PDF 文档 (使用 -p 指定文件名)",
		"help.set_language":    "设置语言 (zh/en)",
		"help.description":     "LaziTex: 零配置 LaTex 编译工具，支持 AI 辅助",

		// 工具描述
		"desc.pdflatex":  "最常用的 PDF 编译器",
		"desc.xelatex":   "支持 Unicode 和现代字体",
		"desc.lualatex":  "Lua 扩展的现代引擎",
		"desc.latex":     "传统 LaTex 编译器",
		"desc.pdftex":    "底层 TeX 引擎",
		"desc.tex":       "原始 TeX 引擎",
		"desc.etex":      "扩展 TeX 引擎",
		"desc.bibtex":    "传统参考文献管理",
		"desc.biber":     "现代参考文献管理",
		"desc.makeindex": "生成索引",
		"desc.xindy":     "多语言索引工具",
		"desc.texindy":   "LaTex 索引包装器",
		"desc.dvipdfmx":  "DVI 转 PDF（支持 CJK）",
		"desc.dvips":     "DVI 转 PostScript",
		"desc.ps2pdf":    "PostScript 转 PDF",
		"desc.dvisvgm":   "DVI 转 SVG",
		"desc.latexmk":   "自动化编译工具",
		"desc.tlmgr":     "TeX Live 包管理器",
		"desc.mpm":       "MiKTeX 包管理器",
		"desc.kpsewhich": "查找 TeX 文件路径",
		"desc.texdoc":    "查看文档",
		"desc.texhash":   "更新文件数据库",
		"desc.updmap":    "更新字体映射",

		// REPL 模式
		"repl.welcome":          "欢迎使用 LaziTex REPL 模式！",
		"repl.help_hint":        "输入 'help' 查看命令，输入 'quit' 或 'exit' 退出",
		"repl.prompt":           "lazitex> ",
		"repl.goodbye":          "再见！",
		"repl.unknown_command":  "未知命令: %s",
		"repl.type_help":        "输入 'help' 查看可用命令",
		"repl.help_title":       "可用命令:",
		"repl.help_desc":        "显示此帮助",
		"repl.version_desc":     "显示版本号",
		"repl.check_desc":       "检查 LaTex 环境",
		"repl.install_desc":     "安装或更新 LaTex 环境",
		"repl.uninstall_desc":   "卸载 LaTex 环境",
		"repl.build_desc":       "构建 LaTex 文档 (使用 -o 指定输出, -s 编译后展示)",
		"repl.build_usage":      "用法: build <文件名.tex> [-o 输出路径] [-s]",
		"repl.preview_desc":     "实时预览 PDF 文档 (使用 -p 指定.tex文件名)",
		"repl.preview_usage":    "用法: preview <文件名.tex>",
		"repl.lang_desc":        "切换语言 (zh/en)",
		"repl.quit_desc":        "退出 REPL",
		"repl.lang_usage":       "用法: lang <zh|en>",
		"repl.lang_current":     "当前语言: ",
		"repl.lang_switched":    "✓ 语言已切换并保存",
		"repl.lang_unsupported": "不支持的语言: %s",
		"repl.lang_available":   "可用语言: zh (中文), en (English)",
		"repl.shell_title":      "REPL 内置快捷命令",
		"repl.help_cd":          "切换当前目录，默认回到用户主目录",
		"repl.help_ls":          "列出当前/指定目录内容",
		"repl.help_pwd":         "显示当前工作目录",
		"repl.help_clear":       "清屏",
		"repl.help_cat":         "查看文件内容",
		"repl.err_init":         "初始化 REPL 失败: %v",
		"repl.err_cd":           "cd 切换目录失败: %v",
		"repl.err_pwd":          "获取当前目录失败: %v",
		"repl.err_ls":           "列出目录失败: %v",
		"repl.err_cat":          "读取文件失败: %v",
		"repl.cat_usage":        "用法: cat <文件名>",
	})

	// 注册英文翻译
	i18n.Register(LangEN, map[string]string{
		// Titles and borders
		"title.env_check":  "LaTex Environment Check Results",
		"label.os":         "🖥️  Operating System",
		"label.distro":     "📦 LaTex Distribution",
		"label.installed":  "✅ Installed",
		"label.tools":      "tools",
		"label.core_tools": "core tools",
		"status.installed": "[Installed]",
		"status.missing":   "[Not Installed]",
		"label.path":       "Path",
		"label.version":    "Version",

		// Tool categories
		"category.compiler":        "🔨 Compilers",
		"category.bibliography":    "📚 Bibliography",
		"category.index":           "📇 Indexing",
		"category.converter":       "🔄 Converters",
		"category.automation":      "⚙️ Automation",
		"category.package_manager": "📦 Package Managers",
		"category.utility":         "🔧 Utilities",

		// Messages
		"msg.unknown_command":    "Unknown command: %s",
		"msg.no_compilers":       "⚠️  Warning: No core LaTex compilers detected",
		"msg.partial_install":    "💡 Note: Some core tools are missing, but basic functionality is available",
		"msg.all_installed":      "🎉 Excellent! All core compilers are installed",
		"msg.install_guide":      "📖 Installation Guide:",
		"msg.unknown_distro":     "Unknown",
		"msg.tip":                "💡 Tip",
		"msg.warning":            "⚠️  Warning",
		"msg.unsupported_os":     "Error: Unsupported operating system '%s'",
		"msg.build_usage":        "Usage: lazitex -b <file.tex> [-o output_path] [-s]",
		"msg.checking_env":       "Checking LaTex environment...",
		"msg.install_failed":     "Installation failed",
		"msg.install_success":    "✅ LaTex environment installed successfully",
		"msg.uninstall_failed":   "Uninstallation failed",
		"msg.uninstall_success":  "✅ LaTex environment uninstalled successfully",
		"msg.building_doc":       "🚀 Building LaTex document: %s",
		"msg.working_dir":        "📁 Working directory: %s",
		"msg.build_success":      "✨ Build successful!",
		"msg.build_failed":       "❌ Build failed",
		"msg.err_abs_path":       "Error: Failed to get absolute path for '%s'",
		"msg.err_mkdir":          "Error: Failed to create output directory '%s'",
		"msg.err_invalid_ext":    "Error: Invalid file type '%s' (only .tex files supported)",
		"msg.err_file_not_found": "Error: File not found '%s'",
		"msg.preview_usage":      "Usage: lazitex -p <file.tex>",
		"msg.preview_failed":     "⚠️  Warning: Preview failed: %v",
		"msg.watching_file":      "👀 Watching file: %s (press Ctrl+C to stop)",
		"msg.watcher_error":      "Error: Watcher unexpectedly crashed: %v",

		// Package management related
		"msg.package_missing":             "🔍 Missing package detected: %s",
		"msg.package_install_prompt":      "💡 Auto-install? [Y/n]: ",
		"msg.package_installing":          "📦 Installing %s...",
		"msg.package_install_success":     "✅ Installation successful!",
		"msg.package_install_failed":      "❌ Installation failed: %v",
		"msg.package_retry_build":         "🔄 Retrying build...",
		"msg.package_no_manager":          "⚠️  No package manager (tlmgr/mpm) found, cannot auto-install",
		"msg.package_need_sudo":           "Administrator privileges required, please enter password...",
		"msg.package_unsupported_manager": "Unsupported package manager: %s",

		// Multiple compilation passes related
		"msg.compile_pass_n":         "Pass %d (resolving cross-references)...",
		"msg.running_bibtex":         "Running bibtex...",
		"msg.running_biber":          "Running biber...",
		"msg.running_makeindex":      "Running makeindex...",
		"msg.running_makeglossaries": "Running makeglossaries...",
		"msg.tool_failed":            "Running %s failed: %v (continuing compilation)",

		// macOS installation messages
		"msg.mac.brew_not_installed":       "Error: Homebrew is not installed. Please install Homebrew first: https://brew.sh",
		"msg.mac.latex_already_installed":  "LaTex is already installed",
		"msg.mac.checking_updates":         "Checking for updates...",
		"msg.mac.update_failed":            "Update failed",
		"msg.mac.update_success":           "✅ Update successful!",
		"msg.mac.installing_basictex":      "Installing BasicTeX via Homebrew...",
		"msg.mac.install_warning":          "This may take a while (approximately 1.5 GB download)...",
		"msg.mac.install_failed":           "Installation failed",
		"msg.mac.updating_tlmgr":           "Updating TeX Live package manager...",
		"msg.mac.update_warning":           "Warning: Error occurred during update",
		"msg.mac.install_success":          "✅ BasicTeX installed successfully!",
		"msg.mac.next_steps":               "Note: You may need to restart your terminal or run: eval $(/usr/libexec/path_helper)",
		"msg.mac.updating_tlmgr_self":      "Updating tlmgr itself...",
		"msg.mac.updating_packages":        "Updating packages (this may take a while)...",
		"msg.mac.all_up_to_date":           "All packages are up to date.",
		"msg.mac.updated_packages":         "Updated packages: %s",
		"msg.mac.invalid_input":            "Invalid input, please enter y (yes) or n (no)",
		"msg.mac.cannot_list_packages":     "Cannot list updatable packages",
		"msg.mac.updatable_packages":       "The following packages can be updated:",
		"msg.mac.confirm_update":           "Do you want to update these packages? [y/n]: ",
		"msg.mac.update_cancelled":         "Update cancelled",
		"msg.mac.install_prompt":           "LaTex environment not detected",
		"msg.mac.install_info":             "Will install BasicTeX via Homebrew (lightweight version, ~1.5 GB)\nInstall command: brew install --cask basictex",
		"msg.mac.confirm_install":          "Do you want to install? [y/n]: ",
		"msg.mac.install_cancelled":        "Installation cancelled",
		"msg.mac.latex_not_installed":      "LaTex environment not detected",
		"msg.mac.uninstall_prompt":         "LaTex is already installed",
		"msg.mac.uninstall_info":           "Will uninstall LaTex environment (including all installed packages)",
		"msg.mac.confirm_uninstall":        "Do you want to uninstall? [y/n]: ",
		"msg.mac.uninstall_cancelled":      "Uninstallation cancelled",
		"msg.mac.uninstalling":             "Uninstalling LaTex environment...",
		"msg.mac.uninstall_success":        "✅ LaTex environment uninstalled successfully!",
		"msg.mac.uninstall_failed":         "Uninstallation failed",
		"msg.mac.uninstall_install_hint":   "Tip: You can use 'lazitex -i' to install LaTex environment with one click",
		"msg.mac.uninstall_mactex_running": "Detected official MacTeX, attempting uninstall script...",
		"msg.mac.uninstall_mactex_manual":  "Detected official MacTeX but uninstall script not found. Please remove manually:",
		"msg.mac.removing_residual":        "Removing residual path: %s",

		// CLI help
		"help.usage":           "Usage:",
		"help.commands":        "Commands:",
		"help.show_help":       "Show this help",
		"help.show_version":    "Show version",
		"help.check_env":       "Check LaTex environment",
		"help.start_repl":      "Start REPL mode",
		"help.start_tui":       "Start terminal UI",
		"help.install_latex":   "Install or update LaTex environment",
		"help.uninstall_latex": "Uninstall LaTex environment",
		"help.build_latex":     "Build LaTex document to PDF (use -o for output, -s to show after build)",
		"help.live_preview":    "Live preview PDF document (use -p for file name)",
		"help.set_language":    "Set language (zh/en)",
		"help.description":     "LaziTex: Zero-config LaTex compilation with AI assistance",

		// Tool descriptions
		"desc.pdflatex":  "Most common PDF compiler",
		"desc.xelatex":   "Unicode and modern font support",
		"desc.lualatex":  "Modern engine with Lua extensions",
		"desc.latex":     "Traditional LaTex compiler",
		"desc.pdftex":    "Low-level TeX engine",
		"desc.tex":       "Original TeX engine",
		"desc.etex":      "Extended TeX engine",
		"desc.bibtex":    "Traditional bibliography management",
		"desc.biber":     "Modern bibliography management",
		"desc.makeindex": "Generate index",
		"desc.xindy":     "Multilingual indexing tool",
		"desc.texindy":   "LaTex index wrapper",
		"desc.dvipdfmx":  "DVI to PDF (CJK support)",
		"desc.dvips":     "DVI to PostScript",
		"desc.ps2pdf":    "PostScript to PDF",
		"desc.dvisvgm":   "DVI to SVG",
		"desc.latexmk":   "Automated compilation tool",
		"desc.tlmgr":     "TeX Live package manager",
		"desc.mpm":       "MiKTeX package manager",
		"desc.kpsewhich": "Find TeX file paths",
		"desc.texdoc":    "View documentation",
		"desc.texhash":   "Update file database",
		"desc.updmap":    "Update font mapping",

		// REPL mode
		"repl.welcome":          "Welcome to LaziTex REPL mode!",
		"repl.help_hint":        "Type 'help' for commands, 'quit' or 'exit' to exit",
		"repl.prompt":           "lazitex> ",
		"repl.goodbye":          "Goodbye!",
		"repl.unknown_command":  "Unknown command: %s",
		"repl.type_help":        "Type 'help' for available commands",
		"repl.help_title":       "Available commands:",
		"repl.help_desc":        "Show this help",
		"repl.version_desc":     "Show version",
		"repl.check_desc":       "Check LaTex environment",
		"repl.install_desc":     "Install or update LaTex environment",
		"repl.uninstall_desc":   "Uninstall LaTex environment",
		"repl.build_desc":       "Build LaTex document (use -o for output, -s to show after build)",
		"repl.build_usage":      "Usage: build <filename.tex> [-o output_path] [-s]",
		"repl.preview_desc":     "Live preview PDF document (use -p for .tex file name)",
		"repl.preview_usage":    "Usage: preview <filename.tex>",
		"repl.lang_desc":        "Switch language (zh/en)",
		"repl.quit_desc":        "Exit REPL",
		"repl.lang_usage":       "Usage: lang <zh|en>",
		"repl.lang_current":     "Current language: ",
		"repl.lang_switched":    "✓ Language switched and saved",
		"repl.lang_unsupported": "Unsupported language: %s",
		"repl.lang_available":   "Available: zh (中文), en (English)",
		"repl.shell_title":      "Shell-like shortcuts",
		"repl.help_cd":          "Change directory (default to home)",
		"repl.help_ls":          "List current/specified directory",
		"repl.help_pwd":         "Print working directory",
		"repl.help_clear":       "Clear screen",
		"repl.help_cat":         "Display file content",
		"repl.err_init":         "Error initializing REPL: %v",
		"repl.err_cd":           "cd failed: %v",
		"repl.err_pwd":          "pwd failed: %v",
		"repl.err_ls":           "ls failed: %v",
		"repl.err_cat":          "cat failed: %v",
		"repl.cat_usage":        "Usage: cat <filename>",
	})
}
