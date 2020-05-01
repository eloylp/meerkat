// +build unit

package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseResources(t *testing.T) {
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
			result := parseResources(c.Input)
			assert.Equal(t, c.Expected, result, "Cannot ensure case '%s', expected output is %v but got %v",
				c.Description, c.Expected, result)
		})
	}
}
