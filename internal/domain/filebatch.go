package domain

import "fmt"

type FileBatch struct {
	filenames []string
	idx       int
}

func NewFileBatch(files []string) (*FileBatch, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files to process")
	}

	return &FileBatch{
		filenames: files,
		idx:       0,
	}, nil
}

func (b *FileBatch) CurrentFile() string {
	if b.idx >= len(b.filenames) {
		return ""
	}
	return b.filenames[b.idx]
}

func (b *FileBatch) NextFile() {
	b.idx++
}

func (b *FileBatch) IsComplete() bool {
	return b.idx >= len(b.filenames)
}

func (b *FileBatch) Progress() int {
	return b.idx
}

func (b *FileBatch) TotalFiles() int {
	return len(b.filenames)
}
