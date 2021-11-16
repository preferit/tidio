package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/gregoryv/cmdline"
	"github.com/preferit/tidio"
)

func main() {
	var (
		cli  = cmdline.NewBasicParser()
		host = cli.Option("-H, --host").String("https://tidio.preferit.se")

		user = cli.Option("-u, --username").String(os.Getenv("USER"))
		pass = cli.Option("-p, --password").String(os.Getenv("PASSWORD"))
		cred = tidio.NewCredentials(user, pass)

		filename = cli.Required("FILE").String("")
	)
	cli.Parse()

	uploadFile(cred, host, filename)
}

func uploadFile(cred *tidio.Credentials, host, filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	api := tidio.NewAPI(host, cred)
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
