package config

import (
	"reflect"
	"testing"
)

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
