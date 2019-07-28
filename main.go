package main

import (
	"go-sentinel/config"
	"go-sentinel/fetch"
	"go-sentinel/www"
	"log"
	"net/http"
	"time"
)

var frameBuffer chan []byte

func init() {
	frameBuffer = make(chan []byte, 10)
}

func main() {
	cfg := config.C()
	go startPolling(cfg.PollInterval, cfg.Resource, frameBuffer)
	server := www.NewServer(cfg.HTTPListenAddress, frameBuffer)
	if err := server.Start(); err != nil {
		log.Fatal(err)
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
