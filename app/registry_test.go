package app

import (
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

func (s *storeMock) Subscribe() (chan io.Reader, string) {
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

	dfr := &dataFlowRegistry{}
	dfr.Add(&dataFlow{
		UUID:      "A1234",
		DataStore: &storeMock{},
		DataPump:  &pumpMock{},
	})
	targetUid := "A12345"
	expectedStore := &storeMock{}
	dfr.Add(&dataFlow{
		UUID:      targetUid,
		DataStore: expectedStore,
		DataPump:  &pumpMock{},
	})

	dfr.Add(&dataFlow{
		UUID:      "A123456",
		DataStore: &storeMock{},
		DataPump:  &pumpMock{},
	})
	result, err := dfr.FindStore(targetUid)
	assert.NoError(t, err)
	assert.Equal(t, expectedStore, result)
}
