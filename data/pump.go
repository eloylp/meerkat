package data

import (
	"io"
	"log"
	"time"

	"github.com/eloylp/meerkat/elements"
)

type fetcher interface {
	Fetch(res string) (io.Reader, error)
}

type Pump struct {
	interval int
	url      string
	fetcher  fetcher
	store    elements.Store
}

func NewDataPump(interval int, url string, fetcher fetcher, store elements.Store) *Pump {
	return &Pump{interval: interval, url: url, fetcher: fetcher, store: store}
}

func (dp *Pump) Start() {
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
