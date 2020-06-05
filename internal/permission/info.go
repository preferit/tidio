package permission

type Info struct {
	uid, gid int
	mode     PermMode
}

func (g *Info) UID() int       { return g.uid }
func (g *Info) GID() int       { return g.gid }
func (g *Info) Mode() PermMode { return g.mode }
