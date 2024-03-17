package fetch_test

import (
	"github.com/legosx/gopro-media-library-verifier/client"
	"github.com/legosx/gopro-media-library-verifier/fetch"
	"github.com/legosx/gopro-media-library-verifier/fetch/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"sort"
	"testing"
)

//go:generate mockgen -destination=./mocks/client.go -package=mocks github.com/legosx/gopro-media-library-verifier/fetch Client

func TestFetcher_GetMedias(t *testing.T) {
	type fields struct {
		client func(mockCtrl *gomock.Controller) fetch.Client
		opts   func(mockCtrl *gomock.Controller) []func(fetcher *fetch.Fetcher)
	}

	type want struct {
		medias []fetch.Media
		err    error
	}

	tests := []struct {
		name string
		fields
		want
	}{
		{
			name: "happy path",
			fields: fields{
				client: func(mockCtrl *gomock.Controller) fetch.Client {
					mock := mocks.NewMockClient(mockCtrl)
					mock.EXPECT().GetPage(gomock.Any(), gomock.Any()).Times(3).DoAndReturn(func(pageNumber, perPage int) (page client.Page, err error) {
						assert.GreaterOrEqual(t, pageNumber, 1)
						assert.LessOrEqual(t, pageNumber, 3)
						assert.Equal(t, 2, perPage)

						mediasPerPages := map[int][]client.Media{
							1: {
								client.NewMedia("file1.mp4", 10),
								client.NewMedia("file2.jpg", 20),
							},
							2: {
								client.NewMedia("file3.mp4", 30),
								client.NewMedia("file4.jpg", 40),
							},
							3: {
								client.NewMedia("file5.mp4", 50),
							},
						}

						return client.NewPage(3, mediasPerPages[pageNumber]), nil
					})

					return mock
				},
				opts: func(mockCtrl *gomock.Controller) []func(fetcher *fetch.Fetcher) {
					return []func(fetcher *fetch.Fetcher){
						fetch.WithPerPage(2),
					}
				},
			},
			want: want{
				medias: []fetch.Media{
					fetch.NewMedia("file1.mp4", 10),
					fetch.NewMedia("file2.jpg", 20),
					fetch.NewMedia("file3.mp4", 30),
					fetch.NewMedia("file4.jpg", 40),
					fetch.NewMedia("file5.mp4", 50),
				},
			},
		},
		{
			name: "sad path, client error on page 1",
			fields: fields{
				client: func(mockCtrl *gomock.Controller) fetch.Client {
					mock := mocks.NewMockClient(mockCtrl)
					mock.EXPECT().GetPage(1, 250).Times(1).DoAndReturn(func(pageNumber, perPage int) (page client.Page, err error) {
						return client.Page{}, assert.AnError
					})

					return mock
				},
				opts: func(mockCtrl *gomock.Controller) []func(fetcher *fetch.Fetcher) {
					return []func(fetcher *fetch.Fetcher){}
				},
			},
			want: want{
				err: errors.Wrap(assert.AnError, "error getting medias"),
			},
		},
		{
			name: "sad path, client error on page 2",
			fields: fields{
				client: func(mockCtrl *gomock.Controller) fetch.Client {
					mock := mocks.NewMockClient(mockCtrl)
					mock.EXPECT().GetPage(gomock.Any(), gomock.Any()).Times(2).DoAndReturn(func(pageNumber, perPage int) (page client.Page, err error) {
						assert.True(t, pageNumber == 1 || pageNumber == 2)

						if pageNumber == 2 {
							return client.Page{}, assert.AnError
						}

						return client.NewPage(2, []client.Media{
							client.NewMedia("file1.mp4", 10),
							client.NewMedia("file2.jpg", 20),
						}), nil
					})

					return mock
				},
				opts: func(mockCtrl *gomock.Controller) []func(fetcher *fetch.Fetcher) {
					return []func(fetcher *fetch.Fetcher){
						fetch.WithPerPage(2),
					}
				},
			},
			want: want{
				err: errors.Wrap(assert.AnError, "error getting medias"),
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

			fetcher := fetch.NewFetcher(tt.fields.client(mockCtrl), tt.fields.opts(mockCtrl)...)

			got, err := fetcher.GetMedias()
			if tt.want.err == nil {
				assert.Equal(t, sortList(tt.want.medias), sortList(got))
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func sortList(list []fetch.Media) []fetch.Media {
	sort.Slice(list, func(i, j int) bool {
		return list[i].FileName() < list[j].FileName()
	})

	return list
}
