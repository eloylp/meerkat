// +build unit

package data_test

import (
	"bytes"
	"github.com/eloylp/meerkat/data"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHTTPFetcher_Fetch(t *testing.T) {
	resURL := "http://iesource.com/camera.jpg"
	body := "OK"
	client := NewTestClient(func(req *http.Request) *http.Response {
		requestedURL := req.URL.String()
		assert.Equal(t, resURL, requestedURL, "Expected camera url is %s but client used %s", resURL, requestedURL)
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}
	})
	fetcher := data.NewHTTPFetcher(client)
	reader, err := fetcher.Fetch(resURL)
	assert.NoError(t, err)
	d, err := ioutil.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, body, string(d), "Expected body is %v got %v", body, reader)
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
