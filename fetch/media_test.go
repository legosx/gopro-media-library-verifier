package fetch_test

import (
	"github.com/legosx/gopro-media-library-verifier/fetch"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMedia(t *testing.T) {
	type fields struct {
		fileName string
		fileSize int64
	}

	type want struct {
		media    fetch.Media
		fileName string
		fileSize int64
	}

	tests := []struct {
		name string
		fields
		want
	}{
		{
			name: "happy path",
			fields: fields{
				fileName: "file1.mp4",
				fileSize: 10,
			},
			want: want{
				media:    fetch.NewMedia("file1.mp4", 10),
				fileName: "file1.mp4",
				fileSize: 10,
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := fetch.NewMedia(tt.fields.fileName, tt.fields.fileSize)
			assert.Equal(t, tt.want.media, got)
			assert.Equal(t, tt.want.fileName, got.FileName())
			assert.Equal(t, tt.want.fileSize, got.FileSize())
		})
	}
}
