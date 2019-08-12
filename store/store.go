package store

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func ticket() int {
	return seededRand.Int()
}

type Store interface {
	AddItem(r io.Reader) error
	Subscribe() (chan io.Reader, int)
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
	items                 []*dataFrame
	subscribers           []subscriber
	maxItems              uint
	maxItemLength         uint64
	subscribersBufferSize uint64
}

func (t *timeLineStore) Length() uint {
	return uint(len(t.items))
}

func NewTimeLineStore(maxItems uint, maxItemLength uint64) *timeLineStore {
	return &timeLineStore{maxItems: maxItems, maxItemLength: maxItemLength, subscribersBufferSize: 10}
}

func (t *timeLineStore) AddItem(r io.Reader) error {

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
	ch := make(chan io.Reader, t.subscribersBufferSize)
	ticket := ticket()
	t.subscribers = append(t.subscribers, subscriber{ch, ticket})
	return ch, ticket
}

func (t *timeLineStore) Unsubscribe(ticket int) error {
	l := len(t.subscribers)
	el := 0
	if l > 0 {
		el = l - 1
	}
	newSubs := make([]subscriber, el)
	var found bool
	for _, s := range t.subscribers {
		if s.ticket != ticket {
			newSubs = append(newSubs, s)
		} else {
			found = true
			close(s.ch)
		}
	}
	if found {
		t.subscribers = newSubs
		return nil
	}
	return newNotFoundError(ticket)
}

func (t *timeLineStore) Reset() {
	for _, s := range t.subscribers {
		close(s.ch)
	}
	t.subscribers = make([]subscriber, 1)
	t.items = make([]*dataFrame, t.maxItems)
}
