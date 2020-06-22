package data

import (
	"fmt"

	"github.com/eloylp/meerkat/elements"
)

type FlowRegistry struct {
	flows []elements.Flow
}

func NewFlowRegistry(flows []elements.Flow) *FlowRegistry {
	return &FlowRegistry{flows: flows}
}

func (dfr *FlowRegistry) Flows() []elements.Flow {
	return dfr.flows
}

func (dfr *FlowRegistry) Add(df elements.Flow) {
	dfr.flows = append(dfr.flows, df)
}

func (dfr *FlowRegistry) FindStore(wfUID string) (elements.Store, error) {
	for _, wf := range dfr.Flows() {
		if wf.UUID() == wfUID {
			return wf.Store(), nil
		}
	}
	return nil, fmt.Errorf("cannot find worklow %v", wfUID)
}
