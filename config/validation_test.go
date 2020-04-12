package config_test

import (
	"github.com/eloylp/meerkat/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		config      config.Config
		mustPass    bool
		description string
	}{
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
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			err := config.Validate(c.config)
			if c.mustPass {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
