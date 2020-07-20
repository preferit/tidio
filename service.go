package tidio

import (
	"io/ioutil"
	"net/http"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/go-timesheet"
	"github.com/gregoryv/rs"
)

func NewService() *Service {
	sys := rs.NewSystem()
	asRoot := rs.Root.Use(sys)
	asRoot.Exec("/bin/mkdir /api")
	asRoot.Exec("/bin/mkdir /api/timesheets")

	// default templates
	res, _ := asRoot.Create("/api/timesheets/202001.timesheet")
	timesheet.Render(res, 2020, 1, 8)
	res.Close()

	return &Service{
		warn: fox.NewSyncLog(ioutil.Discard).Log,
		sys:  sys,
	}
}

type Service struct {
	warn func(...interface{})

	sys *rs.System
}

// SetLogger
func (me *Service) SetLogger(log fox.Logger) {
	me.warn = fox.NewFilterEmpty(log).Log
}

func (me *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	acc := me.authorize()
	asAcc := acc.Use(me.sys)

	// Check resource access permissions
	_, err := asAcc.Stat(r.URL.Path)
	if acc == rs.Anonymous && err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Serve the specific method
	switch r.Method {
	case "GET":
		asAcc.Fexec(w, "/bin/ls", "-json", "-json-name", "resources", r.URL.Path)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// authorize
func (me *Service) authorize() *rs.Account {
	acc := rs.Anonymous
	return acc
}
