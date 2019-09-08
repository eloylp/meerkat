package app

import (
	"errors"
	"fmt"
	"github.com/eloylp/meerkat/store"
)

type dataPump interface {
	Start()
}

type dataFlowRegistry struct {
	Flows []*dataFlow
}

func (dfr *dataFlowRegistry) DataFlows() []*dataFlow {
	return dfr.Flows
}

func (dfr *dataFlowRegistry) Add(df *dataFlow) {
	dfr.Flows = append(dfr.Flows, df)
}

func (dfr *dataFlowRegistry) FindStore(wfUid string) (store.Store, error) {
	for _, wf := range dfr.DataFlows() {
		if wf.UUID == wfUid {
			return wf.DataStore, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Cannot find workflow %v", wfUid))
}

type dataFlow struct {
	UUID      string
	Resource  string
	DataStore store.Store
	DataPump  dataPump
}

func (df *dataFlow) Store() store.Store {
	return df.DataStore
}
func (df *dataFlow) Start() {
	df.DataPump.Start()
}
