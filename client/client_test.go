package client_test

import (
	"bytes"
	"fmt"
	"github.com/legosx/gopro-media-library-verifier/client"
	"github.com/legosx/gopro-media-library-verifier/client/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

//go:generate mockgen -destination=./mocks/http_client.go -package=mocks github.com/legosx/gopro-media-library-verifier/client HTTPClient
//go:generate mockgen -destination=./mocks/http_requester.go -package=mocks github.com/legosx/gopro-media-library-verifier/client HTTPRequester

func TestNewClient(t *testing.T) {
	type fields struct {
		token string
		opts  func(mockCtrl *gomock.Controller) []func(client *client.Client) error
	}

	type want struct {
		err error
	}

	tests := []struct {
		name string
		fields
		want
	}{
		{
			name: "happy path",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					return []func(client *client.Client) error{}
				},
			},
		},
		{
			name: "happy path, with auth check",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						assert.Equal(t, "GET", req.Method)
						assert.Equal(t, "api.gopro.com", req.Header.Get("Authority"))
						assert.Equal(t, "application/vnd.gopro.jk.media+json; version=2.0.0", req.Header.Get("Accept"))
						assert.Equal(t, "Bearer token", req.Header.Get("Authorization"))
						assert.Equal(t, "https://api.gopro.com/notification_center/notifications", req.URL.String())

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
						}, nil
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
						client.WithAuthCheck(),
					}
				},
			},
		},
		{
			name: "sad path, auth check failed",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						assert.Equal(t, "GET", req.Method)
						assert.Equal(t, "api.gopro.com", req.Header.Get("Authority"))
						assert.Equal(t, "application/vnd.gopro.jk.media+json; version=2.0.0", req.Header.Get("Accept"))
						assert.Equal(t, "Bearer token", req.Header.Get("Authorization"))
						assert.Equal(t, "https://api.gopro.com/notification_center/notifications", req.URL.String())

						return &http.Response{
							StatusCode: http.StatusUnauthorized,
							Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
						}, nil
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
						client.WithAuthCheck(),
					}
				},
			},
			want: want{
				err: errors.Wrap(client.NewErrorResponse(&http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
				}), "error checking authentication"),
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

			got, err := client.NewClient(tt.fields.token, tt.fields.opts(mockCtrl)...)
			if tt.want.err == nil {
				assert.NotNil(t, got)
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func TestClient_AuthCheck(t *testing.T) {
	type fields struct {
		token string
		opts  func(mockCtrl *gomock.Controller) []func(client *client.Client) error
	}

	type want struct {
		err error
	}

	tests := []struct {
		name string
		fields
		want
	}{
		{
			name: "happy path",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						assert.Equal(t, "GET", req.Method)
						assert.Equal(t, "api.gopro.com", req.Header.Get("Authority"))
						assert.Equal(t, "application/vnd.gopro.jk.media+json; version=2.0.0", req.Header.Get("Accept"))
						assert.Equal(t, "Bearer token", req.Header.Get("Authorization"))

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
						}, nil
					})

					return []func(client *client.Client) error{client.WithHTTPClient(httpClientMock)}
				},
			},
		},
		{
			name: "sad path",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						assert.Equal(t, "GET", req.Method)
						assert.Equal(t, "api.gopro.com", req.Header.Get("Authority"))
						assert.Equal(t, "application/vnd.gopro.jk.media+json; version=2.0.0", req.Header.Get("Accept"))
						assert.Equal(t, "Bearer token", req.Header.Get("Authorization"))

						return &http.Response{
							StatusCode: http.StatusUnauthorized,
							Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
						}, nil
					})

					return []func(client *client.Client) error{client.WithHTTPClient(httpClientMock)}
				},
			},
			want: want{
				err: errors.Wrap(client.NewErrorResponse(&http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
				}), "error checking authentication"),
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

			c, err := client.NewClient(tt.fields.token, tt.fields.opts(mockCtrl)...)
			assert.NoError(t, err)

			err = c.AuthCheck()
			if tt.want.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func TestClient_GetAllowedExtensions(t *testing.T) {
	type fields struct {
		token string
		opts  func(mockCtrl *gomock.Controller) []func(client *client.Client) error
	}

	type want struct {
		allowedExtensions []string
	}

	tests := []struct {
		name string
		fields
		want
	}{
		{
			name: "happy path",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					return []func(client *client.Client) error{}
				},
			},
			want: want{
				allowedExtensions: []string{".mp4", ".mov", ".360", ".heic", ".jpg", ".jpeg", ".png"},
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

			c, err := client.NewClient(tt.fields.token, tt.fields.opts(mockCtrl)...)
			assert.NoError(t, err)

			got := c.GetAllowedExtensions()
			assert.Equal(t, tt.want.allowedExtensions, got)
		})
	}
}

