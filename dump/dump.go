package dump

type Dumper interface {
	DumpPart(data []byte) error
}
