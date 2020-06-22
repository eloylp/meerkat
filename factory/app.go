package factory

import (
	"net/http"

	"github.com/eloylp/meerkat/config"
	"github.com/eloylp/meerkat/www"
)

func NewHTTPServedApp(cfg config.Config) (*HTTPServedApp, error) {
	dfr, err := NewDataFlowRegistry(cfg)
	if err != nil {
		return nil, err
	}
	s := &http.Server{
		Addr:    cfg.HTTPListenAddress,
		Handler: www.Router(dfr),
	}
	return &HTTPServedApp{
		httpServer:       s,
		dataFlowRegistry: dfr,
	}, nil
}
