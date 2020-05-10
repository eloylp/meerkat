// Package elements represents all the moving parts of this project.
// They are intended to be used in multiple packages so they are all
// grouped here.
package elements

import (
	"io"
)

// DataPump represent something that is able to pull data
// from a source like HTTP, TCP, UNIX SOCKET ...
type DataPump interface {
	Start()
}

// Store can be treated as a data buffer that will fan out data to multiple
// clients at the same time.
type Store interface {
	// AddItem will add to the store a chunk of data
	AddItem(r io.Reader) error
	// Subscribe will return an output channel that will be notified
	// when more data arrives to store. It will also return the associated UUID
	// for later Unsubscribe operation.
	//
	// Depending on the implementation it will dump to the returned channel
	// the entire available data layer until the last frame is reached.
	// If you are not actively consuming this channel, but data continues
	// arriving from the store, the oldest element will be replaced by the
	// new one in the channel an so on.
	Subscribe() (<-chan io.Reader, string)
	// Subscribers will return the current number of active subscribers
	Subscribers() int
	// Unsubscribe will require the UUID obtained via a Subscribe operation to
	// properly clear all resources.
	Unsubscribe(uuid string) error
	// Length will return the length of the stored data.
	Length() int
	// Reset will clear all the data and subscribers and start again.
	Reset()
}

// Dumper its responsible of writing each store reader to the desired
// output.
type Dumper interface {
	WritePart(data io.Reader) error
}
