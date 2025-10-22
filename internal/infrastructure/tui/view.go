package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var s strings.Builder
	s.WriteString("\n")

	switch m.state {
	case FileManageState:
		s.WriteString(m.fileManageView())
	case ProcessingState:
		s.WriteString(m.processingView())
	case EndState:
		s.WriteString(m.endView())
	case ErrorState:
		s.WriteString(m.errorView())
	}

	return s.String()
}

func (m Model) fileManageView() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("📁 File Manager"))
	s.WriteString("\n")

	progress := m.createProgressBar(m.batch.Progress(), m.batch.TotalFiles())
	s.WriteString(progress)
	s.WriteString("\n\n")

	currentFile := fmt.Sprintf("📄 %s", m.batch.CurrentFile())
	s.WriteString(fileStyle.Render(currentFile))
	s.WriteString("\n\n")

	options := []string{
		optionStyle.Render("K") + "eep",
		optionStyle.Render("D") + "elete",
		optionStyle.Render("S") + "kip",
		optionStyle.Render("Q") + "uit",
	}
	optionsLine := strings.Join(options, " "+dividerStyle.String()+" ")

	s.WriteString("❓ Action: ")
	s.WriteString(optionsLine)

	return s.String()
}

func (m Model) processingView() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("⚙️  Processing Files"))
	s.WriteString("\n")

	progress := m.createProgressBar(m.batch.Progress(), m.batch.TotalFiles())
	s.WriteString(progress)
	s.WriteString("\n\n")

	currentFile := fmt.Sprintf("📄 %s", m.batch.CurrentFile())
	s.WriteString(fileStyle.Render(currentFile))
	s.WriteString("\n\n")

	s.WriteString(processingStyle.Render("⏳ Processing..."))

	return s.String()
}

func (m Model) endView() string {
	var s strings.Builder

	s.WriteString(successStyle.Render("🎉 Processing Complete!"))
	s.WriteString("\n\n")

	stats := fmt.Sprintf("✅ Processed %d files", m.batch.TotalFiles())
	s.WriteString(progressStyle.Render(stats))
	s.WriteString("\n\n")

	s.WriteString("👆 Press any key to exit")

	return s.String()
}

func (m Model) errorView() string {
	var s strings.Builder

	s.WriteString(errorStyle.Render("❌ Error Occurred"))
	s.WriteString("\n\n")

	errorMsg := lipgloss.NewStyle().
		Foreground(lipgloss.Color("203")).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("196")).
		PaddingLeft(1).
		Render(m.errMsg)

	s.WriteString(errorMsg)
	s.WriteString("\n\n")

	s.WriteString("👆 Press any key to exit")

	return s.String()
}

func (m Model) createProgressBar(current, total int) string {
	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	width := 30
	filled := int(percentage * float64(width))
	empty := width - filled

	filledBar := strings.Repeat("█", filled)
	emptyBar := strings.Repeat("░", empty)

	filledStyled := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Render(filledBar)

	emptyStyled := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(emptyBar)

	progressText := fmt.Sprintf(" %d/%d (%.1f%%)", current, total, percentage*100)

	progressTextStyled := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Render(progressText)

	return filledStyled + emptyStyled + progressTextStyled
}
