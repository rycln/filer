package tui

import "github.com/rycln/filer/internal/domain"

type state int

const (
	FileManageState state = iota
	ProcessingState
	EndState
	ErrorState
)

type (
	SuccessMsg struct{}
	ErrorMsg   struct{ Err error }
)

type FileManager interface {
	Keep(string) error
	Delete(string) error
}

type Model struct {
	state   state
	errMsg  string
	batch   *domain.FileBatch
	manager FileManager
}

func InitialRootModel(batch *domain.FileBatch, manager FileManager) Model {
	return Model{
		state:   FileManageState,
		batch:   batch,
		manager: manager,
	}
}
