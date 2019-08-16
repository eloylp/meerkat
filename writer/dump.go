package writer

import (
	"io"
)

type Dumper interface {
	WritePart(data io.Reader) error
}
