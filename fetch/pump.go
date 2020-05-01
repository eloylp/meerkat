package fetch

import (
	"io"
	"log"
	"time"

	"github.com/eloylp/meerkat/elements"
)

type fetcher interface {
	Fetch(res string) (io.Reader, error)
}

type DataPump struct {
	interval int
	url      string
	fetcher  fetcher
	store    elements.Store
}

func NewDataPump(interval int, url string, fetcher fetcher, store elements.Store) *DataPump {
	return &DataPump{interval: interval, url: url, fetcher: fetcher, store: store}
}

func (dp *DataPump) Start() {
	for {
		time.Sleep(time.Duration(dp.interval) * time.Second)
		reader, err := dp.fetcher.Fetch(dp.url)
		if err != nil {
			log.Fatal(err)
		}
		if err := dp.store.AddItem(reader); err != nil {
			log.Fatal(err)
		}
	}
}
