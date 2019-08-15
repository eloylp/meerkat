package app

import (
	"fmt"
	"go-sentinel/dump"
	"log"
	"net/http"
	"strings"
)

func (s *server) handleHTMLClient() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-type", "text/html")

		var images string
		for _, dataFlow := range s.dfr.DataFlows() {
			images += fmt.Sprintf(`<img src=%s>`, FrameStreamEndpoint+dataFlow.UUID)
		}

		html := fmt.Sprintf(`<!DOCTYPE html><html><body>%s</body></html>`, images)
		_, _ = w.Write([]byte(html))
	}
}

func (s *server) handleMJPEG() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mimeWriter := dump.NewMJPEGDumper(w)
		contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
		w.Header().Add("Content-Type", contentType)
		DataFlowUUID := strings.TrimPrefix(r.URL.Path, FrameStreamEndpoint)
		store, err := s.dfr.FindStore(DataFlowUUID)
		if err != nil {
			log.Fatal(err)
		}
		readers, ticket := store.Subscribe()
		notify := r.Context().Done()

		go func() {
			<-notify
			if err := store.Unsubscribe(ticket); err != nil {
				log.Fatal(err)
			}
			log.Printf("Client with socket %s left connection", r.RemoteAddr)
		}()

		log.Printf("Started data streaming to client with socket %s", r.RemoteAddr)

		for image := range readers {
			if err := mimeWriter.DumpPart(image); err != nil {
				_, _ = w.Write([]byte("Frame cannot be processed"))
			}
		}
	}
}
