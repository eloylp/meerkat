// +build integration

package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadArgument_URL_ShowsHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	args := []string{Name, "-u", "badAddress"}
	err := run(args, &stdout, &stderr)
	assert.Error(t, err)
	assert.Contains(t, stderr, "meerkat")   // contains program name
	assert.Contains(t, stderr, "-u string") // contains arg help
}
