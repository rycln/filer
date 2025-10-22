package tui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/rycln/filer/internal/domain"
	"github.com/rycln/filer/internal/infrastructure/tui/mocks"
)

func TestInitialModel(t *testing.T) {
	t.Run("should create model with initial state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt", "file2.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		model := InitialModel(batch, mockManager)

		if model.state != FileManageState {
			t.Errorf("Expected initial state FileManageState, got %v", model.state)
		}
		if model.batch != batch {
			t.Error("Batch not set correctly")
		}
		if model.manager != mockManager {
			t.Error("Manager not set correctly")
		}
		if model.errMsg != "" {
			t.Errorf("Expected empty error message, got %s", model.errMsg)
		}
	})
}

func TestModel_Init(t *testing.T) {
	t.Run("should return nil command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		cmd := model.Init()

		if cmd != nil {
			t.Error("Expected nil command from Init")
		}
	})
}

func TestModel_Update_FileManageState(t *testing.T) {
	t.Run("should quit on Ctrl+C", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		msg := tea.KeyMsg{Type: tea.KeyCtrlC}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if cmd == nil {
			t.Error("Expected quit command")
		}
		if updatedModel.state != FileManageState {
			t.Error("State should not change on quit command")
		}
	})

	t.Run("should quit on 'q' key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'q'},
		}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if cmd == nil {
			t.Error("Expected quit command")
		}
		if updatedModel.state != FileManageState {
			t.Error("State should not change on quit command")
		}
	})

	t.Run("should transition to ProcessingState on 'k' key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'k'},
		}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if updatedModel.state != ProcessingState {
			t.Errorf("Expected ProcessingState, got %v", updatedModel.state)
		}
		if cmd == nil {
			t.Error("Expected keep command")
		}
	})

	t.Run("should transition to ProcessingState on 'd' key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'d'},
		}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if updatedModel.state != ProcessingState {
			t.Errorf("Expected ProcessingState, got %v", updatedModel.state)
		}
		if cmd == nil {
			t.Error("Expected delete command")
		}
	})

	t.Run("should skip file on 's' key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt", "file2.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		initialProgress := batch.Progress()

		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'s'},
		}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if updatedModel.batch.Progress() != initialProgress+1 {
			t.Error("Progress should increment after skip")
		}
		if updatedModel.state != FileManageState {
			t.Errorf("Expected FileManageState, got %v", updatedModel.state)
		}
		if cmd != nil {
			t.Error("Expected nil command after skip")
		}
	})

	t.Run("should transition to EndState when skipping last file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'s'},
		}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if updatedModel.state != EndState {
			t.Errorf("Expected EndState, got %v", updatedModel.state)
		}
		if !updatedModel.batch.IsComplete() {
			t.Error("Batch should be complete after skipping last file")
		}
		if cmd != nil {
			t.Error("Expected nil command after skip")
		}
	})
}

func TestModel_Update_ProcessingState(t *testing.T) {
	t.Run("should handle success message and move to next file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt", "file2.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)
		model.state = ProcessingState

		initialProgress := batch.Progress()
		msg := SuccessMsg{}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if updatedModel.batch.Progress() != initialProgress+1 {
			t.Error("Progress should increment after success")
		}
		if updatedModel.state != FileManageState {
			t.Errorf("Expected FileManageState, got %v", updatedModel.state)
		}
		if cmd != nil {
			t.Error("Expected nil command after success")
		}
	})

	t.Run("should transition to EndState after last file success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)
		model.state = ProcessingState

		msg := SuccessMsg{}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if updatedModel.state != EndState {
			t.Errorf("Expected EndState, got %v", updatedModel.state)
		}
		if !updatedModel.batch.IsComplete() {
			t.Error("Batch should be complete after last file")
		}
		if cmd != nil {
			t.Error("Expected nil command after success")
		}
	})

	t.Run("should handle error message and transition to ErrorState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)
		model.state = ProcessingState

		testError := errors.New("test error")
		msg := ErrorMsg{Err: testError}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if updatedModel.state != ErrorState {
			t.Errorf("Expected ErrorState, got %v", updatedModel.state)
		}
		if updatedModel.errMsg != testError.Error() {
			t.Errorf("Expected error message '%s', got '%s'", testError.Error(), updatedModel.errMsg)
		}
		if cmd != nil {
			t.Error("Expected nil command after error")
		}
	})

	t.Run("should quit on Ctrl+C in ProcessingState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)
		model.state = ProcessingState

		msg := tea.KeyMsg{Type: tea.KeyCtrlC}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if cmd == nil {
			t.Error("Expected quit command")
		}
		if updatedModel.state != ProcessingState {
			t.Error("State should not change on quit command")
		}
	})
}

func TestModel_Update_EndState(t *testing.T) {
	t.Run("should quit on any key in EndState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)
		model.state = EndState

		msg := tea.KeyMsg{Type: tea.KeyEnter}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if cmd == nil {
			t.Error("Expected quit command")
		}
		if updatedModel.state != EndState {
			t.Error("State should not change on quit command")
		}
	})
}

