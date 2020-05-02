package data

import (
	"fmt"

	"github.com/eloylp/meerkat/elements"
)

type FlowRegistry struct {
	flows []*Flow
}

func NewFlowRegistry(flows []*Flow) *FlowRegistry {
	return &FlowRegistry{flows: flows}
}

func (dfr *FlowRegistry) Flows() []*Flow {
	return dfr.flows
}

func (dfr *FlowRegistry) Add(df *Flow) {
	dfr.flows = append(dfr.flows, df)
}

func (dfr *FlowRegistry) FindStore(wfUID string) (elements.Store, error) {
	for _, wf := range dfr.Flows() {
		if wf.UUID() == wfUID {
			return wf.dataStore, nil
		}
	}
	return nil, fmt.Errorf("cannot find worklow %v", wfUID)
}
