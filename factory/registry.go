package factory

import (
	"net/http"
	"time"

	"github.com/eloylp/kit/flow/fanout"

	"github.com/google/uuid"

	"github.com/eloylp/meerkat/config"
	"github.com/eloylp/meerkat/data"
)

func NewDataFlowRegistry(cfg config.Config) (*data.FlowRegistry, error) {
	dfr := &data.FlowRegistry{}
	for _, r := range cfg.Resources {
		buffLen := 10 // todo get from config ?
		fo := fanout.NewBufferedFanOut(buffLen, time.Now)
		fetcher := data.NewHTTPFetcher(&http.Client{})
		dataPump := data.NewDataPump(cfg.PollInterval, r, fetcher, fo)
		dfr.Add(data.NewDataFlow(uuid.New().String(), r, fo, dataPump))
	}
	return dfr, nil
}
