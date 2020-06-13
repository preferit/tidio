package tidio

import (
	"fmt"
	"io"
)

type Stateful interface {
	Loader
	Saver
}

type Loader interface {
	Load() error
}

type Saver interface {
	Save() error
}

type None string

func (me None) Open() (io.ReadCloser, error) {
	return nil, fmt.Errorf("%s.Source not set", me)
}
func (me None) Create() (io.WriteCloser, error) {
	return nil, fmt.Errorf("%s.Destination not set", me)
}

type WriteCloser interface {
	io.Closer
	Write([]byte) (int, error)
}

type FilePersistent interface {
	PersistToFile(string)
}
