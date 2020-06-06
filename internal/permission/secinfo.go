package permission

type SecInfo struct {
	uid, gid int
	mode     PermMode
}

// SecInfo returns itself
func (g *SecInfo) SecInfo() SecInfo { return *g }
