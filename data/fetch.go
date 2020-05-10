package data

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type Fetcher struct {
	client *http.Client
}

func NewHTTPFetcher(client *http.Client) *Fetcher {
	return &Fetcher{client: client}
}

func (f *Fetcher) Fetch(res string) (io.Reader, error) {
	r, err := f.client.Get(res)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	return reader, nil
}
