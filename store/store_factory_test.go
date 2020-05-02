// +build unit integration race

package store_test

import (
	"bytes"
	"github.com/eloylp/meerkat/store"
	"strconv"
	"testing"
)

// populatedBufferedStore facilitates construction of an store
// by accepting the number of items to be present and the max
// items that will accept.
func populatedBufferedStore(t *testing.T, items, maxitems int) *store.BufferedStore {
	s := store.NewBufferedStore(maxitems)
	for i := 0; i <= items; i++ {
		data := "d" + strconv.Itoa(i)
		item := bytes.NewReader([]byte(data))
		if err := s.AddItem(item); err != nil {
			t.Fatal(err)
		}
	}
	return s
}
