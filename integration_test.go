package tidio

import (
	"os"
	"strings"
	"testing"
)

func Test_integration(t *testing.T) {
	if !strings.Contains(os.Getenv("group"), "integration") {
		t.SkipNow()
	}
	var (
		api = NewAPI("http://localhost:13001")
		log = Register(api).Buf()
	)
	resp := api.ReadTimesheet("/").MustSend()
	if resp.Status != "200 OK" {
		t.Error(resp.Status, "\n", log.FlushString())
	}

	resp = api.ReadTimesheet("/api/jibberish").MustSend()
	if resp.Status != "401 Unauthorized" {
		t.Errorf("%s\n%s", resp.Status, log.FlushString())
	}
}
