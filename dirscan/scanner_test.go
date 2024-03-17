package dirscan_test

import (
	"github.com/legosx/gopro-media-library-verifier/dirscan"
	"github.com/legosx/gopro-media-library-verifier/dirscan/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"os"
	"sort"
	"testing"
	"time"
)

//go:generate mockgen -destination=./mocks/os.go -package=mocks github.com/legosx/gopro-media-library-verifier/dirscan OS
//go:generate mockgen -destination=./mocks/osfile.go -package=mocks github.com/legosx/gopro-media-library-verifier/dirscan OSFile

func TestScanner_GetFileList(t *testing.T) {
	type fields struct {
		allowedExtensions []string
		os                func(mockCtrl *gomock.Controller) dirscan.OS
	}

	type args struct {
		dirPath string
	}

	type want struct {
		list []dirscan.File
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
			fields: fields{
				allowedExtensions: []string{".mp4", ".jpg"},
				os: func(mockCtrl *gomock.Controller) dirscan.OS {
					mock := mocks.NewMockOS(mockCtrl)
					mock.EXPECT().Stat("/data").Return(nil, nil)
					mock.EXPECT().Stat("/data/dir").Return(nil, nil)

					mockFile := mocks.NewMockOSFile(mockCtrl)
					mockFile.EXPECT().Readdir(-1).Return([]os.FileInfo{
						&fakeFile{name: "file1.mp4", size: 10, mode: 0, isDir: false},
						&fakeFile{name: "file2.jpg", size: 20, mode: 0, isDir: false},
						&fakeFile{name: "file3.mp4", size: 30, mode: 0, isDir: false},
						&fakeFile{name: "file4.jpg", size: 40, mode: 0, isDir: false},
						&fakeFile{name: "file5.mp4", size: 50, mode: 0, isDir: false},
						&fakeFile{name: "file6.non", size: 60, mode: 0, isDir: false},
						&fakeFile{name: "dir", size: 0, mode: os.ModeDir, isDir: true},
					}, nil)
					mockFile.EXPECT().Close().Return(nil)
					mock.EXPECT().Open("/data").Return(mockFile, nil)

					mockFileInner := mocks.NewMockOSFile(mockCtrl)
					mockFileInner.EXPECT().Readdir(-1).Return([]os.FileInfo{
						&fakeFile{name: "file7.mp4", size: 70, mode: 0, isDir: false},
						&fakeFile{name: "file8.jpg", size: 80, mode: 0, isDir: false},
						&fakeFile{name: "file9.non", size: 90, mode: 0, isDir: false},
					}, nil)
					mockFileInner.EXPECT().Close().Return(nil)
					mock.EXPECT().Open("/data/dir").Return(mockFileInner, nil)

					return mock
				},
			},
			args: args{
				dirPath: "/data",
			},
			want: want{
				list: []dirscan.File{
					{Name: "file1.mp4", Path: "/data/file1.mp4", Size: 10},
					{Name: "file2.jpg", Path: "/data/file2.jpg", Size: 20},
					{Name: "file3.mp4", Path: "/data/file3.mp4", Size: 30},
					{Name: "file4.jpg", Path: "/data/file4.jpg", Size: 40},
					{Name: "file5.mp4", Path: "/data/file5.mp4", Size: 50},
					{Name: "file7.mp4", Path: "/data/dir/file7.mp4", Size: 70},
					{Name: "file8.jpg", Path: "/data/dir/file8.jpg", Size: 80},
				},
			},
		},
		{
			name: "sad path, path does not exist",
			fields: fields{
				os: func(mockCtrl *gomock.Controller) dirscan.OS {
					mock := mocks.NewMockOS(mockCtrl)

					err := assert.AnError
					mock.EXPECT().IsNotExist(err).Return(true)
					mock.EXPECT().Stat("/data").Return(nil, err)
					return mock
				},
			},
			args: args{
				dirPath: "/data",
			},
			want: want{
				err: errors.Wrap(assert.AnError, "path does not exist"),
			},
		},
		{
			name: "sad path, can't stat the path",
			fields: fields{
				os: func(mockCtrl *gomock.Controller) dirscan.OS {
					mock := mocks.NewMockOS(mockCtrl)

					err := assert.AnError
					mock.EXPECT().IsNotExist(err).Return(false)
					mock.EXPECT().Stat("/data").Return(nil, err)

					return mock
				},
			},
			args: args{
				dirPath: "/data",
			},
			want: want{
				err: errors.Wrap(assert.AnError, "can't stat the path"),
			},
		},
		{
			name: "sad path, can't open the dir",
			fields: fields{
				os: func(mockCtrl *gomock.Controller) dirscan.OS {
					mock := mocks.NewMockOS(mockCtrl)
					mock.EXPECT().Stat("/data").Return(nil, nil)

					mock.EXPECT().Open("/data").Return(nil, assert.AnError)
					return mock
				},
			},
			args: args{
				dirPath: "/data",
			},
			want: want{
				err: errors.Wrap(assert.AnError, "error opening the directory"),
			},
		},
		{
			name: "sad path, can't close the dir",
			fields: fields{
				os: func(mockCtrl *gomock.Controller) dirscan.OS {
					mock := mocks.NewMockOS(mockCtrl)
					mock.EXPECT().Stat("/data").Return(nil, nil)

					mockFile := mocks.NewMockOSFile(mockCtrl)
					mockFile.EXPECT().Readdir(-1).Return([]os.FileInfo{}, nil)
					mockFile.EXPECT().Close().Return(assert.AnError)
					mock.EXPECT().Open("/data").Return(mockFile, nil)
					return mock
				},
			},
			args: args{
				dirPath: "/data",
			},
			want: want{
				err: errors.Wrap(assert.AnError, "cannot close directory"),
			},
		},
		{
			name: "sad path, can't read the dir",
			fields: fields{
				os: func(mockCtrl *gomock.Controller) dirscan.OS {
					mock := mocks.NewMockOS(mockCtrl)
					mock.EXPECT().Stat("/data").Return(nil, nil)

					mockFile := mocks.NewMockOSFile(mockCtrl)
					mockFile.EXPECT().Readdir(-1).Return([]os.FileInfo{}, assert.AnError)
					mockFile.EXPECT().Close().Return(nil)
					mock.EXPECT().Open("/data").Return(mockFile, nil)
					return mock
				},
			},
			args: args{
				dirPath: "/data",
			},
			want: want{
				err: errors.Wrap(assert.AnError, "error reading the directory"),
			},
		},
		{
			name: "sad path, error getting file list recursively",
			fields: fields{
				allowedExtensions: []string{".mp4", ".jpg"},
				os: func(mockCtrl *gomock.Controller) dirscan.OS {
					mock := mocks.NewMockOS(mockCtrl)
					mock.EXPECT().Stat("/data").Return(nil, nil)

					err := assert.AnError
					mock.EXPECT().Stat("/data/dir").Return(nil, err)

					mock.EXPECT().IsNotExist(err).Return(true)

					mockFile := mocks.NewMockOSFile(mockCtrl)
					mockFile.EXPECT().Readdir(-1).Return([]os.FileInfo{
						&fakeFile{name: "file1.mp4", size: 10, mode: 0, isDir: false},
						&fakeFile{name: "file2.jpg", size: 20, mode: 0, isDir: false},
						&fakeFile{name: "file3.mp4", size: 30, mode: 0, isDir: false},
						&fakeFile{name: "file4.jpg", size: 40, mode: 0, isDir: false},
						&fakeFile{name: "file5.mp4", size: 50, mode: 0, isDir: false},
						&fakeFile{name: "dir", size: 0, mode: os.ModeDir, isDir: true},
					}, nil)
					mockFile.EXPECT().Close().Return(nil)
					mock.EXPECT().Open("/data").Return(mockFile, nil)

					return mock
				},
			},
			args: args{
				dirPath: "/data",
			},
			want: want{
				err: errors.Wrap(
					errors.Wrap(assert.AnError, "path does not exist"),
					"error getting file list recursively",
				),
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

			scanner := dirscan.NewScanner(
				tt.fields.allowedExtensions,
				dirscan.WithOS(tt.fields.os(mockCtrl)),
			)

			got, err := scanner.GetFileList(tt.args.dirPath)
			if tt.want.err == nil {
				assert.Equal(t, sortList(tt.want.list), sortList(got))
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func sortList(list []dirscan.File) []dirscan.File {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list
}

type fakeFile struct {
	name  string
	size  int64
	mode  os.FileMode
	isDir bool
}

func (f *fakeFile) Name() string {
	return f.name
}

func (f *fakeFile) Size() int64 {
	return f.size
}

func (f *fakeFile) Mode() os.FileMode {
	return f.mode
}

func (f *fakeFile) ModTime() time.Time {
	return time.Time{}
}

func (f *fakeFile) IsDir() bool {
	return f.isDir
}

func (f *fakeFile) Sys() interface{} {
	return nil
}
