// cmd/lazitex-cli/modes/tui.go

package modes

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

// model 定义 TUI 的状态
type model struct {
	quitting bool
}

// 结构体model的Init方法
func (m model) Init() tea.Cmd {
	return nil
}

// 初始化model
func initialModel() model {
	return model{quitting: false}
}

// 结构体model的Update方法
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// 结构体model的View方法
func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	return fmt.Sprintf(
		"\n%s\n\n%s\n",
		"Welcome to LaziTex TUI mode!",
		"Press 'q' to quit",
	)
}

// 启动函数
func StartTUI() {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
