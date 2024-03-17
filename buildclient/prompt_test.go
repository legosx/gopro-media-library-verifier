package buildclient_test

import (
	"errors"
	"github.com/legosx/gopro-media-library-verifier/buildclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromptInput(t *testing.T) {
	type args struct {
		label        string
		mask         rune
		hideEntered  bool
		validateFunc func(value string) error
	}

	type want struct {
		value string
		err   error
	}

	tests := []struct {
		name string
		args
		want
	}{
		{
			name: "sad path, interrupt by the test",
			args: args{
				label:       "test",
				mask:        '*',
				hideEntered: true,
				validateFunc: func(value string) error {
					return nil
				},
			},
			want: want{
				err: errors.New("^D"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildclient.PromptInput(tt.args.label, tt.args.mask, tt.args.hideEntered, tt.args.validateFunc)
			if tt.want.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.value, got)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}

		})
	}
}

func TestPromptSelect(t *testing.T) {
	type args struct {
		label string
		items interface{}
	}

	type want struct {
		value string
		err   error
	}

	tests := []struct {
		name string
		args
		want
	}{
		{
			name: "sad path, interrupt by the test",
			args: args{
				label: "test",
				items: []string{"a", "b", "c"},
			},
			want: want{
				err: errors.New("^D"),
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := buildclient.PromptSelect(tt.args.label, tt.args.items)
			if tt.want.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.value, got)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}

		})
	}
}
