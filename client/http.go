package client

import (
	"io"
	"net/http"
)

type HTTPWrapper struct{}

func NewHTTPWrapper() HTTPWrapper {
	return HTTPWrapper{}
}

func (h HTTPWrapper) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}
