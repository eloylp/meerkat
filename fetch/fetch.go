package fetch

import (
	"io/ioutil"
	"net/http"
)

type Fetcher interface {
	Fetch(res string) ([]byte, error)
}

type hTTPFetcher struct {
	client *http.Client
}

func NewHTTPFetcher(client *http.Client) *hTTPFetcher {
	return &hTTPFetcher{client: client}
}

func (f *hTTPFetcher) Fetch(res string) ([]byte, error) {

	r, err := f.client.Get(res)
	if err != nil {
		return nil, err
	}
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	return d, nil

}
