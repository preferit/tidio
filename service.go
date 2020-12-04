package tidio

import (
	"io"
	"os"
	"path"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/rs"
)

func NewService(settings ...ant.Setting) *Service {
	sys := rs.NewSystem()
	asRoot := rs.Root.Use(sys)
	asRoot.Exec("/bin/mkdir /etc/basic")
	asRoot.Exec("/bin/mkdir /api")
	asRoot.Exec("/bin/mkdir /api/timesheets")

	srv := Service{
		sys: sys,
	}
	ant.MustConfigure(&srv, settings...)
	return &srv
}

type Service struct {
	OptionalLogger

	sys  *rs.System
	dest Storage // for persisting the system to

	err error // Set if in unrecoverable state, ie. during setup. Use Error() to read.
}

// Error
func (me *Service) Error() error { return me.err }

// AddAccount creates a system account and stores the secret in
// /etc/basic
func (me *Service) AddAccount(name, secret string) error {
	sh := NewShell(rs.Root.Use(me.sys))
	return CreateAccount(sh, name, secret)
}

func CreateAccount(sh *Shell, name, secret string) error {
	sh.Execf("/bin/mkacc %s", name)
	sh.Execf("/bin/secure -a %s -s %s", name, secret)

	// todo, better define system with /api/home/NAME
	dir := path.Join("/api/timesheets", name)
	sh.Execf("/bin/mkdir %s", dir)

	return sh.Execf("/bin/chown %s %s", name, dir)
}

// RestoreState restores the resource system from the given filename.
func (me *Service) RestoreState() error {
	if me.dest == nil || me.err != nil {
		return me.err
	}
	me.Log("restoring state:", me.dest)
	r, err := me.dest.Open()
	if err != nil {
		return me.failed(err)
	}
	defer r.Close()
	err = me.sys.Import("/", r)
	return me.failed(err)
}

// PersistState persist the system to a configured Storage. Noop if not configured.
func (me *Service) PersistState() error {
	if me.dest == nil || me.err != nil {
		return me.err
	}
	me.Log("persist state: ", me.dest)
	w, err := me.dest.Create()
	if err != nil {
		return me.failed(err)
	}
	defer w.Close()
	return me.sys.Export(w)
}

// failed sets the err field unless already set.
func (me *Service) failed(err error) error {
	if me.err != nil {
		return me.err
	}
	me.err = err
	return err
}

// ----------------------------------------

type Storage interface {
	Create() (io.WriteCloser, error)
	Open() (io.ReadCloser, error)
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

// Open
func (me *FileStorage) Open() (io.ReadCloser, error) {
	return os.Open(me.filename)
}

func (me *FileStorage) String() string { return me.filename }
