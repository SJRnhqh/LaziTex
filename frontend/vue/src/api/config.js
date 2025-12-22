// api/config.js
// config.js 是 LaziTeX 的配置 API，负责管理应用的配置

/**
 * 获取应用配置
 * @returns {Promise<{mode?: string, displayMode?: string}>}
 */
export async function getConfig() {
    const response = await fetch('/api/config')
    
    if (!response.ok) {
        throw new Error(`配置请求失败: ${response.status} ${response.statusText}`)
    }
    
    return await response.json()
}

/**
 * 保存应用配置
 * @param {Object} config - 配置对象 {displayMode?: string, ...}
 * @returns {Promise<void>}
 */
export async function setConfig(config) {
    const response = await fetch('/api/config', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(config)
    })
    
    if (!response.ok) {
        throw new Error(`保存配置失败: ${response.status} ${response.statusText}`)
    }
}