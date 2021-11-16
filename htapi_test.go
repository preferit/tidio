package tidio

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/go-timesheet"
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
		srv  = NewSystem(withJohnAccount)
		ts   = httptest.NewServer(srv.Router())
		api  = NewAPI(ts.URL, asJohn)
		log  = Register(srv, api).Buf()
		path = "/api/timesheets/john/202001.timesheet"
		body = timesheet.Render(2020, 1, 8)
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
		srv = NewSystem(withJohnAccount)
		ts  = httptest.NewServer(srv.Router())
		api = NewAPI(ts.URL)
		log = Register(srv, api).Buf()

		path = "/api/timesheets/john/202001.timesheet"
		body = timesheet.Render(2020, 1, 8)
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
		srv = NewSystem(withJohnAccount)
		ts  = httptest.NewServer(srv.Router())
	)
	defer ts.Close()
	var (
		api  = NewAPI(ts.URL, asJohn)
		path = "/api/timesheets/john/202001.timesheet"
		body = timesheet.Render(2020, 1, 8)
		_    = api.CreateTimesheet(path, body).MustSend()
		resp = api.ReadTimesheet(path).MustSend()
	)
	if resp.StatusCode != 200 {
		t.Error(resp.Status)
	}
}

func TestAPI_ReadTimesheet_noSuchResource(t *testing.T) {
	var (
		srv = NewSystem(withJohnAccount)
		ts  = httptest.NewServer(srv.Router())
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
		srv = NewSystem(withJohnAccount)
		ts  = httptest.NewServer(srv.Router())
	)
	defer ts.Close()
	var (
		api  = NewAPI(ts.URL, asJohn)
		path = "/api/timesheets/john/202001.timesheet"
		body = timesheet.Render(2020, 1, 8)
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
		srv = NewSystem(withJohnAccount)
		ts  = httptest.NewServer(srv.Router())
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
	srv := NewSystem()
	assert := asserter.New(t)
	exp := assert().ResponseFrom(srv.Router())
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
