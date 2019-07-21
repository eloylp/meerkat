package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
)

type Config struct {
	Resource          string
	PollInterval      uint
	HTTPListenAddress string
}

var once sync.Once
var cfg Config
var validators []validator

func init() {
	validators = []validator{
		&ResourceValidator{},
		&PollIntervalValidator{},
		&HTTPListenAddressValidator{},
	}
}

func C() Config {
	once.Do(func() {
		cfg = Config{}
		flag.StringVar(&cfg.Resource, "u", "", "The URL to recover frames from")
		flag.UintVar(&cfg.PollInterval, "i", 1, "The interval to fill the frame buffer")
		flag.StringVar(&cfg.HTTPListenAddress, "l", "0.0.0.0:3000", "Pass the http server listen address for serving results")
		flag.Parse()
		if err := validate(cfg); err != nil {
			flag.PrintDefaults()
			os.Exit(1)
		}
	})
	return cfg
}

type validator interface {
	validate(c Config) error
}

type ResourceValidator struct{}

func (ResourceValidator) validate(c Config) error {
	return stringNotZero("Resource", c.Resource)
}

type PollIntervalValidator struct{}

func (PollIntervalValidator) validate(c Config) error {
	return uintNotZero("Polling interval", c.PollInterval)
}

type HTTPListenAddressValidator struct{}

func (HTTPListenAddressValidator) validate(c Config) error {
	return stringNotZero("HTTP listen address", c.HTTPListenAddress)
}

func stringNotZero(k string, v string) error {
	if v == "" {
		return errors.New(fmt.Sprintf("%s cannot be empty", k))
	}
	return nil
}

func uintNotZero(k string, v uint) error {
	if v == 0 {
		return errors.New(fmt.Sprintf("%s cannot be empty", k))
	}
	return nil
}

func validate(c Config) error {
	for _, v := range validators {
		if err := v.validate(c); err != nil {
			return err
		}
	}
	return nil
}
