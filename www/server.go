package www

import (
	"go-sentinel/store"
	"net/http"
)

const FrameStreamEndpoint = "/data"

type Server interface {
	Start() error
}

type server struct {
	listenAddress string
	store         store.Store
}

func NewServer(listenAddress string, store store.Store) *server {
	return &server{listenAddress: listenAddress, store: store}
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
