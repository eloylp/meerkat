// +build unit integration race

package store_test

import (
	"bytes"
	"github.com/eloylp/meerkat/store"
	"io"
	"testing"
)

func populatedTimeLineStore(t *testing.T) *store.TimeLineStore {
	samples := []io.Reader{
		bytes.NewReader([]byte("d1")),
		bytes.NewReader([]byte("d2")),
		bytes.NewReader([]byte("d3")),
	}
	s := store.NewTimeLineStore(3)
	for _, sample := range samples {
		if err := s.AddItem(sample); err != nil {
			t.Fatal(err)
		}
	}
	return s
}
