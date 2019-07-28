package main

import (
	"go-sentinel/config"
	"go-sentinel/fetch"
	"go-sentinel/www"
	"log"
	"net/http"
)

var frameBuffer chan []byte

func init() {
	frameBuffer = make(chan []byte, 10)
}

func main() {
	cfg := config.C()
	fetcher := fetch.NewHTTPFetcher(&http.Client{})
	dataPump := fetch.NewDataPump(cfg.PollInterval, cfg.Resource, frameBuffer, fetcher)
	go dataPump.Start()
	server := www.NewServer(cfg.HTTPListenAddress, frameBuffer)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
