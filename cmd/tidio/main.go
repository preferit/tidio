// Tidio tidio is a standalone http server for the tidio.Service
package main

import (
	"os"

	"github.com/gregoryv/wolf"
	"github.com/preferit/tidio"
)

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

func main() {
	conf := tidio.Conf
	conf.SetOutput(os.Stderr)

	cmd := wolf.NewOSCmd()
	app := tidio.NewApp()
	code := app.Run(cmd)
	os.Exit(code)
}
