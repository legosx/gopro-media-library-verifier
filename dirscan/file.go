package dirscan

import "os"

type OSFile interface {
	Readdir(n int) ([]os.FileInfo, error)
	Close() error
}

type FileWrapper struct {
	file *os.File
}

func NewFileWrapper(file *os.File) FileWrapper {
	return FileWrapper{file: file}
}

func (f FileWrapper) Readdir(n int) ([]os.FileInfo, error) {
	return f.file.Readdir(n)
}

func (f FileWrapper) Close() error {
	return f.file.Close()
}
