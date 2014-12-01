package output

import (
	"fmt"
	"io"
)

type Outputter interface {
	WritePublic(p []byte) (n int, err error)
	WritePrivate(name string, p []byte) (n int, err error)
}

func NewPrefixed(output io.Writer) Outputter {
	return PrefixedOutputter{output: output}
}

type PrefixedOutputter struct {
	output io.Writer
}

func (out PrefixedOutputter) WritePublic(p []byte) (n int, err error) {
	return out.output.Write(append([]byte("PUBLIC: "), p...))
}
func (out PrefixedOutputter) WritePrivate(name string, p []byte) (n int, err error) {
	return out.output.Write(append([]byte(fmt.Sprintf("PRIVATE[%s]: ", name)), p...))
}
