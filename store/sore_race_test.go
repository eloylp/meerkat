// +build race

package store_test

import (
	"bytes"
	"testing"
	"time"
)

func TestStore_AddItem_supportsrace(t *testing.T) {
	s := populatedTimeLineStore(t)
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
