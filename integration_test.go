package tidio

import (
	"os"
	"strings"
	"testing"
)

func Test_read_root(t *testing.T) {
	api, log := integration(t)
	resp := api.ReadTimesheet("/").MustSend()
	if resp.Status != "200 OK" {
		t.Error(resp.Status, "\n", log.FlushString())
	}
}

func Test_read_unknown(t *testing.T) {
	api, log := integration(t)
	resp := api.ReadTimesheet("/api/jibberish").MustSend()
	got, exp := resp.Status, "401 Unauthorized"
	if got != exp {
		t.Errorf("%s\n%q != %q", log.FlushString(), got, exp)
	}
}

func integration(t *testing.T) (*API, *LogPrinter) {
	t.Helper()
	if !strings.Contains(os.Getenv("group"), "integration") {
		t.SkipNow()
	}
	var (
		api = NewAPI("http://localhost:13001")
		log = Register(api).Buf()
	)
	return api, log
}
