package domain

import "fmt"

// FileBatch manages sequential file processing with progress tracking.
// Tracks current position and completion state for batch operations.
type FileBatch struct {
	filenames []string
	idx       int
}

// NewFileBatch creates a file batch for sequential processing.
// Returns error if files slice is empty.
func NewFileBatch(files []string) (*FileBatch, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files to process")
	}

	return &FileBatch{
		filenames: files,
		idx:       0,
	}, nil
}

// CurrentFile returns the current filename being processed.
// Returns empty string when batch is complete.
func (b *FileBatch) CurrentFile() string {
	if b.idx >= len(b.filenames) {
		return ""
	}
	return b.filenames[b.idx]
}

// NextFile advances to the next file in the batch.
// Increments internal index counter.
func (b *FileBatch) NextFile() {
	b.idx++
}

// IsComplete checks if all files have been processed.
// Returns true when index reaches end of file list.
func (b *FileBatch) IsComplete() bool {
	return b.idx >= len(b.filenames)
}

// Progress returns current processing position.
// Zero-based index of current file.
func (b *FileBatch) Progress() int {
	return b.idx
}

// TotalFiles returns the total number of files in batch.
// Constant value representing initial file count.
func (b *FileBatch) TotalFiles() int {
	return len(b.filenames)
}
