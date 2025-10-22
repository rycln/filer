package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLocal(t *testing.T) {
	t.Run("should create local filesystem with source only", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		local, err := NewLocal(tempDir, "")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if local == nil {
			t.Error("Expected non-nil local filesystem")
		}
		if local.source != tempDir {
			t.Errorf("Expected source %s, got %s", tempDir, local.source)
		}
		if local.target != "" {
			t.Errorf("Expected empty target, got %s", local.target)
		}
	})

	t.Run("should create local filesystem with source and target", func(t *testing.T) {
		tempSource, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempSource)

		tempTarget, err := os.MkdirTemp("", "test_target")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempTarget)

		local, err := NewLocal(tempSource, tempTarget)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if local == nil {
			t.Error("Expected non-nil local filesystem")
		}
		if local.source != tempSource {
			t.Errorf("Expected source %s, got %s", tempSource, local.source)
		}
		if local.target != tempTarget {
			t.Errorf("Expected target %s, got %s", tempTarget, local.target)
		}
	})

	t.Run("should create target directory if it doesn't exist", func(t *testing.T) {
		tempSource, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempSource)

		targetDir := filepath.Join(tempSource, "new_target")

		local, err := NewLocal(tempSource, targetDir)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if local == nil {
			t.Error("Expected non-nil local filesystem")
		}

		// Check that target directory was created
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			t.Error("Target directory was not created")
		}
	})

	t.Run("should return error when target directory creation fails", func(t *testing.T) {
		tempSource, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempSource)

		// Create an invalid target path that should fail
		invalidTarget := "/root/invalid/path/that/should/fail"

		local, err := NewLocal(tempSource, invalidTarget)

		if err == nil {
			t.Error("Expected error for invalid target path")
		}
		if local != nil {
			t.Error("Expected nil local filesystem when error occurs")
		}
	})
}

func TestLocal_KeepFile(t *testing.T) {
	t.Run("should do nothing when target is empty", func(t *testing.T) {
		tempSource, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempSource)

		local, err := NewLocal(tempSource, "")
		if err != nil {
			t.Fatalf("Failed to create local filesystem: %v", err)
		}

		err = local.KeepFile("somefile.txt")
		if err != nil {
			t.Errorf("Expected no error with empty target, got %v", err)
		}
	})

	t.Run("should move file to target directory", func(t *testing.T) {
		tempSource, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempSource)

		tempTarget, err := os.MkdirTemp("", "test_target")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempTarget)

		// Create test file in source
		testFile := "testfile.txt"
		sourceFilePath := filepath.Join(tempSource, testFile)
		err = os.WriteFile(sourceFilePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		local, err := NewLocal(tempSource, tempTarget)
		if err != nil {
			t.Fatalf("Failed to create local filesystem: %v", err)
		}

		err = local.KeepFile(testFile)
		if err != nil {
			t.Errorf("Failed to keep file: %v", err)
		}

		// Check file was moved to target
		targetFilePath := filepath.Join(tempTarget, testFile)
		if _, err := os.Stat(targetFilePath); os.IsNotExist(err) {
			t.Error("File was not moved to target directory")
		}

		// Check file was removed from source
		if _, err := os.Stat(sourceFilePath); !os.IsNotExist(err) {
			t.Error("File was not removed from source directory")
		}
	})

	t.Run("should return error when source file doesn't exist", func(t *testing.T) {
		tempSource, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempSource)

		tempTarget, err := os.MkdirTemp("", "test_target")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempTarget)

		local, err := NewLocal(tempSource, tempTarget)
		if err != nil {
			t.Fatalf("Failed to create local filesystem: %v", err)
		}

		err = local.KeepFile("nonexistent.txt")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})
}

func TestLocal_DeleteFile(t *testing.T) {
	t.Run("should delete existing file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create test file
		testFile := "file_to_delete.txt"
		filePath := filepath.Join(tempDir, testFile)
		err = os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		local, err := NewLocal(tempDir, "")
		if err != nil {
			t.Fatalf("Failed to create local filesystem: %v", err)
		}

		err = local.DeleteFile(testFile)
		if err != nil {
			t.Errorf("Failed to delete file: %v", err)
		}

		// Check file was deleted
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Error("File was not deleted")
		}
	})

	t.Run("should return error when file to delete doesn't exist", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		local, err := NewLocal(tempDir, "")
		if err != nil {
			t.Fatalf("Failed to create local filesystem: %v", err)
		}

		err = local.DeleteFile("nonexistent.txt")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})
}

