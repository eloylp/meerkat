// +build unit

package data_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/meerkat/data"
	"github.com/eloylp/meerkat/data/datatest"
	"github.com/eloylp/meerkat/elements"
)

func TestDataFlowRegistry_FindStore(t *testing.T) {
	df1 := data.NewDataFlow("A1234", "", datatest.NewStoreMock(), datatest.NewPumpMock())
	targetUID := "A12345"
	expectedStore := datatest.NewStoreMock()
	df2 := data.NewDataFlow(targetUID, "", expectedStore, datatest.NewPumpMock())
	df3 := data.NewDataFlow("A123456", "", datatest.NewStoreMock(), datatest.NewPumpMock())
	dfr := data.NewFlowRegistry([]elements.Flow{df1, df2, df3})
	result, err := dfr.FindStore(targetUID)
	assert.NoError(t, err)
	assert.Equal(t, expectedStore, result)
}
