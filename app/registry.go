package app

import (
	"errors"
	"fmt"
	"go-sentinel/fetch"
	"go-sentinel/store"
)

type DataFlowRegistry struct {
	Flows []*DataFlow
}

func (dfr *DataFlowRegistry) DataFlows() []*DataFlow {
	return dfr.Flows
}

func (dfr *DataFlowRegistry) Add(df *DataFlow) {
	dfr.Flows = append(dfr.Flows, df)
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
	DataPump  fetch.DataPump
}

func (df *DataFlow) Store() store.Store {
	return df.DataStore
}
func (df *DataFlow) Start() {
	df.DataPump.Start()
}
