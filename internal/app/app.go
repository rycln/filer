package app

import (
	"os"
	"regexp"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/filer/internal/domain"
	"github.com/rycln/filer/internal/infrastructure/config"
	"github.com/rycln/filer/internal/infrastructure/filesystem"
	"github.com/rycln/filer/internal/infrastructure/tui"
	"github.com/rycln/filer/internal/usecases"
)

type App struct {
	tui *tea.Program
}

func New() (*App, error) {
	cfg := config.NewConfigBuilder().WithFlagParsing().Build()

	filesys, err := filesystem.NewLocal(cfg.Source, cfg.Target)
	if err != nil {
		return nil, err
	}

	filenames, err := filesys.GetFilenames()
	if err != nil {
		return nil, err
	}

	filtered, err := filterAndSortFilenames(filenames, cfg.Pattern)
	if err != nil {
		return nil, err
	}

	batch := domain.NewFileBatch(filtered)
	fileProcessor := usecases.NewFileProcessor(filesys)

	p := tea.NewProgram(tui.InitialModel(batch, fileProcessor))

	return &App{
		tui: p,
	}, nil
}

func filterAndSortFilenames(filenames []string, pattern string) ([]string, error) {
	var filtered []string

	for _, filename := range filenames {
		matched, err := regexp.MatchString(pattern, filename)
		if err != nil {
			return nil, err
		}
		if matched {
			filtered = append(filtered, filename)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i] < filtered[j]
	})

	return filtered, nil
}

func (app *App) Run() error {
	_, err := app.tui.Run()
	if err != nil {
		os.Exit(1)
	}

	return nil
}
