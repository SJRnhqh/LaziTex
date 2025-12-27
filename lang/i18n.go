// lang/i18n.go
// 中英文支持管理模块

package lang

import (
	// 外部包
	"os"
	"strings"

	// 内部包
	cfg "github.com/SJRnhqh/lazitex/config"
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
	// 使用 UpdateConfig 进行部分更新，只更新 language 字段
	return cfg.UpdateConfig(map[string]interface{}{
		"language": string(lang),
	})
}

// LoadLanguagePreference 加载语言偏好
func LoadLanguagePreference() Language {
	jsonConfig := cfg.LoadConfig()

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
		"title.env_check":  "LaTeX 编译环境检测结果",
		"label.os":         "🖥️  操作系统",
		"label.distro":     "📦 LaTeX 发行版",
		"label.installed":  "✨ 已安装",
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
		"msg.no_compilers":       "⚠️  警告: 未检测到任何核心 LaTeX 编译器",
		"msg.partial_install":    "💡 提示: 部分核心工具未安装，但基本功能可用",
		"msg.all_installed":      "🎉 太棒了！所有核心编译器都已安装",
		"msg.install_guide":      "📖 安装指南:",
		"msg.unknown_distro":     "未知",
		"msg.tip":                "💡 提示",
		"msg.warning":            "⚠️  警告",
		"msg.unsupported_os":     "错误: 不支持的操作系统 '%s'",
		"msg.build_usage":        "用法: lazitex -b <文件名.tex> [-o|--output 输出路径] [-s|--show] [-q|--quiet] [-t|--tidy]",
		"msg.checking_env":       "正在检查 LaTeX 环境...",
		"msg.install_failed":     "安装失败",
		"msg.install_success":    "✨ LaTeX 环境安装成功",
		"msg.uninstall_failed":   "卸载失败",
		"msg.uninstall_success":  "✨ LaTeX 环境卸载成功",
		"msg.building_doc":       "🚀 正在构建 LaTeX 文档: %s",
		"msg.working_dir":        "📁 工作目录: %s",
		"msg.output_dir":         "📁 输出目录: %s",
		"msg.build_success":      "✨ 构建成功！",
		"msg.build_failed":       "😢 构建失败",
		"msg.err_abs_path":       "错误: 无法获取文件 '%s' 的绝对路径",
		"msg.err_mkdir":          "错误: 无法创建输出目录 '%s'",
		"msg.err_invalid_ext":    "错误: 无效的文件类型 '%s' (仅支持 .tex 文件)",
		"msg.err_file_not_found": "错误: 找不到文件 '%s'",
		"msg.preview_usage":      "用法: lazitex -p <文件名.tex> [-q|--quiet] [-t|--tidy] [-d|--dev] [:端口号]",
		"msg.show_failed":        "⚠️  警告: 展示失败: %v",
		"msg.watching_file":      "👀 正在实时监听文件: %s (按 Ctrl+C 退出监听)",
		"msg.watcher_error":      "错误: 监听器意外崩溃: %v",

		// Web预览相关
		"msg.server_starting":   "🌐 服务器启动在端口 %d",
		"msg.server_error":      "😢 服务器错误: %v",
		"msg.opening_browser":   "🔗 正在打开浏览器: %s",
		"msg.browser_error":     "⚠️  无法打开浏览器: %v",
		"msg.manual_open":       "💡 请手动打开浏览器访问: %s",
		"msg.config_path_error": "获取配置路径失败",
		"msg.config_save_error": "创建默认配置失败",
		"msg.config_open_error": "打开配置文件失败",
		"msg.config_opening":    "正在打开配置文件",
		"msg.pdf_not_found":     "PDF文件不存在",
		"msg.invalid_port":      "😢错误: 无效的端口号: %s (必须是 1-65535 之间的数字)",
		"msg.dev_mode_vue_url":  "💡 开发模式：Vue 前端地址: %s",

		// LLM 管理相关消息
		// 基础与列表
		"msg.llm.load_config_failed":      "😢 读取配置失败: %v",
		"msg.llm.no_providers":            "📝 暂无 LLM 配置，请先添加一个 LLM 提供商",
		"msg.llm.list_header_id":          "ID",
		"msg.llm.list_header_name":        "名称",
		"msg.llm.list_header_provider":    "提供商",
		"msg.llm.list_header_model":       "模型",
		"msg.llm.list_header_enabled":     "启用",
		"msg.llm.list_header_verified":    "已验证",
		"msg.llm.list_header_verified_at": "验证时间",
		"msg.llm.list_header_active":      "当前",
		"msg.llm.date_format":             "2006-01-02", // 日期格式（Go time 格式字符串）
		// 测试
		"msg.llm.test.test_failed":                     "😢 测试失败: %v",
		"msg.llm.test.not_found":                       "😢 未找到 LLM '%s'",
		"msg.llm.test.available_verified":              "✨ 可用: %s (%s:%s) - 已验证",
		"msg.llm.test.unavailable":                     "😢 不可用: %s (%s:%s)",
		"msg.llm.test.update_config_failed_on_success": "⚠️  测试通过，但更新配置失败: %v",
		"msg.llm.test.update_config_failed_on_failure": "⚠️  测试失败，但更新配置失败: %v",
		"msg.llm.test.multiple_matches":                "找到 %d 个匹配的 LLM，请选择要测试的：",
		// 连接
		"msg.llm.link.not_found":              "😢 LLM Provider '%s' 不存在",
		"msg.llm.link.test_failed":            "😢 连通性测试失败: %v",
		"msg.llm.link.unavailable":            "⚠️  LLM 无法接入: %s (%s)",
		"msg.llm.link.set_active_failed":      "😢 连接失败: %v",
		"msg.llm.link.set_current_failed":     "😢 设置前台失败: %v",
		"msg.llm.link.success":                "✨ %s 已连接",
		"msg.llm.link.activated":              "✨ %s 已激活",
		"msg.llm.link.already_active":         "💡 %s 已经激活",
		"msg.llm.link.already_current":        "💡 %s 已经是当前前台 LLM",
		"msg.llm.link.multiple_matches":       "找到 %d 个匹配的 LLM，请选择要连接的：",
		"msg.llm.link.batch_multiple_skip":    "⚠ %s 匹配到多个 Provider，跳过批量模式下的歧义项",
		"msg.llm.link.batch_already_active":   "💡 %s 已经激活，跳过",
		"msg.llm.link.batch_failed":           "😢 %s 连接失败: %v",
		"msg.llm.link.model_multiple_title":   "模型 '%s' 匹配到多个 LLM",
		"msg.llm.link.registered_title":       "已注册的 LLM Provider",
		"msg.llm.link.no_providers":           "当前没有已注册的 LLM Provider",
		"msg.llm.link.no_available":           "没有可用的 Provider",
		"msg.llm.link.select_prompt":          "请选择：",
		"msg.llm.link.select_number":          "  - 输入数字 (1-%d) 选择特定 LLM 并激活",
		"msg.llm.link.select_all":             "  - 输入 'a' 或 'all' 全部激活",
		"msg.llm.link.select_cancel":          "  - 输入 'q' 或 'quit' 取消",
		"msg.llm.link.select_input":           "选择: ",
		"msg.llm.link.read_input_failed":      "😢 读取输入失败: %v",
		"msg.llm.link.invalid_choice":         "😢 无效的选择: %s",
		"msg.llm.link.activate_failed":        "😢 激活失败: %v",
		"msg.llm.link.activate_success":       "✨ %s 已激活",
		"msg.llm.link.set_current_success":    "✨ 已设为前台",
		"msg.llm.link.activate_all_count":     "共激活 %d 个 Provider",
		"msg.llm.link.status_current":         " [当前前台]",
		"msg.llm.link.status_active":          " [已激活]",
		"msg.llm.link.status_inactive":        " [未激活]",
		"msg.llm.link.provider_format":        "  [%d] %s (Provider: %s, Model: %s)",
		"msg.llm.link.provider_format_simple": "  [%d] %s (ID: %s, Provider: %s)",
		// 取消连接
		"msg.llm.unlink.not_found":            "😢 LLM Provider '%s' 不存在",
		"msg.llm.unlink.not_active":           "⚠️  LLM Provider '%s' 不在激活列表中",
		"msg.llm.unlink.failed":               "😢 取消连接失败: %v",
		"msg.llm.unlink.success":              "✨ %s 已取消连接",
		"msg.llm.unlink.multiple_matches":     "找到 %d 个匹配的 LLM，请选择要取消连接的：",
		"msg.llm.unlink.clear_current_failed": "⚠️  清空当前激活状态失败: %v",
		"msg.llm.unlink.batch_multiple_skip":  "⚠ %s 匹配到多个 Provider，跳过批量模式下的歧义项",
		"msg.llm.unlink.batch_failed":         "😢 %s 取消连接失败: %v",
		"msg.llm.unlink.model_multiple_title": "模型 '%s' 匹配到多个 LLM",
		"msg.llm.unlink.registered_title":     "已注册的 LLM Provider",
		"msg.llm.unlink.active_title":         "已激活的 LLM Provider",
		"msg.llm.unlink.no_providers":         "当前没有已注册的 LLM Provider",
		"msg.llm.unlink.no_active_providers":  "当前没有已激活的 LLM Provider",
		"msg.llm.unlink.select_prompt":        "请选择要取消连接的 Provider：",
		"msg.llm.unlink.select_number":        "  - 输入数字 (1-%d) 选择特定 LLM 并取消连接",
		"msg.llm.unlink.select_all":           "  - 输入 'a' 或 'all' 全部取消连接",
		"msg.llm.unlink.select_cancel":        "  - 输入 'q' 或 'quit' 取消",
		"msg.llm.unlink.select_input":         "选择: ",
		"msg.llm.unlink.unlink_failed":        "😢 取消连接失败: %v",
		"msg.llm.unlink.unlink_success":       "✨ %s 已取消连接",
		"msg.llm.unlink.unlink_all_count":     "共取消连接 %d 个 Provider",
		// 切换
		"msg.llm.switch.not_found":            "😢 未找到 LLM: %s",
		"msg.llm.switch.not_active":           "⚠️  LLM '%s' 未激活，请先使用 link 命令",
		"msg.llm.switch.failed":               "😢 切换失败: %v",
		"msg.llm.switch.success":              "✨ 已切换到 LLM: %s",
		"msg.llm.switch.already_current":      "💡 %s 已经是当前前台 LLM",
		"msg.llm.switch.no_active":            "😢 没有已激活的 LLM Provider",
		"msg.llm.switch.active_title":         "已激活的 LLM Provider",
		"msg.llm.switch.model_multiple_title": "模型 '%s' 匹配到多个已激活的 LLM Provider",
		"msg.llm.switch.select_prompt":        "请选择要切换到前台的 LLM",
		"msg.llm.switch.select_number":        "输入数字 (1-%d) 选择",
		"msg.llm.switch.select_cancel":        "输入 'q' 或 'quit' 取消",
		"msg.llm.switch.select_input":         "请输入: ",
		"msg.llm.switch.batch_not_supported":  "😢 不支持批量切换，一次只能切换到一个 LLM",
		// 删除
		"msg.llm.remove.success":          "🗑️ 已删除 LLM: %s",
		"msg.llm.remove.save_failed":      "⚠️ 保存配置失败: %s",
		"msg.llm.remove.not_found":        "😢 未找到 LLM: %s",
		"msg.llm.remove.cancelled":        "已取消删除",
		"msg.llm.remove.multiple_matches": "找到 %d 个匹配的 LLM Provider，请选择：",
		"msg.llm.remove.confirm_prompt":   "确定要删除该 LLM Provider 吗？[y/N]: ",
		"msg.llm.remove.select_prompt":    "请选择要删除的 LLM Provider (输入序号，或输入 q 取消): ",
		// 通用选择提示
		"msg.llm.select_prompt": "请选择 (输入序号，或输入 q 取消): ",
		"msg.llm.cancelled":     "已取消",
		// 注册相关
		"msg.llm.register.title":                  "=== 注册新的 LLM Provider ===",
		"msg.llm.register.confirm_new":            "是否接入新的 LLM？",
		"msg.llm.register.input_provider":         "请输入 LLM 供应商类型 (例如: ollama, openai)",
		"msg.llm.register.input_baseurl":          "请输入基础URL (例如: http://localhost:11434)",
		"msg.llm.register.unsupported_provider":   "不支持的供应商类型，当前仅支持 ollama",
		"msg.llm.register.confirm_link_existing":  "是否要连接现有的已注册 LLM？",
		"msg.llm.register.input_name":             "请给这个 LLM 配置一个名称 (例如: my-ollama)",
		"msg.llm.register.input_model":            "请输入模型名称 (例如: qwen2.5-coder:7b)",
		"msg.llm.register.test_failed":            "连接测试失败: %v",
		"msg.llm.register.test_success":           "✨ %s",
		"msg.llm.register.save_failed":            "保存失败: %v",
		"msg.llm.register.set_active_failed":      "设置激活失败: %v",
		"msg.llm.register.set_current_failed":     "设置当前 LLM 失败: %v",
		"msg.llm.register.success":                "✨ 注册成功！",
		"msg.llm.register.name_empty":             "名称不能为空，请重新输入",
		"msg.llm.register.name_exists":            "名称 '%s' 已存在，请换一个名称",
		"msg.llm.register.input_empty":            "输入不能为空，请重新输入",
		"msg.llm.register.confirm_invalid":        "请输入 Y 或 n",
		"msg.llm.register.ollama.connect_failed":  "无法连接到 Ollama 服务: %w",
		"msg.llm.register.ollama.status_error":    "Ollama 服务返回错误状态: %d",
		"msg.llm.register.ollama.parse_failed":    "解析响应失败: %w",
		"msg.llm.register.ollama.model_not_found": "模型 '%s' 不存在，请先运行: ollama pull %s",
		"msg.llm.register.ollama.test_success":    "连接成功，模型 '%s' 可用",
		// 智能模式相关
		// Add 命令相关（补充）
		"msg.llm.add.test_success":    "连接测试验证成功，正在注册 LLM: %s (%s)",
		"msg.llm.add.test_failed":     "😢 连接测试失败，取消注册: %v",
		"msg.llm.add.register_failed": "😢 注册失败: %v",
		"msg.llm.add.registered":      "✨ 已注册 LLM: %s (%s)",
		// Provider 基础错误
		"msg.llm.provider.id_exists":   "LLM Provider ID '%s' 已存在",
		"msg.llm.provider.name_exists": "LLM Provider 名称 '%s' 已存在",
		"msg.llm.provider.unsupported": "不支持的 provider: %s",
		// Provider: Ollama
		"msg.llm.provider.ollama.service_not_running":      "Ollama 服务未运行",
		"msg.llm.provider.ollama.model_not_found":          "模型不存在",
		"msg.llm.provider.ollama.api_connect_failed":       "无法连接到 Ollama API: %v",
		"msg.llm.provider.ollama.api_status_error":         "Ollama API 返回错误状态码: %d",
		"msg.llm.provider.ollama.parse_response_failed":    "解析 Ollama 响应失败: %v",
		"msg.llm.provider.ollama.serialize_request_failed": "序列化请求失败: %v",
		// CLI 用法
		"msg.llm.cli_usage": "用法: lazitex -m list | lazitex -m add <model> -p <provider> -n <name> | lazitex -m [remove|link|unlink|test|switch] <id/name/model>",
		// Add 命令相关
		"msg.llm.add.model_required":             "😢 需要提供注册模型名称",
		"msg.llm.add.provider_required":          "😢 需要提供提供商名称",
		"msg.llm.add.name_required":              "😢 需要为注册模型自定义名称",
		"msg.llm.add.provider_missing":           "😢 需要提供提供商",
		"msg.llm.add.unexpected_nonflag":         "😢 意外的非flag参数: %s",
		"msg.llm.add.smart_mode_not_implemented": "⚠️  智能注册/连接功能暂未实现",
		"msg.llm.add.smart_mode_hint":            "请使用: lazitex -m add <model> 注册，或 lazitex -m link <model> 连接",
		"msg.llm.unknown_action":                 "未知操作: %s",
		"msg.dev_mode_backend_url":               "💡 后端 API 地址: %s",
		"msg.prod_mode_embedded":                 "💡 生产模式：前端已嵌入，访问: %s",
		"msg.dev_mode_custom_port_warning":       "⚠️  警告：开发模式下使用了自定义端口 (%d)，但 Vite 代理配置硬编码到 8080\n   将强制使用默认端口 8080 以确保前端可以连接后端",
		"msg.port_8080_occupied":                 "⚠️  端口 8080 已被占用，无法使用开发模式\n   将回退到生产模式，使用自定义端口 %d",
		"msg.checking_vue_dev_server":            "🔍 检测 Vue dev server (%s) 是否可用...",
		"msg.vue_dev_server_available":           "✨ Vue dev server 可用，使用开发模式",
		"msg.vue_dev_server_unavailable":         "⚠️  Vue dev server 不可用，自动回退到生产模式（使用嵌入的前端）",
		"msg.common.yes":                         "是",
		"msg.common.no":                          "否",

		// 包管理相关
		"msg.package_missing":             "🔍 检测到缺失的包: %s",
		"msg.package_install_prompt":      "💡 是否自动安装？[Y/n]: ",
		"msg.package_installing":          "📦 正在安装 %s...",
		"msg.package_install_success":     "✨ 安装成功！",
		"msg.package_install_failed":      "😢 安装失败: %v",
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
		"msg.tidy_success":           "🧹 已清理辅助文件",
		"msg.tidy_warning":           "⚠️  警告：清理辅助文件失败",
		"msg.tidy_remove_failed":     "删除文件 %s 失败: %v",

		// 编译锁和任务管理相关
		"msg.compile_cancelled":        "编译已取消",
		"msg.cancelling_previous_task": "⚠️  取消之前的编译任务...",
		"msg.compile_timeout":          "⏰ 编译超时（超过 %d 秒），已自动取消",

		// macOS 安装相关消息
		"msg.mac.brew_not_installed":       "错误: 未安装 Homebrew。请先安装 Homebrew: https://brew.sh",
		"msg.mac.latex_already_installed":  "检测到 LaTeX 已安装",
		"msg.mac.checking_updates":         "正在检查更新...",
		"msg.mac.update_failed":            "更新失败",
		"msg.mac.update_success":           "✨ 更新成功！",
		"msg.mac.installing_basictex":      "正在通过 Homebrew 安装 BasicTeX...",
		"msg.mac.install_warning":          "这可能需要一些时间（约 1.5 GB 下载）...",
		"msg.mac.install_failed":           "安装失败",
		"msg.mac.updating_tlmgr":           "正在更新 TeX Live 包管理器...",
		"msg.mac.update_warning":           "警告: 更新过程中出现错误",
		"msg.mac.install_success":          "✨ BasicTeX 安装成功！",
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
		"msg.mac.install_prompt":           "未检测到 LaTeX 环境",
		"msg.mac.install_info":             "将使用 Homebrew 安装 BasicTeX (轻量版，约 1.5 GB)\n安装命令: brew install --cask basictex",
		"msg.mac.confirm_install":          "是否要安装? [y/n]: ",
		"msg.mac.install_cancelled":        "安装已取消",
		"msg.mac.latex_not_installed":      "未检测到 LaTeX 环境",
		"msg.mac.uninstall_prompt":         "检测到 LaTeX 已安装",
		"msg.mac.uninstall_info":           "将卸载 LaTeX 环境（包括所有已安装的包）",
		"msg.mac.confirm_uninstall":        "是否要卸载? [y/n]: ",
		"msg.mac.uninstall_cancelled":      "卸载已取消",
		"msg.mac.uninstalling":             "正在卸载 LaTeX 环境...",
		"msg.mac.uninstall_success":        "✨ LaTeX 环境卸载成功！",
		"msg.mac.uninstall_failed":         "卸载失败",
		"msg.mac.uninstall_install_hint":   "提示: 可以使用 'lazitex -i' 一键安装 LaTeX 环境",
		"msg.mac.uninstall_mactex_running": "检测到 MacTeX 官方安装，尝试运行卸载脚本...",
		"msg.mac.uninstall_mactex_manual":  "检测到 MacTeX 官方安装，但未找到卸载脚本，请手动删除：",
		"msg.mac.removing_residual":        "正在移除残留路径: %s",

		// LaTeX 相关消息
		"msg.latex_usage":                 "用法: lazitex -x <-c|--check|-i|--install|-u|--uninstall>",
		"latex.err_no_action":             "😢 错误: 未指定操作",
		"latex.err_multi_actions":         "😢 错误: 指定了多个操作，只能选择一个",
		"latex.err_unknown_flag":          "😢 错误: 未知的标志",
		"latex.err_extra_argument":        "😢 错误: 意外的参数 '%s'，此命令不接受额外参数",
		"msg.latex.not_implemented":       "😢 功能尚未实现（开发中）",
		"msg.latex.linux.not_implemented": "😢 Linux 平台的安装/卸载功能尚未实现（开发中）",
		"msg.latex.win.not_implemented":   "😢 Windows 平台的安装/卸载功能尚未实现（开发中）",

		// Ollama 相关消息
		"msg.ollama.already_installed":              "🦙 Ollama 已安装，版本: %s",
		"msg.ollama.installing":                     "🦙 正在通过 winget 安装 Ollama...",
		"msg.ollama.install_note_gui":               "💡 注意: Windows 版本包含 GUI 界面，但 CLI 工具同样可用",
		"msg.ollama.install_success":                "✨ Ollama 安装成功！版本: %s",
		"msg.ollama.install_failed":                 "😢 winget 安装失败",
		"msg.ollama.install_verify_failed":          "⚠️  安装完成但验证失败，请手动检查",
		"msg.ollama.confirm_install":                "🦙 确定要安装 Ollama 吗？[y/N]: ",
		"msg.ollama.install_cancelled":              "已取消安装",
		"msg.ollama.winget_not_available":           "⚠️  winget 不可用，请先安装 winget 或手动安装 Ollama",
		"msg.ollama.uninstalling":                   "🦙 正在通过 winget 卸载 Ollama...",
		"msg.ollama.uninstall_success":              "✨ Ollama 卸载成功！",
		"msg.ollama.uninstall_failed":               "😢 winget 卸载失败",
		"msg.ollama.uninstall_verify_failed":        "⚠️  卸载完成但验证失败，可能仍有残留文件",
		"msg.ollama.not_installed":                  "🦙 Ollama 未安装",
		"msg.ollama.confirm_uninstall":              "🦙 确定要卸载 Ollama 吗？[y/N]: ",
		"msg.ollama.uninstall_cancelled":            "已取消卸载",
		"msg.ollama.winget_uninstall_not_available": "⚠️  winget 不可用，请手动卸载 Ollama",
		"msg.ollama.service_running":                "✨ Ollama 服务正在运行",
		"msg.ollama.service_not_running":            "⚠️  Ollama 服务未运行",
		"msg.ollama.start_service_hint":             "💡 提示: 请运行 'ollama serve' 启动服务，或重启 Ollama 应用",
		"msg.ollama.path_not_updated":               "⚠️  注意: Ollama 已安装，但命令行工具可能尚未添加到 PATH",
		"msg.ollama.restart_terminal_hint":          "💡 提示: 请重启终端或重新打开命令行窗口，然后再次运行 'ollama -c' 检查",
		"msg.ollama.install_maybe_success":          "⚠️  Ollama 可能已安装，但命令行工具尚未在 PATH 中",
		"msg.ollama.winget_install_prompt":          "💡 winget 是 Windows 的包管理器，用于自动安装 Ollama",
		"msg.ollama.winget_install_confirm":         "是否要自动安装 winget? [Y/n]: ",
		"msg.ollama.winget_install_cancelled":       "已取消安装 winget",
		"msg.ollama.winget_required":                "需要 winget 才能继续安装 Ollama",
		"msg.ollama.installing_winget":              "📦 正在安装 winget...",
		"msg.ollama.winget_store_opened":            "🔗 已打开 Microsoft Store",
		"msg.ollama.winget_store_instruction":       "请在 Microsoft Store 中点击「获取」或「安装」按钮",
		"msg.ollama.winget_wait_prompt":             "安装完成后，请按回车键继续...",
		"msg.ollama.winget_install_success":         "✨ winget 安装成功！",
		"msg.ollama.winget_install_error":           "winget 安装失败",
		"msg.ollama.downloading_winget":             "📥 正在下载 winget 安装包...",
		"msg.ollama.winget_download_failed":         "下载 winget 安装包失败",
		"msg.ollama.installing_winget_package":      "📦 正在安装 winget 安装包...",
		"msg.ollama.winget_install_package_failed":  "安装 winget 安装包失败",
		"msg.ollama.winget_install_verify_failed":   "winget 安装完成但验证失败，请手动检查",
		"msg.ollama.invalid_input":                  "无效输入，请输入 y (是) 或 n (否)",
		"msg.ollama.manager_nil":                    "😢 Ollama 管理器未初始化",
		"msg.ollama.check_failed":                   "检查 Ollama 状态失败",
		"msg.ollama.check_installed":                "Ollama 已安装，版本",
		"msg.ollama.check_installed_no_version":     "Ollama 已安装（无法获取版本）",
		"msg.ollama.status_installed":               "🦙 Ollama 已安装，版本: %s",
		"msg.ollama.status_installed_no_version":    "🦙 Ollama 已安装",
		"msg.ollama.mac.install_method_prompt":      "🦙 检测到已安装 Homebrew，可以选择以下安装方式：",
		"msg.ollama.mac.install_method_option1":     "  1. Homebrew (推荐，快速安装)",
		"msg.ollama.mac.install_method_option2":     "  2. 官方安装脚本 (直接安装)",
		"msg.ollama.mac.install_method_choice":      "请选择 [1/2] (默认: 1): ",
		"msg.ollama_usage":                          "用法: lazitex -o <-c|--check|-i|--install|-u|--uninstall|-s|--status>",
		"ollama.err_no_action":                      "😢 错误: 未指定操作",
		"ollama.err_multi_actions":                  "😢 错误: 指定了多个操作，只能选择一个",
		"ollama.err_unknown_flag":                   "😢 错误: 未知的标志",
		"ollama.err_extra_argument":                 "😢 错误: 意外的参数 '%s'，此命令不接受额外参数",

		// macOS Ollama 特定消息
		"msg.ollama.mac.installing_via_homebrew":   "🦙 正在通过 Homebrew 安装 Ollama...",
		"msg.ollama.mac.installing_via_script":     "🦙 正在通过官方安装脚本安装 Ollama...",
		"msg.ollama.mac.uninstall_prompt":          "🦙 确定要卸载 Ollama 吗？",
		"msg.ollama.mac.uninstall_confirm":         "是否确认卸载? [Y/n]: ",
		"msg.ollama.mac.uninstall_cancelled":       "已取消卸载",
		"msg.ollama.mac.manual_uninstall_hint":     "⚠️  未找到自动卸载脚本，请手动卸载：",
		"msg.ollama.mac.manual_uninstall_steps":    "  请手动执行以下步骤：",
		"msg.ollama.mac.manual_uninstall_step1":    "  1. 删除 ~/.ollama 目录",
		"msg.ollama.mac.manual_uninstall_step2":    "  2. 从 PATH 中移除 ollama 命令（通常在 ~/.zshrc 或 ~/.bash_profile）",
		"msg.ollama.mac.manual_uninstall_step3":    "  3. 删除 /usr/local/bin/ollama 或 ~/.local/bin/ollama（如果存在）",
		"msg.ollama.mac.unknown_install_method":    "⚠️  无法确定 Ollama 的安装方式",
		"msg.ollama.mac.manual_uninstall_prompt":   "请手动卸载 Ollama：",
		"msg.ollama.mac.manual_uninstall_homebrew": "  - Homebrew: brew uninstall ollama",
		"msg.ollama.mac.manual_uninstall_script":   "  - 脚本安装: 删除 ~/.ollama 目录和相关 PATH 配置",
		"msg.ollama.mac.install_note_homebrew":     "💡 注意: 这将通过 Homebrew 安装 Ollama",
		"msg.ollama.mac.install_note_script":       "💡 注意: 这将使用官方安装脚本安装 Ollama",

		// Linux Ollama 特定消息
		"msg.ollama.linux.install_method_prompt":  "🦙 请选择安装方式：",
		"msg.ollama.linux.install_method_option1": "  1. 使用 %s 包管理器安装",
		"msg.ollama.linux.install_method_option2": "  2. 使用官方脚本安装（推荐）",
		"msg.ollama.linux.install_method_choice":  "请输入选项 (1/2，默认 2): ",
		"msg.ollama.linux.installing_via_pm":      "🦙 正在通过 %s 安装 Ollama...",
		"msg.ollama.linux.installing_via_script":  "🦙 正在通过官方脚本安装 Ollama...",
		"msg.ollama.linux.install_note_script":    "💡 注意: 官方脚本安装方式适用于所有 Linux 发行版",
		"msg.ollama.linux.pm_update_failed":       "包列表更新失败",
		"msg.ollama.linux.unsupported_pm":         "不支持的包管理器",
		"msg.ollama.linux.uninstall_prompt":       "🦙 检测到 Ollama 已安装，准备卸载",
		"msg.ollama.linux.uninstall_confirm":      "确定要卸载 Ollama 吗？(y/n，默认 n): ",
		"msg.ollama.linux.uninstall_cancelled":    "卸载已取消",
		"msg.ollama.linux.manual_uninstall_hint":  "⚠️  未找到自动卸载脚本，请手动卸载",
		"msg.ollama.linux.manual_uninstall_steps": "请参考官方文档进行手动卸载：https://ollama.com",

		// 命令行帮助
		"help.usage":        "使用方法:",
		"help.commands":     "命令:",
		"help.show_help":    "显示此帮助",
		"help.show_version": "显示版本号",
		"help.start_repl":   "启动 REPL 模式",
		"help.start_tui":    "启动终端 UI",
		"help.build_latex":  "构建 LaTeX 文档为 PDF",
		"help.live_preview": "实时预览 PDF 文档",
		"help.set_language": "设置语言 (zh/en)",
		"help.config":       "打开配置文件（JSON）",
		"help.latex":        "LaTeX 管理",
		"help.ollama":       "Ollama 管理",
		"help.llm":          "LLM 管理",
		"help.description":  "LaziTex: 零配置 LaTeX 编译工具，支持 AI 辅助",

		// 工具描述
		"desc.pdflatex":  "最常用的 PDF 编译器",
		"desc.xelatex":   "支持 Unicode 和现代字体",
		"desc.lualatex":  "Lua 扩展的现代引擎",
		"desc.latex":     "传统 LaTeX 编译器",
		"desc.pdftex":    "底层 TeX 引擎",
		"desc.tex":       "原始 TeX 引擎",
		"desc.etex":      "扩展 TeX 引擎",
		"desc.bibtex":    "传统参考文献管理",
		"desc.biber":     "现代参考文献管理",
		"desc.makeindex": "生成索引",
		"desc.xindy":     "多语言索引工具",
		"desc.texindy":   "LaTeX 索引包装器",
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
		"repl.ollama":           "Ollama 管理",
		"repl.ollama_usage":     "用法: ollama <-c|--check|-i|--install|-u|--uninstall|-s|--status>",
		"repl.latex_usage":      "用法: latex <-c|--check|-i|--install|-u|--uninstall>",
		"repl.latex":            "LaTeX 管理",
		"repl.build_desc":       "构建 LaTeX 文档",
		"repl.build_usage":      "用法: build <文件名.tex> [-o|--output 输出路径] [-s|--show] [-q|--quiet] [-t|--tidy]",
		"repl.preview_desc":     "实时预览 PDF 文档",
		"repl.preview_usage":    "用法: preview <文件名.tex> [-q|--quiet] [-t|--tidy] [-d|--dev] [:端口号]",
		"repl.lang_desc":        "切换语言 (zh/en)",
		"repl.quit_desc":        "退出 REPL",
		"repl.config_desc":      "打开配置文件（JSON）",
		"repl.llm_desc":         "LLM 管理",
		"repl.llm.link":         "连接并激活 LLM",
		"repl.llm.unlink":       "取消连接 LLM",
		"repl.llm_usage":        "用法: llm list | llm add <model> -p <provider> -n <name> | llm [remove|link|unlink|test|switch] <id/name/model>",
		"repl.lang_usage":       "用法: lang <zh|en>",
		"repl.lang_current":     "当前语言: ",
		"repl.lang_switched":    "✨ 语言已切换并保存",
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
		"title.env_check":  "LaTeX Environment Check Results",
		"label.os":         "🖥️  Operating System",
		"label.distro":     "📦 LaTeX Distribution",
		"label.installed":  "✨ Installed",
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
		"msg.no_compilers":       "⚠️  Warning: No core LaTeX compilers detected",
		"msg.partial_install":    "💡 Note: Some core tools are missing, but basic functionality is available",
		"msg.all_installed":      "🎉 Excellent! All core compilers are installed",
		"msg.install_guide":      "📖 Installation Guide:",
		"msg.unknown_distro":     "Unknown",
		"msg.tip":                "💡 Tip",
		"msg.warning":            "⚠️  Warning",
		"msg.unsupported_os":     "Error: Unsupported operating system '%s'",
		"msg.build_usage":        "Usage: lazitex -b <file.tex> [-o|--output output_path] [-s|--show] [-q|--quiet] [-t|--tidy]",
		"msg.checking_env":       "Checking LaTeX environment...",
		"msg.install_failed":     "Installation failed",
		"msg.install_success":    "✨ LaTeX environment installed successfully",
		"msg.uninstall_failed":   "Uninstallation failed",
		"msg.uninstall_success":  "✨ LaTeX environment uninstalled successfully",
		"msg.building_doc":       "🚀 Building LaTeX document: %s",
		"msg.working_dir":        "📁 Working directory: %s",
		"msg.output_dir":         "📁 Output directory: %s",
		"msg.build_success":      "✨ Build successful!",
		"msg.build_failed":       "😢 Build failed",
		"msg.err_abs_path":       "Error: Failed to get absolute path for '%s'",
		"msg.err_mkdir":          "Error: Failed to create output directory '%s'",
		"msg.err_invalid_ext":    "Error: Invalid file type '%s' (only .tex files supported)",
		"msg.err_file_not_found": "Error: File not found '%s'",
		"msg.preview_usage":      "Usage: lazitex -p <file.tex> [-q|--quiet] [-t|--tidy] [-d|--dev] [:port]",
		"msg.show_failed":        "⚠️  Warning: Show failed: %v",
		"msg.watching_file":      "👀 Watching file: %s (press Ctrl+C to stop)",
		"msg.watcher_error":      "Error: Watcher unexpectedly crashed: %v",

		// Web preview related
		"msg.server_starting":   "🌐 Server starting on port %d",
		"msg.server_error":      "😢 Server error: %v",
		"msg.opening_browser":   "🔗 Opening browser: %s",
		"msg.browser_error":     "⚠️  Failed to open browser: %v",
		"msg.manual_open":       "💡 Please manually open browser: %s",
		"msg.config_path_error": "Failed to get config path",
		"msg.config_save_error": "Failed to create default config",
		"msg.config_open_error": "Failed to open config file",
		"msg.config_opening":    "Opening config file",
		"msg.pdf_not_found":     "PDF file not found",
		"msg.invalid_port":      "😢 Error: Invalid port '%s' (must be a number between 1 and 65535)",
		"msg.dev_mode_vue_url":  "💡 Development mode: Vue frontend URL: %s",

		// LLM management related messages
		// Base & list
		"msg.llm.load_config_failed":      "😢 Failed to load config: %v",
		"msg.llm.no_providers":            "📝 No LLM providers configured, please add one first",
		"msg.llm.list_header_id":          "ID",
		"msg.llm.list_header_name":        "Name",
		"msg.llm.list_header_provider":    "Provider",
		"msg.llm.list_header_model":       "Model",
		"msg.llm.list_header_enabled":     "Enabled",
		"msg.llm.list_header_verified":    "Verified",
		"msg.llm.list_header_verified_at": "VerifiedAt",
		"msg.llm.list_header_active":      "Active",
		"msg.llm.date_format":             "2006-01-02", // Date format (Go time format string)
		// Test
		"msg.llm.test.test_failed":                     "😢 Test failed: %v",
		"msg.llm.test.not_found":                       "😢 LLM '%s' not found",
		"msg.llm.test.available_verified":              "✨ Available: %s (%s:%s) - Verified",
		"msg.llm.test.unavailable":                     "😢 Unavailable: %s (%s:%s)",
		"msg.llm.test.update_config_failed_on_success": "⚠️  Test passed, but failed to update config: %v",
		"msg.llm.test.update_config_failed_on_failure": "⚠️  Test failed, but failed to update config: %v",
		"msg.llm.test.multiple_matches":                "Found %d matching LLMs, please select one to test:",
		// Connect
		"msg.llm.link.not_found":              "😢 LLM Provider '%s' does not exist",
		"msg.llm.link.test_failed":            "😢 Connectivity test failed: %v",
		"msg.llm.link.unavailable":            "⚠️  LLM cannot be connected: %s (%s)",
		"msg.llm.link.set_active_failed":      "😢 Failed to connect: %v",
		"msg.llm.link.set_current_failed":     "😢 Failed to set foreground: %v",
		"msg.llm.link.success":                "✨ %s connected",
		"msg.llm.link.activated":              "✨ %s activated",
		"msg.llm.link.already_active":         "💡 %s is already active",
		"msg.llm.link.already_current":        "💡 %s is already the current foreground LLM",
		"msg.llm.link.multiple_matches":       "Found %d matching LLMs, please select one to link:",
		"msg.llm.link.batch_multiple_skip":    "⚠ %s matches multiple Providers, skipping ambiguous item in batch mode",
		"msg.llm.link.batch_already_active":   "💡 %s is already active, skipping",
		"msg.llm.link.batch_failed":           "😢 %s connection failed: %v",
		"msg.llm.link.model_multiple_title":   "Model '%s' matches multiple LLMs",
		"msg.llm.link.registered_title":       "Registered LLM Providers",
		"msg.llm.link.no_providers":           "No registered LLM Providers",
		"msg.llm.link.no_available":           "No available Providers",
		"msg.llm.link.select_prompt":          "Please select:",
		"msg.llm.link.select_number":          "  - Enter number (1-%d) to select specific LLM and activate",
		"msg.llm.link.select_all":             "  - Enter 'a' or 'all' to activate all",
		"msg.llm.link.select_cancel":          "  - Enter 'q' or 'quit' to cancel",
		"msg.llm.link.select_input":           "Choice: ",
		"msg.llm.link.read_input_failed":      "😢 Failed to read input: %v",
		"msg.llm.link.invalid_choice":         "😢 Invalid choice: %s",
		"msg.llm.link.activate_failed":        "😢 Activation failed: %v",
		"msg.llm.link.activate_success":       "✨ %s activated",
		"msg.llm.link.set_current_success":    "✨ Set as foreground",
		"msg.llm.link.activate_all_count":     "Activated %d Providers",
		"msg.llm.link.status_current":         " [Current Foreground]",
		"msg.llm.link.status_active":          " [Active]",
		"msg.llm.link.status_inactive":        " [Inactive]",
		"msg.llm.link.provider_format":        "  [%d] %s (Provider: %s, Model: %s)",
		"msg.llm.link.provider_format_simple": "  [%d] %s (ID: %s, Provider: %s)",
		// Disconnect
		"msg.llm.unlink.not_found":            "😢 LLM Provider '%s' does not exist",
		"msg.llm.unlink.not_active":           "⚠️  LLM Provider '%s' is not in the active list",
		"msg.llm.unlink.failed":               "😢 Failed to disconnect: %v",
		"msg.llm.unlink.success":              "✨ %s disconnected",
		"msg.llm.unlink.multiple_matches":     "Found %d matching LLMs, please select one to unlink:",
		"msg.llm.unlink.clear_current_failed": "⚠️  Failed to clear current LLM: %v",
		"msg.llm.unlink.batch_multiple_skip":  "⚠ %s matches multiple Providers, skipping ambiguous item in batch mode",
		"msg.llm.unlink.batch_failed":         "😢 %s disconnect failed: %v",
		"msg.llm.unlink.model_multiple_title": "Model '%s' matches multiple LLMs",
		"msg.llm.unlink.registered_title":     "Registered LLM Providers",
		"msg.llm.unlink.active_title":         "Active LLM Providers",
		"msg.llm.unlink.no_providers":         "No registered LLM Providers",
		"msg.llm.unlink.no_active_providers":  "No active LLM Providers",
		"msg.llm.unlink.select_prompt":        "Please select Provider to disconnect:",
		"msg.llm.unlink.select_number":        "  - Enter number (1-%d) to select specific LLM and disconnect",
		"msg.llm.unlink.select_all":           "  - Enter 'a' or 'all' to disconnect all",
		"msg.llm.unlink.select_cancel":        "  - Enter 'q' or 'quit' to cancel",
		"msg.llm.unlink.select_input":         "Choice: ",
		"msg.llm.unlink.unlink_failed":        "😢 Disconnect failed: %v",
		"msg.llm.unlink.unlink_success":       "✨ %s disconnected",
		"msg.llm.unlink.unlink_all_count":     "Disconnected %d Providers",
		// Switch
		"msg.llm.switch.not_found":            "😢 LLM not found: %s",
		"msg.llm.switch.not_active":           "⚠️  LLM '%s' is not active, please use link command first",
		"msg.llm.switch.failed":               "😢 Switch failed: %v",
		"msg.llm.switch.success":              "✨ Switched to LLM: %s",
		"msg.llm.switch.already_current":      "💡 %s is already the current foreground LLM",
		"msg.llm.switch.no_active":            "😢 No active LLM Providers",
		"msg.llm.switch.active_title":         "Active LLM Providers",
		"msg.llm.switch.model_multiple_title": "Model '%s' matches multiple active LLM Providers",
		"msg.llm.switch.select_prompt":        "Please select LLM to switch to foreground:",
		"msg.llm.switch.select_number":        "Enter number (1-%d) to select",
		"msg.llm.switch.select_cancel":        "Enter 'q' or 'quit' to cancel",
		"msg.llm.switch.select_input":         "Choice: ",
		"msg.llm.switch.batch_not_supported":  "😢 Batch switching not supported, can only switch to one LLM at a time",
		// Remove
		"msg.llm.remove.success":          "🗑️ Removed LLM: %s",
		"msg.llm.remove.save_failed":      "⚠️ Failed to save config: %s",
		"msg.llm.remove.not_found":        "😢 LLM not found: %s",
		"msg.llm.remove.cancelled":        "Deletion cancelled",
		"msg.llm.remove.multiple_matches": "Found %d matching LLM Providers, please select:",
		"msg.llm.remove.confirm_prompt":   "Are you sure you want to remove this LLM Provider? [y/N]: ",
		"msg.llm.remove.select_prompt":    "Please select the LLM Provider to remove (enter number, or 'q' to cancel): ",
		// Common selection prompt
		"msg.llm.select_prompt": "Please select (enter number, or 'q' to cancel): ",
		"msg.llm.cancelled":     "Cancelled",
		// Register related
		"msg.llm.register.title":                  "=== Register New LLM Provider ===",
		"msg.llm.register.confirm_new":            "Do you want to add a new LLM?",
		"msg.llm.register.input_provider":         "Enter LLM provider type (e.g.: ollama, openai)",
		"msg.llm.register.input_baseurl":          "Enter base URL (e.g.: http://localhost:11434)",
		"msg.llm.register.unsupported_provider":   "Unsupported provider type, currently only ollama is supported",
		"msg.llm.register.confirm_link_existing":  "Do you want to link an existing registered LLM?",
		"msg.llm.register.input_name":             "Give this LLM a name (e.g.: my-ollama)",
		"msg.llm.register.input_model":            "Enter model name (e.g.: qwen2.5-coder:7b)",
		"msg.llm.register.test_failed":            "Connection test failed: %v",
		"msg.llm.register.test_success":           "✨ %s",
		"msg.llm.register.save_failed":            "Save failed: %v",
		"msg.llm.register.set_active_failed":      "Failed to set active: %v",
		"msg.llm.register.set_current_failed":     "Failed to set current LLM: %v",
		"msg.llm.register.success":                "✨ Registration successful!",
		"msg.llm.register.name_empty":             "Name cannot be empty, please re-enter",
		"msg.llm.register.name_exists":            "Name '%s' already exists, please choose another name",
		"msg.llm.register.input_empty":            "Input cannot be empty, please re-enter",
		"msg.llm.register.confirm_invalid":        "Please enter Y or n",
		"msg.llm.register.ollama.connect_failed":  "Failed to connect to Ollama service: %w",
		"msg.llm.register.ollama.status_error":    "Ollama service returned error status: %d",
		"msg.llm.register.ollama.parse_failed":    "Failed to parse response: %w",
		"msg.llm.register.ollama.model_not_found": "Model '%s' does not exist, please run: ollama pull %s",
		"msg.llm.register.ollama.test_success":    "Connection successful, model '%s' is available",
		// Smart mode related
		// Add command related (supplement)
		"msg.llm.add.test_success":    "Connectivity test verified, registering LLM: %s (%s)",
		"msg.llm.add.test_failed":     "😢 Connectivity test failed, registration cancelled: %v",
		"msg.llm.add.register_failed": "😢 Failed to register: %v",
		"msg.llm.add.registered":      "✨ Registered LLM: %s (%s)",
		// Provider base errors
		"msg.llm.provider.unsupported": "Unsupported provider: %s",
		// Provider: Ollama
		"msg.llm.provider.ollama.service_not_running":      "Ollama service is not running",
		"msg.llm.provider.ollama.model_not_found":          "Model not found",
		"msg.llm.provider.ollama.api_connect_failed":       "Failed to connect to Ollama API: %v",
		"msg.llm.provider.ollama.api_status_error":         "Ollama API returned error status code: %d",
		"msg.llm.provider.ollama.parse_response_failed":    "Failed to parse Ollama response: %v",
		"msg.llm.provider.ollama.serialize_request_failed": "Failed to serialize request: %v",
		// CLI usage
		"msg.llm.cli_usage": "Usage: lazitex -m list | lazitex -m add <model> -p <provider> -n <name> | lazitex -m [remove|link|unlink|test|switch] <id/name/model>",
		// Add command related
		"msg.llm.add.model_required":             "😢 Model name is required",
		"msg.llm.add.provider_required":          "😢 Provider name is required",
		"msg.llm.add.name_required":              "😢 Custom name is required for the registered model",
		"msg.llm.add.provider_missing":           "😢 Provider is required",
		"msg.llm.add.unexpected_nonflag":         "😢 Unexpected non-flag argument: %s",
		"msg.llm.add.smart_mode_not_implemented": "⚠️  Smart registration/linking feature is not implemented yet",
		"msg.llm.add.smart_mode_hint":            "Please use: lazitex -m add <model> to register, or lazitex -m link <model> to connect",
		"msg.llm.unknown_action":                 "Unknown action: %s",
		"msg.dev_mode_backend_url":               "💡 Backend API URL: %s",
		"msg.prod_mode_embedded":                 "💡 Production mode: Frontend embedded, access: %s",
		"msg.dev_mode_custom_port_warning":       "⚠️  Warning: Custom port (%d) used in dev mode, but Vite proxy is hardcoded to 8080\n   Will force use default port 8080 to ensure frontend can connect to backend",
		"msg.port_8080_occupied":                 "⚠️  Port 8080 is already in use, cannot use dev mode\n   Will fallback to production mode with custom port %d",
		"msg.checking_vue_dev_server":            "🔍 Checking if Vue dev server (%s) is available...",
		"msg.vue_dev_server_available":           "✨ Vue dev server is available, using dev mode",
		"msg.vue_dev_server_unavailable":         "⚠️  Vue dev server is not available, auto-fallback to production mode (using embedded frontend)",
		"msg.common.yes":                         "Yes",
		"msg.common.no":                          "No",

		// Package management related
		"msg.package_missing":             "🔍 Missing package detected: %s",
		"msg.package_install_prompt":      "💡 Auto-install? [Y/n]: ",
		"msg.package_installing":          "📦 Installing %s...",
		"msg.package_install_success":     "✨ Installation successful!",
		"msg.package_install_failed":      "😢 Installation failed: %v",
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
		"msg.tidy_success":           "🧹 Cleaned up auxiliary files",
		"msg.tidy_warning":           "⚠️  Warning: Failed to clean auxiliary files",
		"msg.tidy_remove_failed":     "Failed to remove file %s: %v",

		// Compilation lock and task management
		"msg.compile_cancelled":        "Compilation cancelled",
		"msg.cancelling_previous_task": "⚠️  Cancelling previous compilation task...",
		"msg.compile_timeout":          "⏰ Compilation timeout (exceeded %d seconds), automatically cancelled",

		// macOS installation messages
		"msg.mac.brew_not_installed":       "Error: Homebrew is not installed. Please install Homebrew first: https://brew.sh",
		"msg.mac.latex_already_installed":  "LaTeX is already installed",
		"msg.mac.checking_updates":         "Checking for updates...",
		"msg.mac.update_failed":            "Update failed",
		"msg.mac.update_success":           "✨ Update successful!",
		"msg.mac.installing_basictex":      "Installing BasicTeX via Homebrew...",
		"msg.mac.install_warning":          "This may take a while (approximately 1.5 GB download)...",
		"msg.mac.install_failed":           "Installation failed",
		"msg.mac.updating_tlmgr":           "Updating TeX Live package manager...",
		"msg.mac.update_warning":           "Warning: Error occurred during update",
		"msg.mac.install_success":          "✨ BasicTeX installed successfully!",
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
		"msg.mac.install_prompt":           "LaTeX environment not detected",
		"msg.mac.install_info":             "Will install BasicTeX via Homebrew (lightweight version, ~1.5 GB)\nInstall command: brew install --cask basictex",
		"msg.mac.confirm_install":          "Do you want to install? [y/n]: ",
		"msg.mac.install_cancelled":        "Installation cancelled",
		"msg.mac.latex_not_installed":      "LaTeX environment not detected",
		"msg.mac.uninstall_prompt":         "LaTeX is already installed",
		"msg.mac.uninstall_info":           "Will uninstall LaTeX environment (including all installed packages)",
		"msg.mac.confirm_uninstall":        "Do you want to uninstall? [y/n]: ",
		"msg.mac.uninstall_cancelled":      "Uninstallation cancelled",
		"msg.mac.uninstalling":             "Uninstalling LaTeX environment...",
		"msg.mac.uninstall_success":        "✨ LaTeX environment uninstalled successfully!",
		"msg.mac.uninstall_failed":         "Uninstallation failed",
		"msg.mac.uninstall_install_hint":   "Tip: You can use 'lazitex -i' to install LaTeX environment with one click",
		"msg.mac.uninstall_mactex_running": "Detected official MacTeX, attempting uninstall script...",
		"msg.mac.uninstall_mactex_manual":  "Detected official MacTeX but uninstall script not found. Please remove manually:",
		"msg.mac.removing_residual":        "Removing residual path: %s",

		// LaTeX related messages
		"msg.latex_usage":                 "Usage: lazitex -x <-c|--check|-i|--install|-u|--uninstall>",
		"latex.err_no_action":             "😢 Error: No action specified",
		"latex.err_multi_actions":         "😢 Error: Multiple actions specified, only one allowed",
		"latex.err_unknown_flag":          "😢 Error: Unknown flag",
		"latex.err_extra_argument":        "😢 Error: Unexpected argument '%s', this command does not accept extra arguments",
		"msg.latex.not_implemented":       "😢 Feature not implemented (under development)",
		"msg.latex.linux.not_implemented": "😢 Installation/uninstallation on Linux is not implemented (under development)",
		"msg.latex.win.not_implemented":   "😢 Installation/uninstallation on Windows is not implemented (under development)",

		// Ollama related messages
		"msg.ollama.already_installed":              "🦙 Ollama is already installed, version: %s",
		"msg.ollama.installing":                     "🦙 Installing Ollama via winget...",
		"msg.ollama.install_note_gui":               "💡 Note: Windows version includes GUI, but CLI tools are also available",
		"msg.ollama.install_success":                "✨ Ollama installed successfully! Version: %s",
		"msg.ollama.install_failed":                 "😢 winget installation failed",
		"msg.ollama.install_verify_failed":          "⚠️  Installation completed but verification failed, please check manually",
		"msg.ollama.confirm_install":                "🦙 Are you sure you want to install Ollama? [y/N]: ",
		"msg.ollama.install_cancelled":              "Install cancelled",
		"msg.ollama.winget_not_available":           "⚠️  winget is not available, please install winget first or install Ollama manually",
		"msg.ollama.uninstalling":                   "🦙 Uninstalling Ollama via winget...",
		"msg.ollama.uninstall_success":              "✨ Ollama uninstalled successfully!",
		"msg.ollama.uninstall_failed":               "😢 winget uninstallation failed",
		"msg.ollama.uninstall_verify_failed":        "⚠️  Uninstallation completed but verification failed, there may be residual files",
		"msg.ollama.not_installed":                  "🦙 Ollama is not installed",
		"msg.ollama.confirm_uninstall":              "🦙 Are you sure you want to uninstall Ollama? [y/N]: ",
		"msg.ollama.uninstall_cancelled":            "Uninstall cancelled",
		"msg.ollama.winget_uninstall_not_available": "⚠️  winget is not available, please uninstall Ollama manually",
		"msg.ollama.service_running":                "✨ Ollama service is running",
		"msg.ollama.service_not_running":            "⚠️  Ollama service is not running",
		"msg.ollama.start_service_hint":             "💡 Tip: Run 'ollama serve' to start the service, or restart the Ollama application",
		"msg.ollama.path_not_updated":               "⚠️  Note: Ollama is installed, but the command-line tool may not be in PATH yet",
		"msg.ollama.restart_terminal_hint":          "💡 Tip: Please restart your terminal or open a new command window, then run 'ollama -c' again to check",
		"msg.ollama.install_maybe_success":          "⚠️  Ollama may be installed, but the command-line tool is not in PATH yet",
		"msg.ollama.winget_install_prompt":          "💡 winget is Windows package manager, used to automatically install Ollama",
		"msg.ollama.winget_install_confirm":         "Would you like to automatically install winget? [Y/n]: ",
		"msg.ollama.winget_install_cancelled":       "winget installation cancelled",
		"msg.ollama.winget_required":                "winget is required to continue installing Ollama",
		"msg.ollama.installing_winget":              "📦 Installing winget...",
		"msg.ollama.winget_store_opened":            "🔗 Microsoft Store opened",
		"msg.ollama.winget_store_instruction":       "Please click 'Get' or 'Install' button in Microsoft Store",
		"msg.ollama.winget_wait_prompt":             "After installation, please press Enter to continue...",
		"msg.ollama.winget_install_success":         "✨ winget installed successfully!",
		"msg.ollama.winget_install_error":           "winget installation failed",
		"msg.ollama.downloading_winget":             "📥 Downloading winget installer...",
		"msg.ollama.winget_download_failed":         "Failed to download winget installer",
		"msg.ollama.installing_winget_package":      "📦 Installing winget package...",
		"msg.ollama.winget_install_package_failed":  "Failed to install winget package",
		"msg.ollama.winget_install_verify_failed":   "winget installation completed but verification failed, please check manually",
		"msg.ollama.invalid_input":                  "Invalid input, please enter y (yes) or n (no)",
		"msg.ollama.manager_nil":                    "😢 Ollama manager is not initialized",
		"msg.ollama.check_failed":                   "Failed to check Ollama status",
		"msg.ollama.check_installed":                "Ollama is installed, version",
		"msg.ollama.check_installed_no_version":     "Ollama is installed (unable to get version)",
		"msg.ollama.status_installed":               "🦙 Ollama is installed, version: %s",
		"msg.ollama.status_installed_no_version":    "🦙 Ollama is installed",
		"msg.ollama.mac.install_method_prompt":      "🦙 Homebrew detected, choose installation method:",
		"msg.ollama.mac.install_method_option1":     "  1. Homebrew (recommended, fast installation)",
		"msg.ollama.mac.install_method_option2":     "  2. Official install script (direct installation)",
		"msg.ollama.mac.install_method_choice":      "Choose [1/2] (default: 1): ",
		"msg.ollama_usage":                          "Usage: lazitex -o <-c|--check|-i|--install|-u|--uninstall|-s|--status>",
		"ollama.err_no_action":                      "😢 Error: No action specified",
		"ollama.err_multi_actions":                  "😢 Error: Multiple actions specified, only one allowed",
		"ollama.err_unknown_flag":                   "😢 Error: Unknown flag",
		"ollama.err_extra_argument":                 "😢 Error: Unexpected argument '%s', this command does not accept extra arguments",

		// macOS Ollama specific messages
		"msg.ollama.mac.installing_via_homebrew":   "🦙 Installing Ollama via Homebrew...",
		"msg.ollama.mac.installing_via_script":     "🦙 Installing Ollama via official script...",
		"msg.ollama.mac.uninstall_prompt":          "🦙 Are you sure you want to uninstall Ollama?",
		"msg.ollama.mac.uninstall_confirm":         "Confirm uninstall? [Y/n]: ",
		"msg.ollama.mac.uninstall_cancelled":       "Uninstall cancelled",
		"msg.ollama.mac.manual_uninstall_hint":     "⚠️  Automatic uninstall script not found, please uninstall manually:",
		"msg.ollama.mac.manual_uninstall_steps":    "  Please manually perform the following steps:",
		"msg.ollama.mac.manual_uninstall_step1":    "  1. Delete ~/.ollama directory",
		"msg.ollama.mac.manual_uninstall_step2":    "  2. Remove ollama command from PATH (usually in ~/.zshrc or ~/.bash_profile)",
		"msg.ollama.mac.manual_uninstall_step3":    "  3. Delete /usr/local/bin/ollama or ~/.local/bin/ollama (if exists)",
		"msg.ollama.mac.unknown_install_method":    "⚠️  Unable to determine Ollama installation method",
		"msg.ollama.mac.manual_uninstall_prompt":   "Please manually uninstall Ollama:",
		"msg.ollama.mac.manual_uninstall_homebrew": "  - Homebrew: brew uninstall ollama",
		"msg.ollama.mac.manual_uninstall_script":   "  - Script installation: Delete ~/.ollama directory and related PATH configuration",
		"msg.ollama.mac.install_note_homebrew":     "💡 Note: This will install Ollama via Homebrew",
		"msg.ollama.mac.install_note_script":       "💡 Note: This will install Ollama via official install script",

		// Linux Ollama specific messages
		"msg.ollama.linux.install_method_prompt":  "🦙 Choose installation method:",
		"msg.ollama.linux.install_method_option1": "  1. Install via %s package manager",
		"msg.ollama.linux.install_method_option2": "  2. Install via official script (recommended)",
		"msg.ollama.linux.install_method_choice":  "Enter option (1/2, default 2): ",
		"msg.ollama.linux.installing_via_pm":      "🦙 Installing Ollama via %s...",
		"msg.ollama.linux.installing_via_script":  "🦙 Installing Ollama via official script...",
		"msg.ollama.linux.install_note_script":    "💡 Note: Official script installation works on all Linux distributions",
		"msg.ollama.linux.pm_update_failed":       "Package list update failed",
		"msg.ollama.linux.unsupported_pm":         "Unsupported package manager",
		"msg.ollama.linux.uninstall_prompt":       "🦙 Ollama is installed, ready to uninstall",
		"msg.ollama.linux.uninstall_confirm":      "Are you sure you want to uninstall Ollama? (y/n, default n): ",
		"msg.ollama.linux.uninstall_cancelled":    "Uninstall cancelled",
		"msg.ollama.linux.manual_uninstall_hint":  "⚠️  Automatic uninstall script not found, please uninstall manually",
		"msg.ollama.linux.manual_uninstall_steps": "Please refer to official documentation for manual uninstallation: https://ollama.com",

		// CLI help
		"help.usage":        "Usage:",
		"help.commands":     "Commands:",
		"help.show_help":    "Show this help",
		"help.show_version": "Show version",
		"help.start_repl":   "Start REPL mode",
		"help.start_tui":    "Start terminal UI",
		"help.build_latex":  "Build LaTeX document to PDF",
		"help.live_preview": "Live preview PDF document",
		"help.set_language": "Set language (zh/en)",
		"help.config":       "Open config file (JSON)",
		"help.latex":        "LaTeX management",
		"help.ollama":       "Ollama management",
		"help.llm":          "LLM management",
		"help.description":  "LaziTex: Zero-config LaTeX compilation with AI assistance",

		// Tool descriptions
		"desc.pdflatex":  "Most common PDF compiler",
		"desc.xelatex":   "Unicode and modern font support",
		"desc.lualatex":  "Modern engine with Lua extensions",
		"desc.latex":     "Traditional LaTeX compiler",
		"desc.pdftex":    "Low-level TeX engine",
		"desc.tex":       "Original TeX engine",
		"desc.etex":      "Extended TeX engine",
		"desc.bibtex":    "Traditional bibliography management",
		"desc.biber":     "Modern bibliography management",
		"desc.makeindex": "Generate index",
		"desc.xindy":     "Multilingual indexing tool",
		"desc.texindy":   "LaTeX index wrapper",
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
		"repl.ollama":           "Ollama management",
		"repl.ollama_usage":     "Usage: ollama <-c|--check|-i|--install|-u|--uninstall|-s|--status>",
		"repl.latex_usage":      "Usage: latex <-c|--check|-i|--install|-u|--uninstall>",
		"repl.latex":            "LaTeX management",
		"repl.build_desc":       "Build LaTeX document (use -o for output, -s to show after build, -q to quiet mode, -t to tidy mode)",
		"repl.build_usage":      "Usage: build <filename.tex> [-o|--output output_path] [-s|--show] [-q|--quiet] [-t|--tidy]",
		"repl.preview_desc":     "Live preview PDF document (use -q for quiet mode, -t for tidy mode, -d for dev mode, :port for custom port)",
		"repl.preview_usage":    "Usage: preview <filename.tex> [-q|--quiet] [-t|--tidy] [-d|--dev] [:port]",
		"repl.lang_desc":        "Switch language (zh/en)",
		"repl.quit_desc":        "Exit REPL",
		"repl.config_desc":      "Open config file (JSON)",
		"repl.llm_desc":         "LLM management",
		"repl.llm.link":         "Connect and activate LLM",
		"repl.llm.unlink":       "Disconnect LLM",
		"repl.llm_usage":        "Usage: llm list | llm add <model> -p <provider> -n <name> | llm [remove|link|unlink|test|switch] <id/name/model>",
		"repl.lang_usage":       "Usage: lang <zh|en>",
		"repl.lang_current":     "Current language: ",
		"repl.lang_switched":    "✨ Language switched and saved",
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
