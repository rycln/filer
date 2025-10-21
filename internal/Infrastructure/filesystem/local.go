package filesystem

import (
	"fmt"
	"io"
	"os"
)

type Local struct {
	source string
	target string
}

func NewLocal(source, target string) (*Local, error) {
	if target != "" {
		err := os.MkdirAll(target, 0755)
		if err != nil {
			return nil, err
		}
	}

	return &Local{
		source: source,
		target: target,
	}, nil
}

func (l *Local) KeepFile(filename string) error {
	if l.target == "" {
		return nil
	}

	err := moveFileSafe(l.source+"/"+filename, l.target+"/"+filename)
	if err != nil {
		return err
	}

	return nil
}

func moveFileSafe(sourcePath, destPath string) error {
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", sourcePath)
	}

	err := os.Rename(sourcePath, destPath)
	if err == nil {
		return nil
	}

	return copyAndRemove(sourcePath, destPath)
}

func copyAndRemove(sourcePath, destPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		os.Remove(destPath)
		return err
	}

	sourceInfo, _ := sourceFile.Stat()
	destInfo, _ := destFile.Stat()

	if sourceInfo.Size() != destInfo.Size() {
		os.Remove(destPath)
		return fmt.Errorf("the file sizes do not match")
	}

	return os.Remove(sourcePath)
}

func (l *Local) DeleteFile(filename string) error {
	err := os.Remove(l.source + "/" + filename)
	if err != nil {
		return err
	}

	return nil
}
