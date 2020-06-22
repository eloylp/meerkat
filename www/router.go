package www

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/eloylp/meerkat/data"
)

func Router(dfr *data.FlowRegistry) http.Handler {
	h := mux.NewRouter()
	h.HandleFunc(DashboardPath, HandleHTMLClient(dfr)).Methods(http.MethodGet)
	h.HandleFunc(DataStreamPath+"/{id}", HandleMJPEG(dfr)).Methods(http.MethodGet)
	return h
}