func TestClient_GetPage(t *testing.T) {
	type args struct {
		pageNumber int
		perPage    int
	}

	type fields struct {
		token string
		opts  func(mockCtrl *gomock.Controller) []func(client *client.Client) error
	}

	type want struct {
		page client.Page
		err  error
	}

	page2try, page3try := 0, 0

	tests := []struct {
		name string
		fields
		args
		want
	}{
		{
			name: "happy path, get page 1",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).Times(1).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						assert.Equal(t, "GET", req.Method)
						assert.Equal(t, "api.gopro.com", req.Header.Get("Authority"))
						assert.Equal(t, "application/vnd.gopro.jk.media+json; version=2.0.0", req.Header.Get("Accept"))
						assert.Equal(t, "Bearer token", req.Header.Get("Authorization"))

						u, err := url.Parse(req.URL.String())
						assert.NoError(t, err)

						assert.Equal(t, "https", u.Scheme)
						assert.Equal(t, "api.gopro.com", u.Host)
						assert.Equal(t, "/media/search", u.Path)

						q := u.Query()
						assert.Equal(t, "2", q.Get("per_page"))
						assert.Equal(t, "captured_at", q.Get("order_by"))
						assert.Equal(t, "", q.Get("type"))
						assert.Equal(t, "filename,file_size", q.Get("fields"))
						assert.Equal(t, "registered,rendering,pretranscoding,transcoding,failure,ready", q.Get("processing_states"))
						assert.Equal(t, "1", q.Get("page"))

						medias := []string{
							`{"filename": "file1.mp4","file_size": 10}`,
							`{"filename": "file2.jpg","file_size": 20}`,
						}

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       getBody(2, 2, 5, 3, medias),
						}, nil
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
					}
				},
			},
			args: args{pageNumber: 1, perPage: 2},
			want: want{
				page: client.NewPage(
					3,
					[]client.Media{
						client.NewMedia("file1.mp4", 10),
						client.NewMedia("file2.jpg", 20),
					},
				),
			},
		},
		{
			name: "happy path, get page 2 with one retry",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).Times(2).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						u, err := url.Parse(req.URL.String())
						assert.NoError(t, err)

						q := u.Query()
						assert.Equal(t, "2", q.Get("page"))

						medias := []string{
							`{"filename": "file3.mp4","file_size": 30}`,
							`{"filename": "file4.jpg","file_size": 40}`,
						}

						page2try++
						if page2try == 1 {
							medias = []string{
								`{"filename": "file3.mp4","file_size": 30}`,
							}
						}

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       getBody(2, 2, 5, 3, medias),
						}, nil
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
					}
				},
			},
			args: args{pageNumber: 2, perPage: 2},
			want: want{
				page: client.NewPage(
					3,
					[]client.Media{
						client.NewMedia("file3.mp4", 30),
						client.NewMedia("file4.jpg", 40),
					},
				),
			},
		},
		{
			name: "happy path, get page 3 with 5 retries",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).Times(6).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						u, err := url.Parse(req.URL.String())
						assert.NoError(t, err)

						q := u.Query()
						assert.Equal(t, "3", q.Get("page"))

						medias := []string{
							`{"filename": "file7.mp4","file_size": 70}`,
							`{"filename": "file8.jpg","file_size": 80}`,
						}

						page3try++
						if page3try < 6 {
							medias = []string{
								`{"filename": "file7.mp4","file_size": 70}`,
							}
						}

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       getBody(3, 3, 8, 3, medias),
						}, nil
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
					}
				},
			},
			args: args{pageNumber: 3, perPage: 3},
			want: want{
				page: client.NewPage(
					3,
					[]client.Media{
						client.NewMedia("file7.mp4", 70),
						client.NewMedia("file8.jpg", 80),
					},
				),
			},
		},
		{
			name: "sad path, cannot create http request",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {

					httpRequesterMock := mocks.NewMockHTTPRequester(mockCtrl)
					httpRequesterMock.EXPECT().NewRequest(gomock.Any(), gomock.Any(), gomock.Any()).
						Times(1).
						Return(nil, assert.AnError)

					return []func(client *client.Client) error{
						client.WithHTTPClient(mocks.NewMockHTTPClient(mockCtrl)),
						client.WithHTTPRequester(httpRequesterMock),
					}
				},
			},
			args: args{pageNumber: 1, perPage: 2},
			want: want{
				err: errors.Wrap(assert.AnError,
					"error getting page: "+
						"error getting page with retry: "+
						"error getting data from client: "+
						"error creating HTTP request",
				),
			},
		},
		{
			name: "sad path, error from client",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).Times(1).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						return &http.Response{}, assert.AnError
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
					}
				},
			},
			args: args{pageNumber: 1, perPage: 2},
			want: want{
				err: errors.Wrap(assert.AnError,
					"error getting page: "+
						"error getting page with retry: "+
						"error getting data from client: "+
						"error performing HTTP request",
				),
			},
		},
		{
			name: "sad path, error response with empty body",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).Times(1).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusUnauthorized,
							Status:     "401 Unauthorized",
						}, nil
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
					}
				},
			},
			args: args{pageNumber: 0, perPage: 0},
			want: want{
				err: errors.New("error getting page: " +
					"error getting page with retry: " +
					"error getting data from client: " +
					"401 Unauthorized"),
			},
		},
		{
			name: "sad path, error response with body",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).Times(1).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusUnauthorized,
							Status:     "401 Unauthorized",
							Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
						}, nil
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
					}
				},
			},
			args: args{pageNumber: 0, perPage: 0},
			want: want{
				err: errors.New("error getting page: " +
					"error getting page with retry: " +
					"error getting data from client: " +
					"401 Unauthorized"),
			},
		},
		{
			name: "sad path, error response, can not close the body",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).Times(1).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						body := io.NopCloser(bytes.NewBufferString(`{}`))
						err := body.Close()
						assert.NoError(t, err)

						return &http.Response{
							StatusCode: http.StatusUnauthorized,
							Status:     "401 Unauthorized",
							Body: &errorReadCloser{
								Reader: strings.NewReader("response body"),
							},
						}, nil
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
					}
				},
			},
			args: args{pageNumber: 0, perPage: 0},
			want: want{
				err: errors.New("error getting page:" +
					" error getting page with retry: " +
					"error getting data from client: " +
					"error closing HTTP response body: " +
					"401 Unauthorized; Close error"),
			},
		},
		{
			name: "sad path, error response, can not read the body",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).Times(1).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						body := io.NopCloser(bytes.NewBufferString(`{}`))
						err := body.Close()
						assert.NoError(t, err)

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       &errorReader{},
						}, nil
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
					}
				},
			},
			args: args{pageNumber: 0, perPage: 0},
			want: want{
				err: errors.New("error getting page: " +
					"error getting page with retry: " +
					"error getting data from client: " +
					"error reading Data from HTTP response body: " +
					"simulated read error"),
			},
		},
		{
			name: "sad path, invalid json",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).Times(1).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString(`{invalid}`)),
						}, nil
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
					}
				},
			},
			args: args{pageNumber: 0, perPage: 0},
			want: want{
				err: errors.New("error getting page:" +
					" error getting page with retry: " +
					"error decoding JSON: " +
					"invalid character 'i' looking for beginning of object key string, json:" +
					" {invalid}: invalid character 'i' looking for beginning of object key string"),
			},
		},
		{
			name: "sad path, unexpected response from API",
			fields: fields{
				token: "token",
				opts: func(mockCtrl *gomock.Controller) []func(client *client.Client) error {
					httpClientMock := mocks.NewMockHTTPClient(mockCtrl)
					httpClientMock.EXPECT().Do(gomock.Any()).Times(1).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(
								`{"error": "message"}`,
							)),
						}, nil
					})

					return []func(client *client.Client) error{
						client.WithHTTPClient(httpClientMock),
					}
				},
			},
			args: args{pageNumber: 1, perPage: 2},
			want: want{
				err: errors.New("error getting page: error getting page with retry: unexpected response: {\"error\": \"message\"}"),
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

			c, err := client.NewClient(tt.fields.token, tt.fields.opts(mockCtrl)...)
			assert.NoError(t, err)

			got, err := c.GetPage(tt.args.pageNumber, tt.args.perPage)
			if tt.want.err == nil {
				assert.Equal(t, tt.want.page, got)
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func getBody(pageNumber, perPage, totalItems, totalPages int, medias []string) io.ReadCloser {
	return io.NopCloser(bytes.NewBufferString(fmt.Sprintf(
		`{"_pages": {"current_page": %d,"per_page": %d,"total_items": %d,"total_pages": %d},"_embedded": {"media": [%s]}}`,
		pageNumber, perPage, totalItems, totalPages,
		strings.Join(medias, ","),
	)))
}

type errorReadCloser struct {
	io.Reader
}

func (erc *errorReadCloser) Close() error {
	return errors.New("Close error")
}

type errorReader struct{}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

func (erc *errorReader) Close() error {
	return nil
}
