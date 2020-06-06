package permission_test

import (
	"github.com/preferit/tidio/internal/permission"
)

type Blog struct {
	permission.Secured
}

type Entry struct {
	permission.Secured
}

func Example_control_access() {
	blog := &Blog{
		Secured: &permission.Set{},
	}
	entry := &Entry{
		Secured: &permission.Set{},
	}

	actor := permission.NewActor(0, 0)

	if permission.ToCreate(blog, entry, actor) != nil {
		// failed
	}
}
