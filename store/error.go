package store

import (
	"errors"
)

var (
	ErrSubscriberNotFound = errors.New("store: subscriber not found")
)
