// target/linux/installer.go

package linux

import (
	"fmt"
)

// Installer Linux 平台安装器
type Installer struct{}

// NewInstaller 创建 Linux 安装器
func NewInstaller() *Installer {
	return &Installer{}
}

// Install 安装 LaTeX 环境
func (i *Installer) Install() error {
	// TODO: 实现具体的安装逻辑
	return fmt.Errorf("not implemented yet")
}

// Uninstall 卸载 LaTeX 环境
func (i *Installer) Uninstall() error {
	// TODO: 实现具体的卸载逻辑
	return fmt.Errorf("not implemented yet")
}
