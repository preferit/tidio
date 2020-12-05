package main

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/rs"
	"github.com/preferit/tidio"
)

func main() {
	var (
		cli      = cmdline.NewParser(os.Args...)
		filename = cli.Option("--state-file").String("/var/local/tidio/system.state")
		name     = cli.Option("-n, --name").String("")
		secret   = cli.Option("-s, --secret").String("")
		help     = cli.Flag("-h, --help")
	)
	switch {
	case !cli.Ok():
		fmt.Println(cli.Error())
		os.Exit(1)
	case help:
		cli.WriteUsageTo(os.Stderr)
		os.Exit(0)
	}
	err := createAccount(filename, name, secret)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createAccount(filename string, name, secret string) error {
	r, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open state file: %w", err)
	}
	defer r.Close()
	sys := rs.NewSystem()
	err = sys.Import("/", r)
	if err != nil {
		return err
	}

	asRoot := rs.Root.Use(sys)

	sh := tidio.NewShell(asRoot)
	err = tidio.CreateAccount(sh, name, secret)
	if err != nil {
		return err
	}
	return sys.Export(os.Stdout)
}
