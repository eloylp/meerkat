package www

import (
	"net/http"
)

const FrameStreamEndpoint = "/data"

type server struct {
	listenAddress string
	frameBuffer   chan []byte
}

func NewServer(listenAddress string, frameBuffer chan []byte) *server {
	return &server{listenAddress: listenAddress, frameBuffer: frameBuffer}
}

func (s *server) Start() error {
	h := http.NewServeMux()
	h.HandleFunc("/", s.handleHTMLClient())
	h.HandleFunc(FrameStreamEndpoint, s.handleMJPEG(s.frameBuffer))
	if err := http.ListenAndServe(s.listenAddress, h); err != nil {
		return err
	}
	return nil
}
