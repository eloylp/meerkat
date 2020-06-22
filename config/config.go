package config

import (
	"strings"
)

type Config struct {
	Resources         []string
	PollInterval      int
	HTTPListenAddress string
}

func ParseResources(p string) []string {
	p = strings.Replace(p, " ", "", -1)
	return strings.Split(p, ",")
}
