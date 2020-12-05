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
		log = RLog(api).Buf()
	)
	resp := api.ReadTimesheet("/").MustSend()
	if resp.Status != "200 OK" {
		t.Error(resp.Status, "\n", log.FlushString())
	}
}
