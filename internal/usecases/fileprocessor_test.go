package usecases

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/filer/internal/usecases/mocks"
)

func TestNewFileProcessor(t *testing.T) {
	t.Run("should create new file processor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mocks.NewMockFileSystem(ctrl)
		processor := NewFileProcessor(mockFS)

		if processor == nil {
			t.Error("Expected non-nil processor")
		}
	})
}

func TestFileProcessor_Keep(t *testing.T) {
	t.Run("should successfully keep file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mocks.NewMockFileSystem(ctrl)
		processor := NewFileProcessor(mockFS)
		filename := "test.txt"

		mockFS.EXPECT().KeepFile(filename).Return(nil)

		err := processor.Keep(filename)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("should return error when filesystem fails to keep file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mocks.NewMockFileSystem(ctrl)
		processor := NewFileProcessor(mockFS)
		filename := "test.txt"
		expectedErr := errors.New("keep failed")

		mockFS.EXPECT().KeepFile(filename).Return(expectedErr)

		err := processor.Keep(filename)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should handle empty filename", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mocks.NewMockFileSystem(ctrl)
		processor := NewFileProcessor(mockFS)
		filename := ""

		mockFS.EXPECT().KeepFile(filename).Return(nil)

		err := processor.Keep(filename)

		if err != nil {
			t.Errorf("Expected no error with empty filename, got %v", err)
		}
	})
}

func TestFileProcessor_Delete(t *testing.T) {
	t.Run("should successfully delete file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mocks.NewMockFileSystem(ctrl)
		processor := NewFileProcessor(mockFS)
		filename := "test.txt"

		mockFS.EXPECT().DeleteFile(filename).Return(nil)

		err := processor.Delete(filename)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("should return error when filesystem fails to delete file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mocks.NewMockFileSystem(ctrl)
		processor := NewFileProcessor(mockFS)
		filename := "test.txt"
		expectedErr := errors.New("delete failed")

		mockFS.EXPECT().DeleteFile(filename).Return(expectedErr)

		err := processor.Delete(filename)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should handle empty filename for delete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mocks.NewMockFileSystem(ctrl)
		processor := NewFileProcessor(mockFS)
		filename := ""

		mockFS.EXPECT().DeleteFile(filename).Return(nil)

		err := processor.Delete(filename)

		if err != nil {
			t.Errorf("Expected no error with empty filename, got %v", err)
		}
	})
}

func TestFileProcessor_Integration(t *testing.T) {
	t.Run("should call correct filesystem method for each operation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mocks.NewMockFileSystem(ctrl)
		processor := NewFileProcessor(mockFS)
		filename := "document.pdf"

		// Expect keep call
		gomock.InOrder(
			mockFS.EXPECT().KeepFile(filename).Return(nil),
			mockFS.EXPECT().DeleteFile(filename).Return(nil),
		)

		err := processor.Keep(filename)
		if err != nil {
			t.Errorf("Keep failed: %v", err)
		}

		err = processor.Delete(filename)
		if err != nil {
			t.Errorf("Delete failed: %v", err)
		}
	})

	t.Run("should handle different filenames independently", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mocks.NewMockFileSystem(ctrl)
		processor := NewFileProcessor(mockFS)
		filename1 := "file1.txt"
		filename2 := "file2.txt"

		mockFS.EXPECT().KeepFile(filename1).Return(nil)
		mockFS.EXPECT().DeleteFile(filename2).Return(nil)

		err := processor.Keep(filename1)
		if err != nil {
			t.Errorf("Keep file1 failed: %v", err)
		}

		err = processor.Delete(filename2)
		if err != nil {
			t.Errorf("Delete file2 failed: %v", err)
		}
	})
}
