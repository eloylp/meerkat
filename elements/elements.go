// Package elements represents all the moving parts of this project.
// They are intended to be used in multiple packages so they are all
// grouped here.
package elements

import (
	"io"

	"github.com/eloylp/kit/flow/fanout"
)

// DataPump represent something that is able to pull data
// from a source like HTTP, TCP, UNIX SOCKET ...
type DataPump interface {
	Start()
}

// Store can be treated as a data buffer that will fan out data to multiple
// clients at the same time.
type Store interface {
	// Add will add to the store a chunk of data
	Add(elem interface{})
	// Subscribe will return an output channel that will be notified
	// when more data arrives to store. It will also return a cancelFunc
	// and  associated UUID for later Unsubscribe operation
	//
	// If you are not actively consuming this channel, but data continues
	// arriving from the store, the oldest elements will be discarded.
	Subscribe() (<-chan *fanout.Slot, string, fanout.CancelFunc)
	// SubscribersLen will return the current number of active subscribers
	SubscribersLen() int
	// Unsubscribe will require the UUID obtained via a Subscribe operation to
	// properly clear all resources.
	Unsubscribe(uuid string) error
	// Reset will clear all the data and subscribers and start again.
	Reset()
}

// Dumper its responsible of writing each store reader to the desired
// output.
type Dumper interface {
	WritePart(data io.Reader) error
}
