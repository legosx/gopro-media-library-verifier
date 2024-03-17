package client_test

import (
	"github.com/legosx/gopro-media-library-verifier/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPage(t *testing.T) {
	type fields struct {
		totalPages int
		medias     []client.Media
	}

	type want struct {
		page       client.Page
		totalPages int
		medias     []client.Media
	}

	tests := []struct {
		name string
		fields
		want
	}{
		{
			name: "happy path",
			fields: fields{
				totalPages: 1,
				medias:     []client.Media{client.NewMedia("file1.mp4", 10)},
			},
			want: want{
				page:       client.NewPage(1, []client.Media{client.NewMedia("file1.mp4", 10)}),
				totalPages: 1,
				medias:     []client.Media{client.NewMedia("file1.mp4", 10)},
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := client.NewPage(tt.fields.totalPages, tt.fields.medias)
			assert.Equal(t, tt.want.page, got)
			assert.Equal(t, tt.want.totalPages, got.TotalPages())
			assert.Equal(t, tt.want.medias, got.Medias())
		})
	}
}
