package buildclient_test

import (
	"github.com/legosx/gopro-media-library-verifier/buildclient"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestTokenPromptMethodCURL_GetToken(t *testing.T) {
	type fields struct {
		opts func(mockCtrl *gomock.Controller) []func(*buildclient.TokenPromptMethodCURL)
	}

	type want struct {
		token string
		err   error
	}

	tests := []struct {
		name string
		fields
		want
	}{
		{
			name: "happy path, valid token",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(*buildclient.TokenPromptMethodCURL) {
					promptInput := func(label string, mask rune, hideEntered bool, validateFunc func(value string) error) (string, error) {
						return "Bearer valid", nil
					}

					return []func(*buildclient.TokenPromptMethodCURL){
						buildclient.WithCURLPromptInput(promptInput),
					}
				},
			},
			want: want{
				token: "valid",
			},
		},
		{
			name: "happy path, but empty token from bearer",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(*buildclient.TokenPromptMethodCURL) {
					promptInput := func(label string, mask rune, hideEntered bool, validateFunc func(value string) error) (string, error) {
						return "wrong", nil
					}

					return []func(*buildclient.TokenPromptMethodCURL){
						buildclient.WithCURLPromptInput(promptInput),
					}
				},
			},
			want: want{
				err: errors.New("no bearer token found in CURL request"),
			},
		},
		{
			name: "sad path, curl prompt failed",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(*buildclient.TokenPromptMethodCURL) {
					promptInput := func(label string, mask rune, hideEntered bool, validateFunc func(value string) error) (string, error) {
						return "", assert.AnError
					}

					return []func(*buildclient.TokenPromptMethodCURL){
						buildclient.WithCURLPromptInput(promptInput),
					}
				},
			},
			want: want{
				err: errors.Wrap(assert.AnError, "curl prompt failed"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			got, err := buildclient.NewTokenPromptMethodCURL(tt.fields.opts(mockCtrl)...).GetToken()
			if tt.want.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.token, got)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func TestTokenPromptMethodCURL_Validate(t *testing.T) {
	type args struct {
		value string
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
			name: "happy path, valid input",
			args: args{
				value: "Bearer valid",
			},
		},
		{
			name: "sad path, empty input",
			args: args{
				value: "",
			},
			want: want{
				err: errors.New("curl request cannot be empty"),
			},
		},
		{
			name: "sad path, no bearer token found",
			args: args{
				value: "no token",
			},
			want: want{
				err: errors.New("no bearer token found in CURL request"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			err := buildclient.NewTokenPromptMethodCURL().Validate(tt.args.value)
			if tt.want.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}
