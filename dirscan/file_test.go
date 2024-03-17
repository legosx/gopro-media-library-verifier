package dirscan_test

import (
	"github.com/legosx/gopro-media-library-verifier/dirscan"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"os"
	"testing"
)

func TestFileWrapper_Readdir(t *testing.T) {
	type fields struct {
		file func(mockCtrl *gomock.Controller) *os.File
	}

	type args struct {
		n int
	}

	type want struct {
		name string
		err  error
	}

	tests := []struct {
		name string
		fields
		args
		want
	}{
		{
			name: "happy path",
			args: args{
				n: -1,
			},
			fields: fields{
				file: func(mockCtrl *gomock.Controller) *os.File {
					file, err := os.Open("../cmd")
					assert.NoError(t, err)

					return file
				},
			},
			want: want{
				name: "root.go",
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			got, err := dirscan.NewFileWrapper(tt.fields.file(mockCtrl)).Readdir(tt.args.n)
			if err == nil {
				assert.True(t, isNameExistInFileInfoList(got, "root.go"))
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func TestFileWrapper_Close(t *testing.T) {
	type fields struct {
		file func(mockCtrl *gomock.Controller) *os.File
	}

	tests := []struct {
		name string
		fields
	}{
		{
			name: "happy path",
			fields: fields{
				file: func(mockCtrl *gomock.Controller) *os.File {
					file, err := os.Open("../go.mod")
					assert.NoError(t, err)

					return file
				},
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			got := dirscan.NewFileWrapper(tt.fields.file(mockCtrl)).Close()
			assert.NoError(t, got)
		})
	}
}

func isNameExistInFileInfoList(fileInfoList []os.FileInfo, name string) bool {
	for _, fileInfo := range fileInfoList {
		if fileInfo.Name() == name {
			return true
		}
	}

	return false
}
