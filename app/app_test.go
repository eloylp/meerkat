package app

import (
	"io"
	"testing"
)

type storeMock struct {
	MLength uint
}

func (s storeMock) AddItem(r io.Reader) error {
	panic("implement me")
}

func (s storeMock) Subscribe() (chan io.Reader, string) {
	panic("implement me")
}

func (s storeMock) Subscribers() uint {
	panic("implement me")
}

func (s storeMock) Unsubscribe(ticket string) error {
	panic("implement me")
}

func (s storeMock) Length() uint {
	return s.MLength
}

func (s storeMock) Reset() {
	panic("implement me")
}

type pumpMock struct {
}

func (p pumpMock) Start() {
	panic("implement me")
}

func TestDataFlowRegistry_FindStore(t *testing.T) {

	dfr := &dataFlowRegistry{}
	dfr.Add(&dataFlow{
		UUID:      "A1234",
		DataStore: &storeMock{},
		DataPump:  &pumpMock{},
	})
	targetUid := "A12345"
	var expectedStoreLength uint = 12
	dfr.Add(&dataFlow{
		UUID:      targetUid,
		DataStore: &storeMock{MLength: expectedStoreLength},
		DataPump:  &pumpMock{},
	})

	dfr.Add(&dataFlow{
		UUID:      "A123456",
		DataStore: &storeMock{},
		DataPump:  &pumpMock{},
	})

	result, err := dfr.FindStore(targetUid)
	if err != nil {
		t.Fatal(err)
	}

	length := result.Length()
	if length != expectedStoreLength {
		t.Errorf("Expected store has length %v, got store with length %v", expectedStoreLength, length)
	}

}
