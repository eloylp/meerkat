package app

import (
	"go-sentinel/config"
	"go-sentinel/fetch"
	"go-sentinel/store"
	"go-sentinel/unique"
	"net/http"
)

type App struct {
	httpServer       *server
	dataFlowRegistry *dataFlowRegistry
}

func NewApp(cfg config.Config) *App {

	dfr := &dataFlowRegistry{}

	for _, r := range cfg.Resources {
		dataStore := store.NewTimeLineStore(10)
		fetcher := fetch.NewHTTPFetcher(&http.Client{})
		dataPump := fetch.NewDataPump(cfg.PollInterval, r, fetcher, dataStore)
		dfr.Add(&dataFlow{
			UUID:      unique.UUID4(),
			Resource:  r,
			DataStore: dataStore,
			DataPump:  dataPump,
		})
	}

	return &App{
		httpServer:       newHTTPServer(cfg.HTTPListenAddress, dfr),
		dataFlowRegistry: dfr,
	}
}

func (a *App) Start() error {
	for _, dataFlow := range a.dataFlowRegistry.DataFlows() {
		go dataFlow.Start()
	}
	if err := a.httpServer.Start(); err != nil {
		return err
	}
	return nil
}
