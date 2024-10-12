package client

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	url_endpoint       = "https://api.gopro.com/"
	path_notifications = "notification_center/notifications"
	path_media_search  = "media/search"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPRequester interface {
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
}

type Client struct {
	token         string
	httpClient    HTTPClient
	httpRequester HTTPRequester
}

func NewClient(token string, opts ...func(client *Client) error) (client *Client, err error) {
	client = &Client{
		token:         token,
		httpClient:    &http.Client{},
		httpRequester: NewHTTPWrapper(),
	}

	for _, opt := range opts {
		if err = opt(client); err != nil {
			return &Client{}, err
		}
	}

	return client, nil
}

func WithAuthCheck() func(c *Client) error {
	return func(c *Client) error {
		return c.AuthCheck()
	}
}

func WithHTTPClient(httpClient HTTPClient) func(c *Client) error {
	return func(c *Client) error {
		c.httpClient = httpClient

		return nil
	}
}

func WithHTTPRequester(httpRequester HTTPRequester) func(c *Client) error {
	return func(c *Client) error {
		c.httpRequester = httpRequester

		return nil
	}
}

func (c Client) AuthCheck() (err error) {
	if _, err = c.get(path_notifications, map[string]string{}); err != nil {
		return errors.Wrap(err, "error checking authentication")
	}

	return nil
}

func (c Client) GetAllowedExtensions() []string {
	return []string{".mp4", ".mov", ".360", ".heic", ".jpg", ".jpeg", ".png"}
}

func (c Client) GetPage(pageNumber, perPage int) (page Page, err error) {
	response, err := c.getPageWithRetry(pageNumber, perPage, 10)
	if err != nil {
		return Page{}, errors.Wrap(err, "error getting page")
	}

	var medias []Media

	for _, media := range response.Embedded.Media {
		medias = append(medias, NewMedia(media.FileName, media.FileSize))
	}

	return NewPage(response.Pages.TotalPages, medias), nil
}

// Sometimes it doesn't return all items from the first try
func (c Client) getPageWithRetry(pageNumber, perPage, maxRetries int) (page *page, err error) {
	for retry := 0; retry < maxRetries; retry++ {
		if page, err = c.getPage(pageNumber, perPage); err != nil {
			return nil, errors.Wrap(err, "error getting page with retry")
		}

		pageMediaCount := len(page.Embedded.Media)
		isExpectedItems := pageMediaCount == perPage

		if isLastPage := pageNumber == page.Pages.TotalPages; isLastPage {
			isExpectedItems = ((page.Pages.CurrentPage-1)*perPage)+pageMediaCount == page.Pages.TotalItems
		}

		if isExpectedItems {
			break
		}
	}

	return page, nil
}

func (c Client) getPage(pageNumber, perPage int) (page *page, err error) {
	body, err := c.get(path_media_search, map[string]string{
		"fields":            strings.Join(c.getDefaultFields(), ","),
		"processing_states": strings.Join(c.getDefaultProcessingStates(), ","),
		"order_by":          "captured_at",
		"per_page":          strconv.Itoa(perPage),
		"page":              strconv.Itoa(pageNumber),
		"type":              strings.Join(c.getDefaultTypes(), ","),
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting data from client")
	}

	if err = json.Unmarshal(body, &page); err != nil {
		return nil, errors.Wrapf(err, "error decoding JSON: %+v, json: %s", err, string(body))
	}

	if page.Pages.CurrentPage == 0 {
		return nil, errors.Errorf("unexpected response: %s", string(body))
	}

	return page, nil
}

func (c Client) get(path string, queryParameters map[string]string) (body []byte, err error) {
	headers := map[string]string{
		"Authority":     "api.gopro.com",
		"Accept":        "application/vnd.gopro.jk.media+json; version=2.0.0",
		"Authorization": "Bearer " + c.token,
	}

	url := url_endpoint + path

	delimiter := "?"
	for key, value := range queryParameters {
		url = url + delimiter + key + "=" + value
		delimiter = "&"
	}

	// Create an HTTP request with the specified URL and headers
	req, err := c.httpRequester.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating HTTP request")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Perform the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error performing HTTP request")
	}
	defer func() {
		if body := resp.Body; body != nil {
			if innerErr := body.Close(); innerErr != nil {
				err = errors.Wrap(multierr.Append(err, innerErr), "error closing HTTP response body")
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, NewErrorResponse(resp)
	}

	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, errors.Wrap(err, "error reading Data from HTTP response body")
	}

	return body, nil
}

func (c Client) getDefaultProcessingStates() []string {
	return []string{
		"registered",
		"rendering",
		"pretranscoding",
		"transcoding",
		"failure",
		"ready",
	}
}

func (c Client) getDefaultFields() []string {
	return []string{"filename", "file_size"}
}

func (c Client) getDefaultTypes() []string {
	return []string{
		"Burst",
		"BurstVideo",
		"Continuous",
		"LoopedVideo",
		"Photo",
		"TimeLapse",
		"TimeLapseVideo",
		"Video",
		"MultiClipEdit",
	}
}
