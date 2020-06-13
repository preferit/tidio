package tidio

import (
	"fmt"
	"io"
	"os"
)

type Destination interface {
	Create() (io.WriteCloser, error)
}

// ----------------------------------------

type BrokenDestination struct{}

func (BrokenDestination) Create() (io.WriteCloser, error) {
	return nil, fmt.Errorf("broken destination")
}

// ----------------------------------------

type FileDestination string

func (me FileDestination) Create() (io.WriteCloser, error) {
	return os.Create(string(me))
}

// ----------------------------------------

func NopDestination() *nopDest {
	return &nopDest{}
}

type nopDest struct{}

func (me *nopDest) Create() (io.WriteCloser, error) {
	return me, nil
}

func (*nopDest) Write(b []byte) (int, error) {
	return len(b), nil
}

func (*nopDest) Close() error { return nil }
