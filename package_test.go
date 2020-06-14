package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/draw/shape/design"
	"github.com/preferit/tidio/internal"
)

func Test_stateful(t *testing.T) {
	var (
		assert = asserter.New(t)
		_, bad = assert().Errors()
		e      = internal.None("")
	)
	_, err := e.Open()
	bad(err)
	_, err = e.Create()
	bad(err)
}

func Test_routing_request(t *testing.T) {
	var (
		d       = design.NewSequenceDiagram()
		router  = d.Add("Router")
		service = d.AddStruct(Service{})
		auth    = d.AddInterface((*Accounts)(nil))
	)
	d.Link(router, service, "RoleByKey()")
	d.Link(service, auth, "FindAccountByKey()")
	d.SaveAs("/tmp/aquiring_role.svg")
}

type BrokenDestination struct{}

func (BrokenDestination) Create() (io.WriteCloser, error) {
	return nil, fmt.Errorf("broken destination")
}

type BrokenSource struct{}

func (BrokenSource) Open() (io.ReadCloser, error) {
	return nil, fmt.Errorf("broken source")
}

type StringSource string

func (me StringSource) Open() (io.ReadCloser, error) {
	return ioutil.NopCloser(strings.NewReader(string(me))), nil
}

// ----------------------------------------

func Nowhere() *nowhere {
	return &nowhere{}
}

type nowhere struct{}

func (me *nowhere) Create() (io.WriteCloser, error) {
	return me, nil
}

func (*nowhere) Write(b []byte) (int, error) {
	return len(b), nil
}

func (*nowhere) Close() error { return nil }
