package factory

import (
	"net/http"
	"time"

	httpserver "github.com/eloylp/go-serve/www"

	"github.com/eloylp/meerkat/config"
	"github.com/eloylp/meerkat/flow"
	"github.com/eloylp/meerkat/www"
)

type HTTPServedApp struct {
	httpServer       *http.Server
	dataFlowRegistry *flow.DataFlowRegistry
}

func (a *HTTPServedApp) Start() error {
	for _, dataFlow := range a.dataFlowRegistry.DataFlows() {
		go dataFlow.Start()
	}
	httpserver.Shutdown(a.httpServer, 20*time.Second)
	if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func NewHTTPServedApp(cfg config.Config) (*HTTPServedApp, error) {
	dfr, err := NewDataFlowRegistry(cfg)
	if err != nil {
		return nil, err
	}
	h := http.NewServeMux()
	h.HandleFunc(www.DashboardPath, www.HandleHTMLClient(dfr))
	h.HandleFunc(www.DataStreamPath, www.HandleMJPEG(dfr))
	s := &http.Server{
		Addr:    cfg.HTTPListenAddress,
		Handler: h,
	}
	return &HTTPServedApp{
		httpServer:       s,
		dataFlowRegistry: dfr,
	}, nil
}
