package tidio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/stamp"
)

func NewService() *Service {
	return &Service{
		Resources: NewMemResources(),
		Accounts:  NewMemAccounts(),
		warn:      fox.NewSyncLog(ioutil.Discard).Log,
	}
}

type Service struct {
	Stateful
	Resources
	Accounts

	warn func(...interface{})
}

func (me *Service) Load() error {
	err := errors{
		me.Resources.Load(),
		me.Accounts.Load(),
	}
	return err.First()
}

func (me *Service) Save() error {
	err := errors{
		me.Resources.Save(),
		me.Accounts.Save(),
	}
	return err.First()
}

// SetDataDir sets directory where state is persisted
func (s *Service) SetDataDir(dir string) {
	s.Resources.PersistToFile(path.Join(dir, "timesheets.json"))
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
	account.Resources = me.Resources
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

// ----------------------------------------

func (me *Service) FindResource(resource *Resource) error {
	switch {
	case resource.Path == "":
		NewStats(resource)
	default:
		return fmt.Errorf("FindResource: %s not found", resource.Path)
	}
	return nil
}

func NewStats(resource *Resource) {
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
	resource.Seal.Mode = 04444
	resource.ReadCloser = r
}
