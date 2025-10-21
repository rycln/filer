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
	KeepFile(filename string) error
	DeleteFile(filename string) error
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
