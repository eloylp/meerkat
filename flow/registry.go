package flow

import (
	"errors"
	"fmt"
	"github.com/eloylp/meerkat/store"
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

func (dfr *DataFlowRegistry) FindStore(wfUid string) (store.Store, error) {
	for _, wf := range dfr.DataFlows() {
		if wf.UUID() == wfUid {
			return wf.dataStore, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Cannot find workflow %v", wfUid))
}
