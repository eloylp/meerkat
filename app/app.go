package app

import (
	"go-sentinel/config"
	"go-sentinel/fetch"
	"go-sentinel/store"
	"go-sentinel/www"
	"net/http"
)

type App struct {
	hTTPServer www.Server
	dataPump   fetch.DataPump
}

func NewApp(config config.Config) *App {

	dataStore := store.NewTimeLineStore(10)
	fetcher := fetch.NewHTTPFetcher(&http.Client{})
	return &App{
		hTTPServer: www.NewServer(config.HTTPListenAddress, dataStore),
		dataPump:   fetch.NewDataPump(config.PollInterval, config.Resource, fetcher, dataStore),
	}
}

func (a *App) Start() error {
	go a.dataPump.Start()
	if err := a.hTTPServer.Start(); err != nil {
		return err
	}
	return nil
}
