package config

import (
	"fmt"
)

var validators = []validator{
	&ResourceValidator{},
	&PollIntervalValidator{},
	&HTTPListenAddressValidator{},
}

type validator interface {
	validate(c Config) error
}

type ResourceValidator struct{}

func (ResourceValidator) validate(c Config) error {
	for _, r := range c.Resources {
		if err := stringNotZero("Resource", r); err != nil {
			return err
		}
	}
	return nil
}

type PollIntervalValidator struct{}

func (PollIntervalValidator) validate(c Config) error {
	return notZero("Polling interval", c.PollInterval)
}

type HTTPListenAddressValidator struct{}

func (HTTPListenAddressValidator) validate(c Config) error {
	return stringNotZero("HTTP listen address", c.HTTPListenAddress)
}

func stringNotZero(k, v string) error {
	if v == "" {
		return fmt.Errorf("%s cannot be empty", k)
	}
	return nil
}

func notZero(k string, v int) error {
	if v == 0 {
		return fmt.Errorf("%s cannot be empty", k)
	}
	return nil
}

func Validate(c Config) error {
	for _, v := range validators {
		if err := v.validate(c); err != nil {
			return err
		}
	}
	return nil
}
