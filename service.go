package tidio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/gregoryv/nugo"
	"github.com/gregoryv/stamp"
)

func NewService() *Service {
	return &Service{
		Timesheets: NewMemSheets(),
		Accounts:   NewMemAccounts(),
	}
}

type Service struct {
	Stateful
	Timesheets
	Accounts
}

func (s *Service) Load() error {
	err := errors{
		s.Timesheets.Load(),
		s.Accounts.Load(),
	}
	return err.First()
}

func (s *Service) Save() error {
	err := errors{
		s.Timesheets.Save(),
		s.Accounts.Save(),
	}
	return err.First()
}

// SetDataDir sets directory where state is persisted
func (s *Service) SetDataDir(dir string) {
	s.Timesheets.PersistToFile(path.Join(dir, "timesheets.json"))
	s.Accounts.PersistToFile(path.Join(dir, "accounts.json"))
}

func (s *Service) AccountByKey(key string) (*Account, bool) {
	if key == "" {
		return nil, false
	}
	var account Account
	if err := s.FindAccountByKey(&account, key); err != nil {
		return nil, false
	}
	account.Timesheets = s.Timesheets
	return &account, true
}

func (me *Service) FindAccount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		account, ok := me.AccountByKey(key)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "account", account))
		next.ServeHTTP(w, r)
	})
}

func (me *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// only serve /api requests
	path := r.URL.Path
	prefix := "/api"
	if startsWith(path, prefix) {
		path = path[len(prefix):]
	}
	resource := Resource{
		Path: path,
	}
	if err := me.FindResource(&resource); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	publicResource := true
	if !publicResource {
		key := r.Header.Get("Authorization")
		account, found := me.AccountByKey(key)
		if !found {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_ = account
	}
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resource.Content)
	resource.Content.Close()

}

func (me *Service) FindResource(resource *Resource) error {
	switch {
	case resource.Path == "":
		r, w := io.Pipe()
		go func() {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"revision": stamp.InUse().Revision,
				"version":  stamp.InUse().ChangelogVersion,
				"resources": []string{
					"/api/timesheets/",
				},
			})
			w.Close()
		}()
		resource.Content = r
		return nil
	}
	return fmt.Errorf("todo")
}

func startsWith(s, prefix string) bool {
	return strings.Index(s, prefix) == 0
}

func (me *Service) WriteResource(r *Resource, user *Account) error {
	return nil
}

// ----------------------------------------

type Resource struct {
	nugo.Seal
	Path    string
	Content io.ReadCloser
}
