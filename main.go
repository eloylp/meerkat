package main

import (
	"go-sentinel/config"
	"go-sentinel/fetch"
	"go-sentinel/store"
	"go-sentinel/www"
	"log"
	"net/http"
)

func main() {
	cfg := config.C()
	fetcher := fetch.NewHTTPFetcher(&http.Client{})
	timeLineStore := store.NewTimeLineStore(10, 10000)
	dataPump := fetch.NewDataPump(cfg.PollInterval, cfg.Resource, fetcher, timeLineStore)
	go dataPump.Start()
	server := www.NewServer(cfg.HTTPListenAddress, timeLineStore)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
