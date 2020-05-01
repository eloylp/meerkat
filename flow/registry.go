package flow

import (
	"fmt"

	"github.com/eloylp/meerkat/elements"
)

type DataFlowRegistry struct {
	flows []*DataFlow
}

func NewDataFlowRegistry(flows []*DataFlow) *DataFlowRegistry {
	return &DataFlowRegistry{flows: flows}
}

func (dfr *DataFlowRegistry) DataFlows() []*DataFlow {
	return dfr.flows
}

func (dfr *DataFlowRegistry) Add(df *DataFlow) {
	dfr.flows = append(dfr.flows, df)
}

func (dfr *DataFlowRegistry) FindStore(wfUID string) (elements.Store, error) {
	for _, wf := range dfr.DataFlows() {
		if wf.UUID() == wfUID {
			return wf.dataStore, nil
		}
	}
	return nil, fmt.Errorf("cannot find worklow %v", wfUID)
}
