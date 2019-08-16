package store

import (
	"fmt"
)

func newNotFoundError(uuid string) *NotFoundError {
	return &NotFoundError{uuid: uuid}
}

type NotFoundError struct {
	uuid string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Cannot found subscriber with UUID %v", e.uuid)
}
