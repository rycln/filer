package app

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/filer/internal/domain"
	"github.com/rycln/filer/internal/infrastructure/config"
	"github.com/rycln/filer/internal/infrastructure/filesystem"
	"github.com/rycln/filer/internal/infrastructure/filter"
	"github.com/rycln/filer/internal/infrastructure/tui"
	"github.com/rycln/filer/internal/usecases"
)

type App struct {
	tui *tea.Program
}

func New() (*App, error) {
	cfg, err := config.NewConfigBuilder().WithFlagParsing().Build()
	if err != nil {
		return nil, err
	}

	filesys, err := filesystem.NewLocal(cfg.Source, cfg.Target)
	if err != nil {
		return nil, err
	}

	filenames, err := filesys.GetFilenames()
	if err != nil {
		return nil, err
	}

	fileFilter := filter.NewRegexpFilter(cfg.Pattern)
	filtered, err := fileFilter.Filter(filenames)
	if err != nil {
		return nil, err
	}

	batch, err := domain.NewFileBatch(filtered)
	if err != nil {
		return nil, err
	}
	fileProcessor := usecases.NewFileProcessor(filesys)

	p := tea.NewProgram(tui.InitialModel(batch, fileProcessor))

	return &App{
		tui: p,
	}, nil
}

func (app *App) Run() error {
	_, err := app.tui.Run()
	if err != nil {
		os.Exit(1)
	}

	return nil
}
