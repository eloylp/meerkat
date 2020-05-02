package factory

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/eloylp/meerkat/config"
	"github.com/eloylp/meerkat/data"
	"github.com/eloylp/meerkat/store"
)

func NewDataFlowRegistry(cfg config.Config) (*data.FlowRegistry, error) {
	dfr := &data.FlowRegistry{}
	for _, r := range cfg.Resources {
		maxItems := 10
		dataStore := store.NewTimeLineStore(maxItems)
		fetcher := data.NewHTTPFetcher(&http.Client{})
		dataPump := data.NewDataPump(cfg.PollInterval, r, fetcher, dataStore)
		dfr.Add(data.NewDataFlow(uuid.New().String(), r, dataStore, dataPump))
	}
	return dfr, nil
}
