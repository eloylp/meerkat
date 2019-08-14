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
	_, ticket := s.Subscribe()
	if ticket <= 0 {
		t.Errorf("Ticket number must be above 0, got %v", ticket)
	}
}

func TestTimeLineStore_Unsubscribe(t *testing.T) {
	s := populatedTimeLineStore(t)
	_, _ = s.Subscribe()
	ch2, ticket2 := s.Subscribe()
	if err := s.Unsubscribe(ticket2); err != nil {
		t.Error(err)
	}

	var expectedSubscribers uint = 1
	subscribersNumResult := s.SubscribersNum()
	if subscribersNumResult != expectedSubscribers {
		t.Errorf("Expected subscribers after unsubscribe is %v got %v", expectedSubscribers, subscribersNumResult)
	}
	_, ok := <-ch2
	if ok {
		t.Error("Channel is not closed after unsubscribe")
	}
}

func TestTimeLineStore_Unsubscribe_NotFound(t *testing.T) {
	s := populatedTimeLineStore(t)
	_, ticket := s.Subscribe()
	if err := s.Unsubscribe(ticket); err != nil {
		t.Error(err)
	}
	err := s.Unsubscribe(ticket)
	switch err.(type) {
	case *NotFoundError:
		break
	default:
		t.Errorf("Expected not found error but got %v", err)
	}
}

func TestTimeLineStore(t *testing.T) {

	s := populatedTimeLineStore(t)
	listenCh, _ := s.Subscribe()

	s.Reset()
	var count uint
	count++
	for item := range listenCh {
		expected := "d" + fmt.Sprint(count)
		item, err := ioutil.ReadAll(item)
		if err != nil {
			t.Fatal(err)
		}
		got := string(item)
		if expected != got {
			t.Fatalf("Error listening subscribed frames. Expected frame was %s but got %s", expected, got)
		}
		count++
	}

}

func populatedTimeLineStore(t *testing.T) *timeLineStore {
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
			go s.AddItem(bytes.NewReader([]byte("d")))
		}
	}
}
