package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func main() {
	req, _ := http.NewRequest(
		"POST",
		"http://localhost:13001/api/timesheets/john/202001.timesheet",
		strings.NewReader("x"),
	)
	req.Header.Set("Content-Type", "text/plain")
	req.SetBasicAuth("john", "secret")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	dump, _ := httputil.DumpResponse(resp, false)
	io.Copy(os.Stdout, bytes.NewReader(dump))
}
