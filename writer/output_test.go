// +build unit

package writer_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/eloylp/meerkat/writer"
)

func TestMJPEGDumper_Boundary(t *testing.T) {
	const expectedBoundaryLength = 60
	w := new(bytes.Buffer)
	d := writer.NewMJPEGWriter(w)
	bLength := len(d.Boundary())
	if bLength != expectedBoundaryLength {
		t.Errorf("expected boundary is %v and was %v", expectedBoundaryLength, bLength)
	}
}

func TestNewMJPEGDumper(t *testing.T) {
	w := new(bytes.Buffer)
	d := writer.NewMJPEGWriter(w)
	data := []byte("Data")
	dataReader := bytes.NewReader(data)

	if err := d.WritePart(dataReader); err != nil {
		t.Error(err)
	}
	expectedPart := "--" + d.Boundary() + " Content-Type: image/jpeg  Data"
	writePartString := w.String()
	re := regexp.MustCompile(`\r?\n`)
	writePartSanitized := re.ReplaceAllString(writePartString, " ")
	if writePartSanitized != expectedPart {
		t.Errorf("expected string part is \n %s \ngot\n %s", expectedPart, writePartSanitized)
	}
}
