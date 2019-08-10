package store

import (
	"time"
)

type Store interface {
	Start(input chan []byte) chan error
	Add(data []byte) error
	Shutdown() error
	CleanUp() error
}

type DataFrame struct {
	timeStamp time.Time
	data      []byte
}

type timeLineStore struct {
	store []DataFrame
}
