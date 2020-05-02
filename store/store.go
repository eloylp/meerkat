package store

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"sync"
	"time"

	guuid "github.com/google/uuid"
)

type subscriber struct {
	ch   chan io.Reader
	UUID string
}

type dataFrame struct {
	timeStamp time.Time
	length    int
	data      []byte
}

type BufferedStore struct {
	items              []*dataFrame
	subscribers        []subscriber
	maxItems           int
	subscriberBuffSize int
	L                  sync.RWMutex
}

func (t *BufferedStore) Length() int {
	t.L.RLock()
	defer t.L.RUnlock()
	return len(t.items)
}

func (t *BufferedStore) Subscribers() int {
	t.L.RLock()
	defer t.L.RUnlock()
	return len(t.subscribers)
}

func NewBufferedStore(maxItems, subscriberBuffSize int) *BufferedStore {
	return &BufferedStore{maxItems: maxItems, subscriberBuffSize: subscriberBuffSize}
}

func (t *BufferedStore) AddItem(r io.Reader) error {
	t.L.Lock()
	defer t.L.Unlock()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	frame := &dataFrame{
		timeStamp: time.Now(),
		length:    len(data),
		data:      data,
	}

	t.store(frame)
	t.gcItems()
	t.publish(frame)
	return nil
}

func (t *BufferedStore) store(frame *dataFrame) {
	t.items = append(t.items, frame)
}

func (t *BufferedStore) gcItems() {
	l := len(t.items)
	if l > t.maxItems {
		i := make([]*dataFrame, t.maxItems)
		e := t.items[1:l]
		copy(i, e)
		t.items = i
	}
}

func (t *BufferedStore) publish(df *dataFrame) {
	for _, s := range t.subscribers {
		s.ch <- bytes.NewReader(df.data)
	}
}

func (t *BufferedStore) Subscribe() (<-chan io.Reader, string) {
	t.L.Lock()
	defer t.L.Unlock()
	ch := make(chan io.Reader, t.subscriberBuffSize)
	uuid := guuid.New().String()
	t.subscribers = append(t.subscribers, subscriber{ch, uuid})
	for _, frame := range t.items {
		if frame != nil {
			ch <- bytes.NewReader(frame.data)
		}
	}
	return ch, uuid
}

func (t *BufferedStore) Unsubscribe(uuid string) error {
	if !t.exists(uuid) {
		return ErrSubscriberNotFound
	}
	t.L.Lock()
	defer t.L.Unlock()
	ec := t.estimateCapacity()
	newSubs := make([]subscriber, 0, ec)
	for _, s := range t.subscribers {
		if s.UUID == uuid {
			close(s.ch)
		} else {
			newSubs = append(newSubs, s)
		}
	}
	t.subscribers = newSubs
	log.Printf("Unsubscribed UUID from store %v, pending uuids %v", uuid, len(t.subscribers))
	return nil
}

func (t *BufferedStore) exists(uuid string) bool {
	t.L.RLock()
	defer t.L.RUnlock()
	for _, s := range t.subscribers {
		if s.UUID == uuid {
			return true
		}
	}
	return false
}

func (t *BufferedStore) estimateCapacity() int {
	cl := len(t.subscribers)
	ec := 0
	if cl > 0 {
		ec = cl - 1
	}
	return ec
}

func (t *BufferedStore) Reset() {
	t.L.Lock()
	defer t.L.Unlock()
	for _, s := range t.subscribers {
		close(s.ch)
	}
	t.subscribers = nil
	t.items = nil
}
