package store

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"
)

type Store interface {
	AddItem(r io.Reader) error
	Subscribe() chan io.Reader
	Length() uint
	Close()
}

type dataFrame struct {
	timeStamp time.Time
	length    uint64
	data      []byte
}

type timeLineStore struct {
	items                 []*dataFrame
	subscribers           []chan io.Reader
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
		s <- bytes.NewReader(df.data)
	}
}
func (t *timeLineStore) Subscribe() chan io.Reader {
	ch := make(chan io.Reader, t.subscribersBufferSize)
	t.subscribers = append(t.subscribers, ch)
	return ch
}

func (t *timeLineStore) Close() {
	for _, s := range t.subscribers {
		close(s)
	}
}
