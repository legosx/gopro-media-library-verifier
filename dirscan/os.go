package dirscan

import "os"

type OSWrapper struct {
}

func NewOSWrapper() OSWrapper {
	return OSWrapper{}
}

func (o OSWrapper) Stat(name string) (fileInfo os.FileInfo, err error) {
	return os.Stat(name)
}

func (o OSWrapper) Open(name string) (OSFile, error) {
	osFile, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return NewFileWrapper(osFile), nil
}

func (o OSWrapper) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
