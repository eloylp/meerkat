package store_test

import (
	"bytes"
	"fmt"
	"go-sentinel/store"
	"io"
	"io/ioutil"
	"testing"
)

func TestTimeLineStore(t *testing.T) {

	samples := []io.Reader{
		bytes.NewReader([]byte("d1")),
		bytes.NewReader([]byte("d2")),
		bytes.NewReader([]byte("d3")),
	}
	s := store.NewTimeLineStore(3, 10000)
	listenCh := s.Subscribe()
	for _, sample := range samples {
		if err := s.AddItem(sample); err != nil {
			t.Fatal(err)
		}
	}
	s.Close()
	var count uint
	count++
	for item := range listenCh {
		expected := "d" + fmt.Sprint(count)
		item, err := ioutil.ReadAll(item)
		if err != nil {
			t.Fatal(err)
		}
		got := string(item)
		if expected != got {
			t.Fatalf("Error listening subscribed frames. Expected frame was %s but got %s", expected, got)
		}
		count++
	}

}
