// Tidio tidio is a standalone http server for the tidio.Service
package main

import (
	"os"

	"github.com/gregoryv/wolf"
)

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

func main() {
	cmd := wolf.NewOSCmd()
	code := NewTidio(cmd).Run()
	os.Exit(code)
}
