package app

import (
	"fmt"
	"go-sentinel/writer"
	"log"
	"net/http"
	"strings"
)

func (s *server) handleHTMLClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-type", "text/html")

		var images string
		for _, dataFlow := range s.dfr.DataFlows() {
			images += fmt.Sprintf(`<img src=%s>`, DataStreamPath+dataFlow.UUID)
		}

		html := fmt.Sprintf(`<!DOCTYPE html><html><body>%s</body></html>`, images)
		_, _ = w.Write([]byte(html))
	}
}

func (s *server) handleMJPEG() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mimeWriter := writer.NewMJPEGWriter(w)
		contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
		w.Header().Add("Content-Type", contentType)
		DataFlowUUID := strings.TrimPrefix(r.URL.Path, DataStreamPath)
		store, err := s.dfr.FindStore(DataFlowUUID)
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
			if err := mimeWriter.WritePart(reader); err != nil {
				_, _ = w.Write([]byte("Frame cannot be processed"))
			}
		}
	}
}
