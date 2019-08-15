package app

import (
	"errors"
	"fmt"
	"go-sentinel/config"
	"go-sentinel/fetch"
	"go-sentinel/store"
	"go-sentinel/www"
	"net/http"
)

type Server interface {
	Start() error
}

type App struct {
	hTTPServer Server
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

type DataFlowRegistry struct {
	DataFlows []*DataFlow
}

func (dfr *DataFlowRegistry) Add(df *DataFlow) {
	dfr.DataFlows = append(dfr.DataFlows, df)
}

func (dfr *DataFlowRegistry) FindStore(wfUid string) (store.Store, error) {
	for _, wf := range dfr.DataFlows {
		if wf.UUID == wfUid {
			return wf.Store, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Cannot find workflow %v", wfUid))
}

type DataFlow struct {
	UUID  string
	Store store.Store
	Pump  fetch.DataPump
}
