package factory

import (
	"github.com/eloylp/meerkat/config"
	"github.com/eloylp/meerkat/www"
	"net/http"
)

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
