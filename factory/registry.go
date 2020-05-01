package factory

import (
	"net/http"

	"github.com/eloylp/meerkat/config"
	"github.com/eloylp/meerkat/fetch"
	"github.com/eloylp/meerkat/flow"
	"github.com/eloylp/meerkat/store"
	"github.com/eloylp/meerkat/unique"
)

func NewDataFlowRegistry(cfg config.Config) (*flow.DataFlowRegistry, error) {
	dfr := &flow.DataFlowRegistry{}
	for _, r := range cfg.Resources {
		maxItems := 10
		dataStore := store.NewTimeLineStore(maxItems)
		fetcher := fetch.NewHTTPFetcher(&http.Client{})
		dataPump := fetch.NewDataPump(cfg.PollInterval, r, fetcher, dataStore)
		dfr.Add(flow.NewDataFlow(unique.UUID4(), r, dataStore, dataPump))
	}
	return dfr, nil
}
