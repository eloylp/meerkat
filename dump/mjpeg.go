package dump

import (
	"io"
	"mime/multipart"
	"net/textproto"
)

type mJPEGDumper struct {
	mimeWriter *multipart.Writer
}

func (d *mJPEGDumper) Boundary() string {
	return d.mimeWriter.Boundary()
}

func NewMJPEGDumper(w io.Writer) *mJPEGDumper {
	return &mJPEGDumper{mimeWriter: multipart.NewWriter(w)}
}

func (d *mJPEGDumper) DumpPart(data []byte) error {
	h := make(textproto.MIMEHeader)
	h.Add("Content-Type", "image/jpeg")
	w, err := d.mimeWriter.CreatePart(h)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}
