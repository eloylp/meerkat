package app

import (
	"net/http"
)

const DataStreamPath = "/data/"
const DashboardPath = "/"

type server struct {
	listenAddress string
	dfr           *dataFlowRegistry
}

func newHTTPServer(listenAddress string, dfr *dataFlowRegistry) *server {
	return &server{listenAddress: listenAddress, dfr: dfr}
}

func (s *server) Start() error {
	h := http.NewServeMux()
	h.HandleFunc(DashboardPath, s.handleHTMLClient())
	h.HandleFunc(DataStreamPath, s.handleMJPEG())
	if err := http.ListenAndServe(s.listenAddress, h); err != nil {
		return err
	}
	return nil
}
