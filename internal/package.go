package internal

import (
	"fmt"
	"io"
)

type None string

func (me None) Open() (io.ReadCloser, error) {
	return nil, fmt.Errorf("%s.Source not set", me)
}
func (me None) Create() (io.WriteCloser, error) {
	return nil, fmt.Errorf("%s.Destination not set", me)
}
