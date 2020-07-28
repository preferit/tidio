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
	w, _ := asRoot.Create("/api/timesheets/202001.timesheet")
	timesheet.Render(w, 2020, 1, 8)
	w.Close()

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
	cmd := rs.NewCmd("/bin/mkacc", name)
	if err := asRoot.Run(cmd); err != nil {
		return err
	}
	cmd = rs.NewCmd("/bin/secure", "-a", name, "-s", secret)
	asRoot.Run(cmd)

	dir := path.Join("/api/timesheets", name)
	cmd = rs.NewCmd("/bin/mkdir", dir)
	return asRoot.Run(cmd)
}

// SetLogger
func (me *Service) SetLogger(log fox.Logger) {
	me.warn = fox.NewFilterEmpty(log).Log
}

func (me *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	acc, err := me.authenticate(r.Header.Get("Authorization"))
	if err != nil {
		me.warn(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	asAcc := acc.Use(me.sys)

	// Check resource access permissions
	_, err = asAcc.Stat(r.URL.Path)
	if acc == rs.Anonymous && err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "plain/text")
		w.Write([]byte(err.Error()))
		return
	}

	// Serve the specific method
	switch r.Method {
	case "GET":
		// todo if url is a resource return it's content
		cmd := rs.NewCmd("/bin/ls", "-json", "-json-name", "resources", r.URL.Path)
		cmd.Out = w
		asAcc.Run(cmd)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (me *Service) authenticate(h string) (*rs.Account, error) {
	if h == "" {
		return rs.Anonymous, nil
	}
	basic, err := ParseBasicAuth(h)
	if err != nil {
		return rs.Anonymous, err
	}
	var (
		name   = basic.AccountName
		secret = basic.Secret
	)
	asRoot := rs.Root.Use(me.sys)
	cmd := rs.NewCmd("/bin/secure", "-c", "-a", name, "-s", secret)
	if err := asRoot.Run(cmd); err != nil {
		return rs.Anonymous, err
	}
	var acc rs.Account
	err = asRoot.LoadAccount(&acc, name)
	return &acc, err
}
