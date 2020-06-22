package data

import (
	"github.com/eloylp/meerkat/elements"
)

type Flow struct {
	uUID      string
	resource  string
	dataStore elements.Store
	dataPump  elements.DataPump
}

func NewDataFlow(uuid, resource string, dataStore elements.Store, dataPump elements.DataPump) *Flow {
	return &Flow{uUID: uuid, resource: resource, dataStore: dataStore, dataPump: dataPump}
}

func (df *Flow) UUID() string {
	return df.uUID
}

func (df *Flow) Store() elements.Store {
	return df.dataStore
}
func (df *Flow) Start() {
	df.dataPump.Start()
}
