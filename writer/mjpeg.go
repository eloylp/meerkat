package writer

import (
	"io"
	"mime/multipart"
	"net/textproto"
)

type mJPEGWriter struct {
	mimeWriter *multipart.Writer
}

func (d *mJPEGWriter) Boundary() string {
	return d.mimeWriter.Boundary()
}

func NewMJPEGWriter(w io.Writer) *mJPEGWriter {
	return &mJPEGWriter{mimeWriter: multipart.NewWriter(w)}
}

func (d *mJPEGWriter) WritePart(data io.Reader) error {
	h := make(textproto.MIMEHeader)
	h.Add("Content-Type", "image/jpeg")
	w, err := d.mimeWriter.CreatePart(h)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, data); err != nil {
		return err
	}
	return nil
}
