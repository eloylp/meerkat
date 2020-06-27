package factory

import (
	"net/http"
	"sync"
	"time"

	"github.com/eloylp/kit/shutdown"

	"github.com/eloylp/meerkat/data"
)

type HTTPServedApp struct {
	httpServer       *http.Server
	dataFlowRegistry *data.FlowRegistry
}

func (a *HTTPServedApp) Start() error {
	for _, dataFlow := range a.dataFlowRegistry.Flows() {
		go dataFlow.Start()
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	shutDownTimeout := 20 * time.Second
	shutdown.WithOSSignals(a.httpServer, shutDownTimeout, wg, nil)
	if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	wg.Wait()
	return nil
}
