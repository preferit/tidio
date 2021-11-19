package tidio

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/worklog"
)

var (
	asJohn          = NewCredentials("john", "secret")
	withJohnAccount = asJohn
)

func init() {
	Conf.SetDebug(true)
}

func TestAPI_CreateTimesheet_asJohn(t *testing.T) {
	var (
		sys  = NewSystem(withJohnAccount)
		ts   = httptest.NewServer(NewRouter(sys))
		api  = NewAPI(ts.URL, asJohn)
		log  = Register(sys, api).Buf()
		path = "/api/timesheets/john/202001.timesheet"
		body = worklog.Render(2020, 1, 8)
		req  = api.CreateTimesheet(path, body)
	)
	defer ts.Close()

	resp := req.MustSend()
	if resp.StatusCode != 201 {
		t.Error(resp.Status, "\n", log.FlushString())
	}
}

func TestAPI_CreateTimesheet_asAnonymous(t *testing.T) {
	var (
		sys = NewSystem(withJohnAccount)
		ts  = httptest.NewServer(NewRouter(sys))
		api = NewAPI(ts.URL)
		log = Register(sys, api).Buf()

		path = "/api/timesheets/john/202001.timesheet"
		body = worklog.Render(2020, 1, 8)
		req  = api.CreateTimesheet(path, body)
	)
	defer ts.Close()

	resp := req.MustSend()
	if resp.StatusCode != 401 {
		t.Error(resp.Status, "\n", log.FlushString())
	}
}

func TestAPI_ReadTimesheet_asJohn(t *testing.T) {
	var (
		sys = NewSystem(withJohnAccount)
		ts  = httptest.NewServer(NewRouter(sys))
	)
	defer ts.Close()
	var (
		api  = NewAPI(ts.URL, asJohn)
		path = "/api/timesheets/john/202001.timesheet"
		body = worklog.Render(2020, 1, 8)
		_    = api.CreateTimesheet(path, body).MustSend()
		resp = api.ReadTimesheet(path).MustSend()
	)
	if resp.StatusCode != 200 {
		t.Error(resp.Status)
	}
}

func TestAPI_ReadTimesheet_noSuchResource(t *testing.T) {
	var (
		sys = NewSystem(withJohnAccount)
		ts  = httptest.NewServer(NewRouter(sys))
	)
	defer ts.Close()
	var (
		api  = NewAPI(ts.URL, asJohn)
		resp = api.ReadTimesheet("/api/timesheets/john/nosuch").MustSend()
	)
	if resp.StatusCode != 404 {
		t.Error(resp.Status)
	}
}

func TestAPI_ReadTimesheet_asAnonymous(t *testing.T) {
	var (
		sys = NewSystem(withJohnAccount)
		ts  = httptest.NewServer(NewRouter(sys))
	)
	defer ts.Close()
	var (
		api  = NewAPI(ts.URL, asJohn)
		path = "/api/timesheets/john/202001.timesheet"
		body = worklog.Render(2020, 1, 8)
		_    = api.CreateTimesheet(path, body).MustSend()
	)
	api = NewAPI(ts.URL) // anonymous
	resp := api.ReadTimesheet(path).MustSend()
	if resp.StatusCode == 200 {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		t.Error("should fail\n", buf.String())
	}
}

func Test_hacks(t *testing.T) {
	var (
		sys = NewSystem(withJohnAccount)
		ts  = httptest.NewServer(NewRouter(sys))
	)
	defer ts.Close()

	t.Run("malformed basic auth", func(t *testing.T) {
		api := NewAPI(ts.URL)
		log := Log(api).Buf()
		r := api.CreateTimesheet("", nil).Request
		r.Header.Set("Authorization", "Basi")
		resp := api.MustSend()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Error("should fail:", resp.Status)
			t.Log(log.FlushString())
		}
	})

	t.Run("empty basic auth", func(t *testing.T) {
		api := NewAPI(ts.URL, &Credentials{})
		resp := api.CreateTimesheet("", nil).MustSend()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Error("should fail:", resp.Status)
		}
	})
}

func TestAPI_Send_nil_request(t *testing.T) {
	var (
		api = NewAPI("http://localhost")
	)
	if _, err := api.Send(); err == nil {
		t.Error("should fail")
	}

	defer func() {
		if e := recover(); e == nil {
			t.Error("should panic")
		}
	}()
	api.MustSend()
}

func TestAPI_Send_failing_response(t *testing.T) {
	var (
		api = NewAPI("http://_1234nosuchhost.net")
		log = Register(api).Buf()
	)
	api.Request, _ = http.NewRequest("GET", "/", nil)
	if _, err := api.Send(); err == nil {
		t.Error("should fail\n", log.FlushString())
	}
}

func TestAPI_warnings(t *testing.T) {
	var (
		api = NewAPI("http://localhost")
	)
	api.warn(nil)
	api.warn(fmt.Errorf("failed"))
}

func TestAPI_Auth_nil_request(t *testing.T) {
	api := NewAPI("http://localhost")
	api.Auth(nil)
}

func Test_defaults(t *testing.T) {
	sys := NewSystem()
	assert := asserter.New(t)
	exp := assert().ResponseFrom(NewRouter(sys))
	exp.StatusCode(200, "GET", "/")
	exp.StatusCode(200, "GET", "/api")
	exp.StatusCode(405, "X", "/api")
	exp.StatusCode(http.StatusUnauthorized, "GET", "/api/timesheets")
	exp.StatusCode(http.StatusUnauthorized, "GET", "/api/timesheets",
		http.Header{
			"Authorization": []string{"Basic invalid-encoding"},
		},
	)
	exp.BodyIs(`{"resources":[{"name": "timesheets"}]}`, "GET", "/api")
}

func TestBasicAuth_Set(t *testing.T) {
	t.Run("on nil", func(t *testing.T) {
		c := NewBasicAuth(&Credentials{})
		if err := c.Set(nil); err == nil {
			t.Error("should fail")
		}
	})

	t.Run("on *http.Request", func(t *testing.T) {
		c := NewBasicAuth(&Credentials{})
		r, _ := http.NewRequest("GET", "http://example.com", nil)
		if err := c.Set(r); err != nil {
			t.Error("should work:", err)
		}
	})
}

func Test_loggers(t *testing.T) {
	log := Register(t).Buf()
	defer Unregister(t)

	somefunc(t) // should log
	Conf.Unregister(t)
	somefunc(t) // no logger registered

	got := log.FlushString()
	if strings.Count(got, "hello") != 1 {
		t.Errorf("cached log\n%s", got)
		t.Error("writes", log.writes)
	}
}

func Test_Register_panics(t *testing.T) {
	defer catchPanic(t)
	Register()
}

func somefunc(t *testing.T) {
	Log(t).Info("hello")
	Log(t).Info("world")
}

// ----------------------------------------

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
