package tui

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
	total   int
	idx     int
	errMsg  string
	names   []string
	manager FileManager
}

func InitialRootModel(names []string, manager FileManager) Model {
	return Model{
		state:   FileManageState,
		total:   len(names),
		names:   names,
		manager: manager,
	}
}
