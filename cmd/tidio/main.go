// Command tidio is a standalone http server for the tidio.Service
package main

import (
	"os"

	"github.com/gregoryv/fox"
)

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

func main() {
	c := NewCommand()
	c.Logger = fox.NewSyncLog(os.Stderr).FilterEmpty()
	if err := c.run(os.Args...); err != nil {
		c.Log(err)
		os.Exit(1)
	}
}
