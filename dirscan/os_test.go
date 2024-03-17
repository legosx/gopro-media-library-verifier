package dirscan_test

import (
	"github.com/legosx/gopro-media-library-verifier/dirscan"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestOSWrapper_Stat(t *testing.T) {
	type args struct {
		name string
	}

	type fileInfo struct {
		name  string
		isDir bool
	}

	type want struct {
		fileInfo fileInfo
		err      error
	}

	tests := []struct {
		name string
		args
		want
	}{
		{
			name: "happy path, file",
			args: args{
				name: "../go.mod",
			},
			want: want{
				fileInfo: fileInfo{
					name:  "go.mod",
					isDir: false,
				},
				err: nil,
			},
		},
		{
			name: "happy path, directory",
			args: args{
				name: "../cmd",
			},
			want: want{
				fileInfo: fileInfo{
					name:  "cmd",
					isDir: true,
				},
				err: nil,
			},
		},
		{
			name: "sad path, no file",
			args: args{
				name: "does-not-exist",
			},
			want: want{
				err: errors.New("stat does-not-exist: no such file or directory"),
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := dirscan.NewOSWrapper().Stat(tt.args.name)
			if tt.want.err == nil {
				assert.Equal(t, tt.want.fileInfo.name, got.Name())
				assert.Equal(t, tt.want.fileInfo.isDir, got.IsDir())
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func TestOSWrapper_Open(t *testing.T) {
	type args struct {
		name string
	}

	type want struct {
		err error
	}

	tests := []struct {
		name string
		args
		want
	}{
		{
			name: "happy path",
			args: args{
				name: "../go.mod",
			},
		},
		{
			name: "sad path, no file",
			args: args{
				name: "does-not-exist",
			},
			want: want{
				err: errors.New("open does-not-exist: no such file or directory"),
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := dirscan.NewOSWrapper().Open(tt.args.name)
			if tt.want.err == nil {
				assert.NotNil(t, got)
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func TestOSWrapper_IsNotExist(t *testing.T) {
	type args struct {
		err func() error
	}

	type want struct {
		IsNotExist bool
	}

	tests := []struct {
		name string
		args
		want
	}{
		{
			name: "happy path",
			args: args{err: func() error {
				_, errFileNotExist := os.Stat("does-not-exist")
				return errFileNotExist
			}},
			want: want{
				IsNotExist: true,
			},
		},
		{
			name: "sad path, another error",
			args: args{err: func() error {
				return assert.AnError
			}},
			want: want{
				IsNotExist: false,
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := dirscan.NewOSWrapper().IsNotExist(tt.args.err())

			assert.Equal(t, tt.want.IsNotExist, got)
		})
	}
}
