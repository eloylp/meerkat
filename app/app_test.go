package app_test

import (
	"go-sentinel/app"
	"io"
	"testing"
)

type StoreMock struct {
	MLength uint
}

func (s StoreMock) AddItem(r io.Reader) error {
	panic("implement me")
}

func (s StoreMock) Subscribe() (chan io.Reader, string) {
	panic("implement me")
}

func (s StoreMock) Subscribers() uint {
	panic("implement me")
}

func (s StoreMock) Unsubscribe(ticket string) error {
	panic("implement me")
}

func (s StoreMock) Length() uint {
	return s.MLength
}

func (s StoreMock) Reset() {
	panic("implement me")
}

type PumpMock struct {
}

func (p PumpMock) Start() {
	panic("implement me")
}

func TestDataFlowRegistry_FindStore(t *testing.T) {

	dfr := &app.DataFlowRegistry{}
	dfr.Add(&app.DataFlow{
		UUID:      "A1234",
		DataStore: &StoreMock{},
		DataPump:  &PumpMock{},
	})
	targetUid := "A12345"
	var expectedStoreLength uint = 12
	dfr.Add(&app.DataFlow{
		UUID:      targetUid,
		DataStore: &StoreMock{MLength: expectedStoreLength},
		DataPump:  &PumpMock{},
	})

	dfr.Add(&app.DataFlow{
		UUID:      "A123456",
		DataStore: &StoreMock{},
		DataPump:  &PumpMock{},
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
