package tidio

import (
	"github.com/gregoryv/nugo"
)

type Resource struct {
	nugo.Seal
	Path   string
	Entity interface{}
}

type Account struct {
	nugo.Ring
}

func (me *Account) NewResource(path string, entity interface{}) *Resource {
	return &Resource{
		Seal:   me.Seal(),
		Path:   path,
		Entity: entity,
	}
}

type Folder struct {
	nugo.Seal
	Path string
}
