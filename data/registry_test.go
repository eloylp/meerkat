// +build unit

package data_test

import (
	"github.com/eloylp/kit/flow/fanout"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/eloylp/meerkat/data"
)

type storeMock struct {
	mock.Mock
}

func (s *storeMock) Add(elem interface{}) {
	s.Called(elem)
}

func (s *storeMock) Subscribe() (<-chan *fanout.Slot, string, fanout.CancelFunc) { //nolint:gocritic
	args := s.Called()
	return args.Get(0).(chan *fanout.Slot), args.String(1), args.Get(2).(fanout.CancelFunc)
}

func (s *storeMock) SubscribersLen() int {
	args := s.Called()
	return args.Int(0)
}

func (s *storeMock) Unsubscribe(ticket string) error {
	args := s.Called(ticket)
	return args.Error(0)
}

func (s *storeMock) Reset() {
	s.Called()
}

type pumpMock struct {
	mock.Mock
}

func (p *pumpMock) Start() {
	p.Called()
}

func TestDataFlowRegistry_FindStore(t *testing.T) {
	df1 := data.NewDataFlow("A1234", "", &storeMock{}, &pumpMock{})
	targetUID := "A12345"
	expectedStore := &storeMock{}
	df2 := data.NewDataFlow(targetUID, "", expectedStore, &pumpMock{})
	df3 := data.NewDataFlow("A123456", "", &storeMock{}, &pumpMock{})
	dfr := data.NewFlowRegistry([]*data.Flow{df1, df2, df3})
	result, err := dfr.FindStore(targetUID)
	assert.NoError(t, err)
	assert.Equal(t, expectedStore, result)
}
