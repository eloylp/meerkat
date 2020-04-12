package config_test

import (
	"github.com/eloylp/meerkat/config"
	"testing"
)

type sample struct {
	config      config.Config
	mustPass    bool
	description string
}

func TestValidate(t *testing.T) {
	cases := validateCases()
	for _, c := range cases {
		err := config.Validate(c.config)
		if c.mustPass && err != nil {
			t.Errorf("Pass expectation: %s", c.description)
		}
		if !c.mustPass && err == nil {
			t.Errorf("Fail expectation: %s", c.description)
		}
	}
}

func validateCases() []sample {
	return []sample{
		{config: config.Config{
			PollInterval:      1,
			Resources:         []string{"http://example.com/camdump.jpeg"},
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    true,
			description: "Must accept complete config with url as resource",
		},
		{config: config.Config{
			PollInterval:      1,
			Resources:         []string{"/var/motion/dump.jpeg"},
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    true,
			description: "Must accept complete config with path as resource",
		},
		{config: config.Config{
			PollInterval:      0,
			Resources:         []string{"/var/motion/dump.jpeg"},
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    false,
			description: "Must not accept a config with 0 poll interval setting",
		},
		{config: config.Config{
			PollInterval:      1,
			Resources:         []string{""},
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    false,
			description: "Must not accept a config with no resource setting",
		},
		{config: config.Config{
			PollInterval:      1,
			Resources:         []string{"/var/motion/dump.jpeg"},
			HTTPListenAddress: "",
		},
			mustPass:    false,
			description: "Must not accept a config with no http listen address setting",
		},
	}
}
