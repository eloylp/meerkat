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

type dataPump struct {
	interval int
	url      string
	fetcher  fetcher
	store    elements.Store
}

func NewDataPump(interval int, url string, fetcher fetcher, store elements.Store) *dataPump {
	return &dataPump{interval: interval, url: url, fetcher: fetcher, store: store}
}

func (dp *dataPump) Start() {
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
