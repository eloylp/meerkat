package main

import (
	"fmt"
	"go-sentinel/config"
	"go-sentinel/fetch"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

var frameBuffer chan []byte

const FrameStreamEndpoint = "/data"

func init() {
	frameBuffer = make(chan []byte, 10)
}

func main() {

	cfg := config.C()
	go startPolling(cfg.PollInterval, cfg.Resource, frameBuffer)
	h := http.NewServeMux()
	h.HandleFunc(FrameStreamEndpoint, MJPEG(frameBuffer))
	h.HandleFunc("/", HTMLClient())
	if err := http.ListenAndServe(cfg.HTTPListenAddress, h); err != nil {
		log.Fatal(err)
	}
}

func HTMLClient() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-type", "text/html")
		html := fmt.Sprintf(`<!DOCTYPE html><html><body><img src="%s"></body></html>`, FrameStreamEndpoint)
		_, _ = w.Write([]byte(html))
	}
}

func MJPEG(imageBuffer chan []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mimeWriter := multipart.NewWriter(w)
		contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
		w.Header().Add("Content-Type", contentType)

		for image := range imageBuffer {
			writeFrame(mimeWriter, image)
		}
	}
}

func writeFrame(mimeWriter *multipart.Writer, image []byte) {
	partHeader := make(textproto.MIMEHeader)
	partHeader.Add("Content-Type", "image/jpeg")
	partWriter, partErr := mimeWriter.CreatePart(partHeader)
	if partErr != nil {
		log.Fatal(partErr.Error())
	}
	if _, writeErr := partWriter.Write(image); writeErr != nil {
		log.Fatal(writeErr.Error())
	}
}

func startPolling(interval uint, url string, frames chan []byte) {
	f := fetch.NewHTTPFetcher(&http.Client{})
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		frame, err := f.Fetch(url)
		if err != nil {
			log.Fatal(err)
		}
		frames <- frame
	}
}
