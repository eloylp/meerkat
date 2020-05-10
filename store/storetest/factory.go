package storetest

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/eloylp/meerkat/store"
)

// PopulatedBufferedStore facilitates construction of an store
// by accepting the number of items to be present, and
// store.BufferedStore needed constructor args.
func PopulatedBufferedStore(t testing.TB, items, maxItems, subscriberBuffSize int) *store.BufferedStore {
	s := store.NewBufferedStore(maxItems, subscriberBuffSize)
	for i := 0; i < items; i++ {
		data := "d" + strconv.Itoa(i)
		item := bytes.NewReader([]byte(data))
		if err := s.AddItem(item); err != nil {
			t.Fatal(err)
		}
	}
	return s
}
