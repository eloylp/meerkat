package www

import (
	"github.com/eloylp/meerkat/flow"
	"net/http"
)

const DataStreamPath = "/data/"
const DashboardPath = "/"

type HTTPServer struct {
	listenAddress string
	dfr           *flow.DataFlowRegistry
}

func NewHTTPServer(listenAddress string, dfr *flow.DataFlowRegistry) *HTTPServer {
	return &HTTPServer{listenAddress: listenAddress, dfr: dfr}
}

func (s *HTTPServer) Start() error {
	h := http.NewServeMux()
	h.HandleFunc(DashboardPath, HandleHTMLClient(s.dfr))
	h.HandleFunc(DataStreamPath, HandleMJPEG(s.dfr))
	if err := http.ListenAndServe(s.listenAddress, h); err != nil {
		return err
	}
	return nil
}
