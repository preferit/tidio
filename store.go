package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/gregoryv/fox"
)

func NewStore(dir string) *Store {
	// DO NOT buffer this channel or multiple write operations will be
	// commited together.
	ch := make(chan string)
	store := &Store{
		Logger:   fox.NewSyncLog(ioutil.Discard),
		writeOps: ch, // synchronize all write operations
		dir:      dir,
	}
	go func() {
		for {
			warn(store.Commit(<-ch))
		}
	}()
	return store
}

type Store struct {
	fox.Logger
	writeOps chan string
	dir      string
}

func (s *Store) Commit(msg string) error {
	warn(exec.Command("git", "-C", s.dir, "add", ".").Run())
	return exec.Command("git", "-C", s.dir, "commit", "-m", msg).Run()
}

func (s *Store) IsInitiated() bool {
	stat, err := os.Stat(path.Join(s.dir, ".git"))
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func (s *Store) Init() error {
	return exec.Command("git", "-C", s.dir, "init").Run()
}

func (s *Store) Do(op WriteOp) error {
	filename := path.Join(s.dir, op.file)
	os.MkdirAll(path.Dir(filename), 0755)
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = io.Copy(fh, op.data)
	s.writeOps <- fmt.Sprintf("%s write %s", op.account, op.file)
	op.data.Close()
	return err
}

type WriteOp struct {
	account string
	file    string
	data    io.ReadCloser
}

func (s *Store) ReadFile(w io.Writer, file string) error {
	filename := path.Join(s.dir, file)
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = io.Copy(w, fh)
	return err
}

func (s *Store) Glob(user, pattern string) []string {
	found, err := filepath.Glob(path.Join(s.dir, user, pattern))
	warn(err)
	return found
}
