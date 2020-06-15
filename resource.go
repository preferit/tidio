package tidio

import (
	"github.com/gregoryv/nugo"
)

type Resource struct {
	nugo.Seal
	path   string
	entity interface{}
}

type Account struct {
	nugo.Ring
}