func TestLocal_GetFilenames(t *testing.T) {
	t.Run("should return empty list for empty directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		local, err := NewLocal(tempDir, "")
		if err != nil {
			t.Fatalf("Failed to create local filesystem: %v", err)
		}

		filenames, err := local.GetFilenames()
		if err != nil {
			t.Errorf("Failed to get filenames: %v", err)
		}
		if len(filenames) != 0 {
			t.Errorf("Expected empty filenames list, got %v", filenames)
		}
	})

	t.Run("should return only files, not directories", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create files
		files := []string{"file1.txt", "file2.txt", "file3.go"}
		for _, file := range files {
			filePath := filepath.Join(tempDir, file)
			err = os.WriteFile(filePath, []byte("content"), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}

		// Create a subdirectory
		subDir := filepath.Join(tempDir, "subdir")
		err = os.MkdirAll(subDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create subdirectory: %v", err)
		}

		local, err := NewLocal(tempDir, "")
		if err != nil {
			t.Fatalf("Failed to create local filesystem: %v", err)
		}

		filenames, err := local.GetFilenames()
		if err != nil {
			t.Errorf("Failed to get filenames: %v", err)
		}

		if len(filenames) != len(files) {
			t.Errorf("Expected %d files, got %d", len(files), len(filenames))
		}

		// Check that all expected files are present
		for _, expectedFile := range files {
			found := false
			for _, actualFile := range filenames {
				if actualFile == expectedFile {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected file %s not found in result", expectedFile)
			}
		}
	})

	t.Run("should return error for non-existent source directory", func(t *testing.T) {
		local, err := NewLocal("/nonexistent/path/12345", "")
		if err != nil {
			t.Fatalf("Failed to create local filesystem: %v", err)
		}

		_, err = local.GetFilenames()
		if err == nil {
			t.Error("Expected error for non-existent source directory")
		}
	})
}

func Test_moveFileSafe(t *testing.T) {
	t.Run("should return error when source file doesn't exist", func(t *testing.T) {
		err := moveFileSafe("/nonexistent/source.txt", "/some/target.txt")
		if err == nil {
			t.Error("Expected error for non-existent source file")
		}
	})

	t.Run("should successfully rename file within same filesystem", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_move")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		sourcePath := filepath.Join(tempDir, "source.txt")
		targetPath := filepath.Join(tempDir, "target.txt")

		err = os.WriteFile(sourcePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err = moveFileSafe(sourcePath, targetPath)
		if err != nil {
			t.Errorf("Failed to move file: %v", err)
		}

		// Check file was moved
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			t.Error("File was not moved to target path")
		}
		if _, err := os.Stat(sourcePath); !os.IsNotExist(err) {
			t.Error("File still exists at source path")
		}
	})
}

func Test_copyAndRemove(t *testing.T) {
	t.Run("should copy file and remove original", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_copy")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		sourcePath := filepath.Join(tempDir, "source.txt")
		targetPath := filepath.Join(tempDir, "target.txt")

		content := []byte("test content for copy")
		err = os.WriteFile(sourcePath, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err = copyAndRemove(sourcePath, targetPath)
		if err != nil {
			t.Errorf("Failed to copy and remove file: %v", err)
		}

		// Check file was copied
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			t.Error("File was not copied to target path")
		}

		// Check original was removed
		if _, err := os.Stat(sourcePath); !os.IsNotExist(err) {
			t.Error("Original file was not removed")
		}

		// Check content is correct
		copiedContent, err := os.ReadFile(targetPath)
		if err != nil {
			t.Fatalf("Failed to read copied file: %v", err)
		}
		if string(copiedContent) != string(content) {
			t.Error("Copied file content doesn't match original")
		}
	})

	t.Run("should return error when source file cannot be opened", func(t *testing.T) {
		err := copyAndRemove("/nonexistent/source.txt", "/some/target.txt")
		if err == nil {
			t.Error("Expected error for non-existent source file")
		}
	})

	t.Run("should return error when target file cannot be created", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_copy")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		sourcePath := filepath.Join(tempDir, "source.txt")
		err = os.WriteFile(sourcePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Use invalid target path
		invalidTarget := "/root/invalid/target.txt"

		err = copyAndRemove(sourcePath, invalidTarget)
		if err == nil {
			t.Error("Expected error for invalid target path")
		}
	})
}
