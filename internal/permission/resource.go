package permission

type Resource interface {
	UID() int
	GID() int
	Mode() PermMode
}