func TestModel_Update_ErrorState(t *testing.T) {
	t.Run("should quit on any key in ErrorState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)
		model.state = ErrorState
		model.errMsg = "test error"

		msg := tea.KeyMsg{Type: tea.KeySpace}
		updatedTeaModel, cmd := model.Update(msg)
		updatedModel := updatedTeaModel.(Model)

		if cmd == nil {
			t.Error("Expected quit command")
		}
		if updatedModel.state != ErrorState {
			t.Error("State should not change on quit command")
		}
	})
}

func TestModel_Keep_Delete_Commands(t *testing.T) {
	t.Run("keep command should return SuccessMsg on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		mockManager.EXPECT().Keep("file1.txt").Return(nil)

		cmd := model.keep()
		msg := cmd()

		switch msg.(type) {
		case SuccessMsg:
			// Expected
		case ErrorMsg:
			t.Error("Expected SuccessMsg, got ErrorMsg")
		default:
			t.Errorf("Expected SuccessMsg, got %T", msg)
		}
	})

	t.Run("keep command should return ErrorMsg on failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		expectedErr := errors.New("keep failed")
		mockManager.EXPECT().Keep("file1.txt").Return(expectedErr)

		cmd := model.keep()
		msg := cmd()

		switch msg := msg.(type) {
		case ErrorMsg:
			if msg.Err != expectedErr {
				t.Errorf("Expected error %v, got %v", expectedErr, msg.Err)
			}
		case SuccessMsg:
			t.Error("Expected ErrorMsg, got SuccessMsg")
		default:
			t.Errorf("Expected ErrorMsg, got %T", msg)
		}
	})

	t.Run("delete command should return SuccessMsg on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		mockManager.EXPECT().Delete("file1.txt").Return(nil)

		cmd := model.delete()
		msg := cmd()

		switch msg.(type) {
		case SuccessMsg:
			// Expected
		case ErrorMsg:
			t.Error("Expected SuccessMsg, got ErrorMsg")
		default:
			t.Errorf("Expected SuccessMsg, got %T", msg)
		}
	})

	t.Run("delete command should return ErrorMsg on failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		expectedErr := errors.New("delete failed")
		mockManager.EXPECT().Delete("file1.txt").Return(expectedErr)

		cmd := model.delete()
		msg := cmd()

		switch msg := msg.(type) {
		case ErrorMsg:
			if msg.Err != expectedErr {
				t.Errorf("Expected error %v, got %v", expectedErr, msg.Err)
			}
		case SuccessMsg:
			t.Error("Expected ErrorMsg, got SuccessMsg")
		default:
			t.Errorf("Expected ErrorMsg, got %T", msg)
		}
	})
}

func TestModel_View(t *testing.T) {
	t.Run("should render file manage view", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt", "file2.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		view := model.View()

		if !strings.Contains(view, "File Manager") {
			t.Error("View should contain 'File Manager' title")
		}
		if !strings.Contains(view, "file1.txt") {
			t.Error("View should contain current filename")
		}
		if !strings.Contains(view, "Keep") || !strings.Contains(view, "Delete") {
			t.Error("View should contain action options")
		}
	})

	t.Run("should render processing view", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)
		model.state = ProcessingState

		view := model.View()

		if !strings.Contains(view, "Processing Files") {
			t.Error("View should contain 'Processing Files' title")
		}
		if !strings.Contains(view, "Processing...") {
			t.Error("View should contain processing indicator")
		}
	})

	t.Run("should render end view", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)
		model.state = EndState

		view := model.View()

		if !strings.Contains(view, "Processing Complete") {
			t.Error("View should contain completion message")
		}
		if !strings.Contains(view, "Processed 1 files") {
			t.Error("View should contain file count")
		}
	})

	t.Run("should render error view", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)
		model.state = ErrorState
		model.errMsg = "test error message"

		view := model.View()

		if !strings.Contains(view, "Error Occurred") {
			t.Error("View should contain error title")
		}
		if !strings.Contains(view, "test error message") {
			t.Error("View should contain error message")
		}
	})
}

func TestModel_CreateProgressBar(t *testing.T) {
	t.Run("should create progress bar for partial completion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt", "file2.txt", "file3.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		progressBar := model.createProgressBar(3, 10)

		if !strings.Contains(progressBar, "3/10") {
			t.Error("Progress bar should contain current/total count")
		}
		if !strings.Contains(progressBar, "30.0%") {
			t.Error("Progress bar should contain percentage")
		}
	})

	t.Run("should create progress bar for complete state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := mocks.NewMockFileManager(ctrl)
		batch, err := domain.NewFileBatch([]string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt"})
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}
		model := InitialModel(batch, mockManager)

		progressBar := model.createProgressBar(5, 5)

		if !strings.Contains(progressBar, "5/5") {
			t.Error("Progress bar should contain current/total count")
		}
		if !strings.Contains(progressBar, "100.0%") {
			t.Error("Progress bar should show 100% completion")
		}
	})
}
