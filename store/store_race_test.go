// +build race

package store_test

import (
	"bytes"
	"github.com/eloylp/meerkat/store/storetest"
	"testing"
	"time"
)

func TestStore_AddItem_SupportsRace(t *testing.T) {
	s := storetest.PopulatedBufferedStore(t, 3, 3, 5)
	subs, _ := s.Subscribe()

	go func() {
		for {
			<-subs
		}
	}()
	timer := time.NewTimer(time.Second * 10)
loop:
	for {
		select {
		case <-timer.C:
			break loop
		default:
			go s.AddItem(bytes.NewReader([]byte("d"))) //nolint:errcheck
		}
	}
}
