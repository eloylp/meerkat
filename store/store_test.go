// +build unit

package store_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/eloylp/meerkat/store"

	"github.com/stretchr/testify/assert"
)

func TestBufferedStore_Subscribe_ElementsAreSentToSubscribers(t *testing.T) {
	items := 3
	maxItems := 3
	maxSubsBuffSize := 10
	s := populatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	ch, _ := s.Subscribe()
	err := s.AddItem(bytes.NewReader([]byte("d3"))) // Extra (4th) data that is removed by excess
	assert.NoError(t, err)
	var chOk bool
	var itemR io.Reader
	for i := 0; i < items+1; i++ {
		itemR, chOk = <-ch
		item, err := ioutil.ReadAll(itemR)
		assert.NoError(t, err)
		want := "d" + fmt.Sprint(i)
		got := string(item)
		// We check that all data is present, even removed by excess
		assert.Equal(t, want, got, "error listening subscribed frames,wanted frame was %q but got %q", want, got)
	}
	assert.True(t, chOk, "channel remains open for future consumes")
}

func TestBufferedStore_Subscribe_ReturnValues(t *testing.T) {
	items := 3
	maxItems := 3
	maxSubsBuffSize := 10
	s := populatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	ch, uuid := s.Subscribe()
	assert.NotEmpty(t, uuid, "want a uuid not an empty string")
	assert.NotNil(t, ch, "want a channel")
	_, ok := <-ch
	assert.True(t, ok, "want a open channel")
}

func TestBufferedStore_Unsubscribe(t *testing.T) {
	items := 3
	maxItems := 3
	maxSubsBuffSize := 10
	s := populatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	// Adds one extra subscriber for test hardening.
	_, _ = s.Subscribe()
	ch, uuid := s.Subscribe()
	err := s.Unsubscribe(uuid)
	assert.NoError(t, err)
	assert.Equal(t, 1, s.Subscribers())

	// exaust channel
	var count int
	for range ch {
		if count == 2 {
			break
		}
		count++
	}
	_, ok := <-ch
	assert.False(t, ok, "want channel closed after unsubscribe and consumed")
}

func TestBufferedStore_Unsubscribe_NotFound(t *testing.T) {
	items := 3
	maxItems := 3
	maxSubsBuffSize := 10
	s := populatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	err := s.Unsubscribe("A1234")
	assert.IsType(t, store.ErrSubscriberNotFound, err, "wanted store.ErrSubscriberNotFound got %T", err)
}

func TestBufferedStore_Reset(t *testing.T) {
	items := 0
	maxItems := 3
	maxSubsBuffSize := 10
	s := populatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	ch, _ := s.Subscribe()
	err := s.AddItem(bytes.NewReader([]byte("dd")))
	assert.NoError(t, err)
	s.Reset()
	assert.Equal(t, 0, s.Subscribers(), "no subscribers expected after reset")
	assert.Equal(t, 0, s.Length(), "no items expected after reset")
	// Check channel is closed afer concumption
	<-ch
	_, ok := <-ch
	assert.False(t, ok)
}

func TestNewBufferedStore_AddItem_OldItemsClear(t *testing.T) {
	items := 3
	maxItems := 3
	maxSubsBuffSize := 10
	s := populatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	err := s.AddItem(bytes.NewReader([]byte("d4")))
	assert.NoError(t, err)
	want := 3
	got := s.Length()
	assert.Equal(t, want, got, "want %v resultant items got %v", want, got)
}
