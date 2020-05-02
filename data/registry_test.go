// +build unit

package data_test

import (
	"github.com/eloylp/meerkat/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"testing"
)

type storeMock struct {
	mock.Mock
}

func (s *storeMock) AddItem(r io.Reader) error {
	args := s.Called(r)
	return args.Error(0)
}

func (s *storeMock) Subscribe() (<-chan io.Reader, string) { //nolint:gocritic
	args := s.Called()
	return args.Get(0).(chan io.Reader), args.String(1)
}

func (s *storeMock) Subscribers() int {
	args := s.Called()
	return args.Int(0)
}

func (s *storeMock) Unsubscribe(ticket string) error {
	args := s.Called(ticket)
	return args.Error(0)
}

func (s *storeMock) Length() int {
	args := s.Called()
	return args.Int(0)
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
