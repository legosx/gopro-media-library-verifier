package fetch

import (
	"github.com/legosx/gopro-media-library-verifier/client"
	"github.com/pkg/errors"
)

type Client interface {
	GetPage(pageNumber, perPage int) (page client.Page, err error)
}

type Fetcher struct {
	client  Client
	perPage int
}

func NewFetcher(client Client, opts ...func(fetcher *Fetcher)) Fetcher {
	f := Fetcher{
		client:  client,
		perPage: 250,
	}

	for _, opt := range opts {
		opt(&f)
	}

	return f
}

func WithPerPage(perPage int) func(f *Fetcher) {
	return func(f *Fetcher) {
		f.perPage = perPage
	}
}

func (f Fetcher) GetMedias() (medias []Media, err error) {
	medias = []Media{}

	handleResult := func(result getPageResult) error {
		if err = result.err; err != nil {
			return errors.Wrap(err, "error getting medias")
		}

		if page := result.page; len(page.Medias()) > 0 {
			medias = append(medias, f.convertClientMedias(page.Medias())...)
		}

		return nil
	}

	maxConcurrentCalls := 10

	result := newGetPageResult(f.client.GetPage(1, f.perPage))
	if err = handleResult(result); err != nil {
		return []Media{}, err
	}

	totalPages := result.page.TotalPages()

	resultCh := make(chan getPageResult, totalPages)
	sem := make(chan struct{}, maxConcurrentCalls)
	for pageNumber := 2; pageNumber <= totalPages; pageNumber++ {
		go func(pageNumber int) {
			sem <- struct{}{}
			defer func() { <-sem }()

			resultCh <- newGetPageResult(f.client.GetPage(pageNumber, f.perPage))
		}(pageNumber)
	}

	for i := 2; i <= totalPages; i++ {
		if err = handleResult(<-resultCh); err != nil {
			return []Media{}, err
		}
	}

	return medias, nil
}

func (f Fetcher) convertClientMedias(medias []client.Media) (convertedMedias []Media) {
	for _, media := range medias {
		convertedMedias = append(convertedMedias, NewMedia(media.FileName(), media.FileSize()))
	}

	return convertedMedias
}

type getPageResult struct {
	page client.Page
	err  error
}

func newGetPageResult(page client.Page, err error) getPageResult {
	return getPageResult{
		page: page,
		err:  err,
	}
}
