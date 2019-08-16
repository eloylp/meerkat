package fetch_test

import (
	"bytes"
	"go-sentinel/fetch"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHTTPFetcher_Fetch(t *testing.T) {

	resUrl := "http://iesource.com/camera.jpg"
	body := "OK"
	client := NewTestClient(func(req *http.Request) *http.Response {
		requestedUrl := req.URL.String()
		if requestedUrl != resUrl {
			t.Errorf("Expected camera url is %s but client used %s", resUrl, requestedUrl)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}
	})
	fetcher := fetch.NewHTTPFetcher(client)
	reader, err := fetcher.Fetch(resUrl)
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Error(err)
	}
	if string(d) != body {
		t.Errorf("Expected body is %v got %v", body, reader)
	}
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
