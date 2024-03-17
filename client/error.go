package client

import "net/http"

type ErrorResponse struct {
	resp *http.Response
}

func NewErrorResponse(resp *http.Response) error {
	return ErrorResponse{resp: resp}
}

func (e ErrorResponse) Error() string {
	return e.resp.Status
}
