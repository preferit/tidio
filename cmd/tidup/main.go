package main

import (
	"fmt"
	"os"
	"path"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/fox"
	"github.com/preferit/tidio"
)

func main() {
	var (
		cli      = cmdline.New(os.Args...)
		host     = cli.Option("--host").String("https://tidio.preferit.se")
		help     = cli.Flag("-h, --help")
		user     = cli.Option("-u, --username").String(os.Getenv("USER"))
		pass     = cli.Option("-p, --password").String(os.Getenv("PASSWORD"))
		cred     = tidio.NewCredentials(user, pass)
		filename = cli.NeedArg("FILE", 0).String()
	)

	switch {
	case !cli.Ok():
		fmt.Println("Try --help for more information")
	case help:
		cli.WriteUsageTo(os.Stdout)
	default:
		uploadFile(cred, host, filename)
	}
}

func uploadFile(cred *tidio.Credentials, host, filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	api := tidio.NewAPI(
		host,
		cred,
		fox.Logging{
			fox.NewSyncLog(os.Stderr),
		},
	)
	// todo optional path
	path := path.Join("/api/timesheets/john/", filename)
	resp := api.CreateTimesheet(path, fh).MustSend()
	fmt.Println(api.Request.Header)
	fmt.Println(resp.Status)
	return nil
}
