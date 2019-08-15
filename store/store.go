package store

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Todo replace this with uuid (heres also a data race)
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func ticket() int {
	return seededRand.Int()
}

type Store interface {
	AddItem(r io.Reader) error
	Subscribe() (chan io.Reader, int)
	Subscribers() uint
	Unsubscribe(ticket int) error
	Length() uint
	Reset()
}

type subscriber struct {
	ch     chan io.Reader
	ticket int
}

type dataFrame struct {
	timeStamp time.Time
	length    uint64
	data      []byte
}

type timeLineStore struct {
	items              []*dataFrame
	subscribers        []subscriber
	maxItems           uint
	subscriberBuffSize uint64
	L                  sync.RWMutex
}

func (t *timeLineStore) Length() uint {
	t.L.RLock()
	defer t.L.RUnlock()
	return uint(len(t.items))
}

func (t *timeLineStore) Subscribers() uint {
	t.L.RLock()
	defer t.L.RUnlock()
	return uint(len(t.subscribers))
}

func NewTimeLineStore(maxItems uint) *timeLineStore {
	return &timeLineStore{maxItems: maxItems, subscriberBuffSize: 10}
}

func (t *timeLineStore) AddItem(r io.Reader) error {
	t.L.Lock()
	defer t.L.Unlock()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	frame := &dataFrame{
		timeStamp: time.Now(),
		length:    uint64(len(data)),
		data:      data,
	}

	t.store(frame)
	t.gcItems()
	t.publish(frame)
	return nil
}

func (t *timeLineStore) store(frame *dataFrame) {
	t.items = append(t.items, frame)
}

func (t *timeLineStore) gcItems() {
	l := uint(len(t.items))
	if l > t.maxItems {
		i := make([]*dataFrame, t.maxItems)
		e := t.items[1:l]
		copy(i, e)
		t.items = i
	}
}

func (t *timeLineStore) publish(df *dataFrame) {
	for _, s := range t.subscribers {
		s.ch <- bytes.NewReader(df.data)
	}
}

func (t *timeLineStore) Subscribe() (chan io.Reader, int) {
	t.L.Lock()
	defer t.L.Unlock()
	ch := make(chan io.Reader, t.subscriberBuffSize)
	ticket := ticket()
	t.subscribers = append(t.subscribers, subscriber{ch, ticket})
	return ch, ticket
}

func (t *timeLineStore) Unsubscribe(ticket int) error {
	if !t.exists(ticket) {
		return newNotFoundError(ticket)
	}
	t.L.Lock()
	defer t.L.Unlock()
	ec := t.estimateCapacity()
	newSubs := make([]subscriber, 0, ec)
	for _, s := range t.subscribers {
		if s.ticket == ticket {
			close(s.ch)
		} else {
			newSubs = append(newSubs, s)
		}
	}
	t.subscribers = newSubs
	log.Printf("Unsubscribed ticket from store %v, pending tickets %v", ticket, len(t.subscribers))
	return nil
}

func (t *timeLineStore) exists(ticket int) bool {
	t.L.RLock()
	defer t.L.RUnlock()
	for _, s := range t.subscribers {
		if s.ticket == ticket {
			return true
		}
	}
	return false
}

func (t *timeLineStore) estimateCapacity() int {
	cl := len(t.subscribers)
	ec := 0
	if cl > 0 {
		ec = cl - 1
	}
	return ec
}

func (t *timeLineStore) Reset() {
	t.L.Lock()
	defer t.L.Unlock()
	for _, s := range t.subscribers {
		close(s.ch)
	}
	t.subscribers = make([]subscriber, 1)
	t.items = make([]*dataFrame, t.maxItems)
}
