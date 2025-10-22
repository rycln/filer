package domain

import (
	"testing"
)

func TestNewFileBatch(t *testing.T) {
	t.Run("should return error for empty files", func(t *testing.T) {
		files := []string{}
		batch, err := NewFileBatch(files)

		if err == nil {
			t.Error("Expected error for empty files")
		}
		if batch != nil {
			t.Error("Expected nil batch when error occurs")
		}
		expectedErr := "no files to process"
		if err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got '%v'", expectedErr, err)
		}
	})

	t.Run("should create batch with single file", func(t *testing.T) {
		files := []string{"file1.txt"}
		batch, err := NewFileBatch(files)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if batch == nil {
			t.Error("Expected non-nil batch")
		}
		if batch.TotalFiles() != 1 {
			t.Errorf("Expected total files 1, got %d", batch.TotalFiles())
		}
		if batch.IsComplete() {
			t.Error("Expected batch to not be complete")
		}
	})

	t.Run("should create batch with multiple files", func(t *testing.T) {
		files := []string{"file1.txt", "file2.txt", "file3.txt"}
		batch, err := NewFileBatch(files)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if batch == nil {
			t.Error("Expected non-nil batch")
		}
		if batch.TotalFiles() != 3 {
			t.Errorf("Expected total files 3, got %d", batch.TotalFiles())
		}
		if batch.IsComplete() {
			t.Error("Expected batch to not be complete")
		}
	})
}

func TestFileBatch_CurrentFile(t *testing.T) {
	t.Run("should return first file initially", func(t *testing.T) {
		files := []string{"file1.txt", "file2.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		current := batch.CurrentFile()

		if current != "file1.txt" {
			t.Errorf("Expected 'file1.txt', got '%s'", current)
		}
	})

	t.Run("should return empty string after completion", func(t *testing.T) {
		files := []string{"file1.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		batch.NextFile()

		current := batch.CurrentFile()
		if current != "" {
			t.Errorf("Expected empty string, got '%s'", current)
		}
	})

	t.Run("should return correct file after multiple next operations", func(t *testing.T) {
		files := []string{"file1.txt", "file2.txt", "file3.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		batch.NextFile()
		current := batch.CurrentFile()
		if current != "file2.txt" {
			t.Errorf("Expected 'file2.txt', got '%s'", current)
		}

		batch.NextFile()
		current = batch.CurrentFile()
		if current != "file3.txt" {
			t.Errorf("Expected 'file3.txt', got '%s'", current)
		}
	})
}

func TestFileBatch_NextFile(t *testing.T) {
	t.Run("should move to next file in sequence", func(t *testing.T) {
		files := []string{"file1.txt", "file2.txt", "file3.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		if batch.CurrentFile() != "file1.txt" {
			t.Error("Should start with first file")
		}

		batch.NextFile()
		if batch.CurrentFile() != "file2.txt" {
			t.Error("Should move to second file")
		}

		batch.NextFile()
		if batch.CurrentFile() != "file3.txt" {
			t.Error("Should move to third file")
		}
	})

	t.Run("should handle multiple next operations beyond bounds", func(t *testing.T) {
		files := []string{"file1.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		batch.NextFile()
		batch.NextFile()
		batch.NextFile()

		if !batch.IsComplete() {
			t.Error("Expected batch to be complete")
		}
		if batch.CurrentFile() != "" {
			t.Errorf("Expected empty current file, got '%s'", batch.CurrentFile())
		}
	})
}

func TestFileBatch_IsComplete(t *testing.T) {
	t.Run("should not be complete for non-empty batch initially", func(t *testing.T) {
		files := []string{"file1.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		if batch.IsComplete() {
			t.Error("Expected non-empty batch to not be complete initially")
		}
	})

	t.Run("should become complete after processing all files", func(t *testing.T) {
		files := []string{"file1.txt", "file2.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		if batch.IsComplete() {
			t.Error("Should not be complete initially")
		}

		batch.NextFile()
		if batch.IsComplete() {
			t.Error("Should not be complete after first next")
		}

		batch.NextFile()
		if !batch.IsComplete() {
			t.Error("Should be complete after processing all files")
		}
	})

	t.Run("should remain complete after additional next operations", func(t *testing.T) {
		files := []string{"file1.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		batch.NextFile()
		batch.NextFile()

		if !batch.IsComplete() {
			t.Error("Should remain complete after additional next operations")
		}
	})
}

func TestFileBatch_Progress(t *testing.T) {
	t.Run("should start at zero progress", func(t *testing.T) {
		files := []string{"file1.txt", "file2.txt", "file3.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		if batch.Progress() != 0 {
			t.Errorf("Expected progress 0, got %d", batch.Progress())
		}
	})

	t.Run("should increment progress with next operations", func(t *testing.T) {
		files := []string{"file1.txt", "file2.txt", "file3.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		batch.NextFile()
		if batch.Progress() != 1 {
			t.Errorf("Expected progress 1, got %d", batch.Progress())
		}

		batch.NextFile()
		if batch.Progress() != 2 {
			t.Errorf("Expected progress 2, got %d", batch.Progress())
		}
	})

	t.Run("should handle progress beyond file count", func(t *testing.T) {
		files := []string{"file1.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		batch.NextFile()
		batch.NextFile()
		batch.NextFile()

		if batch.Progress() != 3 {
			t.Errorf("Expected progress 3, got %d", batch.Progress())
		}
	})
}

func TestFileBatch_TotalFiles(t *testing.T) {
	t.Run("should return correct count for single file", func(t *testing.T) {
		files := []string{"file1.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		if batch.TotalFiles() != 1 {
			t.Errorf("Expected 1 total file, got %d", batch.TotalFiles())
		}
	})

	t.Run("should return correct count for multiple files", func(t *testing.T) {
		files := []string{"a.txt", "b.txt", "c.txt", "d.txt", "e.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		if batch.TotalFiles() != 5 {
			t.Errorf("Expected 5 total files, got %d", batch.TotalFiles())
		}
	})

	t.Run("should maintain consistent total files count after operations", func(t *testing.T) {
		files := []string{"file1.txt", "file2.txt"}
		batch, err := NewFileBatch(files)
		if err != nil {
			t.Fatalf("Failed to create batch: %v", err)
		}

		initialTotal := batch.TotalFiles()

		batch.NextFile()
		batch.NextFile()
		batch.NextFile()

		if batch.TotalFiles() != initialTotal {
			t.Errorf("Total files should remain %d, got %d", initialTotal, batch.TotalFiles())
		}
	})
}
