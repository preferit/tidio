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
	root := Account{Username: "root"}
	s.AddUser(&root)
	return s
}

type Service struct {
	warn func(...interface{})

	next    sync.Mutex
	store   []*Resource
	nextUID uint
	nextGID uint
}

func (me *Service) AddUser(account *Account) error {
	if account.Username == "" {
		return fmt.Errorf("AddUser: empty username")
	}
	account.Ring = nugo.NewRing(me.nextIDs())
	filename := fmt.Sprintf("%s.account", account.Username)
	me.AddResource(&Resource{
		Seal:   account.Seal(),
		Path:   filename,
		Entity: account,
	})
	// home dir
	me.AddResource(&Resource{
		Seal: account.Seal(),
		Path: "/home/" + account.Username,
	})
	return nil
}

func (me *Service) Mkdir(dir string, a *Account) error {
	var parent Resource
	if err := me.Find(&parent, path.Dir(dir)); err != nil {
		return err
	}
	perm := nugo.Permission{}

	r := Resource{
		Seal: a.Seal(),
		Path: dir,
	}
	if err := perm.ToCreate(&parent, &r, a); err != nil {
		return err
	}
	me.AddResource(&r)
	return nil
}

// AddResource adds a resource with no premission control
func (me *Service) AddResource(r *Resource) {
	me.store = append(me.store, r)
}

func (me *Service) Find(resource *Resource, filename string) error {
	for _, r := range me.store {
		if r.Path == filename {
			*resource = *r
			return nil
		}
	}
	return fmt.Errorf("Find: %s not found", filename)
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

// ----------------------------------------

func (me *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
