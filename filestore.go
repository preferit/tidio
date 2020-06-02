package tidio

import (
	"fmt"
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
			if err := store.Commit(<-ch); err != nil {
				store.Log(err)
			}
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
	err := exec.Command("git", "-C", s.dir, "add", ".").Run()
	if err != nil {
		return err
	}
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
	err := ioutil.WriteFile(filename, data, perm)
	s.writeOp <- fmt.Sprintf("write %s", file)
	return err
}
