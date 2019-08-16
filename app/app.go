package app

import (
	"go-sentinel/config"
	"go-sentinel/fetch"
	"go-sentinel/store"
	"go-sentinel/unique"
	"net/http"
)

type Server interface {
	Start() error
}

type App struct {
	hTTPServer       Server
	DataFlowRegistry *DataFlowRegistry
}

func NewApp(config config.Config) *App {

	dfr := &DataFlowRegistry{}

	for _, r := range config.Resources {
		dataStore := store.NewTimeLineStore(10)
		fetcher := fetch.NewHTTPFetcher(&http.Client{})
		dataPump := fetch.NewDataPump(config.PollInterval, r, fetcher, dataStore)
		dfr.Add(&DataFlow{
			UUID:      unique.UUID4(),
			Resource:  r,
			DataStore: dataStore,
			DataPump:  dataPump,
		})
	}

	return &App{
		hTTPServer:       newServer(config.HTTPListenAddress, dfr),
		DataFlowRegistry: dfr,
	}
}

func (a *App) Start() error {
	for _, dataFlow := range a.DataFlowRegistry.DataFlows() {
		go dataFlow.Start()
	}
	if err := a.hTTPServer.Start(); err != nil {
		return err
	}
	return nil
}
