// target/linux/latex.go
// Linux 平台特定 LaTeX 管理器实现

package linux

import (
	// 内部包
	core "github.com/SJRnhqh/lazitex/core"
)

// LaTeXManager Linux 平台 LaTeX 管理器
type LaTeXManager struct {
	checker   *Checker
	installer *Installer
}

// NewLaTeXManager 创建 Linux LaTeX 管理器
func NewLaTeXManager() *LaTeXManager {
	return &LaTeXManager{
		checker:   NewChecker(),
		installer: NewInstaller(),
	}
}

// 实现 EnvironmentChecker 接口
func (m *LaTeXManager) GetPlatformName() string {
	return m.checker.GetPlatformName()
}

func (m *LaTeXManager) GetSearchPaths() []string {
	return m.checker.GetSearchPaths()
}

func (m *LaTeXManager) DetectDistribution(tools map[string]core.CompilerInfo) string {
	return m.checker.DetectDistribution(tools)
}

func (m *LaTeXManager) GetInstallGuide() string {
	return m.checker.GetInstallGuide()
}

func (m *LaTeXManager) PostCheck(env *core.LaTeXEnvironment) {
	m.checker.PostCheck(env)
}

// 实现 EnvironmentInstaller 接口
func (m *LaTeXManager) Install() error {
	return m.installer.Install()
}

func (m *LaTeXManager) Uninstall() error {
	return m.installer.Uninstall()
}