package tidio

import (
	"github.com/gregoryv/nugo"
)

type Resource struct {
	nugo.Seal
	path   string
	Entity interface{}
}

type Account struct {
	nugo.Ring
}

func (me *Account) NewResource(path string, entity interface{}) *Resource {
	return &Resource{
		Seal:   me.Seal(),
		path:   path,
		Entity: entity,
	}
}
