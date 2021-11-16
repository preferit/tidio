// Package tidio provides a web System for timesheet reports.
package tidio

import (
	"io"
	"os"
	"path"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/nexus"
	"github.com/gregoryv/rs"
)

func NewSystem(settings ...ant.Setting) *System {
	sys := rs.NewSystem()
	asRoot := rs.Root.Use(sys)
	asRoot.Exec("/bin/mkdir /etc/basic")
	asRoot.Exec("/bin/mkdir /api")
	asRoot.Exec("/bin/mkdir /api/timesheets")

	srv := &System{
		sys: sys,
	}
	Register(srv)
	ant.MustConfigure(srv, settings...)
	return srv
}

type System struct {
	nexus.Failure

	sys  *rs.System
	dest Storage // for persisting the system to
}

func (me *System) UseFileStorage(filename string) error {
	me.dest = NewFileStorage(filename)

	_, err := os.Stat(filename)
	switch {
	case os.IsNotExist(err):
		return me.PersistState()
	default:
		return me.RestoreState()
	}
	return me.Error()
}

// AddAccount creates a system account and stores the secret in
// /etc/basic
func (me *System) AddAccount(name, secret string) error {
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
func (me *System) RestoreState() error {
	if me.dest == nil || !me.Ok() {
		return me.Error()
	}
	Log(me).Info("restoring state:", me.dest)
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
func (me *System) PersistState() error {
	if me.dest == nil || !me.Ok() {
		return me.Error()
	}
	w, err := me.dest.Create()
	if err != nil {
		return me.Fail(err)
	}
	defer w.Close()
	Log(me).Info("persist state: ", me.dest)
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

// ----------------------------------------

// NewShell returns a new shell. Once an error occurs the shell no
// longer functions, similar to bash -e flag.
func NewShell(acc *rs.Syscall) *Shell {
	return &Shell{account: acc}
}

type Shell struct {
	account *rs.Syscall
	err     error
}

// Execf
func (me *Shell) Execf(format string, args ...interface{}) error {
	if me.err != nil {
		return me.err
	}
	me.err = me.account.Execf(format, args...)
	return me.err
}

// ----------------------------------------

// Credentials provides ways to authenticate a requests via header
// manipulation. Zero value credentials means anonymous.
func NewCredentials(account, secret string) *Credentials {
	return &Credentials{
		account: account,
		secret:  secret,
	}
}

type Credentials struct {
	account string
	secret  string
}

func (me *Credentials) Set(v interface{}) error {
	switch v := v.(type) {
	case usesCredentials:
		v.SetCredentials(me)
	case *System:
		v.AddAccount(me.account, me.secret)
	default:
		return ant.SetFailed(v, me)
	}
	return nil
}

type usesCredentials interface {
	SetCredentials(*Credentials)
}
