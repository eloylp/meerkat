package config

import (
	"flag"
	"sync"
)

type config struct {
	Resource          string
	PollInterval      uint
	HTTPListenAddress string
}

var once sync.Once
var cfg *config

func C() *config {
	once.Do(func() {
		cfg = &config{}
		flag.StringVar(&cfg.Resource, "u", "", "The URL to recover frames from")
		flag.UintVar(&cfg.PollInterval, "i", 1, "The interval to fill the frame buffer")
		flag.StringVar(&cfg.HTTPListenAddress, "l", "0.0.0.0:3000", "Pass the http server listen address for serving results")
		flag.Parse()
	})
	return cfg
}

func validate(c *config) error {
	return nil
}
