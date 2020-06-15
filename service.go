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

func (me *Service) Load() error {
	err := errors{
		me.Timesheets.Load(),
		me.Accounts.Load(),
	}
	return err.First()
}

func (me *Service) Save() error {
	err := errors{
		me.Timesheets.Save(),
		me.Accounts.Save(),
	}
	return err.First()
}

// SetDataDir sets directory where state is persisted
func (s *Service) SetDataDir(dir string) {
	s.Timesheets.PersistToFile(path.Join(dir, "timesheets.json"))
	s.Accounts.PersistToFile(path.Join(dir, "accounts.json"))
}

func (me *Service) AccountByKey(account *Account, key string) error {
	if err := me.FindAccountByKey(account, key); err != nil {
		return err
	}
	me.fillAccount(account)
	return nil
}

func (me *Service) fillAccount(account *Account) {
	account.Timesheets = me.Timesheets
}

func (me *Service) FindAccount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		var account Account // todo maybe default to noname
		if err := me.AccountByKey(&account, key); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "account", &account))
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
	var account Account
	if !publicResource {
		key := r.Header.Get("Authorization")
		if err := me.AccountByKey(&account, key); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
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
