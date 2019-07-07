package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

var imageBuffer chan []byte

func init() {
	imageBuffer = make(chan []byte, 10)
}

type config struct {
	CameraUrl           string
	PollingIntervalSecs uint
	HTTPListenAddress   string
}

func main() {

	cfg := cfg()
	go startCameraPolling(cfg.PollingIntervalSecs, cfg.CameraUrl, imageBuffer)
	h := http.NewServeMux()
	h.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {

		mimeWriter := multipart.NewWriter(w)
		contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
		w.Header().Add("Content-Type", contentType)

		for image := range imageBuffer {

			partHeader := make(textproto.MIMEHeader)
			partHeader.Add("Content-Type", "image/jpeg")

			partWriter, partErr := mimeWriter.CreatePart(partHeader)
			if nil != partErr {
				log.Fatal(partErr.Error())
			}

			if _, writeErr := partWriter.Write(image); nil != writeErr {
				log.Fatal(writeErr.Error())
			}
		}
	})
	h.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-type", "text/html")
		html := `<!DOCTYPE html><html><body><img src="/data"></body></html>`
		_, _ = w.Write([]byte(html))
	})
	if err := http.ListenAndServe(cfg.HTTPListenAddress, h); err != nil {
		log.Fatal(err)
	}
}

func cfg() *config {
	c := &config{}
	flag.StringVar(&c.CameraUrl, "u", "", "Pass the camera url")
	flag.UintVar(&c.PollingIntervalSecs, "i", 1, "Pass the camera interval")
	flag.StringVar(&c.HTTPListenAddress, "l", "0.0.0.0:3000", "Pass the http server for serving results")
	flag.Parse()
	return c
}

func startCameraPolling(interval uint, cameraUrl string, images chan []byte) {

	for {
		time.Sleep(time.Duration(interval) * time.Second)
		resp, err := http.Get(cameraUrl)
		if err != nil {
			log.Fatal(err)
		}
		jpeg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
		images <- jpeg
	}
}
