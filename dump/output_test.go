package dump_test

import (
	"bytes"
	"go-sentinel/dump"
	"regexp"
	"testing"
)

func TestMJPEGDumper_Boundary(t *testing.T) {
	const expectedBoundaryLength = 60
	w := new(bytes.Buffer)
	d := dump.NewMJPEGDumper(w)
	bLength := len(d.Boundary())
	if bLength != expectedBoundaryLength {
		t.Errorf("Expected boundary is %v and was %v", expectedBoundaryLength, bLength)
	}
}

func TestNewMJPEGDumper(t *testing.T) {
	w := new(bytes.Buffer)
	d := dump.NewMJPEGDumper(w)
	data := []byte("Data")
	dataReader := bytes.NewReader(data)

	if err := d.DumpPart(dataReader); err != nil {
		t.Error(err)
	}
	expectedPart := "--" + d.Boundary() + " Content-Type: image/jpeg  Data"
	writePartString := string(w.Bytes())
	re := regexp.MustCompile(`\r?\n`)
	writePartSanitized := re.ReplaceAllString(writePartString, " ")
	if writePartSanitized != expectedPart {
		t.Errorf("Expected string part is \n %s \ngot\n %s", expectedPart, writePartSanitized)
	}
}
