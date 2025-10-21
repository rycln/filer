package domain

type FileBatch struct {
	filenames []string
	idx       int
}

func NewFileBatch(files []string) *FileBatch {
	return &FileBatch{
		filenames: files,
		idx:       0,
	}
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
	return b.idx + 1
}

func (b *FileBatch) TotalFiles() int {
	return len(b.filenames)
}
