package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case FileManageState:
		handleFileManageState(m, msg)
	case ProcessingState:
		handleProcessingeState(m, msg)
	case EndState:
		handleEndState(m, msg)
	case ErrorState:
		handleErrorState(m, msg)
	}

	return m, nil
}

func handleFileManageState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyRunes:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "k":
				m.state = ProcessingState
				return m, m.keep()
			case "d":
				m.state = ProcessingState
				return m, m.delete()
			case "s":
				m.idx++
				return m, nil
			}
		}
	}

	return m, nil
}

func (m *Model) keep() tea.Cmd {
	return func() tea.Msg {
		err := m.manager.KeepFile(m.names[m.idx])
		if err != nil {
			return ErrorMsg{
				Err: err,
			}
		}

		return SuccessMsg{}
	}
}

func (m *Model) delete() tea.Cmd {
	return func() tea.Msg {
		err := m.manager.DeleteFile(m.names[m.idx])
		if err != nil {
			return ErrorMsg{
				Err: err,
			}
		}

		return SuccessMsg{}
	}
}

func handleProcessingeState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyRunes:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			}
		}
	case ErrorMsg:
		m.errMsg = msg.Err.Error()
		m.state = ErrorState
	case SuccessMsg:
		m.idx++
		m.state = FileManageState
	}

	return m, nil
}

func handleEndState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	}

	return m, nil
}

func handleErrorState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	}

	return m, nil
}
