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

var frameBuffer chan []byte

const FrameStreamEndpoint = "/data"

func init() {
	frameBuffer = make(chan []byte, 10)
}

type config struct {
	CameraUrl           string
	PollingIntervalSecs uint
	HTTPListenAddress   string
}

func main() {

	cfg := cfg()
	go startCameraPolling(cfg.PollingIntervalSecs, cfg.CameraUrl, frameBuffer)
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

func cfg() *config {
	c := &config{}
	flag.StringVar(&c.CameraUrl, "u", "", "Pass the camera url")
	flag.UintVar(&c.PollingIntervalSecs, "i", 1, "Pass the camera interval")
	flag.StringVar(&c.HTTPListenAddress, "l", "0.0.0.0:3000", "Pass the http server for serving results")
	flag.Parse()
	return c
}

func startCameraPolling(interval uint, cameraUrl string, frames chan []byte) {

	for {
		time.Sleep(time.Duration(interval) * time.Second)
		frame, err := pullFrameFromURL(cameraUrl)
		if err != nil {
			log.Fatal(err)
		}
		frames <- frame
	}
}

func pullFrameFromURL(u string) ([]byte, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	jpeg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	return jpeg, nil
}
