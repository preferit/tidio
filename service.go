package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"

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

	srv := &Service{
		sys:    sys,
		Router: NewRouter(sys),
	}
	srv.SetLogger(fox.NewSyncLog(ioutil.Discard))
	return srv
}

type Service struct {
	fox.Logger
	warn func(...interface{})

	sys *rs.System
	*Router
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
	me.Logger = log
	me.warn = fox.NewFilterEmpty(log).Log
}

// RestoreState restores the resource system from the given filename.
// Restoring the state replaces current system.
func (me *Service) RestoreState(filename string) error {
	me.Log("restore state: ", filename)
	r, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open state file: %w", err)
	}
	defer r.Close()
	me.sys, err = rs.Import(r)
	//me.warn(err)
	return err
}

// AutoPersist enables automatic persistence of system state to given filename.
func (me *Service) AutoPersist(dest Storage) {
	last := me.sys.LastModified()
	go func() {
		for {
			// todo decouple and use events
			modified := me.sys.LastModified()
			if !modified.After(last) {
				time.Sleep(time.Second)
				continue
			}
			last = modified
			err := me.PersistState(dest)
			me.warn(err)
		}
	}()
}

// PersistState
func (me *Service) PersistState(dest Storage) error {
	me.Log("persist state")
	w, err := dest.Create()
	if err != nil {
		return err
	}
	defer w.Close()
	return me.sys.Export(w)
}

type Storage interface {
	Create() (io.WriteCloser, error)
}

func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{filename: filename}
}

type FileStorage struct {
	filename string
}

// Create
func (me *FileStorage) Create() (io.WriteCloser, error) {
	return os.Create(me.filename)
}
