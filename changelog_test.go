package tidio

import (
	"testing"

	. "github.com/gregoryv/web"
)

func Test_generate_changelog(t *testing.T) {
	err := NewPage(
		Html(Body(NewChangelog())),
	).SaveAs("changelog.md")
	if err != nil {
		t.Fatal(err)
	}

}
