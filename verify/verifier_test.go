package verify_test

import (
	"github.com/legosx/gopro-media-library-verifier/dirscan"
	"github.com/legosx/gopro-media-library-verifier/fetch"
	"github.com/legosx/gopro-media-library-verifier/verify"
	"github.com/legosx/gopro-media-library-verifier/verify/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

//go:generate mockgen -destination=./mocks/fetcher.go -package=mocks github.com/legosx/gopro-media-library-verifier/verify Fetcher
//go:generate mockgen -destination=./mocks/scanner.go -package=mocks github.com/legosx/gopro-media-library-verifier/verify Scanner

func TestVerifier_IdentifyMissingFiles(t *testing.T) {
	type fields struct {
		fetcher func(mockCtrl *gomock.Controller) verify.Fetcher
		scanner func(mockCtrl *gomock.Controller) verify.Scanner
	}

	type args struct {
		path string
	}

	type want struct {
		filePaths []string
		err       error
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
				fetcher: func(mockCtrl *gomock.Controller) verify.Fetcher {
					mock := mocks.NewMockFetcher(mockCtrl)
					mock.
						EXPECT().
						GetMedias().
						Return(
							[]fetch.Media{
								fetch.NewMedia("file1.mp4", 1000),
								fetch.NewMedia("file2.jpg", 2000),
							},
							nil,
						)

					return mock
				},
				scanner: func(mockCtrl *gomock.Controller) verify.Scanner {
					mock := mocks.NewMockScanner(mockCtrl)
					mock.
						EXPECT().
						GetFileList("/dir").
						Return(
							[]dirscan.File{
								{Name: "file1.mp4", Path: "/dir/file1.mp4", Size: 1000},
								{Name: "file2.jpg", Path: "/dir/file2.jpg", Size: 2000},
								{Name: "file3.mp4", Path: "/dir/file3.mp4", Size: 3000},
								{Name: "file4.jpg", Path: "/dir/file4.jpg", Size: 4000},
								{Name: "file5.mp4", Path: "/dir/file5.mp4", Size: 5000},
							},
							nil,
						)

					return mock
				},
			},
			args: args{
				path: "/dir",
			},
			want: want{
				filePaths: []string{
					"/dir/file3.mp4",
					"/dir/file4.jpg",
					"/dir/file5.mp4",
				},
				err: nil,
			},
		},
		{
			name: "sad path, scanner.GetFileList error",
			fields: fields{
				fetcher: func(mockCtrl *gomock.Controller) verify.Fetcher {
					return mocks.NewMockFetcher(mockCtrl)
				},
				scanner: func(mockCtrl *gomock.Controller) verify.Scanner {
					mock := mocks.NewMockScanner(mockCtrl)
					mock.
						EXPECT().
						GetFileList("/dir").
						Return([]dirscan.File{}, assert.AnError)

					return mock
				},
			},
			args: args{
				path: "/dir",
			},
			want: want{
				filePaths: []string{},
				err:       errors.Wrap(assert.AnError, "error getting local files"),
			},
		},
		{
			name: "sad path, fetcher.GetMedias error",
			fields: fields{
				fetcher: func(mockCtrl *gomock.Controller) verify.Fetcher {
					mock := mocks.NewMockFetcher(mockCtrl)
					mock.
						EXPECT().
						GetMedias().
						Return([]fetch.Media{}, assert.AnError)

					return mock
				},
				scanner: func(mockCtrl *gomock.Controller) verify.Scanner {
					mock := mocks.NewMockScanner(mockCtrl)
					mock.
						EXPECT().
						GetFileList("/dir").
						Return(
							[]dirscan.File{
								{Name: "file1.mp4", Path: "/dir/file1.mp4", Size: 1000},
								{Name: "file2.jpg", Path: "/dir/file2.jpg", Size: 2000},
								{Name: "file3.mp4", Path: "/dir/file3.mp4", Size: 3000},
								{Name: "file4.jpg", Path: "/dir/file4.jpg", Size: 4000},
								{Name: "file5.mp4", Path: "/dir/file5.mp4", Size: 5000},
							},
							nil,
						)

					return mock
				},
			},
			args: args{
				path: "/dir",
			},
			want: want{
				filePaths: []string{},
				err: errors.Wrap(
					errors.Wrap(assert.AnError, "error getting remote medias"),
					"error getting remote files",
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

			verifier := verify.NewVerifier(
				tt.fields.fetcher(mockCtrl),
				tt.fields.scanner(mockCtrl),
			)
			got, err := verifier.IdentifyMissingFiles(tt.args.path)

			assert.Equal(t, tt.want.filePaths, got)
			if tt.want.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}
