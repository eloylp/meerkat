package dump

import (
	"io"
)

type Dumper interface {
	DumpPart(data io.Reader) error
}
