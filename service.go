package tidio

import (
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/go-timesheet"
	"github.com/gregoryv/rs"
)

func NewService() *Service {
	sys := rs.NewSystem()
	asRoot := rs.Root.Use(sys)
	asRoot.Exec("/bin/mkdir /etc/basic")
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

// AddAccount creates a system account and stores the secret in
// /etc/basic
func (me *Service) AddAccount(name, secret string) error {
	asRoot := rs.Root.Use(me.sys)
	w := ioutil.Discard
	if err := asRoot.Fexec(w, "/bin/mkacc", name); err != nil {
		return err
	}
	key := NewKey(secret, path.Join("/etc/accounts", name+".acc"))
	return asRoot.Save(path.Join("/etc/basic", name+".key"), &key)
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
