package www_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/meerkat/data"
	"github.com/eloylp/meerkat/data/datatest"
	"github.com/eloylp/meerkat/elements"
	"github.com/eloylp/meerkat/www"
)

func TestHandleHTMLClient(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, www.DashboardPath, nil)
	h := router()
	h.ServeHTTP(rec, req)
	gotContentType := rec.Result().Header.Get("Content-Type")
	assert.Contains(t, gotContentType, "text/html")
	body := rec.Result().Body
	bytes, err := ioutil.ReadAll(body)
	assert.NoError(t, err)
	want := fmt.Sprintf(`<img src=%s>`, www.DataStreamPath+"/A1234")
	assert.Contains(t, string(bytes), want)
}

func router() http.Handler {
	return www.Router(dataFlowRegistry())
}

func dataFlowRegistry() *data.FlowRegistry {
	store := datatest.NewStoreMock()
	pump := datatest.NewPumpMock()
	flow := data.NewDataFlow("A1234", "example.com/picture.png", store, pump)
	dfr := data.NewFlowRegistry([]elements.Flow{flow})
	return dfr
}
