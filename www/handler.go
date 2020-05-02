package www

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/eloylp/meerkat/data"

	"github.com/eloylp/meerkat/writer"
)

func HandleHTMLClient(dfr *data.FlowRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-type", "text/h")

		var img string
		for _, df := range dfr.Flows() {
			img += fmt.Sprintf(`<img src=%s>`, DataStreamPath+df.UUID())
		}

		doc := fmt.Sprintf(`<!DOCTYPE html><html><body>%s</body></html>`, img)
		_, _ = w.Write([]byte(doc))
	}
}

func HandleMJPEG(dfr *data.FlowRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mJPEGWriter := writer.NewMJPEGWriter(w)
		contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mJPEGWriter.Boundary())
		w.Header().Add("Content-Type", contentType)
		dataFlowUUID := strings.TrimPrefix(r.URL.Path, DataStreamPath)
		store, err := dfr.FindStore(dataFlowUUID)
		if err != nil {
			log.Fatal(err)
		}
		readers, uuid := store.Subscribe()
		notify := r.Context().Done()

		go func() {
			<-notify
			if err := store.Unsubscribe(uuid); err != nil {
				log.Fatal(err)
			}
			log.Printf("Client with socket %s left connection", r.RemoteAddr)
		}()

		log.Printf("Started data streaming to client with socket %s", r.RemoteAddr)

		for reader := range readers {
			if err := mJPEGWriter.WritePart(reader); err != nil {
				_, _ = w.Write([]byte("Frame cannot be processed"))
			}
		}
	}
}
