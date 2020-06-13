package tidio

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_stateful(t *testing.T) {
	var (
		assert = asserter.New(t)
		_, bad = assert().Errors()
		e      = None("")
	)
	_, err := e.Open()
	bad(err)
	_, err = e.Create()
	bad(err)
}
