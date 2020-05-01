package factory

import (
	"github.com/eloylp/go-serve/www"
	"github.com/eloylp/meerkat/flow"
	"net/http"
	"time"
)

type HTTPServedApp struct {
	httpServer       *http.Server
	dataFlowRegistry *flow.DataFlowRegistry
}

func (a *HTTPServedApp) Start() error {
	for _, dataFlow := range a.dataFlowRegistry.DataFlows() {
		go dataFlow.Start()
	}
	shutDownTimeout := 20 * time.Second
	www.Shutdown(a.httpServer, shutDownTimeout)
	if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
