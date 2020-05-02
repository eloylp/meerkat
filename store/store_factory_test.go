// +build unit integration race

package store_test

import (
	"bytes"
	"github.com/eloylp/meerkat/store"
	"strconv"
	"testing"
)

// populated BufferedStore facilitates construction of an store
// by accepting the number of items to be present, and
// store.BufferedStore needed constructor args.
func populatedBufferedStore(t *testing.T, items, maxitems, subscriberBuffSize int) *store.BufferedStore {
	s := store.NewBufferedStore(maxitems, subscriberBuffSize)
	for i := 0; i <= items; i++ {
		data := "d" + strconv.Itoa(i)
		item := bytes.NewReader([]byte(data))
		if err := s.AddItem(item); err != nil {
			t.Fatal(err)
		}
	}
	return s
}
