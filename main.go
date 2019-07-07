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

var images chan []byte

func main() {

	var cameraUrl string
	flag.StringVar(&cameraUrl, "u", "", "Pass the camera url")
	var interval uint
	flag.UintVar(&interval, "i", 1, "Pass the camera interval")
	var listenAddress string
	flag.StringVar(&listenAddress, "l", "0.0.0.0:3000", "Pass the http server for serving results")
	flag.Parse()
	images = make(chan []byte, 10)
	go startCameraPolling(interval, cameraUrl, images)
	h := http.NewServeMux()
	h.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {

		mimeWriter := multipart.NewWriter(w)
		contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
		w.Header().Add("Content-Type", contentType)

		for image := range images {

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
	if err := http.ListenAndServe(listenAddress, h); err != nil {
		log.Fatal(err)
	}
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
