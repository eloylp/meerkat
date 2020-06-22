package datatest

import (
	"github.com/stretchr/testify/mock"
)

type PumpMock struct {
	mock.Mock
}

func NewPumpMock() *PumpMock {
	return &PumpMock{}
}

func (p *PumpMock) Start() {
	p.Called()
}
