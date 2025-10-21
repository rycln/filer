package usecases

type FileSystem interface {
	KeepFile(string) error
	DeleteFile(string) error
}

type FileProcessor struct {
	fs FileSystem
}

func NewFileProcessor(fs FileSystem) *FileProcessor {
	return &FileProcessor{
		fs: fs,
	}
}

func (p *FileProcessor) Keep(filename string) error {
	return p.fs.KeepFile(filename)
}

func (p *FileProcessor) Delete(filename string) error {
	return p.fs.DeleteFile(filename)
}
