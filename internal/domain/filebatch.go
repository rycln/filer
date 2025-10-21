package domain

type FileBatch struct {
	files []string
	idx   int
}

func NewFileBatch(files []string) *FileBatch {
	return &FileBatch{
		files: files,
		idx:   0,
	}
}

func (b *FileBatch) CurrentFile() string {
	if b.idx >= len(b.files) {
		return ""
	}
	return b.files[b.idx]
}

func (b *FileBatch) NextFile() {
	b.idx++
}

func (b *FileBatch) IsComplete() bool {
	return b.idx >= len(b.files)
}

func (b *FileBatch) Progress() int {
	return b.idx + 1
}

func (b *FileBatch) TotalFiles() int {
	return len(b.files)
}
