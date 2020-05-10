// +build unit

package store_test

import (
	"bytes"
	"fmt"
	"github.com/eloylp/meerkat/store/storetest"
	"io"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	"github.com/eloylp/meerkat/store"

	"github.com/stretchr/testify/assert"
)

func TestBufferedStore_Subscribe_ElementsAreSentToSubscribers(t *testing.T) {
	items := 3
	maxItems := 3
	maxSubsBuffSize := 10
	s := storetest.PopulatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	ch, _ := s.Subscribe()
	err := s.AddItem(bytes.NewReader([]byte("d3"))) // Extra (4th) data that is removed by excess
	assert.NoError(t, err)
	var chOk bool
	var itemR io.Reader
	for i := 0; i < items+1; i++ {
		itemR, chOk = <-ch
		item, rErr := ioutil.ReadAll(itemR)
		assert.NoError(t, rErr)
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
	s := storetest.PopulatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
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
	s := storetest.PopulatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	// Adds one extra subscriber for test hardening.
	_, _ = s.Subscribe()
	ch, uuid := s.Subscribe()
	err := s.Unsubscribe(uuid)
	assert.NoError(t, err)
	assert.Equal(t, 1, s.Subscribers())

	// exhaust channel
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
	s := storetest.PopulatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	err := s.Unsubscribe("A1234")
	assert.IsType(t, store.ErrSubscriberNotFound, err, "wanted store.ErrSubscriberNotFound got %T", err)
}

func TestBufferedStore_Reset(t *testing.T) {
	items := 0
	maxItems := 3
	maxSubsBuffSize := 10
	s := storetest.PopulatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	ch, _ := s.Subscribe()
	err := s.AddItem(bytes.NewReader([]byte("dd")))
	assert.NoError(t, err)
	s.Reset()
	assert.Equal(t, 0, s.Subscribers(), "no subscribers expected after reset")
	assert.Equal(t, 0, s.Length(), "no items expected after reset")
	// Check channel is closed after consumption
	<-ch
	_, ok := <-ch
	assert.False(t, ok)
}

func TestNewBufferedStore_AddItem_OldItemsClear(t *testing.T) {
	items := 3
	maxItems := 3
	maxSubsBuffSize := 10
	s := storetest.PopulatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	err := s.AddItem(bytes.NewReader([]byte("d4")))
	assert.NoError(t, err)
	want := 3
	got := s.Length()
	assert.Equal(t, want, got, "want %v resultant items got %v", want, got)
}

func TestBufferedStore_AddItem_NoActiveSubscriberDoesntBlock(t *testing.T) {
	items := 3 // "d + index" items (continue reading comments ...)
	maxItems := 3
	maxSubsBuffSize := 10
	s := storetest.PopulatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	ch, _ := s.Subscribe()      // this subscriber will try to block the entire system
	limitValueForBlocking := 11 // So we will exceed by one (limit value of maxSubsBuffSize)

	testEnd := make(chan struct{}, 1)
	go func() {
		for i := 0; i < limitValueForBlocking; i++ {
			// "dn + index" will mark new data segments that may override old ones in factory "d + index".
			err := s.AddItem(bytes.NewReader([]byte("dn" + strconv.Itoa(i))))
			assert.NoError(t, err)
		}
		testEnd <- struct{}{}
	}()
	select {
	case <-time.NewTimer(2 * time.Second).C: // Will break test if its blocking.
		t.Error("exceeded wait time. May subscribers are blocking the buffer")
	case <-testEnd:
		t.Log("successfully inserted items without active subscribers")
	}
	s.Reset() // this will close the subscriber channel, allowing us to follow with the next check.

	// Now check that all discarded elements in subscriber buffer are old, from the factory.
	// This will be done by checking all data in subscriber channels contains "dn + index".
	for e := range ch {
		content, err := ioutil.ReadAll(e)
		assert.NoError(t, err)
		assert.Contains(t, string(content), "dn",
			"The are data in subs channel that needs to be discarded")
	}
}
