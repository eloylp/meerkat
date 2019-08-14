package fetch

import (
	"io"
	"net/http"
)

type Fetcher interface {
	Fetch(res string) (io.Reader, error)
}

type hTTPFetcher struct {
	client *http.Client
}

func NewHTTPFetcher(client *http.Client) *hTTPFetcher {
	return &hTTPFetcher{client: client}
}

func (f *hTTPFetcher) Fetch(res string) (io.Reader, error) {
	r, err := f.client.Get(res)
	if err != nil {
		return nil, err
	}
	return r.Body, nil
}
