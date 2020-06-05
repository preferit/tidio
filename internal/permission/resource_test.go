package permission

import "testing"

func Test_resource(t *testing.T) {
	var e Resource
	e = &thing{}
	e.UID()
	e.GID()
	e.Mode()
}

type thing struct {
	uid, gid int
	mode     PermMode
}

func (g *thing) UID() int       { return g.uid }
func (g *thing) GID() int       { return g.gid }
func (g *thing) Mode() PermMode { return g.mode }
