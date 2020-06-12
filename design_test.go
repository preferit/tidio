package tidio

import (
	"testing"

	"github.com/gregoryv/draw/shape/design"
)

func Test_aquiring_role(t *testing.T) {
	var (
		d       = design.NewSequenceDiagram()
		router  = d.Add("Router")
		service = d.AddStruct(Service{})
		auth    = d.AddInterface((*Accounts)(nil))
	)
	d.Link(router, service, "RoleByKey()")
	d.Link(service, auth, "LoadAccountByKey()")
	d.SaveAs("/tmp/aquiring_role.svg")
}
