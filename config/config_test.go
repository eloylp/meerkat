package config

import (
	"testing"
)

type sample struct {
	config      *config
	mustPass    bool
	description string
}

func TestValidate(t *testing.T) {
	cases := validateCases()
	for _, c := range cases {
		err := validate(c.config)
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
		{config: &config{
			PollInterval:      1,
			Resource:          "http://example.com/camdump.jpeg",
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    true,
			description: "Must accept complete config with url as resource",
		},
		{config: &config{
			PollInterval:      1,
			Resource:          "/var/motion/dump.jpeg",
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    true,
			description: "Must accept complete config with path as resource",
		},
		{config: &config{
			PollInterval:      0,
			Resource:          "/var/motion/dump.jpeg",
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    false,
			description: "Must not accept a config with 0 poll interval setting",
		},
		{config: &config{
			PollInterval:      1,
			Resource:          "",
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    false,
			description: "Must not accept a config with no resource setting",
		},
		{config: &config{
			PollInterval:      1,
			Resource:          "/var/motion/dump.jpeg",
			HTTPListenAddress: "",
		},
			mustPass:    false,
			description: "Must not accept a config with no http listen address setting",
		},
	}
}
