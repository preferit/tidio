package tidio

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

type Store struct {
	dir string
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
	return ioutil.WriteFile(path.Join(s.dir, file), data, perm)
}
