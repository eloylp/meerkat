// +build unit

package store_test

import (
	"bytes"
	"fmt"
	"github.com/eloylp/meerkat/store"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"testing"
)

func TestBufferedStore_Subscribe(t *testing.T) {
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

func TestBufferedStore_Unsubscribe_notfound(t *testing.T) {
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

func TestBufferedStore(t *testing.T) {
	items := 3
	maxItems := 3
	maxSubsBuffSize := 10
	s := populatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	listenCh, _ := s.Subscribe()

	s.Reset()
	var dataCount int
	var itemCount int
	for item := range listenCh {
		itemCount++
		expected := "d" + fmt.Sprint(dataCount)
		item, err := ioutil.ReadAll(item)
		if err != nil {
			t.Fatal(err)
		}
		got := string(item)
		if expected != got {
			t.Fatalf("Error listening subscribed frames. Expected frame was %s but got %s", expected, got)
		}
		dataCount++
	}
	if itemCount != 3 {
		t.Fatal("Not elements in channel !!")
	}
}

func TestNewBufferedStore_OldItemsClear(t *testing.T) {
	items := 3
	maxItems := 3
	maxSubsBuffSize := 10
	s := populatedBufferedStore(t, items, maxItems, maxSubsBuffSize)
	subs, _ := s.Subscribe()
	if err := s.AddItem(bytes.NewReader([]byte("d4"))); err != nil {
		t.Fatal(err)
	}
	expectedSize := 3
	size := s.Length()
	if size != 3 {
		t.Errorf("Expected resultant items is %v but %v obtained", expectedSize, size)
	}
	s.Reset()
	var lastItemR io.Reader
	for item := range subs {
		lastItemR = item
	}
	data, err := ioutil.ReadAll(lastItemR)
	if err != nil {
		t.Fatal(err)
	}
	lastItem := string(data)
	expectedLastItem := "d4"
	if lastItem != expectedLastItem {
		t.Errorf("Expected last item is %s but got %s", expectedLastItem, lastItem)
	}
}
