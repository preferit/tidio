package tidio

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/fox"
	"github.com/gregoryv/go-timesheet"
)

var (
	withJohnAccount = InitialAccount{"john", "secret"}
	asJohn          = NewCredentials("john", "secret")
)

func TestAPI_CreateTimesheet_asJohn(t *testing.T) {
	var (
		log = fox.Logging{t}
		srv = NewService(log, withJohnAccount)
		ts  = httptest.NewServer(srv.Router())
	)
	defer ts.Close()
	var (
		api  = NewAPI(ts.URL, asJohn, log)
		path = "/api/timesheets/john/202001.timesheet"
		body = timesheet.Render(2020, 1, 8)
		resp = api.CreateTimesheet(path, body).MustSend()
	)
	if resp.StatusCode != 201 {
		t.Error(resp.Status)
	}
}

func TestAPI_CreateTimesheet_asAnonymous(t *testing.T) {
	var (
		log = fox.Logging{t}
		srv = NewService(log, withJohnAccount)
		ts  = httptest.NewServer(srv.Router())
	)
	defer ts.Close()
	var (
		api  = NewAPI(ts.URL, log)
		path = "/api/timesheets/john/202001.timesheet"
		body = timesheet.Render(2020, 1, 8)
		resp = api.CreateTimesheet(path, body).MustSend()
	)
	if resp.StatusCode != 401 {
		t.Error(resp.Status)
	}
}

func TestAPI_ReadTimesheet_asJohn(t *testing.T) {
	var (
		log = fox.Logging{t}
		srv = NewService(log, withJohnAccount)
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
		log = fox.Logging{t}
		srv = NewService(log, withJohnAccount)
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
		log = fox.Logging{t}
		srv = NewService(log, withJohnAccount)
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
		log = fox.Logging{t}
		srv = NewService(log, withJohnAccount)
		ts  = httptest.NewServer(srv.Router())
	)
	defer ts.Close()

	t.Run("malformed basic auth", func(t *testing.T) {
		api := NewAPI(ts.URL)
		r := api.CreateTimesheet("", nil).Request
		r.Header.Set("Authorization", "Basi")
		resp := api.MustSend()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Error("should fail:", resp.Status)
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
		log = fox.Logging{t}
		api = NewAPI("http://localhost", log)
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
		log = fox.Logging{t}
		api = NewAPI("http://_1234nosuchhost.net", log)
	)
	api.Request, _ = http.NewRequest("GET", "/", nil)
	if _, err := api.Send(); err == nil {
		t.Error("should fail")
	}
}

func TestAPI_warnings(t *testing.T) {
	var (
		log = fox.Logging{t}
		api = NewAPI("http://localhost", log)
	)
	api.warn(nil)
	api.warn(fmt.Errorf("failed"))
}

func TestAPI_Auth_nil_request(t *testing.T) {
	api := NewAPI("http://localhost")
	api.Auth(nil)
}

func Test_defaults(t *testing.T) {
	srv := NewService()
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
