package www

import (
	"fmt"
	"go-sentinel/dump"
	"net/http"
)

func (s *server) handleHTMLClient() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-type", "text/html")
		html := fmt.Sprintf(`<!DOCTYPE html><html><body><img src="%s"></body></html>`, FrameStreamEndpoint)
		_, _ = w.Write([]byte(html))
	}
}

func (s *server) handleMJPEG() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mimeWriter := dump.NewMJPEGDumper(w)
		contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
		w.Header().Add("Content-Type", contentType)

		readers, _ := s.store.Subscribe()

		for image := range readers {
			if err := mimeWriter.DumpPart(image); err != nil {
				_, _ = w.Write([]byte("Frame cannot be processed"))
			}
		}
	}
}
