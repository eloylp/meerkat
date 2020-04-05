package flow

import (
	"github.com/eloylp/meerkat/store"
)

type dataPump interface {
	Start()
}

type DataFlow struct {
	uUID      string
	resource  string
	dataStore store.Store
	dataPump  dataPump
}

func NewDataFlow(UUID string, resource string, dataStore store.Store, dataPump dataPump) *DataFlow {
	return &DataFlow{uUID: UUID, resource: resource, dataStore: dataStore, dataPump: dataPump}
}

func (df *DataFlow) UUID() string {
	return df.uUID
}

func (df *DataFlow) Store() store.Store {
	return df.dataStore
}
func (df *DataFlow) Start() {
	df.dataPump.Start()
}
