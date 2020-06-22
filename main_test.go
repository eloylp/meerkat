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
	args := []string{Name}
	err := run(args, &stdout, &stderr)
	gotStderr := stderr.String()
	assert.Error(t, err)
	assert.Contains(t, gotStderr, "-u string") // contains arg help
}
