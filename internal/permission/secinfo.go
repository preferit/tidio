package permission

type Set struct {
	UID, GID uint
	Mode     PermMode
}

// SecInfo returns itself
func (g *Set) PermSet() Set { return *g }
