package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/go-timesheet"
	"github.com/gregoryv/rs"
)

func NewService(settings ...Setting) *Service {
	sys := rs.NewSystem()
	asRoot := rs.Root.Use(sys)
	asRoot.Exec("/bin/mkdir /etc/basic")
	asRoot.Exec("/bin/mkdir /api")
	asRoot.Exec("/bin/mkdir /api/timesheets")

	srv := &Service{
		sys: sys,
	}
	srv.SetLogger(fox.NewSyncLog(ioutil.Discard))
	for _, setting := range settings {
		err := setting.Set(srv)
		if err != nil {
			panic(err)
		}
	}
	return srv
}

type Service struct {
	fox.Logger
	warn func(...interface{})

	sys *rs.System
}

// InitResources
func (me *Service) InitResources() error {
	// default templates
	asRoot := rs.Root.Use(me.sys)
	w, err := asRoot.Create("/api/timesheets/202001.timesheet")
	if err != nil {
		return err
	}
	timesheet.RenderTo(w, 2020, 1, 8)
	w.Close()
	asRoot.Exec("/bin/chmod 05555 /api/timesheets/202001.timesheet")
	return nil
}

func (me *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := &Response{
		sys: me.sys,
	}
	err := resp.Build(r)
	if err != nil {
		me.warn(err)
		resp.WriteError(w, err)
		return
	}
	resp.Send(w)
}

// AddAccount creates a system account and stores the secret in
// /etc/basic
func (me *Service) AddAccount(name, secret string) error {
	asRoot := NewShell(rs.Root.Use(me.sys))
	asRoot.Execf("/bin/mkacc %s", name)
	asRoot.Execf("/bin/secure -a %s -s %s", name, secret)
	asRoot.Execf("/bin/mkdir /api/timesheets")

	dir := path.Join("/api/timesheets", name)
	asRoot.Execf("/bin/mkdir %s", dir)

	return asRoot.Execf("/bin/chown %s %s", name, dir)
}

// SetLogger
func (me *Service) SetLogger(log fox.Logger) {
	me.Logger = log
	me.warn = fox.NewFilterEmpty(log).Log
}

// RestoreState restores the resource system from the given filename.
func (me *Service) RestoreState(filename string) error {
	me.Log("restore state: ", filename)
	r, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open state file: %w", err)
	}
	defer r.Close()
	err = me.sys.Import("/", r)
	//me.warn(err)
	return err
}

// AutoPersist enables automatic persistence of system state to given filename.
func (me *Service) AutoPersist(dest Storage, every time.Duration) {
	last := me.sys.LastModified()
	go func() {
		for {
			// todo decouple and use events
			modified := me.sys.LastModified()
			if !modified.After(last) {
				time.Sleep(every)
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
	me.Log("persist state: ", dest)
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

func (me *FileStorage) String() string { return me.filename }
