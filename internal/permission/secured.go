package permission

type Secured interface {
	PermSet() Set
}
