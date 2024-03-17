package buildclient_test

import (
	"github.com/legosx/gopro-media-library-verifier/buildclient"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestTokenPromptMethodInput_GetToken(t *testing.T) {
	type fields struct {
		opts func(mockCtrl *gomock.Controller) []func(*buildclient.TokenPromptMethodInput)
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
				opts: func(mockCtrl *gomock.Controller) []func(*buildclient.TokenPromptMethodInput) {
					promptInput := func(label string, mask rune, hideEntered bool, validateFunc func(value string) error) (string, error) {
						return "valid", nil
					}

					return []func(*buildclient.TokenPromptMethodInput){
						buildclient.WithInputPromptInput(promptInput),
					}
				},
			},
			want: want{
				token: "valid",
			},
		},
		{
			name: "sad path, input prompt failed",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(*buildclient.TokenPromptMethodInput) {
					promptInput := func(label string, mask rune, hideEntered bool, validateFunc func(value string) error) (string, error) {
						return "", assert.AnError
					}

					return []func(*buildclient.TokenPromptMethodInput){
						buildclient.WithInputPromptInput(promptInput),
					}
				},
			},
			want: want{
				err: errors.Wrap(assert.AnError, "input prompt failed"),
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

			got, err := buildclient.NewTokenPromptMethodInput(tt.fields.opts(mockCtrl)...).GetToken()
			if tt.want.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.token, got)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}

		})
	}
}

func TestTokenPromptMethodInput_Validate(t *testing.T) {
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
				value: "valid",
			},
		},
		{
			name: "sad path, empty input",
			args: args{
				value: "",
			},
			want: want{
				err: errors.New("token cannot be empty"),
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

			err := buildclient.NewTokenPromptMethodInput().Validate(tt.args.value)
			if tt.want.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}
