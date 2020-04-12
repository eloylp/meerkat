package fetch

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

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
	// TODO , investigate if is needed to read the entire body to close it.
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
