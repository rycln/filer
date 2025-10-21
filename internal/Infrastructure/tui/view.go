package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	switch m.state {
	case FileManageState:
		return m.fileManageView()
	case ProcessingState:
		return m.processingView()
	case EndState:
		return m.endView()
	case ErrorState:
		return m.errorView()
	}

	return ""
}

func (m Model) fileManageView() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("[%d/%d]\n\n", m.batch.Progress(), m.batch.TotalFiles()))
	s.WriteString(fmt.Sprintf("File: %s\n\n", m.batch.CurrentFile()))
	s.WriteString("[K]eep, [D]elete, [S]kip, [Q]uit? ")

	return s.String()
}

func (m Model) processingView() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("[%d/%d]\n\n", m.batch.Progress(), m.batch.TotalFiles()))
	s.WriteString(fmt.Sprintf("File: %s\n\n", m.batch.CurrentFile()))
	s.WriteString("Processing...\n")

	return s.String()
}

func (m Model) endView() string {
	return fmt.Sprintf("Completed! Processed %d files.\nPress any key to exit...\n", m.batch.TotalFiles())
}

func (m Model) errorView() string {
	var s strings.Builder

	s.WriteString("Error occurred:\n\n")
	s.WriteString(fmt.Sprintf("%s\n\n", m.errMsg))
	s.WriteString("Press any key to exit...\n")

	return s.String()
}
