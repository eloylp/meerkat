package datatest

import (
	"github.com/eloylp/kit/flow/fanout"
	"github.com/stretchr/testify/mock"
)

type StoreMock struct {
	mock.Mock
}

func NewStoreMock() *StoreMock {
	return &StoreMock{}
}

func (s *StoreMock) Add(elem interface{}) {
	s.Called(elem)
}

func (s *StoreMock) Subscribe() (<-chan *fanout.Slot, string, fanout.CancelFunc) { //nolint:gocritic
	args := s.Called()
	return args.Get(0).(chan *fanout.Slot), args.String(1), args.Get(2).(fanout.CancelFunc)
}

func (s *StoreMock) SubscribersLen() int {
	args := s.Called()
	return args.Int(0)
}

func (s *StoreMock) Unsubscribe(ticket string) error {
	args := s.Called(ticket)
	return args.Error(0)
}

func (s *StoreMock) Reset() {
	s.Called()
}
