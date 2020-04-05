package flow

import (
	"github.com/eloylp/meerkat/store"
)

type dataPump interface {
	Start()
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
