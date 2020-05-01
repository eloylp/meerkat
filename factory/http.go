package factory

import (
	"github.com/eloylp/go-serve/www"
	"github.com/eloylp/meerkat/data"
	"net/http"
	"time"
)

type HTTPServedApp struct {
	httpServer       *http.Server
	dataFlowRegistry *data.FlowRegistry
}

func (a *HTTPServedApp) Start() error {
	for _, dataFlow := range a.dataFlowRegistry.Flows() {
		go dataFlow.Start()
	}
	shutDownTimeout := 20 * time.Second
	www.Shutdown(a.httpServer, shutDownTimeout)
	if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
