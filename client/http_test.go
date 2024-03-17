package client_test

import (
	"github.com/legosx/gopro-media-library-verifier/client"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestHttpWrapper_NewRequest(t *testing.T) {
	type args struct {
		method string
		url    string
		body   io.Reader
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
				method: "GET",
				url:    "https://test.com/",
				body:   io.NopCloser(nil),
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := client.NewHTTPWrapper().NewRequest(tt.args.method, tt.args.url, tt.args.body)
			if tt.want.err == nil {
				assert.NotNil(t, got)
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}
