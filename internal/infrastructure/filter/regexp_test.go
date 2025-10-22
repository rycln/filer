package filter

import (
	"regexp"
	"testing"
)

func TestNewRegexpFilter(t *testing.T) {
	t.Run("should create filter with empty pattern", func(t *testing.T) {
		filter := NewRegexpFilter("")

		if filter == nil {
			t.Error("Expected non-nil filter")
		}
		if filter.pattern != "" {
			t.Errorf("Expected empty pattern, got %s", filter.pattern)
		}
	})

	t.Run("should create filter with valid pattern", func(t *testing.T) {
		pattern := ".*\\.txt"
		filter := NewRegexpFilter(pattern)

		if filter == nil {
			t.Error("Expected non-nil filter")
		}
		if filter.pattern != pattern {
			t.Errorf("Expected pattern %s, got %s", pattern, filter.pattern)
		}
	})

	t.Run("should create filter with complex pattern", func(t *testing.T) {
		pattern := "^test_[0-9]+\\.go$"
		filter := NewRegexpFilter(pattern)

		if filter == nil {
			t.Error("Expected non-nil filter")
		}
		if filter.pattern != pattern {
			t.Errorf("Expected pattern %s, got %s", pattern, filter.pattern)
		}
	})
}

func TestRegexpFilter_Filter(t *testing.T) {
	t.Run("should return all filenames when pattern is empty", func(t *testing.T) {
		filter := NewRegexpFilter("")
		filenames := []string{"file1.txt", "file2.jpg", "file3.go"}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != len(filenames) {
			t.Errorf("Expected %d files, got %d", len(filenames), len(result))
		}
		for i, filename := range filenames {
			if result[i] != filename {
				t.Errorf("Expected filename %s, got %s", filename, result[i])
			}
		}
	})

	t.Run("should return empty slice when no filenames match", func(t *testing.T) {
		filter := NewRegexpFilter("\\.md$")
		filenames := []string{"file1.txt", "file2.jpg", "file3.go"}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("should filter txt files only", func(t *testing.T) {
		filter := NewRegexpFilter("\\.txt$")
		filenames := []string{"file1.txt", "file2.jpg", "file3.go", "doc.txt"}
		expected := []string{"doc.txt", "file1.txt"}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != len(expected) {
			t.Errorf("Expected %d files, got %d", len(expected), len(result))
		}
		for i, filename := range expected {
			if result[i] != filename {
				t.Errorf("Expected filename %s, got %s", filename, result[i])
			}
		}
	})

	t.Run("should filter files with prefix", func(t *testing.T) {
		filter := NewRegexpFilter("^test_")
		filenames := []string{"test_file.go", "production.go", "test_data.txt", "main.go"}
		expected := []string{"test_data.txt", "test_file.go"}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != len(expected) {
			t.Errorf("Expected %d files, got %d", len(expected), len(result))
		}
		for i, filename := range expected {
			if result[i] != filename {
				t.Errorf("Expected filename %s, got %s", filename, result[i])
			}
		}
	})

	t.Run("should filter files with numbers", func(t *testing.T) {
		filter := NewRegexpFilter("^[a-z]+_[0-9]+\\.txt$")
		filenames := []string{"file_1.txt", "file_123.txt", "test.go", "data_45.txt", "invalid_name.txt"}
		expected := []string{"data_45.txt", "file_1.txt", "file_123.txt"}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != len(expected) {
			t.Errorf("Expected %d files, got %d", len(expected), len(result))
		}
		for i, filename := range expected {
			if result[i] != filename {
				t.Errorf("Expected filename %s, got %s", filename, result[i])
			}
		}
	})

	t.Run("should return sorted results", func(t *testing.T) {
		filter := NewRegexpFilter("\\.go$")
		filenames := []string{"z.go", "a.go", "m.go", "b.go"}
		expected := []string{"a.go", "b.go", "m.go", "z.go"}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != len(expected) {
			t.Errorf("Expected %d files, got %d", len(expected), len(result))
		}
		for i, filename := range expected {
			if result[i] != filename {
				t.Errorf("Expected filename %s, got %s", filename, result[i])
			}
		}
	})

	t.Run("should handle empty input slice", func(t *testing.T) {
		filter := NewRegexpFilter("\\.txt$")
		filenames := []string{}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("should return error for invalid regex pattern", func(t *testing.T) {
		filter := NewRegexpFilter("[invalid")
		filenames := []string{"file1.txt"}

		result, err := filter.Filter(filenames)

		if err == nil {
			t.Error("Expected error for invalid regex pattern")
		}
		if result != nil {
			t.Errorf("Expected nil result when error occurs, got %v", result)
		}

		_, compileErr := regexp.Compile("[invalid")
		if err.Error() != compileErr.Error() {
			t.Errorf("Expected regexp compilation error, got %v", err)
		}
	})

	t.Run("should filter case sensitive by default", func(t *testing.T) {
		filter := NewRegexpFilter("^[a-z]+\\.txt$")
		filenames := []string{"file.txt", "File.txt", "FILE.TXT"}
		expected := []string{"file.txt"}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != len(expected) {
			t.Errorf("Expected %d files, got %d", len(expected), len(result))
		}
		for i, filename := range expected {
			if result[i] != filename {
				t.Errorf("Expected filename %s, got %s", filename, result[i])
			}
		}
	})

	t.Run("should filter files with special characters", func(t *testing.T) {
		filter := NewRegexpFilter("^file\\-with\\-dash_[0-9]+\\.txt$")
		filenames := []string{"file-with-dash_1.txt", "file_with_dash_1.txt", "file-with-dash_123.txt"}
		expected := []string{"file-with-dash_1.txt", "file-with-dash_123.txt"}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != len(expected) {
			t.Errorf("Expected %d files, got %d", len(expected), len(result))
		}
		for i, filename := range expected {
			if result[i] != filename {
				t.Errorf("Expected filename %s, got %s", filename, result[i])
			}
		}
	})
}

func TestRegexpFilter_EdgeCases(t *testing.T) {
	t.Run("should handle dot files", func(t *testing.T) {
		filter := NewRegexpFilter("^\\.")
		filenames := []string{".gitignore", ".env", "normal.txt", ".hidden"}
		expected := []string{".env", ".gitignore", ".hidden"}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != len(expected) {
			t.Errorf("Expected %d files, got %d", len(expected), len(result))
		}
		for i, filename := range expected {
			if result[i] != filename {
				t.Errorf("Expected filename %s, got %s", filename, result[i])
			}
		}
	})

	t.Run("should handle files with multiple extensions", func(t *testing.T) {
		filter := NewRegexpFilter("\\.tar\\.gz$")
		filenames := []string{"archive.tar.gz", "backup.tar.gz", "file.gz", "archive.tar"}
		expected := []string{"archive.tar.gz", "backup.tar.gz"}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != len(expected) {
			t.Errorf("Expected %d files, got %d", len(expected), len(result))
		}
		for i, filename := range expected {
			if result[i] != filename {
				t.Errorf("Expected filename %s, got %s", filename, result[i])
			}
		}
	})

	t.Run("should maintain stability with duplicate filenames", func(t *testing.T) {
		filter := NewRegexpFilter("\\.txt$")
		filenames := []string{"file.txt", "file.txt", "image.jpg", "file.txt"}
		expected := []string{"file.txt", "file.txt", "file.txt"}

		result, err := filter.Filter(filenames)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != len(expected) {
			t.Errorf("Expected %d files, got %d", len(expected), len(result))
		}
		for i, filename := range expected {
			if result[i] != filename {
				t.Errorf("Expected filename %s, got %s", filename, result[i])
			}
		}
	})
}
