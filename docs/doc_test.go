package docs

import (
	"testing"
)

func Test_generate_index(t *testing.T) {
	page := NewIndex()
	page.SaveAs("index.html")

}
