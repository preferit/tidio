package tidio

import (
	"io/ioutil"
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
	asRoot.Run(cmd)

	cmd = rs.NewCmd("/bin/chown", name, dir)
	return asRoot.Run(cmd)
}

// SetLogger
func (me *Service) SetLogger(log fox.Logger) {
	me.warn = fox.NewFilterEmpty(log).Log
}
