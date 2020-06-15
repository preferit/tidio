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
	Username string
}
