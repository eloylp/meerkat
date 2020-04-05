package factory

import (
	"github.com/eloylp/meerkat/config"
	"github.com/eloylp/meerkat/fetch"
	"github.com/eloylp/meerkat/flow"
	"github.com/eloylp/meerkat/store"
	"github.com/eloylp/meerkat/unique"
	"github.com/eloylp/meerkat/www"
	"net/http"
)

type App struct {
	httpServer       *www.HTTPServer
	dataFlowRegistry *flow.DataFlowRegistry
}

func NewApp(cfg config.Config) *App {

	dfr := &flow.DataFlowRegistry{}

	for _, r := range cfg.Resources {
		dataStore := store.NewTimeLineStore(10)
		fetcher := fetch.NewHTTPFetcher(&http.Client{})
		dataPump := fetch.NewDataPump(cfg.PollInterval, r, fetcher, dataStore)
		dfr.Add(flow.NewDataFlow(unique.UUID4(), r, dataStore, dataPump))
	}

	return &App{
		httpServer:       www.NewHTTPServer(cfg.HTTPListenAddress, dfr),
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
