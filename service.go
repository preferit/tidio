package tidio

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"sync"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/nugo"
)

func NewService() *Service {
	s := &Service{
		warn:    fox.NewSyncLog(ioutil.Discard).Log,
		store:   make([]*Resource, 0),
		nextUID: 100,
		nextGID: 100,
	}
	s.NewAccount(&Account{}, "root")
	return s
}

type Service struct {
	warn func(...interface{})

	next    sync.Mutex
	store   []*Resource
	nextUID uint
	nextGID uint
}

func (me *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (me *Service) NewAccount(a *Account, username string) error {
	if username == "" {
		return fmt.Errorf("NewAccount: empty username")
	}
	a.Ring = nugo.NewRing(me.nextIDs())
	me.Add(a, path.Join("accounts", username), a)
	return nil
}

func (me *Service) Add(as *Account, path string, entity interface{}) {
	me.store = append(me.store, &Resource{
		Seal:   as.Seal(),
		path:   path,
		entity: entity,
	})
}

func (me *Service) WriteTo(w io.Writer) error {
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"resources": me.store,
	})
}

func (me *Service) nextIDs() (uid uint, gid uint) {
	me.next.Lock()
	uid = me.nextUID
	gid = me.nextGID
	me.nextUID++
	me.nextGID++
	me.next.Unlock()
	return
}
