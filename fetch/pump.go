package fetch

import (
	"io/ioutil"
	"log"
	"time"
)

type dataPump struct {
	interval uint
	url      string
	frames   chan []byte
	fetcher  Fetcher
}

func NewDataPump(interval uint, url string, frames chan []byte, fetcher Fetcher) *dataPump {
	return &dataPump{interval: interval, url: url, frames: frames, fetcher: fetcher}
}

func (dp *dataPump) Start() {
	for {
		time.Sleep(time.Duration(dp.interval) * time.Second)
		reader, err := dp.fetcher.Fetch(dp.url)
		if err != nil {
			log.Fatal(err)
		}
		d, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}
		dp.frames <- d
	}
}
