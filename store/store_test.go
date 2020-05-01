// +build unit

package store

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
	"time"
)

func TestTimeLineStore_Subscribe(t *testing.T) {
	s := populatedTimeLineStore(t)
	_, uuid := s.Subscribe()
	if uuid == "" {
		t.Errorf("Ticket number must be above 0, got %v", uuid)
	}
}

func TestTimeLineStore_Unsubscribe(t *testing.T) {
	s := populatedTimeLineStore(t)
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

func TestTimeLineStore_Unsubscribe_NotFound(t *testing.T) {
	s := populatedTimeLineStore(t)
	_, uuid := s.Subscribe()
	if err := s.Unsubscribe(uuid); err != nil {
		t.Error(err)
	}
	err := s.Unsubscribe(uuid)
	switch err.(type) {
	case *NotFoundError:
		break
	default:
		t.Errorf("Expected not found error but got %v", err)
	}
}

func TestTimeLineStore_Reset(t *testing.T) {
	s := populatedTimeLineStore(t)
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
func TestTimeLineStore(t *testing.T) {
	s := populatedTimeLineStore(t)
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

func populatedTimeLineStore(t *testing.T) *TimeLineStore {
	samples := []io.Reader{
		bytes.NewReader([]byte("d1")),
		bytes.NewReader([]byte("d2")),
		bytes.NewReader([]byte("d3")),
	}
	s := NewTimeLineStore(3)
	for _, sample := range samples {
		if err := s.AddItem(sample); err != nil {
			t.Fatal(err)
		}
	}
	return s
}

func TestNewTimeLineStore_OldItemsClear(t *testing.T) {
	s := populatedTimeLineStore(t)
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

func TestNewTimeLineStore_DataRace(t *testing.T) {
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
