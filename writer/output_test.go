// +build unit

package writer_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/meerkat/writer"
)

func TestMJPEGDumper_Boundary(t *testing.T) {
	const wantedBoundaryLength = 60
	w := bytes.NewBuffer(nil)
	d := writer.NewMJPEGWriter(w)
	gotBoundaryLength := len(d.Boundary())
	assert.Equal(t, wantedBoundaryLength, gotBoundaryLength)
}

func TestNewMJPEGDumper(t *testing.T) {
	w := bytes.NewBuffer(nil)
	d := writer.NewMJPEGWriter(w)
	data := []byte("Data")
	reader := bytes.NewReader(data)
	if err := d.WritePart(reader); err != nil {
		t.Error(err)
	}
	wantedPart := "--" + d.Boundary() + " Content-Type: image/jpeg  Data"
	writePartString := w.String()
	re := regexp.MustCompile(`\r?\n`)
	gotSanitizedPart := re.ReplaceAllString(writePartString, " ")
	assert.Equal(t, wantedPart, gotSanitizedPart, "expected string part is \n %s \ngot\n %s", wantedPart, gotSanitizedPart)
}
