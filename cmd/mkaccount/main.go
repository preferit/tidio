package main

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/fox"
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
	log := fox.NewSyncLog(os.Stdout)
	trace, cleanup := tidio.NewTrace(log)
	defer cleanup("createAccount", name, "****")

	r, err := os.Open(filename)
	if err != nil {
		trace.Log(err)
		return fmt.Errorf("open state file: %w", err)
	}
	defer r.Close()
	sys := rs.NewSystem()
	err = sys.Import("/", r)
	if err != nil {
		trace.Log(err)
		return err
	}

	asRoot := rs.Root.Use(sys)
	asRoot.SetAuditer(trace)
	sh := tidio.NewShell(asRoot)
	err = tidio.CreateAccount(sh, name, secret)
	if err != nil {
		return err
	}
	return sys.Export(os.Stdout)
}
