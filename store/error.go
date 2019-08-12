package store

import (
	"fmt"
)

func newNotFoundError(ticket int) *NotFoundError {
	return &NotFoundError{ticket: ticket}
}

type NotFoundError struct {
	ticket int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Cannot found subscriber with ticket %v", e.ticket)
}
