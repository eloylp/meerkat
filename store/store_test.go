// +build unit

package store_test

import (
	"bytes"
	"fmt"
	"github.com/eloylp/meerkat/store"
	"io"
	"io/ioutil"
	"testing"
)

func TestBufferedStore_Subscribe(t *testing.T) {
	s := populatedBufferedStore(t)
	_, uuid := s.Subscribe()
	if uuid == "" {
		t.Errorf("Ticket number must be above 0, got %v", uuid)
	}
}

func TestBufferedStore_Unsubscribe(t *testing.T) {
	s := populatedBufferedStore(t)
	_, _ = s.Subscribe()
	ch2, uuid2 := s.Subscribe()
	if err := s.Unsubscribe(uuid2); err != nil {
		t.Error(err)
	}

	expectedSubscribers := 1
	subscribersNumResult := s.Subscribers()
	if subscribersNumResult != expectedSubscribers {
		t.Errorf("Expected subscribers after unsubscribe is %v got %v", expectedSubscribers, subscribersNumResult)
	}
	var count int
	for range ch2 {
		count++
	}
	if count != 3 {
		t.Errorf("Exhausted channel must have last three elements, consumed %v", count)
	}
}

func TestBufferedStore_Unsubscribe_NotFound(t *testing.T) {
	s := populatedBufferedStore(t)
	_, uuid := s.Subscribe()
	if err := s.Unsubscribe(uuid); err != nil {
		t.Error(err)
	}
	err := s.Unsubscribe(uuid)
	switch err.(type) {
	case *store.NotFoundError:
		break
	default:
		t.Errorf("Expected not found error but got %v", err)
	}
}

func TestBufferedStore_Reset(t *testing.T) {
	s := populatedBufferedStore(t)
	s.Reset()
	listenCh, _ := s.Subscribe()
	if err := s.AddItem(bytes.NewReader([]byte("dd"))); err != nil {
		t.Fatal(err)
	}
	s.Reset()
	var count int
	for item := range listenCh {
		count++
		item, err := ioutil.ReadAll(item)
		if err != nil {
			t.Fatal(err)
		}
		expected := "dd"
		if expected != string(item) {
			t.Fatalf("Expected nil values after reset")
		}
	}
	if count != 1 {
		t.Fatal("Only one element was expected in channel")
	}
}
func TestBufferedStore(t *testing.T) {
	s := populatedBufferedStore(t)
	listenCh, _ := s.Subscribe()

	s.Reset()
	var dataCount int
	var itemCount int
	dataCount++
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
	s := populatedBufferedStore(t)
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
