package config

import (
	"reflect"
	"testing"
)

type sample struct {
	config      Config
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
		{config: Config{
			PollInterval:      1,
			Resources:         []string{"http://example.com/camdump.jpeg"},
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    true,
			description: "Must accept complete config with url as resource",
		},
		{config: Config{
			PollInterval:      1,
			Resources:         []string{"/var/motion/dump.jpeg"},
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    true,
			description: "Must accept complete config with path as resource",
		},
		{config: Config{
			PollInterval:      0,
			Resources:         []string{"/var/motion/dump.jpeg"},
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    false,
			description: "Must not accept a config with 0 poll interval setting",
		},
		{config: Config{
			PollInterval:      1,
			Resources:         []string{""},
			HTTPListenAddress: "0.0.0.0:8080",
		},
			mustPass:    false,
			description: "Must not accept a config with no resource setting",
		},
		{config: Config{
			PollInterval:      1,
			Resources:         []string{"/var/motion/dump.jpeg"},
			HTTPListenAddress: "",
		},
			mustPass:    false,
			description: "Must not accept a config with no http listen address setting",
		},
	}
}

func Test_parseResources(t *testing.T) {

	type sample struct {
		Input       string
		Expected    []string
		Description string
	}

	s := []sample{
		{
			Input:       "http://example.com/image.jpg",
			Expected:    []string{"http://example.com/image.jpg"},
			Description: "Must accept single URL resource",
		},
		{
			Input:       "http://example.com/image.jpg,http://example.com/image2.jpg",
			Expected:    []string{"http://example.com/image.jpg", "http://example.com/image2.jpg"},
			Description: "Must accept multiple URL resource",
		},
		{
			Input:       "  http://example.com/image.jpg,  http://example.com/image2.jpg  ",
			Expected:    []string{"http://example.com/image.jpg", "http://example.com/image2.jpg"},
			Description: "Must accept multiple URL resource",
		},
	}

	for _, c := range s {
		result := parseResources(c.Input)
		if !reflect.DeepEqual(result, c.Expected) {
			t.Errorf("Cannot ensure case '%s', expected output is %v but got %v", c.Description, c.Expected, result)
		}
	}
}
