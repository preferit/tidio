package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/fox"
	"github.com/preferit/tidio"
)

func main() {
	var (
		cli  = cmdline.NewParser(os.Args...)
		host = cli.Option("--host").String("https://tidio.preferit.se")
		help = cli.Flag("-h, --help")

		user = cli.Option("-u, --username").String(os.Getenv("USER"))
		pass = cli.Option("-p, --password").String(os.Getenv("PASSWORD"))
		cred = tidio.NewCredentials(user, pass)

		filename = cli.Required("FILE").String()
	)

	switch {
	case help:
		cli.WriteUsageTo(os.Stdout)

	case !cli.Ok():
		fmt.Println(cli.Error())
		fmt.Println("Try -h for more information")

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
	if resp.StatusCode >= 400 {
		fmt.Println("--- body ---")
		io.Copy(os.Stderr, resp.Body)
	}
	return nil
}
