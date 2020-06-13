package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Source interface {
	Open() (io.ReadCloser, error)
}

// ----------------------------------------

type StringSource string

func (me StringSource) Open() (io.ReadCloser, error) {
	return ioutil.NopCloser(strings.NewReader(string(me))), nil
}

// ----------------------------------------

type BrokenSource struct{}

func (BrokenSource) Open() (io.ReadCloser, error) {
	return nil, fmt.Errorf("broken source")
}

// ----------------------------------------

type FileSource string

func (me FileSource) Open() (io.ReadCloser, error) {
	return os.Open(string(me))
}
