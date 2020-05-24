package www

import (
	"github.com/eloylp/meerkat/data"
	"net/http"
)

func Router(dfr *data.FlowRegistry) http.Handler {
	h := http.NewServeMux()
	h.HandleFunc(DashboardPath, HandleHTMLClient(dfr))
	h.HandleFunc(DataStreamPath, HandleMJPEG(dfr))
	return h
}
