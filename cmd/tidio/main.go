// Tidio tidio is a standalone http server for the tidio.Service
package main

import (
	"os"

	"github.com/gregoryv/wolf"
	"github.com/preferit/tidio"
)

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

func main() {
	cmd := wolf.NewOSCmd()
	tidio.Conf.SetOutput(os.Stderr)
	app := tidio.NewApp()
	code := app.Run(cmd)
	os.Exit(code)
}
