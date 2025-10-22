package tui

import "github.com/rycln/filer/internal/domain"

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// state represents TUI application states.
// Controls UI rendering and user input handling.
type state int

const (
	FileManageState state = iota // Main file management interface
	ProcessingState              // File operation in progress
	EndState                     // Processing completed
	ErrorState                   // Error display state
)

// SuccessMsg indicates successful file operation.
// Used for state transitions after keep/delete.
type SuccessMsg struct{}

// ErrorMsg wraps file operation errors.
// Carries error details for error state.
type ErrorMsg struct{ Err error }

// FileManager defines file operations for TUI.
// Abstraction for keep/delete business logic.
type FileManager interface {
	Keep(string) error
	Delete(string) error
}

// Model represents TUI application state.
// Manages UI state, file batch and business logic.
type Model struct {
	state   state
	errMsg  string
	batch   *domain.FileBatch
	manager FileManager
}

// InitialModel creates TUI model with file batch.
// Starts in FileManageState for user interaction.
func InitialModel(batch *domain.FileBatch, manager FileManager) Model {
	return Model{
		state:   FileManageState,
		batch:   batch,
		manager: manager,
	}
}
