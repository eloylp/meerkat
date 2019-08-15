package app

import (
	"net/http"
)

const FrameStreamEndpoint = "/data/"

type server struct {
	listenAddress string
	dfr           *DataFlowRegistry
}

func newServer(listenAddress string, dfr *DataFlowRegistry) *server {
	return &server{listenAddress: listenAddress, dfr: dfr}
}

func (s *server) Start() error {
	h := http.NewServeMux()
	h.HandleFunc("/", s.handleHTMLClient())
	h.HandleFunc(FrameStreamEndpoint, s.handleMJPEG())
	if err := http.ListenAndServe(s.listenAddress, h); err != nil {
		return err
	}
	return nil
}
