package config

import (
	"flag"
	"os"
	"strings"
	"sync"
)

type Config struct {
	Resources         []string
	PollInterval      int
	HTTPListenAddress string
}

var once sync.Once
var cfg Config
var validators []validator

func C() Config {
	once.Do(func() {
		cfg = Config{}
		var resources string
		flag.StringVar(&resources, "u", "", "The comma separated URLS to recover frames from")
		flag.IntVar(&cfg.PollInterval, "i", 1, "The interval to fill the frame buffer")
		flag.StringVar(&cfg.HTTPListenAddress, "l", "0.0.0.0:3000", "Pass the http server listen address for serving results")
		flag.Parse()
		cfg.Resources = parseResources(resources)
		if err := Validate(cfg); err != nil {
			flag.PrintDefaults()
			os.Exit(1)
		}
	})
	return cfg
}

func parseResources(p string) []string {
	p = strings.Replace(p, " ", "", -1)
	return strings.Split(p, ",")
}
