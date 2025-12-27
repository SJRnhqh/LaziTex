<template>
    <div ref="editorContainer" class="monaco-editor-container">
        <div v-if="loading" class="loading">加载编辑器中...</div>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'

// 定义 props：接收父组件传入的内容
const props = defineProps({
    modelValue: {
        type: String,
        default: ''
    }
})

// 定义 emit：向父组件发送内容变化
const emit = defineEmits(['update:modelValue'])

const editorContainer = ref(null)
const loading = ref(true)
let editor = null
let monaco = null

onMounted(async () => {
    try {
        // 动态导入 Monaco Editor（避免构建时的静态分析问题）
        const monacoModule = await import('monaco-editor')
        monaco = monacoModule.default || monacoModule
        
        // 注册 LaTeX 语言
        monaco.languages.register({ id: 'latex' })
        
        // 设置 LaTeX 的语法高亮规则
        monaco.languages.setMonarchTokensProvider('latex', {
        tokenizer: {
            root: [
            // LaTeX 命令：\command
            [/\\[a-zA-Z@]+/, 'keyword'],
            // 环境：\begin{...} \end{...}
            [/\\begin\{[^}]+\}/, 'keyword'],
            [/\\end\{[^}]+\}/, 'keyword'],
            // 注释：%
            [/%[^\n]*/, 'comment'],
            // 字符串：{...}
            [/\{[^}]*\}/, 'string'],
            // 数字
            [/\d+/, 'number'],
            ]
        }
})
    
    // 创建编辑器实例
    editor = monaco.editor.create(editorContainer.value, {
        value: props.modelValue || '',
        language: 'latex',
        theme: 'vs-dark',
        automaticLayout: true,
        fontSize: 14,
        minimap: { enabled: false },
        wordWrap: 'on',
    })
    
    // 监听内容变化
    editor.onDidChangeModelContent(() => {
        const value = editor.getValue()
        emit('update:modelValue', value)
    })
    
    loading.value = false
    } catch (error) {
        console.error('Failed to load Monaco Editor:', error)
        loading.value = false
    }
})

// 监听 props 变化，同步到编辑器
watch(() => props.modelValue, (newValue) => {
    if (editor && editor.getValue() !== newValue) {
        editor.setValue(newValue || '')
    }
})

onUnmounted(() => {
  // 清理编辑器
    if (editor) {
        editor.dispose()
    }
})
</script>

<style scoped>
.monaco-editor-container {
    width: 100%;
    height: 100%;
    position: relative;
}

.loading {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    color: #D8DEE9;
    font-size: 14px;
}
</style>