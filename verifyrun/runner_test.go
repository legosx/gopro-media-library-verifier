package verifyrun_test

import (
	"crypto/rand"
	"github.com/legosx/gopro-media-library-verifier/buildclient"
	"github.com/legosx/gopro-media-library-verifier/client"
	"github.com/legosx/gopro-media-library-verifier/dirscan"
	"github.com/legosx/gopro-media-library-verifier/fetch"
	"github.com/legosx/gopro-media-library-verifier/verifyrun"
	"github.com/legosx/gopro-media-library-verifier/verifyrun/mocks"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"os"
	"path/filepath"
	"testing"
)

//go:generate mockgen -destination=./mocks/verifier.go -package=mocks github.com/legosx/gopro-media-library-verifier/verifyrun Verifier

func TestRunner_Run(t *testing.T) {
	type fields struct {
		path              string
		outputFilePath    func() string
		tokenPromptMethod verifyrun.TokenPromptMethod
		opts              func(mockCtrl *gomock.Controller) []func(*verifyrun.Runner)
	}

	type want struct {
		err    error
		output string
	}

	tests := []struct {
		name string
		fields
		want
	}{
		{
			name: "happy path",
			fields: fields{
				path:              "test",
				tokenPromptMethod: verifyrun.TokenPromptMethodInput,
				outputFilePath:    func() string { return "" },
				opts: func(mockCtrl *gomock.Controller) []func(*verifyrun.Runner) {
					buildClient := func(opts ...func(builder *buildclient.Builder)) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					buildVerifier := func(fetcher fetch.Fetcher, scanner dirscan.Scanner) verifyrun.Verifier {
						verifier := mocks.NewMockVerifier(mockCtrl)

						verifier.EXPECT().IdentifyMissingFiles("test").Return([]string{"test/file1.mp4"}, nil)

						return verifier
					}

					return []func(*verifyrun.Runner){
						verifyrun.WithBuildClient(buildClient),
						verifyrun.WithBuildVerifier(buildVerifier),
					}
				},
			},
		},
		{
			name: "sad path, buildClient fails",
			fields: fields{
				path:              "test",
				tokenPromptMethod: verifyrun.TokenPromptMethodCURL,
				outputFilePath:    func() string { return "" },
				opts: func(mockCtrl *gomock.Controller) []func(*verifyrun.Runner) {
					buildClient := func(opts ...func(builder *buildclient.Builder)) (c *client.Client, err error) {
						return nil, assert.AnError
					}

					return []func(*verifyrun.Runner){
						verifyrun.WithBuildClient(buildClient),
					}
				},
			},
			want: want{
				err: assert.AnError,
			},
		},
		{
			name: "sad path, invalid token prompt method",
			fields: fields{
				path:              "test",
				tokenPromptMethod: "invalid",
				outputFilePath:    func() string { return "" },
				opts: func(mockCtrl *gomock.Controller) []func(*verifyrun.Runner) {
					return []func(*verifyrun.Runner){}
				},
			},
			want: want{
				err: errors.New("invalid token prompt method: invalid"),
			},
		},
		{
			name: "happy path, no token prompt methods specified",
			fields: fields{
				path:           "test",
				outputFilePath: func() string { return "" },
				opts: func(mockCtrl *gomock.Controller) []func(*verifyrun.Runner) {
					buildClient := func(opts ...func(builder *buildclient.Builder)) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					buildVerifier := func(fetcher fetch.Fetcher, scanner dirscan.Scanner) verifyrun.Verifier {
						verifier := mocks.NewMockVerifier(mockCtrl)

						verifier.EXPECT().IdentifyMissingFiles("test").Return([]string{}, nil)

						return verifier
					}

					return []func(*verifyrun.Runner){
						verifyrun.WithBuildClient(buildClient),
						verifyrun.WithBuildVerifier(buildVerifier),
					}
				},
			},
		},
		{
			name: "sad path, verifier IdentifyMissingFiles fails",
			fields: fields{
				path:           "test",
				outputFilePath: func() string { return "" },
				opts: func(mockCtrl *gomock.Controller) []func(*verifyrun.Runner) {
					buildClient := func(opts ...func(builder *buildclient.Builder)) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					buildVerifier := func(fetcher fetch.Fetcher, scanner dirscan.Scanner) verifyrun.Verifier {
						verifier := mocks.NewMockVerifier(mockCtrl)

						verifier.EXPECT().IdentifyMissingFiles("test").Return([]string{}, assert.AnError)

						return verifier
					}

					return []func(*verifyrun.Runner){
						verifyrun.WithBuildClient(buildClient),
						verifyrun.WithBuildVerifier(buildVerifier),
					}
				},
			},
			want: want{
				err: assert.AnError,
			},
		},
		{
			name: "happy path, outputFilePath specified",
			fields: fields{
				path: "test",
				outputFilePath: func() string {
					path, err := createRandomOutputFilePath()
					assert.NoError(t, err)

					return path
				},
				opts: func(mockCtrl *gomock.Controller) []func(*verifyrun.Runner) {
					buildClient := func(opts ...func(builder *buildclient.Builder)) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					buildVerifier := func(fetcher fetch.Fetcher, scanner dirscan.Scanner) verifyrun.Verifier {
						verifier := mocks.NewMockVerifier(mockCtrl)

						verifier.EXPECT().IdentifyMissingFiles("test").Return([]string{"test/file1.mp4"}, nil)

						return verifier
					}

					return []func(*verifyrun.Runner){
						verifyrun.WithBuildClient(buildClient),
						verifyrun.WithBuildVerifier(buildVerifier),
					}
				},
			},
			want: want{
				output: "test/file1.mp4\n",
			},
		},
		{
			name: "happy path, but can't write to outputFilePath",
			fields: fields{
				path: "test",
				outputFilePath: func() string {
					return "/d/o/e/s/not/exist/output.txt"
				},
				opts: func(mockCtrl *gomock.Controller) []func(*verifyrun.Runner) {
					buildClient := func(opts ...func(builder *buildclient.Builder)) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					buildVerifier := func(fetcher fetch.Fetcher, scanner dirscan.Scanner) verifyrun.Verifier {
						verifier := mocks.NewMockVerifier(mockCtrl)

						verifier.EXPECT().IdentifyMissingFiles("test").Return([]string{"test/file1.mp4"}, nil)

						return verifier
					}

					return []func(*verifyrun.Runner){
						verifyrun.WithBuildClient(buildClient),
						verifyrun.WithBuildVerifier(buildVerifier),
					}
				},
			},
			want: want{
				err: errors.New("open /d/o/e/s/not/exist/output.txt: no such file or directory"),
			},
		},
		{
			name: "happy path, real verifier",
			fields: fields{
				path:              "/d/o/e/s/not/exist",
				tokenPromptMethod: verifyrun.TokenPromptMethodInput,
				outputFilePath:    func() string { return "" },
				opts: func(mockCtrl *gomock.Controller) []func(*verifyrun.Runner) {
					buildClient := func(opts ...func(builder *buildclient.Builder)) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					return []func(*verifyrun.Runner){
						verifyrun.WithBuildClient(buildClient),
					}
				},
			},
			want: want{
				err: errors.New("error getting local files: path does not exist: stat /d/o/e/s/not/exist: no such file or directory"),
			},
		},
		{
			name: "sad path, real client fails",
			fields: fields{
				path:              "test",
				tokenPromptMethod: verifyrun.TokenPromptMethodInput,
				outputFilePath:    func() string { return "" },
				opts: func(mockCtrl *gomock.Controller) []func(*verifyrun.Runner) {
					return []func(*verifyrun.Runner){}
				},
			},
			want: want{
				err: errors.New("failed to get client from user prompt: token prompt failed: input prompt failed: ^D"),
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

			outputFilePath := tt.outputFilePath()

			err := verifyrun.NewRunner(tt.path, outputFilePath, tt.tokenPromptMethod, tt.fields.opts(mockCtrl)...).Run()
			if tt.want.err == nil {
				assert.NoError(t, err)

				if tt.want.output != "" {
					output, err := os.ReadFile(outputFilePath)
					assert.NoError(t, err)
					assert.Equal(t, tt.want.output, string(output))
				}
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func TestRunner_Init(t *testing.T) {
	type args struct {
		cmd *cobra.Command
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
				cmd: &cobra.Command{},
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := verifyrun.Init(tt.args.cmd)
			if tt.want.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func randomString() string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 30)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b)
}

func createRandomOutputFilePath() (configFilePath string, err error) {
	path := filepath.Join(os.TempDir(), randomString())
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}

	return filepath.Join(path, "output.txt"), nil
}
