package writer

import (
	"io"
	"mime/multipart"
	"net/textproto"
)

type MJPEGWriter struct {
	mimeWriter *multipart.Writer
}

func (mw *MJPEGWriter) Boundary() string {
	return mw.mimeWriter.Boundary()
}

func NewMJPEGWriter(w io.Writer) *MJPEGWriter {
	return &MJPEGWriter{mimeWriter: multipart.NewWriter(w)}
}

func (mw *MJPEGWriter) WritePart(data io.Reader) error {
	h := make(textproto.MIMEHeader)
	h.Add("Content-Type", "image/jpeg")
	w, err := mw.mimeWriter.CreatePart(h)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, data); err != nil {
		return err
	}
	return nil
}
