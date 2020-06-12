package tidio

import (
	"io"
	"os"
)

type Stateful interface {
	WriteState(WriteOpener) error
	ReadState(ReadOpener) error
}

type WriteOpener func() (io.WriteCloser, error)
type ReadOpener func() (io.ReadCloser, error)

func toFile(filename string) WriteOpener {
	return func() (io.WriteCloser, error) {
		return os.Create(filename)
	}
}

func fromFile(filename string) ReadOpener {
	return func() (io.ReadCloser, error) {
		return os.Open(filename)
	}
}
