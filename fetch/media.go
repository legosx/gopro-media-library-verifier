package fetch

type Media struct {
	fileName string
	fileSize int64
}

func NewMedia(fileName string, fileSize int64) Media {
	return Media{
		fileName: fileName,
		fileSize: fileSize,
	}
}

func (m Media) FileName() string {
	return m.fileName
}

func (m Media) FileSize() int64 {
	return m.fileSize
}
