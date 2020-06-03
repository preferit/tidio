package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/gregoryv/fox"
)

func NewStore(dir string) *Store {
	ch := make(chan string)
	store := &Store{
		Logger:  fox.NewSyncLog(ioutil.Discard),
		writeOp: ch, // synchronize all write operations
		dir:     dir,
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
	writeOp chan string
	dir     string
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

func (s *Store) WriteFile(file string, data []byte, perm os.FileMode) error {
	filename := path.Join(s.dir, file)
	os.MkdirAll(path.Dir(filename), 0755)
	err := ioutil.WriteFile(filename, data, perm)
	s.writeOp <- fmt.Sprintf("write %s", file)
	return err
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
