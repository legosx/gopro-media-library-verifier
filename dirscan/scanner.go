package dirscan

import (
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Scanner struct {
	allowedExtensions []string
	os                OS
}

func NewScanner(allowedExtensions []string, opts ...func(scanner *Scanner)) (scanner Scanner) {
	scanner = Scanner{
		allowedExtensions: allowedExtensions,
		os:                NewOSWrapper(),
	}

	for _, opt := range opts {
		opt(&scanner)
	}

	return scanner
}

type OS interface {
	Stat(name string) (os.FileInfo, error)
	Open(name string) (OSFile, error)
	IsNotExist(err error) bool
}

func WithOS(os OS) func(s *Scanner) {
	return func(s *Scanner) {
		s.os = os
	}
}

type File struct {
	Name string
	Size int64
	Path string
}

func (s Scanner) GetFileList(dirPath string) (list []File, err error) {
	if _, err := s.os.Stat(dirPath); err != nil {
		if s.os.IsNotExist(err) {
			return []File{}, errors.Wrap(err, "path does not exist")
		}

		return []File{}, errors.Wrap(err, "can't stat the path")
	}

	// Open the directory
	dir, err := s.os.Open(dirPath)
	if err != nil {
		return []File{}, errors.Wrap(err, "error opening the directory")
	}
	defer func() {
		if innerErr := dir.Close(); innerErr != nil {
			err = errors.Wrap(multierr.Append(err, innerErr), "cannot close directory")
		}
	}()

	// Read the directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return []File{}, errors.Wrap(err, "error reading the directory")
	}

	list = make([]File, 0)

	// Iterate through the files
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			innerList, err := s.GetFileList(path.Join(dirPath, fileInfo.Name()))
			if err != nil {
				return []File{}, errors.Wrap(err, "error getting file list recursively")
			}

			list = append(list, innerList...)
		}

		if !fileInfo.Mode().IsRegular() {
			continue
		}

		fileName := strings.ToLower(fileInfo.Name())
		if !s.isFileNameAllowed(fileName) {
			continue
		}

		list = append(list, File{
			Name: fileInfo.Name(),
			Size: fileInfo.Size(),
			Path: path.Join(dirPath, fileInfo.Name()),
		})
	}

	return list, nil
}

func (s Scanner) isFileNameAllowed(filename string) bool {
	extension := strings.ToLower(filepath.Ext(filename))

	for _, v := range s.allowedExtensions {
		if v == strings.ToLower(extension) {
			return true
		}
	}

	return false
}
