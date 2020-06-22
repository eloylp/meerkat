// +build unit

package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/meerkat/config"
)

func Test_ParseResources(t *testing.T) {
	s := []struct {
		Input       string
		Expected    []string
		Description string
	}{
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
		t.Run(c.Description, func(t *testing.T) {
			result := config.ParseResources(c.Input)
			assert.Equal(t, c.Expected, result, "Cannot ensure case '%s', expected output is %v but got %v",
				c.Description, c.Expected, result)
		})
	}
}
