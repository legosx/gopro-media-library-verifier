package verify

import (
	"fmt"
	"github.com/legosx/gopro-media-library-verifier/dirscan"
	"github.com/legosx/gopro-media-library-verifier/fetch"
	"github.com/pkg/errors"
)

type Fetcher interface {
	GetMedias() (medias []fetch.Media, err error)
}

type Scanner interface {
	GetFileList(dirPath string) (list []dirscan.File, err error)
}

type Verifier struct {
	fetcher Fetcher
	scanner Scanner
}

func NewVerifier(mediaFetcher Fetcher, scanner Scanner) Verifier {
	return Verifier{fetcher: mediaFetcher, scanner: scanner}
}

func (v Verifier) IdentifyMissingFiles(path string) (filePaths []string, err error) {
	fmt.Printf("\nIdentifying files that are not yet uploaded to cloud from\n%s\nbased on: fileName, fileSize\n", path)

	localFiles, err := v.scanner.GetFileList(path)
	if err != nil {
		return []string{}, errors.Wrap(err, "error getting local files")
	}

	remoteFiles, err := v.getRemoteFiles()
	if err != nil {
		return []string{}, errors.Wrap(err, "error getting remote files")
	}

	return v.getFilePathsOfMissingFiles(localFiles, remoteFiles), nil
}

func (v Verifier) getRemoteFiles() (remoteFiles []dirscan.File, err error) {
	remoteMedias, err := v.fetcher.GetMedias()
	if err != nil {
		return []dirscan.File{}, errors.Wrap(err, "error getting remote medias")
	}

	return v.convertMediasToFiles(remoteMedias), nil
}

func (v Verifier) getFilePathsOfMissingFiles(localFiles, remoteFiles []dirscan.File) (filePaths []string) {
	filePaths = []string{}

	for _, localFile := range localFiles {
		if v.fileExists(localFile, remoteFiles) {
			continue
		}

		filePaths = append(filePaths, localFile.Path)
	}

	return filePaths
}

func (v Verifier) fileExists(lookupFile dirscan.File, files []dirscan.File) (exists bool) {
	for _, file := range files {
		if lookupFile.Name == file.Name && lookupFile.Size == file.Size {
			return true
		}
	}

	return false
}

func (v Verifier) convertMediasToFiles(medias []fetch.Media) (files []dirscan.File) {
	files = []dirscan.File{}

	for _, m := range medias {
		files = append(files, dirscan.File{
			Name: m.FileName(),
			Size: m.FileSize(),
		})
	}

	return files
}
