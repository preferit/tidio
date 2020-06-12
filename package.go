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

var warn = fox.NewSyncLog(os.Stdout).FilterEmpty().Log

var ErrForbidden = fmt.Errorf("forbidden")

func checkTimesheetFilename(name string) error {
	format := `\d\d\d\d\d\d\.timesheet`
	if ok, _ := regexp.MatchString(format, name); !ok {
		return fmt.Errorf("bad filename: expected format %s", format)
	}
	return nil
}

type Stateful interface {
	WriteState(io.WriteCloser, error) error
	ReadState(ReadOpener) error
}

type ReadOpener func() (io.ReadCloser, error)
