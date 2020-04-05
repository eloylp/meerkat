package app

import (
	"errors"
	"fmt"
	"github.com/eloylp/meerkat/store"
)

type dataPump interface {
	Start()
}

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
		if wf.UUID == wfUid {
			return wf.DataStore, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Cannot find workflow %v", wfUid))
}

type DataFlow struct {
	UUID      string
	Resource  string
	DataStore store.Store
	DataPump  dataPump
}

func NewDataFlow(UUID string, resource string, dataStore store.Store, dataPump dataPump) *DataFlow {
	return &DataFlow{UUID: UUID, Resource: resource, DataStore: dataStore, DataPump: dataPump}
}

func (df *DataFlow) Store() store.Store {
	return df.DataStore
}
func (df *DataFlow) Start() {
	df.DataPump.Start()
}
