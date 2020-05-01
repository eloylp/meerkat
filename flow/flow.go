package flow

import (
	"github.com/eloylp/meerkat/elements"
)

type DataFlow struct {
	uUID      string
	resource  string
	dataStore elements.Store
	dataPump  elements.DataPump
}

func NewDataFlow(uuid string, resource string, dataStore elements.Store, dataPump elements.DataPump) *DataFlow {
	return &DataFlow{uUID: uuid, resource: resource, dataStore: dataStore, dataPump: dataPump}
}

func (df *DataFlow) UUID() string {
	return df.uUID
}

func (df *DataFlow) Store() elements.Store {
	return df.dataStore
}
func (df *DataFlow) Start() {
	df.dataPump.Start()
}
