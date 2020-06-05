package perm

type Resource interface {
	UID() int
	GID() int
	Mode() PermMode
}
