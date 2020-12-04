package tidio

import (
	"io"
	"os"
	"path"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/nexus"
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
	nexus.Failure

	sys  *rs.System
	dest Storage // for persisting the system to
}

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
	if me.dest == nil || !me.Ok() {
		return me.Error()
	}
	me.Log("restoring state:", me.dest)
	r, err := me.dest.Open()
	if err != nil {
		return me.Fail(err)
	}
	defer r.Close()
	err = me.sys.Import("/", r)
	return me.Fail(err)
}

// PersistState persist the system to a configured Storage. Noop if
// not configured.
func (me *Service) PersistState() error {
	if me.dest == nil || !me.Ok() {
		return me.Error()
	}
	w, err := me.dest.Create()
	if err != nil {
		return me.Fail(err)
	}
	defer w.Close()
	me.Log("persist state: ", me.dest)
	return me.sys.Export(w)
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

func (me *FileStorage) Create() (io.WriteCloser, error) {
	return os.Create(me.filename)
}

func (me *FileStorage) Open() (io.ReadCloser, error) {
	return os.Open(me.filename)
}

func (me *FileStorage) String() string { return me.filename }
