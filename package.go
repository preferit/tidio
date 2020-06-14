/*
Package tidio provides domain logic for the tidio web service.

*/
package tidio

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/gregoryv/fox"
)

type Stateful interface {
	Loader
	Saver
}
type Loader interface{ Load() error }
type Saver interface{ Save() error }

type Source interface{ Open() (io.ReadCloser, error) }
type Destination interface {
	Create() (io.WriteCloser, error)
}

type FilePersistent interface {
	PersistToFile(string)
}

// ----------------------------------------

type FileSource string

func (me FileSource) Open() (io.ReadCloser, error) {
	return os.Open(string(me))
}

type FileDestination string

func (me FileDestination) Create() (io.WriteCloser, error) {
	return os.Create(string(me))
}

// ----------------------------------------

var warn = fox.NewSyncLog(os.Stdout).FilterEmpty().Log

var ErrForbidden = fmt.Errorf("forbidden")

func checkTimesheetFilename(name string) error {
	format := `\d\d\d\d\d\d\.timesheet`
	if ok, _ := regexp.MatchString(format, name); !ok {
		return fmt.Errorf("bad filename: expected format %s", format)
	}
	return nil
}

type errors []error

func (me errors) First() error {
	for _, err := range me {
		if err != nil {
			return err
		}
	}
	return nil
}
