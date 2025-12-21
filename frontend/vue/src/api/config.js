// api/config.js
// config.js 是 LaziTeX 的配置 API，负责管理应用的配置

/**
 * 获取应用配置
 * @returns {Promise<{mode: string}>}
 */
export async function getConfig() {
    const response = await fetch('/api/config')
    
    if (!response.ok) {
        throw new Error(`配置请求失败: ${response.status} ${response.statusText}`)
    }
    
    return await response.json()
}